package builder
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
  "sync"
  "context"

  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
  
  "github.com/unikraft/wayfinder/api/proto"
  "github.com/unikraft/wayfinder/internal/models"
  "github.com/unikraft/wayfinder/internal/strutils"
  "github.com/unikraft/wayfinder/modules/container"
  "github.com/unikraft/wayfinder/pkg/common/errors"
)

type build struct {
  sync.RWMutex
  container    *container.Container
  err           error
  runtime       time.Duration
}

type Service struct {
  p      *provider
  builds map[string]*build
}

func (s *Service) CreateBuild(ctx context.Context, req *proto.CreateBuildRequest) (*proto.CreateBuildResponse, error) {
  if req.PermutationId <= 0 {
    return &proto.CreateBuildResponse{
      Success: false,
    }, errors.NewMissingParameterError("permutation_id")
  }

  if req.Image == "" {
    return &proto.CreateBuildResponse{
      Success: false,
    }, errors.NewMissingParameterError("image")
  }

  if req.Commands == "" {
    return &proto.CreateBuildResponse{
      Success: false,
    }, errors.NewMissingParameterError("commands")
  }

  if len(req.Cores) == 0 {
    return &proto.CreateBuildResponse{
      Success: false,
    }, errors.NewMissingParameterError("cores")
  }

  if len(req.EnvVars) == 0 {
    return &proto.CreateBuildResponse{
      Success: false,
    }, errors.NewMissingParameterError("envvars")
  }

  // Create a new entry in the database for this build
  buildModel, err := s.p.DB.Repos().Builds().CreateBuildForPermutation(&models.Build{
    PermutationId: uint(req.PermutationId),
    Cores: strutils.JoinUint64(req.Cores, ","),
  })
  if err != nil {
    return &proto.CreateBuildResponse{
      Success: false,
    }, status.Errorf(codes.Internal, "could not create build in database: %s", err)
  }

  uuid := buildModel.UUID.String()

  // This is likely never going to occur
  if _, ok := s.builds[uuid]; ok {
    return &proto.CreateBuildResponse{
      Success: false,
    }, status.Errorf(codes.Internal, "build with uuid=%s already exists: %s", uuid)
  }

  builder, err := s.p.Container.NewContainer(uuid)
  if err != nil {
    s.p.DB.Repos().Builds().SetStatusKilledByBuildUuid(uuid)
    return &proto.CreateBuildResponse{
      Success: false,
    }, status.Errorf(codes.Internal, "cannot create build container: %s", err)
  }
  
  if err := builder.PullAndAttachImage(req.Image); err != nil {
    s.p.DB.Repos().Builds().SetStatusKilledByBuildUuid(uuid)
    return &proto.CreateBuildResponse{
      Success: false,
    }, status.Errorf(codes.Internal, "cannot pull image: %s", err)
  }

  if err := builder.SetCommands(req.Commands); err != nil {
    s.p.DB.Repos().Builds().SetStatusKilledByBuildUuid(uuid)
    return &proto.CreateBuildResponse{
      Success: false,
    }, status.Errorf(codes.Internal, "cannot set build container command: %s", err)
  }

  builder.SetCores(req.Cores)

  if len(req.Devices) > 0 {
    builder.SetDevices(req.Devices)
  }

  if len(req.Capabilities) > 0 {
    builder.SetCapabilities(req.Capabilities)
  }

  if len(req.Workdir) > 0 {
    builder.SetWorkdir(req.Workdir)
  }

  var envVars []string
  for _, env := range req.EnvVars {
    envVars = append(envVars, fmt.Sprintf("%s=%s", env.Name, env.Value))
  }

  builder.AddEnvVars(envVars)

  // Save the build
  s.builds[uuid] = &build{
    container: builder,
  }

  return &proto.CreateBuildResponse{
    Success: true,
    Uuid:    uuid,
  }, nil
}

