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

  "github.com/adjust/rmq/v4"
  
  "github.com/unikraft/wayfinder/spec"
  "github.com/unikraft/wayfinder/api/proto"
  "github.com/erda-project/erda-infra/base/logs"
  "github.com/unikraft/wayfinder/internal/coremap"
)

type TaskConsumer struct {
  p       *provider
  Log      logs.Logger
}

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

  // Check if we received the full job specification
  err := json.Unmarshal([]byte(taskBytes), &task)
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
  for _, coreId := range coresToFree {
    retries := 0
    c.p.CoreMap().ReleaseCore(coreId)
    for {
      if activity, _ := c.p.CoreMap().GetCoreActivity(coreId); activity == nil {
        break;
      }
      if retries > 10 {
        return fmt.Errorf("could not release core %d", coreId)
      }
      time.Sleep(1 * time.Millisecond)
      c.p.CoreMap().ReleaseCore(coreId)
      retries++
    }
  }
  return nil
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
  // if err := s.p.DB.Repos().Job().FindJob(uint(req.JobId), job); err != nil {
  //   return &proto.CreateBuildResponse{
  //     Success: false,
  //     Status:  proto.BuildStatus_BUILD_FAILED,
  //   }, status.Errorf(codes.NotFound, "job with id=%d not found", req.JobId)
  // }

  // permutation := &models.Permutation{}
  // if err := c.p.DB.DB().Where("id = ?", perm.Id).First(&permutation).Error; err != nil {
  //   return fmt.Errorf("permutation with id=%d not found", perm.Id)
  // }

  build := &build{}

  requiredNumBuildCores := int(task.Build.Cores)
  var buildCoreIds []uint64
  var buildCores []*coremap.Core
  
  // Wait until we have some free cores.
  for {
    if len(buildCores) >= requiredNumBuildCores {
      break
    }

    // TODO: Logic on whether the build needs to be numa/socket/cache/core aware
    cores := c.p.CoreMap().FindAllFreeCoresAcrossAllNumaNodes()
    for _, core := range cores {
      // Immediately reserve this core
      if err := c.p.CoreMap().SetCoreActivity(core.Id(), &build); err != nil {
        continue
      }

      buildCoreIds = append(buildCoreIds, core.Id())
      buildCores = append(buildCores, core)

      // Also check here in case we received more cores than requested
      if len(buildCores) >= requiredNumBuildCores {
        break
      }
    }

    time.Sleep(c.p.Cfg.GraceTime)
  }

  var buildEnvVars []*proto.BuildEnvVar
  for _, param := range task.CurrentPerm.Params {
    buildEnvVars = append(buildEnvVars, &proto.BuildEnvVar{
      Name: param.Name,
      Value: param.Value,
    })
  }

  fmt.Printf("params = %#v\n", buildEnvVars)

  // Create the build
  c.Log.Infof("creating build container for permutation_id=%d", permutation.Id)
  createBuildResp, err := c.p.Builder.CreateBuild(context.TODO(), &proto.CreateBuildRequest{
    PermutationId: int64(permutation.Id),
    Image:         task.Build.Image,
    Devices:       task.Build.Devices,
    Capabilities:  task.Build.Capabilities,
    Cores:         buildCoreIds,
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

  // Save the outputs from the build
  c.Log.Infof("saving outputs from permutation_id=%d", permutation.Id)
  saveBuildOutputsResp, err := c.p.Builder.SaveBuildOutputsToDisk(context.TODO(), &proto.SaveBuildOutputsToDiskRequest{
    Uuid:       build.uuid,
    Outputs:   &proto.BuildOutputs{
      Kernel:   task.Build.Outputs.Kernel,
      InitRd:   task.Build.Outputs.InitRd,
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

  // We add an extra core to be used for the QEMU VMM.
  // TODO: Implement
  requiredNumTestCores := 1 + task.Test.Kernel.Cores + task.Test.BenchTool.Cores

  var testCoreIds []uint64
  var testCores []*coremap.Core

  for {
    if uint64(len(testCores)) >= requiredNumTestCores {
      break
    }

    // TODO: Logic on whether the build needs to be numa/socket/cache/core aware
    cores := c.p.CoreMap().FindAllFreeCoresAcrossAllNumaNodes()
    for _, core := range cores {
      // Immediately reserve this core
      if err := c.p.CoreMap().SetCoreActivity(core.Id(), &test); err != nil {
        continue
      }

      testCoreIds = append(testCoreIds, core.Id())
      testCores = append(testCores, core)

      // Also check here in case we received more cores than requested
      if uint64(len(testCores)) >= requiredNumTestCores {
        break
      }
    }

    time.Sleep(c.p.Cfg.GraceTime)
  }

  // Pop a core from our list of reserved cores.  Niiave!  VMM, VM and Bench
  // tool locality matter!
  var x uint64
  var kernelCores []uint64
  benchToolCores := append([]uint64(nil), testCoreIds...) // mem copy trick
  for i := uint64(0); i < task.Test.Kernel.Cores; i++ {
    // Pop from benchToolCores and push to kernelCores
    x, benchToolCores = benchToolCores[len(benchToolCores)-1],
                        benchToolCores[:len(benchToolCores)-1]
    kernelCores = append(kernelCores, x)
  }

  // Create the test
  c.Log.Infof("creating test for permutation_id=%d", permutation.Id)
  createTestResp, err := c.p.Tester.CreateTest(context.TODO(), &proto.CreateTestRequest{
    PermutationId:  int64(permutation.Id),
    Kernel:        &proto.TestKernel{
      Image:        saveBuildOutputsResp.Outputs.Kernel,
      InitRd:       saveBuildOutputsResp.Outputs.InitRd,
      Cores:        kernelCores,
      Args:         task.Test.Kernel.Args,
      Memory:       task.Test.Kernel.Memory,
    },
    BenchTool:     &proto.TestBenchTool{
      Image:        task.Test.BenchTool.Image,
      Devices:      task.Test.BenchTool.Devices,
      Capabilities: task.Test.BenchTool.Capabilities,
      Commands:     task.Test.BenchTool.Commands,
      Cores:        benchToolCores,
      StartDelay:   task.Test.BenchTool.StartDelay,
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
