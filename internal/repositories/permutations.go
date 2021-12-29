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
  "strconv"
  "gorm.io/gorm"

  "github.com/unikraft/wayfinder/spec"
  "github.com/unikraft/wayfinder/internal/models"
)

// PermutationsRepository uses gorm.DB for querying the database
type PermutationsRepository struct {
  db *gorm.DB
}

// NewPermutationsRepository returns a default PermutationsRepository which uses
// gorm.DB for querying the database
func NewPermutationsRepository(db *gorm.DB) *PermutationsRepository {
  return &PermutationsRepository{db}
}

// CreateParam adds a new param-value combination row to the Params table in the
// database
func (repo *PermutationsRepository) CreateParam(param *models.Param) (*models.Param, error) {
  if err := repo.db.Create(param).Error; err != nil {
    return nil, err
  }
  return param, nil
}

// FindOrCreateFromJobSpec is a multi-table function which creates the
// desired permutation as well as all the parameters needed for this
// permutation if does not exist.
func (r *PermutationsRepository) FindOrCreateFromJobSpec(job *spec.JobSpec) (*models.Permutation, error) {
  var err error
  permutation := &models.Permutation{}
  
  // Have we seen this permutation before?
  result := r.db.Where("job_id = ? and checksum = ?", &job.Id, &job.CurrentPerm.Checksum).First(&permutation);
  if result.RowsAffected == 1 {
    return permutation, nil
  }

  permutation.JobId = job.Id
  permutation.Checksum = job.CurrentPerm.Checksum

  // Populate the list of parameters (and their values) for this permutation
  for _, param := range job.CurrentPerm.Params {
    p := &models.Param{
      Name: param.Name,
      Type: param.Type,
    }
    switch param.Type {
      case "str":
        p.ValueStr = param.Value
        if err != nil {
          return nil, fmt.Errorf("could not parse param integer: %s", err)
        }
      case "int":
        p.ValueInt, err = strconv.ParseInt(param.Value, 10, 64)
      default:
        return nil, fmt.Errorf("unknown parameter type: %s", param.Type)
    }

    if err = r.db.Where(&p).FirstOrCreate(&p).Error; err != nil {
      return nil, fmt.Errorf("could not find or create parameter: %s", err)
    }

    permutation.Params = append(permutation.Params, *p)
  }

  // Create a new permutation entry
  if err = r.db.Create(&permutation).Error; err != nil {
    return nil, err
  }

  return permutation, nil
}


// UpdatePermutation updates only the Data field using Key as selector.
func (s *PermutationsRepository) UpdatePermutation(permutation *models.Permutation) (*models.Permutation, error) {
  if err := s.db.Model(permutation).Where("id = ?", permutation.Id).Updates(permutation).Error; err != nil {
    return nil, err
  }
  return permutation, nil
}
