package models

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
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/unikraft/wayfinder/api/proto"
)

// Permutation type that extends gorm.Model
type Permutation struct {
	Base

	UUID uuid.UUID `gorm:"type:char(36)"`

	JobId    uint                       `gorm:"column:job_id"`
	Checksum string                     `gorm:"column:checksum"`
	Params   []Param                    `gorm:"many2many:permutation_params" json:"params"`
	Results  []Result                   `gorm:"foreignKey:permutation_id"    json:"results"`
	Status   proto.JobPermutationStatus `gorm:"column:status"                json:"status"`
	Builds   []Build                    `gorm:"foreignKey:permutation_id"    json:"build_times"`
	Tests    []Test                     `gorm:"foreignKey:permutation_id"    json:"test_times"`
}

func (u *Permutation) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.New()
	return nil
}
