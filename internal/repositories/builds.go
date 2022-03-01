package repositories
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
  "time"
  "gorm.io/gorm"

  "github.com/unikraft/wayfinder/api/proto"
  "github.com/unikraft/wayfinder/internal/models"
  "github.com/unikraft/wayfinder/internal/version"
)

// BuildsRepository uses gorm.DB for querying the database
type BuildsRepository struct {
  db *gorm.DB
}

// NewBuildsRepository returns a BuildsRepository which uses
// gorm.DB for querying the database
func NewBuildsRepository(db *gorm.DB) *BuildsRepository {
  return &BuildsRepository{db}
}

// CreateBuildForPermutation adds a new test row to the Tests table in the
// database
func (repo *BuildsRepository) CreateBuildForPermutation(build *models.Build) (*models.Build, error) {
  build.Status = proto.BuildStatus_BUILD_CREATED
  build.WayfinderVersion = version.String()

  if err := repo.db.Create(build).Error; err != nil {
    return nil, err
  }
  return build, nil
}

func (repo *BuildsRepository) DeleteBuild(build *models.Build) error {  
  if err := repo.db.Delete(build).Error; err != nil {
    return err
  }

  return nil
}


func (repo *BuildsRepository) DeleteBuildByBuildUuid(uuid string) error {
  build := models.Build{}

  if err := repo.db.Where("uuid = ?", uuid).First(&build).Error; err != nil {
    return err
  }
  
  if err := repo.db.Delete(build).Error; err != nil {
    return err
  }

  return nil
}

// SetStatusByBuildUuid sets the status of the build to the desired status by
// the Build's UUID.
func (repo *BuildsRepository) SetStatusByBuildUuid(uuid string, status proto.BuildStatus) error {
  build := &models.Build{}

  if err := repo.db.Where("uuid = ?", uuid).First(&build).Error; err != nil {
    return err
  }

  build.Status = status;

  if err := repo.db.Save(build).Error; err != nil {
    return err
  }

  return nil;
}

// SetStatusCreatedByBuildId sets the state of the build to "created"
func (repo *BuildsRepository) SetStatusCreatedByBuildUuid(uuid string) error {
  return repo.SetStatusByBuildUuid(uuid, proto.BuildStatus_BUILD_CREATED)
}

// SetStatusRunningByBuildUuid sets the state of the build to "running"
func (repo *BuildsRepository) SetStatusRunningByBuildUuid(uuid string) error {
  return repo.SetStatusByBuildUuid(uuid, proto.BuildStatus_BUILD_RUNNING)
}

// SetStatusPausedByBuildUuid sets the state of the build to "paused"
func (repo *BuildsRepository) SetStatusPausedByBuildUuid(uuid string) error {
  return repo.SetStatusByBuildUuid(uuid, proto.BuildStatus_BUILD_PAUSED)
}

// SetStatusSuccessByBuildUuid sets the state of the build to "success"
func (repo *BuildsRepository) SetStatusSuccessByBuildUuid(uuid string) error {
  return repo.SetStatusByBuildUuid(uuid, proto.BuildStatus_BUILD_SUCCESS)
}

// SetStatusKilledByBuildUuid sets the state of the build to "killed"
func (repo *BuildsRepository) SetStatusKilledByBuildUuid(uuid string) error {
  return repo.SetStatusByBuildUuid(uuid, proto.BuildStatus_BUILD_KILLED)
}

// SetStatusFailedByBuildUuid sets the state of the build to "failed"
func (repo *BuildsRepository) SetStatusFailedByBuildUuid(uuid string) error {
  return repo.SetStatusByBuildUuid(uuid, proto.BuildStatus_BUILD_FAILED)
}

// SetRuntimeByBuildUuid sets the state of the build to "created"
func (repo *BuildsRepository) SetRuntimeByBuildUuid(uuid string, runtime time.Duration) error {
  build := &models.Build{}

  if err := repo.db.Where("uuid = ?", uuid).First(&build).Error; err != nil {
    return err
  }

  build.Runtime = runtime;

  if err := repo.db.Save(build).Error; err != nil {
    return err
  }

  return nil;
}

// SetKernelPathByBuildUuid sets the location on disk of the kernel binary
func (repo *BuildsRepository) SetKernelPathByBuildUuid(uuid, path string) error {
  build := &models.Build{}

  if err := repo.db.Where("uuid = ?", uuid).First(&build).Error; err != nil {
    return err
  }

  build.KernelPath = path;

  if err := repo.db.Save(build).Error; err != nil {
    return err
  }

  return nil;
}

// SetInitRdPathByBuildUuid sets the location on disk of the initrd file
func (repo *BuildsRepository) SetInitRdPathByBuildUuid(uuid, path string) error {
  build := &models.Build{}

  if err := repo.db.Where("uuid = ?", uuid).First(&build).Error; err != nil {
    return err
  }

  build.InitRdPath = path;

  if err := repo.db.Save(build).Error; err != nil {
    return err
  }

  return nil;
}


// AddDiskPathByBuildUuid sets the location on disk of the initrd file
func (repo *BuildsRepository) AddDiskPathByBuildUuid(uuid string, buildOutputDisk *models.BuildOutputDisk) (*models.BuildOutputDisk, error) {
  build := &models.Build{}

  if err := repo.db.Where("uuid = ?", uuid).First(&build).Error; err != nil {
    return nil, err
  }

  buildOutputDisk.BuildId = build.ID

  if err := repo.db.Create(buildOutputDisk).Error; err != nil {
    return nil, err
  }

  return buildOutputDisk, nil
}
