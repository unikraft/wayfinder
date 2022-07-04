package tester

// SPDX-License-Identifier: BSD-3-Clause
//
// Authors: Alexander Jung <alex@unikraft.io>
//
// Copyright (c) 2021, Unikraft UG.  All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/unikraft/wayfinder/api/proto"
	"github.com/unikraft/wayfinder/internal/models"
	"github.com/unikraft/wayfinder/internal/strutils"
	"github.com/unikraft/wayfinder/modules/container"
	"github.com/unikraft/wayfinder/modules/libvirt"
	"github.com/unikraft/wayfinder/pkg/common/errors"
)

type test struct {
	sync.RWMutex
	domain     *libvirt.Domain
	startDelay time.Duration
	container  *container.Container
	err        error
	runtime    time.Duration
}

var (
	pidCount = 0
	// This variable contains a map of fake PIDs and a corresponding UUID for a
	// test and can be used to uniquely assign a numerical value to UUID.
	pidUuidMap = make(map[int]string, 0)
)

type Service struct {
	p     *provider
	tests map[string]*test
}

func (s *Service) CreateTest(ctx context.Context, req *proto.CreateTestRequest) (*proto.CreateTestResponse, error) {
	if req.PermutationId <= 0 {
		return &proto.CreateTestResponse{
			Success: false,
		}, errors.NewMissingParameterError("permutation_id")
	}

	// TODO: Auto-populate kernel + benchtool from config if not provided.
	if req.Kernel == nil {
		return &proto.CreateTestResponse{
			Success: false,
		}, errors.NewMissingParameterError("kernel")
	}

	if req.Kernel.Image == "" {
		return &proto.CreateTestResponse{
			Success: false,
		}, errors.NewMissingParameterError("kernel.image")
	}

	if len(req.Kernel.Cores) <= 0 {
		return &proto.CreateTestResponse{
			Success: false,
		}, errors.NewMissingParameterError("kernel.cores")
	}

	if req.Kernel.Memory == "" {
		return &proto.CreateTestResponse{
			Success: false,
		}, errors.NewMissingParameterError("kernel.memory")
	}

	if req.BenchTool == nil {
		return &proto.CreateTestResponse{
			Success: false,
		}, errors.NewMissingParameterError("benchtool")
	}

	if req.BenchTool.Image == "" {
		return &proto.CreateTestResponse{
			Success: false,
		}, errors.NewMissingParameterError("benchtool.image")
	}

	if req.BenchTool.Commands == "" {
		return &proto.CreateTestResponse{
			Success: false,
		}, errors.NewMissingParameterError("benchtool.commands")
	}

	if len(req.BenchTool.Cores) == 0 {
		return &proto.CreateTestResponse{
			Success: false,
		}, errors.NewMissingParameterError("benchtool.cores")
	}

	// Create a new entry in the database for this test
	testModel, err := s.p.DB.Repos().Tests().CreateTestForPermutation(&models.Test{
		PermutationId:  uint(req.PermutationId),
		VMMCores:       strutils.JoinUint64(req.VmmCores, ","),
		KernelCores:    strutils.JoinUint64(req.Kernel.Cores, ","),
		BenchToolCores: strutils.JoinUint64(req.BenchTool.Cores, ","),
	})
	if err != nil {
		return &proto.CreateTestResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "could not create test in database: %s", err)
	}

	uuid := testModel.UUID.String()

	if _, ok := s.tests[uuid]; ok {
		return &proto.CreateTestResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "test with uuid=%s already exists", uuid)
	}

	// Create a new fake PID to be used for the bridge and for counting the
	// domains
	pidCount += 1
	pid := pidCount%4 + 1

	pidUuidMap[pid] = uuid

	// Regex Match all values in order for memory
	r, _ := regexp.Compile("^[1-9]+[0-9]*")
	memoryValue := r.FindStringSubmatch(req.Kernel.Memory)[0]
	memoryUnit := "MiB"

	// Regex to match unit type, we ignore difference between 'b' and 'B'
	r, _ = regexp.Compile("(MiB|M|MB|Mib|Mb)$")
	if r.MatchString(req.Kernel.Memory) {
		memoryUnit = "MiB"
	}

	r, _ = regexp.Compile("(GiB|G|GB|Gib|Gb)$")
	if r.MatchString(req.Kernel.Memory) {
		memoryUnit = "GiB"
	}

	// Convert memoryValue to uint
	memoryValueUint, _ := strconv.ParseUint(memoryValue, 10, 32)

	domain, err := s.p.Libvirt.NewDomain(
		pid,
		uuid,
		req.Kernel.Image,
		req.Kernel.InitRd,
		req.Kernel.Args,
		req.Kernel.Disks,
		req.Kernel.Cores,
		uint(memoryValueUint),
		memoryUnit,
		req.BenchTool.Monitors,
	)
	if err != nil {
		s.p.DB.Repos().Tests().SetStatusKernelFailedStartupByTestUuid(uuid)
		return &proto.CreateTestResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "cannot create test VM: %s", err)
	}

	// Initialize the domain
	if err = domain.Init(); err != nil {
		s.p.DB.Repos().Tests().SetStatusKernelFailedStartupByTestUuid(uuid)
		return &proto.CreateTestResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "could not get initialise domain: %s", err)
	}

	if len(req.VmmCores) > 0 {
		if err = domain.PinVMMToCores(req.VmmCores); err != nil {
			s.p.DB.Repos().Tests().SetStatusWayfinderFailedInternal(uuid)
			return &proto.CreateTestResponse{
				Success: false,
			}, status.Errorf(codes.Internal, "could not pin VMM to cores: %s", err)
		}
	}

	// Initialize the benchmark tool's container
	container, err := s.p.Container.NewContainer(uuid)
	if err != nil {
		s.p.DB.Repos().Tests().SetStatusBenchToolFailedStartupByTestUuid(uuid)
		return &proto.CreateTestResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "cannot create test container: %s", err)
	}

	if err := container.PullAndAttachImage(req.BenchTool.Image); err != nil {
		s.p.DB.Repos().Tests().SetStatusBenchToolFailedStartupByTestUuid(uuid)
		return &proto.CreateTestResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "cannot pull test container's image: %s", err)
	}

	if err := container.SetCommands(req.BenchTool.Commands); err != nil {
		s.p.DB.Repos().Tests().SetStatusBenchToolFailedStartupByTestUuid(uuid)
		return &proto.CreateTestResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "cannot set test container's command: %s", err)
	}

	var envVars []string

	if domain.IP() == nil {
		s.p.DB.Repos().Tests().SetStatusKernelFailedNetworkByTestUuid(uuid)
		return &proto.CreateTestResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "could not get domain IP: %s", err)
	}

	envVars = append(envVars, fmt.Sprintf("WAYFINDER_DOMAIN_IP_ADDR=%s", domain.IP().String()))

	if len(req.BenchTool.EnvVars) > 0 {
		for _, env := range req.BenchTool.EnvVars {
			envVars = append(envVars, fmt.Sprintf("%s=%s", env.Name, env.Value))
		}
	}

	envVars = append(envVars, fmt.Sprintf("WAYFINDER_TOTAL_CORES=%d", len(req.BenchTool.Cores)))

	benchCores := ""
	for i, core := range req.BenchTool.Cores {
		if i == 0 {
			benchCores = fmt.Sprintf("%d", core)
		} else {
			benchCores = fmt.Sprintf("%s,%d", benchCores, core)
		}
	}
	envVars = append(envVars, fmt.Sprintf("WAYFINDER_CORE_ID0=%s", benchCores))

	container.AddEnvVars(envVars)
	container.SetCores(req.BenchTool.Cores)

	if len(req.BenchTool.Devices) > 0 {
		container.SetDevices(req.BenchTool.Devices)
	}

	if len(req.BenchTool.Capabilities) > 0 {
		container.SetCapabilities(req.BenchTool.Capabilities)
	}

	startDelay := s.p.Cfg.DefaultStartDelay
	if req.BenchTool.StartDelay > 0 {
		// TODO: Update spec to use time.Duration
		startDelay = time.Second * time.Duration(req.BenchTool.StartDelay)
	}

	// Save the build
	s.tests[uuid] = &test{
		domain:     domain,
		container:  container,
		startDelay: startDelay,
	}

	return &proto.CreateTestResponse{
		Success: true,
		Uuid:    uuid,
	}, nil
}

