package scheduler
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
  "fmt"
  "time"
  "context"
  "encoding/json"
  "compress/zlib"
  "bytes"
  "io"

  "github.com/adjust/rmq/v4"
  
  "github.com/unikraft/wayfinder/spec"
  "github.com/unikraft/wayfinder/api/proto"
  "github.com/erda-project/erda-infra/base/logs"
  "github.com/unikraft/wayfinder/internal/coremap"
  "github.com/unikraft/wayfinder/internal/strutils"
)

type TaskConsumer struct {
  p       *provider
  Log      logs.Logger
}

const (
  stageBuild = iota
  stageTest
)

type stage int

type build struct {
  uuid string
}

type test struct {
  uuid string
}

func NewTaskConsumer(p *provider) *TaskConsumer {
  return &TaskConsumer{p: p}
}

func (c *TaskConsumer) Consume(delivery rmq.Delivery) {
  taskBytes := delivery.Payload()
  task := spec.JobSpec{}

  var compressed bytes.Buffer
  compressed = *bytes.NewBuffer([]byte(taskBytes))
  var decompressed bytes.Buffer
  r, _ := zlib.NewReader(&compressed)
  io.Copy(&decompressed, r)
  r.Close()

  // Check if we received the full job specification
  err := json.Unmarshal(decompressed.Bytes(), &task)
  if err != nil {
    c.p.Log.Errorf("could not unmarshal job: %s", err)
    if err := delivery.Reject(); err != nil {
      c.p.Log.Errorf("failed to reject job: %s", err)
    }
  }

  // Create a new logger for this task
  c.Log = c.p.Log.Sub(task.CurrentPerm.Checksum)

  // Start the task
  err = c.StartTask(&task)

  if err != nil {
    c.Log.Errorf("could not complete permutation: %s", err)
    if err = delivery.Reject(); err != nil {
      c.Log.Errorf("failed to reject permutation: %s", err)
    }
  } else {
    delivery.Ack()
  }
}

// Busy-waits to release all cores. This happens because sometimes the
// containers don't have time to exit, after signaling that work is done.
// Most of the time this should complete in a single pass.
func (c *TaskConsumer) releaseCoresById(coresToFree []uint64) error {
  c.Log.Debugf("releasing cores: %s", strutils.JoinUint64(coresToFree, ","))

  for _, coreId := range coresToFree {
    retries := 0
    freecore:for {
      err := c.p.CoreMap().ReleaseCore(coreId)
      if err != nil {
        if retries > 10 {
          return fmt.Errorf("could not release core with id=%d", coreId)
        }

        retries++
        time.Sleep(1 * time.Second)
        continue
      }

      break freecore
    }
  }
  return nil
}

// Returns the restriction level asked for when it corresponds to the isolated stage
// Otherwise, reserve cores with no restriction
func (c *TaskConsumer) calculateCoremap(taskStage stage, level proto.JobIsolLevel,
                        split proto.JobIsolSplit) coremap.CoreRestriction {
  if (split == proto.JobIsolSplit_JOB_ISOL_SPLIT_BOTH) ||
    (taskStage == stageBuild && split == proto.JobIsolSplit_JOB_ISOL_SPLIT_BUILDS) ||
    (taskStage == stageTest && split == proto.JobIsolSplit_JOB_ISOL_SPLIT_TESTS) {
    return coremap.CoreRestriction(level)
  }

  return coremap.CoreOptionNoRestriction
}

// Busy-waits to reserve a core.  This method is required for builds and tests
// which have allocated core requirements.  An arbitrary number of cores which
// are required can be entered as input, as well as a pointer to an interface
// which will be assigned to the core, and the list of core IDs which are then
// reserved will be returned.
func (c *TaskConsumer) busyWaitForCores(requiredNumCores int, activity interface{}, taskStage stage,
                        level proto.JobIsolLevel, split proto.JobIsolSplit) ([]uint64) {
  var buildCoreIds []uint64
  var buildCores []*coremap.Core
  
  // Wait until we have some free cores.
  for {

    c.Log.Debugf("Waiting for %d cores...", requiredNumCores)

    restriction := c.calculateCoremap(taskStage, level, split)
    cores := c.p.CoreMap().FindFreeCores(restriction)
    for _, core := range cores {
      // Immediately reserve this core
      if err := c.p.CoreMap().SetCoreActivity(core.Id(), activity); err != nil {
        // To see core mapping during allocation, use:
        // c.p.CoreMap().Print()
        continue
      }

      // To see core mapping during allocation, use:
      // c.p.CoreMap().Print()

      buildCoreIds = append(buildCoreIds, core.Id())
      buildCores = append(buildCores, core)

      // Also check here in case we received more cores than requested
      if len(buildCores) >= requiredNumCores {
        break
      }
    }

    if len(buildCores) >= requiredNumCores {
      break
    } else {
      time.Sleep(c.p.Cfg.GraceTime)
    }
  }

  c.Log.Debugf("Reserved cores: %s", strutils.JoinUint64(buildCoreIds, ", "))

  return buildCoreIds;
}

