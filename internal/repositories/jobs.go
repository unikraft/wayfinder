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
  "fmt"
  "gorm.io/gorm"

  "github.com/unikraft/wayfinder/pkg/sys"
  "github.com/unikraft/wayfinder/api/proto"
  "github.com/unikraft/wayfinder/internal/base64"
  "github.com/unikraft/wayfinder/internal/models"
)

// JobsRepository uses gorm.DB for querying the database
type JobsRepository struct {
  db *gorm.DB
}

// NewJobsRepository returns a default JobsRepository which uses
// gorm.DB for querying the database
func NewJobsRepository(db *gorm.DB) *JobsRepository {
  return &JobsRepository{db}
}

// CreateJob adds a new Job row to the Jobs table in the database
func (repo *JobsRepository) CreateJob(job *models.Job) (*models.Job, error) {
  // If the host is not set, we can look it up
  if job.HostId == 0 {
    host := &models.Host{}

    dmiUuid, err := sys.GetSysDmiUUID()
    if err != nil {
      return nil, err
    }

    if err := repo.db.Where("dmi_uuid = ?", dmiUuid).First(&host).Error; err != nil {
      return nil, fmt.Errorf("unknown host: %s", dmiUuid)
    }

    job.HostId = host.Id
  }

  if !base64.IsBase64(job.Config) {
    job.Config = base64.Encode(job.Config)
  }

  if err := repo.db.Create(job).Error; err != nil {
    return nil, err
  }
  return job, nil
}

// FindJob finds a job and decodes its config
func (repo *JobsRepository) FindJob(id uint, job *models.Job) error {
  if err := repo.db.Where("id = ?", id).First(&job).Error; err != nil {
    return err
  }

  if base64.IsBase64(job.Config) {
    decoded, err := base64.Decode(job.Config)
    if err != nil {
      return err
    }
    job.Config = decoded
  }

  return nil
}

// SetStatusJobById sets the state of the job to the desired state by the Job's
// ID.
func (repo *JobsRepository) SetStatusByJobId(id uint, status proto.JobStatus) error {
  job := &models.Job{}

  if err := repo.db.Where("id = ?", id).First(&job).Error; err != nil {
    return err
  }

  job.Status = status;

  if err := repo.db.Save(job).Error; err != nil {
    return err
  }

  return nil;
}

// SetStatusCreatedByJobId sets the state of the job to "created"
func (repo *JobsRepository) SetStatusCreatedByJobId(id uint) error {
  return repo.SetStatusByJobId(id, proto.JobStatus_JOB_STATUS_CREATED)
}

// SetStatusNotRunningByJobId sets the state of the job to "running"
func (repo *JobsRepository) SetStatusRunningByJobId(id uint) error {
  return repo.SetStatusByJobId(id, proto.JobStatus_JOB_STATUS_RUNNING)
}

// SetStatusPausedByJobId sets the state of the job to "paused"
func (repo *JobsRepository) SetStatusPausedByJobId(id uint) error {
  return repo.SetStatusByJobId(id, proto.JobStatus_JOB_STATUS_PAUSED)
}

// SetStatusFailedByJobId sets the state of the job to "failed"
func (repo *JobsRepository) SetStatusFailedByJobId(id uint) error {
  return repo.SetStatusByJobId(id, proto.JobStatus_JOB_STATUS_FAILED)
}

// SetStatusSuccessByJobId sets the state of the job to "success"
func (repo *JobsRepository) SetStatusSuccessByJobId(id uint) error {
  return repo.SetStatusByJobId(id, proto.JobStatus_JOB_STATUS_SUCCESS)
}