func (s *Service) StartBuild(ctx context.Context, req *proto.StartBuildRequest) (*proto.StartBuildResponse, error) {
  if req.Uuid == "" {
    return &proto.StartBuildResponse{
      Success: false,
    }, errors.NewMissingParameterError("uuid")
  }

  // TODO: Look up builds with jobId/permId rather than via container ID (this
  // is for the proto/API).

  build, ok := s.builds[req.Uuid]
  if !ok {
    return &proto.StartBuildResponse{
      Success: false,
    }, status.Errorf(codes.NotFound, "cannot find build with id=%s", req.Uuid)
  }

  // Lock the build from other requests
  build.RLock()
  defer build.RUnlock()

  err := build.container.Init()
  if err != nil {
    s.p.DB.Repos().Builds().SetStatusFailedByBuildUuid(req.Uuid)

    return &proto.StartBuildResponse{
      Success: false,
    }, status.Errorf(codes.Internal, "could not initialize build container: %s", err)
  }

  s.p.DB.Repos().Builds().SetStatusRunningByBuildUuid(req.Uuid)

  // Start a thread which oversees this build
  go func(){
    runtime, err := build.container.StartAndWait()

    build.RLock()
    defer build.RUnlock()
    
    build.runtime = runtime

    // TODO: Putting note here not sure where to put it.  There's a bug where
    // the build time for a job with which has the permutation as another
    // saves to this.  This is a relational database model error.
    s.p.DB.Repos().Builds().SetRuntimeByBuildUuid(req.Uuid, runtime)

    if err != nil {
      build.err = err
      s.p.DB.Repos().Builds().SetStatusFailedByBuildUuid(req.Uuid)
    } else {
      s.p.DB.Repos().Builds().SetStatusSuccessByBuildUuid(req.Uuid)
    }
  }()

  return &proto.StartBuildResponse{
    Success: true,
  }, nil
}

func (s *Service) GetBuildStatus(ctx context.Context, req *proto.GetBuildStatusRequest) (*proto.GetBuildStatusResponse, error) {
  if req.Uuid == "" {
    return &proto.GetBuildStatusResponse{
      Success: false,
    }, errors.NewMissingParameterError("uuid")
  }

  buildModel := &models.Build{}
  if err := s.p.DB.DB().Where("uuid = ?", req.Uuid).First(&buildModel).Error; err != nil {
    return &proto.GetBuildStatusResponse{
      Success: false,
    }, status.Errorf(codes.NotFound, "build with uuid=%d not found", req.Uuid)
  }

  // TODO: Look up builds with jobId/permId rather than via container ID (this
  // is for the proto/API).

  build, ok := s.builds[req.Uuid]
  if !ok {
    return &proto.GetBuildStatusResponse{
      Success: false,
    }, status.Errorf(codes.NotFound, "cannot find build with id=%s", req.Uuid)
  }

  return &proto.GetBuildStatusResponse{
    Success: true,
    Status:  buildModel.Status,
    Runtime: build.runtime.String(),
  }, nil
}