// Formats the environment variables for the test
// The duration is added to the environment variables
func (c *TaskConsumer) packEnvVars(task *spec.JobSpec) []*proto.TestEnvVar {
  var envVars []*proto.TestEnvVar

  var duration string
  if task.Test.BenchTool.Duration == 0 {
    duration = "30"
  } else {
    duration = fmt.Sprint(task.Test.BenchTool.Duration)
  }


  envVars = append(envVars, &proto.TestEnvVar{
    Name: "DURATION",
    Value: duration,
  })

  for name, value := range task.Test.BenchTool.Environment {
    envVars = append(envVars, &proto.TestEnvVar{
      Name: name,
      Value: value,
    })
  }

  return envVars
}

func (c *TaskConsumer) StartTask(task *spec.JobSpec) error {
  // TODO: Implement a *real* scheduler.  For now this consumer accepts any task
  // it receives and attempts to complete it.  This method should essentially
  // determine which core/socket to schedule the task on based on the request
  // as well as "intelligently" reject permutations.

  // Find or create this permutation in the database
  permutation, err := c.p.DB.Repos().Permutations().FindOrCreateFromJobSpec(task)
  if err != nil {
    return err
  }

  // if req.JobId <= 0 {
  //   return &proto.CreateBuildResponse{
  //     Success: false,
  //     Status:  proto.BuildStatus_BUILD_FAILED,
  //   }, errors.NewMissingParameterError("job_id")
  // }

  // if req.PermutationId <= 0 {
  //   return &proto.CreateBuildResponse{
  //     Success: false,
  //     Status:  proto.BuildStatus_BUILD_FAILED,
  //   }, errors.NewMissingParameterError("permutation_id")
  // }

  // job := &models.Job{}
  // if err := s.p.DB.Repos().Job().FindJob(req.JobId, job); err != nil {
  //   return &proto.CreateBuildResponse{
  //     Success: false,
  //     Status:  proto.BuildStatus_BUILD_FAILED,
  //   }, status.Errorf(codes.NotFound, "job with id=%d not found", req.JobId)
  // }

  // permutation := &models.Permutation{}
  // if err := c.p.DB.DB().Where("id = ?", perm.Id).First(&permutation).Error; err != nil {
  //   return fmt.Errorf("permutation with id=%d not found", perm.Id)
  // }

  build := build{}

  buildCoreIds := c.busyWaitForCores(int(task.Build.Cores), &build, stageBuild, task.IsolLevel, task.IsolSplit)

  var buildEnvVars []*proto.BuildEnvVar
  for _, param := range task.CurrentPerm.Params {
    buildEnvVars = append(buildEnvVars, &proto.BuildEnvVar{
      Name: param.Name,
      Value: param.Value,
    })
  }

  // Create the build
  c.Log.Infof("creating build container for permutation_id=%d", permutation.Id)
  createBuildResp, err := c.p.Builder.CreateBuild(context.TODO(), &proto.CreateBuildRequest{
    PermutationId: int64(permutation.Id),
    Image:         task.Build.Image,
    Devices:       task.Build.Devices,
    Capabilities:  task.Build.Capabilities,
    Cores:         buildCoreIds,
    Workdir:       task.Build.Workdir,
    Commands:      task.Build.Commands,
    EnvVars:       buildEnvVars,
  })
  if err != nil {
    c.releaseCoresById(buildCoreIds)
    return fmt.Errorf("could not create build: %s", err)
  }

  build.uuid = createBuildResp.Uuid

  // Start the build
  c.Log.Infof("starting build container for permutation_id=%d", permutation.Id)
  _, err = c.p.Builder.StartBuild(context.TODO(), &proto.StartBuildRequest{
    Uuid: build.uuid,
  })
  if err != nil {
    c.releaseCoresById(buildCoreIds)
    return fmt.Errorf("could not start build: %s", err)
  }

  // Wait for the build to complete
  c.Log.Infof("waiting for build container for permutation_id=%d to complete...", permutation.Id)
  var numRetries uint64 = 0;
  buildstatus:for {
    statusBuildResp, err := c.p.Builder.GetBuildStatus(context.TODO(), &proto.GetBuildStatusRequest{
      Uuid: build.uuid,
    })

    // Retry 5 times waiting a second each, fail if no response
    if err != nil {
      numRetries++
      if numRetries >= 5 {
        c.releaseCoresById(buildCoreIds)
        return fmt.Errorf("could not get build status: %s", err)
      }
      time.Sleep(time.Second)
      continue
    }

    switch statusBuildResp.Status {
      case proto.BuildStatus_BUILD_SUCCESS,
           proto.BuildStatus_BUILD_KILLED,
           proto.BuildStatus_BUILD_FAILED:
        
        _, err = time.ParseDuration(statusBuildResp.Runtime)
        if err != nil {
          c.releaseCoresById(buildCoreIds)
          return fmt.Errorf("could not convert test runtime into duration: %s", err)
        }

        break buildstatus;
    }

    time.Sleep(c.p.Cfg.GraceTime)
  }

  c.Log.Infof("build container for permutation_id=%d exited!", permutation.Id)

  disks := []*proto.BuildOutputDiskImage{}

  // TODO: This is where spec/proto map, this should be done in proto or in spec
  // and not here.  This, or spec should be derived from proto.  The problem is
  // named enums.
  for _, disk := range task.Build.Outputs.Disks {
    outputDisk := &proto.BuildOutputDiskImage{
      Path: disk.Path,
      Name: disk.Name,
    }
    switch disk.Type {
    case spec.OutputDiskImageTypeRaw:
      outputDisk.Type = proto.BuildOutputDiskImageType_BUILD_OUTPUT_DISK_RAW
    }

    disks = append(disks, outputDisk)
  }

  // Save the outputs from the build
  c.Log.Infof("saving outputs from permutation_id=%d", permutation.Id)
  saveBuildOutputsResp, err := c.p.Builder.SaveBuildOutputsToDisk(context.TODO(), &proto.SaveBuildOutputsToDiskRequest{
    Uuid:       build.uuid,
    Outputs:   &proto.BuildOutputs{
      Kernel:   task.Build.Outputs.Kernel,
      InitRd:   task.Build.Outputs.InitRd,
      Disks:    disks,
    },
  })
  if err != nil {
    c.releaseCoresById(buildCoreIds)
    return fmt.Errorf("could not save build outputs: %s", err)
  }

  // Destroy the build
  c.Log.Infof("destroying the build environment for permutation_id=%d", permutation.Id)
  _, err = c.p.Builder.DestroyBuild(context.TODO(), &proto.DestroyBuildRequest{
    Uuid: build.uuid,
  })
  if err != nil {
    c.releaseCoresById(buildCoreIds)
    return fmt.Errorf("could not destroy build environment%s", err)
  }

  // Free up cores
  err = c.releaseCoresById(buildCoreIds)
  if (err != nil) {
    return fmt.Errorf("could not release cores: %s", err)
  }

  //
  // Test
  //
  test := test{}

  vmmCoreIds := c.busyWaitForCores(1, &test, stageTest, task.IsolLevel, task.IsolSplit)
  benchToolCoreIds := c.busyWaitForCores(int(task.Test.BenchTool.Cores), &test, stageTest, task.IsolLevel, task.IsolSplit)
  kernelCoreIds := c.busyWaitForCores(int(task.Test.Kernel.Cores), &test, stageTest, task.IsolLevel, task.IsolSplit)
  testCoreIds := append(vmmCoreIds, benchToolCoreIds...)
  testCoreIds = append(testCoreIds, kernelCoreIds...)

  // Create the test
  c.Log.Infof("creating test for permutation_id=%d", permutation.Id)
  createTestResp, err := c.p.Tester.CreateTest(context.TODO(), &proto.CreateTestRequest{
    PermutationId:  int64(permutation.Id),
    VmmCores:       vmmCoreIds,
    Kernel:        &proto.TestKernel{
      Image:        saveBuildOutputsResp.Outputs.Kernel,
      InitRd:       saveBuildOutputsResp.Outputs.InitRd,
      Disks:        saveBuildOutputsResp.Outputs.Disks,
      Cores:        kernelCoreIds,
      Args:         task.Test.Kernel.Args,
      Memory:       task.Test.Kernel.Memory,
    },
    BenchTool:     &proto.TestBenchTool{
      Image:        task.Test.BenchTool.Image,
      Devices:      task.Test.BenchTool.Devices,
      Capabilities: task.Test.BenchTool.Capabilities,
      Commands:     task.Test.BenchTool.Commands,
      Cores:        benchToolCoreIds,
      StartDelay:   task.Test.BenchTool.StartDelay,
      EnvVars:      c.packEnvVars(task),
    },
  })
  if err != nil {
    c.releaseCoresById(testCoreIds)
    return fmt.Errorf("could not create test: %s", err)
  }

  test.uuid = createTestResp.Uuid

  var results []*proto.TestResult
  var resultType proto.TestResultType 
  for _, result := range task.Test.Results {

    // TODO: This switch should be abstract in case we need it elsewhere.
    // The proto definitions should probably have it and the spec should be
    // merged with the proto.
    switch result.Type {
      case "int",
           "integer":
        resultType = proto.TestResultType_TEST_RESULT_INT
      case "str",
           "string":
        resultType = proto.TestResultType_TEST_RESULT_STR
      case "bool":
        resultType = proto.TestResultType_TEST_RESULT_BOOL
      case "float":
        resultType = proto.TestResultType_TEST_RESULT_FLOAT
    }

    results = append(results, &proto.TestResult{
      Name: result.Name,
      Path: result.Path,
      Type: resultType,
    })
  }

  // Start the test
  c.Log.Infof("starting test for permutation_id=%d", permutation.Id)
  _, err = c.p.Tester.StartTest(context.TODO(), &proto.StartTestRequest{
    Uuid:    test.uuid,
    Results: results,
  })
  if err != nil {
    c.releaseCoresById(testCoreIds)
    return fmt.Errorf("could not start test: %s", err)
  }

  // Wait for the test to complete
  numRetries = 0
  c.Log.Infof("waiting for test for permutation_id=%d to complete...", permutation.Id)
  teststatus:for {
    statusTestResp, err := c.p.Tester.GetTestStatus(context.TODO(), &proto.GetTestStatusRequest{
      Uuid: test.uuid,
    })
    // Retry 5 times waiting a second each, fail if no response
    if err != nil {
      numRetries++
      if numRetries >= 5 {
        c.releaseCoresById(testCoreIds)
        return fmt.Errorf("could not get test status: %s", err)
      }
      time.Sleep(time.Second)
      continue
    }

    switch statusTestResp.Status {
      case proto.TestStatus_TEST_SUCCESS,
           proto.TestStatus_TEST_KILLED,
           proto.TestStatus_TEST_KERNEL_FAILED,
           proto.TestStatus_TEST_BENCHTOOL_FAILED:
        
        _, err = time.ParseDuration(statusTestResp.Runtime)
        if err != nil {
          c.releaseCoresById(testCoreIds)
          return fmt.Errorf("could not convert test runtime into duration: %s", err)
        }

        break teststatus;
    }

    time.Sleep(c.p.Cfg.GraceTime)
  }

  c.Log.Infof("test for permutation_id=%d to complete!", permutation.Id)

  // Destroy the test
  c.Log.Infof("destroying test permutation_id=%d", permutation.Id)
  _, err = c.p.Tester.DestroyTest(context.TODO(), &proto.DestroyTestRequest{
    Uuid: test.uuid,
  })
  if err != nil {
    c.releaseCoresById(testCoreIds)
    return fmt.Errorf("could not destroy test: %s", err)
  }

  err = c.releaseCoresById(testCoreIds)
  if err != nil {
    return fmt.Errorf("could not release cores: %s", err)
  }
  return nil
}
