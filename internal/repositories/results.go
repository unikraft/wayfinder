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

  "github.com/unikraft/wayfinder/api/proto"
  "github.com/unikraft/wayfinder/internal/models"
)

// ResultsRepository uses gorm.DB for querying the database
type ResultsRepository struct {
  db *gorm.DB
}

// NewResultsRepository returns a ResultsRepository which uses
// gorm.DB for querying the database
func NewResultsRepository(db *gorm.DB) *ResultsRepository {
  return &ResultsRepository{db}
}

// SaveResultIntByTestUuid adds a new integer-based test result
func (r *ResultsRepository) SaveResultIntByTestUuid(testUuid, name string, value int64) (*models.Result, error) {
  // Look up the test
  test := &models.Test{}
  if r.db.Where("uuid = ?", &testUuid).First(&test).RowsAffected != 1 {
    return nil, fmt.Errorf("could not find test with uuid: %s", testUuid)
  }

  result := &models.Result{
    TestId:        test.Id,
    PermutationId: test.PermutationId,
    Name:          name,
    Type:          proto.TestResultType_TEST_RESULT_INT,
    ValueInt:      value,
  }

  if err := r.db.Create(result).Error; err != nil {
    return nil, err
  }

  return result, nil
}


// SaveResultStrByTestUuid adds a new integer-based test result
func (r *ResultsRepository) SaveResultStrByTestUuid(testUuid, name, value string) (*models.Result, error) {
  // Look up the test
  test := &models.Test{}
  if r.db.Where("uuid = ?", &testUuid).First(&test).RowsAffected != 1 {
    return nil, fmt.Errorf("could not find test with uuid: %s", testUuid)
  }

  result := &models.Result{
    TestId:        test.Id,
    PermutationId: test.PermutationId,
    Name:          name,
    Type:          proto.TestResultType_TEST_RESULT_STR,
    ValueStr:      value,
  }

  if err := r.db.Create(result).Error; err != nil {
    return nil, err
  }

  return result, nil
}


// SaveResultFloatByTestUuid adds a new integer-based test result
func (r *ResultsRepository) SaveResultFloatByTestUuid(testUuid, name string, value float64) (*models.Result, error) {
  // Look up the test
  test := &models.Test{}
  if r.db.Where("uuid = ?", &testUuid).First(&test).RowsAffected != 1 {
    return nil, fmt.Errorf("could not find test with uuid: %s", testUuid)
  }

  result := &models.Result{
    TestId:        test.Id,
    PermutationId: test.PermutationId,
    Name:          name,
    Type:          proto.TestResultType_TEST_RESULT_FLOAT,
    ValueFloat:    value,
  }

  if err := r.db.Create(result).Error; err != nil {
    return nil, err
  }

  return result, nil
}


// SaveResultBoolByTestUuid adds a new integer-based test result
func (r *ResultsRepository) SaveResultBoolByTestUuid(testUuid, name string, value bool) (*models.Result, error) {
  // Look up the test
  test := &models.Test{}
  if r.db.Where("uuid = ?", &testUuid).First(&test).RowsAffected != 1 {
    return nil, fmt.Errorf("could not find test with uuid: %s", testUuid)
  }

  result := &models.Result{
    TestId:        test.Id,
    PermutationId: test.PermutationId,
    Name:          name,
    Type:          proto.TestResultType_TEST_RESULT_BOOL,
    ValueBool:     value,
  }

  if err := r.db.Create(result).Error; err != nil {
    return nil, err
  }

  return result, nil
}

// Extract all results for a given job
func (r *ResultsRepository) FindResults(jobId uint, offset, limit int) ([]*models.Result, error) {
  var results []*models.Result
  r.db.Offset(offset).Limit(limit).Where("job_id = ?", jobId).
      Joins("JOIN permutations ON results.permutation_id = permutations.id").
      Preload("results").Select("results.*").Find(&results)

  return results, nil
}