func (s *Service) StartTest(ctx context.Context, req *proto.StartTestRequest) (*proto.StartTestResponse, error) {
	if req.Uuid == "" {
		s.p.DB.Repos().Tests().SetStatusWayfinderFailedInternal(req.Uuid)
		return &proto.StartTestResponse{
			Success: false,
		}, errors.NewMissingParameterError("uuid")
	}

	// TODO: Look up tests with jobId/permId rather than via container ID (this
	// is for the proto/API).

	test, ok := s.tests[req.Uuid]
	if !ok {
		s.p.DB.Repos().Tests().SetStatusWayfinderFailedInternal(req.Uuid)
		return &proto.StartTestResponse{
			Success: false,
		}, status.Errorf(codes.NotFound, "cannot find test with id=%s", req.Uuid)
	}

	err := test.container.Init()
	if err != nil {
		test.err = err
		s.p.DB.Repos().Tests().SetStatusBenchToolFailedStartupByTestUuid(req.Uuid)

		return &proto.StartTestResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "could not initialize test container: %s", err)
	}

	s.p.DB.Repos().Tests().SetStatusRunningByTestUuid(req.Uuid)

	// Start a thread which oversees the running VM domain
	go func() {
		// Start the domain
		err := test.domain.Start()
		if err != nil {
			test.err = err
			s.p.DB.Repos().Tests().SetStatusKernelFailedStartupByTestUuid(req.Uuid)
		}
	}()

	collectData := make(chan bool)

	// Start measuring the domain's resources, after the start delay
	// NOTE: The benchtool may have a warm up time too, the resources should
	// only be started to measure after this warm up time.  How do we calculate
	// the warm up time? Is it a spec option?
	go func() {
		pushMetrics := <-collectData

		for pushMetrics {
			time.Sleep(s.p.Cfg.MetricsFreq)
			if err := s.p.PushMetrics(req.Uuid, req.JobId, test.domain.GetResourceMeasurements()); err != nil {
				s.p.Log.Errorf("could not push metrics: %s", err)
			}
			select {
			case pushMetrics, ok = <-collectData:
				if !ok {
					pushMetrics = false
				}
			default:
				continue
			}
		}
	}()

	// Start a threat which oversees the running benchmark tool's container
	go func() {
		// Give the container a delay in its startup, so we can wait for kernel
		// to boot up correctly
		s.p.Log.Debugf("Waiting %d seconds before starting benchtool...", test.startDelay/1000000000)
		time.Sleep(time.Duration(test.startDelay))

		// The other thread can start pushing metrics now
		collectData <- true

		// Now start the container.  On return, the benchtool has exited.
		runtime, err := test.container.StartAndWait()
		collectData <- false

		test.RLock()
		defer test.RUnlock()

		// Save the total runtime of the benchtool and thus the experiment.
		test.runtime = runtime
		s.p.DB.Repos().Tests().SetRuntimeByTestUuid(req.Uuid, runtime)

		if err != nil {
			test.err = err
			s.p.DB.Repos().Tests().SetStatusBenchToolFailedByTestUuid(req.Uuid)
		} else {
			s.p.DB.Repos().Tests().SetStatusSuccessByTestUuid(req.Uuid)

			// Save all the known results from the test
			if len(req.Results) > 0 {
				for _, result := range req.Results {
					switch result.Type {
					case proto.TestResultType_TEST_RESULT_INT:
						value, err := test.container.GetResultInt(result.Path)
						if err != nil {
							s.p.Log.Errorf("could not get test result value: %s", err)
							continue
						}

						_, err = s.p.DB.Repos().Results().SaveResultIntByTestUuid(req.Uuid, result.Name, value)
						if err != nil {
							s.p.Log.Errorf("could not save test result: %s", err)
						}

					case proto.TestResultType_TEST_RESULT_STR:
						value, err := test.container.GetResultStr(result.Path)
						if err != nil {
							s.p.Log.Errorf("could not get test result value: %s", err)
							continue
						}

						_, err = s.p.DB.Repos().Results().SaveResultStrByTestUuid(req.Uuid, result.Name, value)
						if err != nil {
							s.p.Log.Errorf("could not save test result: %s", err)
						}

					case proto.TestResultType_TEST_RESULT_FLOAT:
						value, err := test.container.GetResultFloat(result.Path)
						if err != nil {
							s.p.Log.Errorf("could not get test result value: %s", err)
							continue
						}

						_, err = s.p.DB.Repos().Results().SaveResultFloatByTestUuid(req.Uuid, result.Name, value)
						if err != nil {
							s.p.Log.Errorf("could not save test result: %s", err)
						}

					case proto.TestResultType_TEST_RESULT_BOOL:
						value, err := test.container.GetResultBool(result.Path)
						if err != nil {
							s.p.Log.Errorf("could not get test result value: %s", err)
							continue
						}

						_, err = s.p.DB.Repos().Results().SaveResultBoolByTestUuid(req.Uuid, result.Name, value)
						if err != nil {
							s.p.Log.Errorf("could not save test result: %s", err)
						}
					}
				}
			}
		}
	}()

	return &proto.StartTestResponse{
		Success: true,
	}, nil
}

