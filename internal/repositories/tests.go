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

// TestsRepository uses gorm.DB for querying the database
type TestsRepository struct {
	db *gorm.DB
}

// NewTestsRepository returns a TestsRepository which uses
// gorm.DB for querying the database
func NewTestsRepository(db *gorm.DB) *TestsRepository {
	return &TestsRepository{db}
}

// CreateTestForPermutation adds a new test row to the Tests table in the
// database
func (repo *TestsRepository) CreateTestForPermutation(test *models.Test) (*models.Test, error) {
	test.Status = proto.TestStatus_TEST_CREATED
	test.WayfinderVersion = version.String()

	if err := repo.db.Create(test).Error; err != nil {
		return nil, err
	}
	return test, nil
}

// Delete a given test. If purge is used, the entry is deleted permanently
func (repo *TestsRepository) DeleteTest(test *models.Test, purge bool) error {
	var deleteType *gorm.DB

	if purge {
		deleteType = repo.db.Unscoped()
	} else {
		deleteType = repo.db
	}

	if err := deleteType.Delete(&models.Test{}, "uuid = ?", test.UUID).Error; err != nil {
		return err
	}

	return nil
}

// Deletes a test from the tests table. If purge is used,
// the entry is deleted permanently
func (repo *TestsRepository) DeleteTestByTestUuid(uuid string, purge bool) error {
	var deleteType *gorm.DB

	if purge {
		deleteType = repo.db.Unscoped()
	} else {
		deleteType = repo.db
	}

	if err := deleteType.Delete(&models.Test{}, "uuid = ?", uuid).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TestsRepository) DeleteTestsByPermutationId(permutationId int64, purge bool) error {
	var deleteType *gorm.DB

	if purge {
		deleteType = repo.db.Unscoped()
	} else {
		deleteType = repo.db
	}

	if err := deleteType.Delete(&models.Test{}, "permutation_id = ?", permutationId).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TestsRepository) FindTestByTestUuid(uuid string) (*models.Test, error) {
	test := models.Test{}

	if err := repo.db.Where("uuid = ?", uuid).First(&test).Error; err != nil {
		return nil, err
	}

	return &test, nil
}

// SetStatusByTestUuid sets the status of the test to the desired status by
// the Test's UUID.
func (repo *TestsRepository) SetStatusByTestUuid(uuid string, status proto.TestStatus) error {
	test := &models.Test{}

	if err := repo.db.Where("uuid = ?", uuid).First(&test).Error; err != nil {
		return err
	}

	test.Status = status

	if err := repo.db.Save(test).Error; err != nil {
		return err
	}

	return nil
}

// SetStatusCreatedByTestId sets the state of the test to "created"
func (repo *TestsRepository) SetStatusCreatedByTestUuid(uuid string) error {
	return repo.SetStatusByTestUuid(uuid, proto.TestStatus_TEST_CREATED)
}

// SetStatusRunningByTestUuid sets the state of the test to "running"
func (repo *TestsRepository) SetStatusRunningByTestUuid(uuid string) error {
	return repo.SetStatusByTestUuid(uuid, proto.TestStatus_TEST_RUNNING)
}

// SetStatusPausedByTestUuid sets the state of the test to "paused"
func (repo *TestsRepository) SetStatusPausedByTestUuid(uuid string) error {
	return repo.SetStatusByTestUuid(uuid, proto.TestStatus_TEST_PAUSED)
}

// SetStatusKilledByTestUuid sets the state of the test to "killed"
func (repo *TestsRepository) SetStatusKilledByTestUuid(uuid string) error {
	return repo.SetStatusByTestUuid(uuid, proto.TestStatus_TEST_KILLED)
}

// SetStatusKernelFailedByTestUuid sets the state of the test to "failed"
func (repo *TestsRepository) SetStatusKernelFailedByTestUuid(uuid string) error {
	return repo.SetStatusByTestUuid(uuid, proto.TestStatus_TEST_KERNEL_FAILED)
}

// SetStatusBenchToolFailedByTestUuid sets the state of the test to "failed"
func (repo *TestsRepository) SetStatusBenchToolFailedByTestUuid(uuid string) error {
	return repo.SetStatusByTestUuid(uuid, proto.TestStatus_TEST_BENCHTOOL_FAILED)
}

// SetStatusSuccessByTestUuid sets the state of the test to "success"
func (repo *TestsRepository) SetStatusSuccessByTestUuid(uuid string) error {
	return repo.SetStatusByTestUuid(uuid, proto.TestStatus_TEST_SUCCESS)
}

// SetRuntimeByTestUuid sets the state of the test to "created"
func (repo *TestsRepository) SetRuntimeByTestUuid(uuid string, runtime time.Duration) error {
	test := &models.Test{}

	if err := repo.db.Where("uuid = ?", uuid).First(&test).Error; err != nil {
		return err
	}

	test.Runtime = runtime

	if err := repo.db.Save(test).Error; err != nil {
		return err
	}

	return nil
}