func (s *Service) SaveBuildOutputsToDisk(ctx context.Context, req *proto.SaveBuildOutputsToDiskRequest) (*proto.SaveBuildOutputsToDiskResponse, error) {
  if req.Uuid == "" {
    return &proto.SaveBuildOutputsToDiskResponse{
      Success: false,
    }, errors.NewMissingParameterError("uuid")
  }

  if req.Outputs == nil {
    return &proto.SaveBuildOutputsToDiskResponse{
      Success: false,
    }, errors.NewMissingParameterError("outputs")
  }

  if req.Outputs.Kernel == "" {
    return &proto.SaveBuildOutputsToDiskResponse{
      Success: false,
    }, errors.NewMissingParameterError("outputs.kernel")
  }

  // TODO: Look up builds with jobId/permId rather than via container ID (this
  // is for the proto/API).

  build, ok := s.builds[req.Uuid]
  if !ok {
    return &proto.SaveBuildOutputsToDiskResponse{
      Success: false,
    }, status.Errorf(codes.NotFound, "cannot find build with id=%s", req.Uuid)
  }
  
  // Save kernel
  kernel, err := build.container.SaveOutput(s.p.Cfg.OutputDir, req.Outputs.Kernel)
  if err != nil {
    return &proto.SaveBuildOutputsToDiskResponse{
      Success: false,
    }, status.Errorf(codes.Internal, "build kernel could not be saved: %s", err)
  }

  s.p.Log.Infof("Saved: %s", kernel)
  if err := s.p.DB.Repos().Builds().SetKernelPathByBuildUuid(req.Uuid, kernel); err != nil {
    return &proto.SaveBuildOutputsToDiskResponse{
      Success: false,
    }, status.Errorf(codes.Internal, "build kernel path could not be saved database: %s", err)
  }

  initrd := ""

  if req.Outputs.InitRd != "" {
    initrd, err = build.container.SaveOutput(s.p.Cfg.OutputDir, req.Outputs.InitRd)
    if err != nil {
      return &proto.SaveBuildOutputsToDiskResponse{
        Success: false,
      }, status.Errorf(codes.Internal, "build could not be saved: %s", err)
    }

    if err := s.p.DB.Repos().Builds().SetInitRdPathByBuildUuid(req.Uuid, initrd); err != nil {
      return &proto.SaveBuildOutputsToDiskResponse{
        Success: false,
      }, status.Errorf(codes.Internal, "build initrd path could not be saved database: %s", err)
    }

    s.p.Log.Infof("Saved: %s", initrd)
  }

  disks := []*proto.BuildOutputDiskImage{}

  if len(req.Outputs.Disks) > 0 {
    for _, disk := range req.Outputs.Disks {
      diskPath, err := build.container.SaveOutput(s.p.Cfg.OutputDir, disk.Path)
      if err != nil {
        return &proto.SaveBuildOutputsToDiskResponse{
          Success: false,
        }, status.Errorf(codes.Internal, "build disk image could not be saved: %s", err)
      }

      if _, err := s.p.DB.Repos().Builds().AddDiskPathByBuildUuid(req.Uuid, &models.BuildOutputDisk{
        Type: disk.Type,
        Path: diskPath,
      }); err != nil {
        return &proto.SaveBuildOutputsToDiskResponse{
          Success: false,
        }, status.Errorf(codes.Internal, "build disk image path could not be saved database: %s", err)
      }

      disk.Path = diskPath
      disks = append(disks, disk)

      s.p.Log.Infof("Saved disk image: %s", diskPath)
    }
  }

  return &proto.SaveBuildOutputsToDiskResponse{
    Success: true,
    Outputs: &proto.BuildOutputs{
      Kernel: kernel,
      InitRd: initrd,
      Disks:  disks,
    },
  }, nil
}

func (s *Service) DestroyBuild(ctx context.Context, req *proto.DestroyBuildRequest) (*proto.DestroyBuildResponse, error) {
  if req.Uuid == "" {
    return &proto.DestroyBuildResponse{
      Success: false,
    }, errors.NewMissingParameterError("uuid")
  }

  // TODO: Look up builds with jobId/permId rather than via container ID (this
  // is for the proto/API).

  build, ok := s.builds[req.Uuid]
  if !ok {
    return &proto.DestroyBuildResponse{
      Success: false,
    }, status.Errorf(codes.NotFound, "cannot find build with id=%s", req.Uuid)
  }
  
  // TODO: Actually we should try and garbage collect this later too.

  if err := build.container.Destroy(); err != nil {
    return &proto.DestroyBuildResponse{
      Success: false,
    }, status.Errorf(codes.Internal, "cannot destroy build container: %s", err)
  }

  // wipe the build from memory
  delete(s.builds, req.Uuid)

  return &proto.DestroyBuildResponse{
    Success: true,
  }, nil
}