func (s *Service) GetTestStatus(ctx context.Context, req *proto.GetTestStatusRequest) (*proto.GetTestStatusResponse, error) {
	if req.Uuid == "" {
		return &proto.GetTestStatusResponse{
			Success: false,
		}, errors.NewMissingParameterError("uuid")
	}

	testModel := &models.Test{}
	if err := s.p.DB.DB().Where("uuid = ?", req.Uuid).First(&testModel).Error; err != nil {
		return &proto.GetTestStatusResponse{
			Success: false,
		}, status.Errorf(codes.NotFound, "test with uuid=%s not found", req.Uuid)
	}

	// TODO: Look up tests with jobId/permId rather than via test ID (this
	// is for the proto/API).

	test, ok := s.tests[req.Uuid]
	if !ok {
		return &proto.GetTestStatusResponse{
			Success: false,
		}, status.Errorf(codes.NotFound, "cannot find test with uuid=%s", req.Uuid)
	}

	return &proto.GetTestStatusResponse{
		Success: true,
		Status:  testModel.Status,
		Runtime: test.runtime.String(),
	}, nil
}

func (s *Service) DestroyTest(ctx context.Context, req *proto.DestroyTestRequest) (*proto.DestroyTestResponse, error) {
	if req.Uuid == "" {
		return &proto.DestroyTestResponse{
			Success: false,
		}, errors.NewMissingParameterError("uuid")
	}

	// TODO: Look up builds with jobId/permId rather than via container ID (this
	// is for the proto/API).

	test, ok := s.tests[req.Uuid]
	if !ok {
		return &proto.DestroyTestResponse{
			Success: false,
		}, status.Errorf(codes.NotFound, "cannot find test with uuid=%s", req.Uuid)
	}

	// TODO: Actually we should try and garbage collect this later too.

	// var err error
	// // Try 3 times to kill domain
	// for i := 3; i > 0; i-- {
	//   fmt.Printf("Trying to kill domain...\n")

	//   time.Sleep(1 * time.Second)
	// }

	if err := test.domain.Destroy(); err != nil {
		return &proto.DestroyTestResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "cannot destroy test kernel domain: %s", err)
	}

	if err := test.container.Destroy(); err != nil {
		return &proto.DestroyTestResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "cannot destroy test benchtool container: %s", err)
	}

	// wipe the test from memory
	delete(pidUuidMap, test.domain.FakePid())
	delete(s.tests, req.Uuid)

	return &proto.DestroyTestResponse{
		Success: true,
	}, nil
}
