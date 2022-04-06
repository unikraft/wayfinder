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
	"gorm.io/gorm"

	"github.com/unikraft/wayfinder/internal/models"
)

// ParamsRepository uses gorm.DB for querying the database
type ParamsRepository struct {
	db *gorm.DB
}

// NewParamsRepository returns a default ParamsRepository which uses
// gorm.DB for querying the database
func NewParamsRepository(db *gorm.DB) *ParamsRepository {
	return &ParamsRepository{db}
}

// CreateParam adds a new param-value combination row to the Params table in the
// database
func (repo *ParamsRepository) CreateParam(param *models.Param) (*models.Param, error) {
	if err := repo.db.Create(param).Error; err != nil {
		return nil, err
	}
	return param, nil
}

// CreateParamInt adds a new param-value combination row to the Params table in
// the database
func (repo *ParamsRepository) CreateParamIntForJobId(id uint, key string, val int) (*models.Param, error) {
	job := &models.Job{}

	if err := repo.db.Where("id = ?", id).First(&job).Error; err != nil {
		return nil, err
	}

	// param := &models.Param{
	//   JobId: id,
	// }
	// return param, nil

	return nil, nil
}

// CreateParamStr adds a new param-value combination row to the Params table in
// the database
func (repo *ParamsRepository) CreateParamStrForJobId(id uint, key, val string) (*models.Param, error) {
	job := &models.Job{}

	if err := repo.db.Where("id = ?", id).First(&job).Error; err != nil {
		return nil, err
	}

	// param := &models.Param{
	//   JobId: id,
	//   Job:
	// }
	// return param, nil

	return nil, nil
}
