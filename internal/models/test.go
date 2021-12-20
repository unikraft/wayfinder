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
  "time"

  "gorm.io/gorm"
  "github.com/google/uuid"

  "github.com/unikraft/wayfinder/api/proto"
)

// Test type that extends gorm.Model
type Test struct {
  Base

  UUID             uuid.UUID        `gorm:"type:char(36)"`

  PermutationId    uint             `gorm:"column:permutation_id;"    json:"permutation_id"`

  Status           proto.TestStatus `gorm:"column:status"             json:"status"`
  Runtime          time.Duration    `gorm:"column:runtime;default:0"  json:"runtime"`
  WayfinderVersion string           `gorm:"column:wayfinder_version" json:"wayfinder_version"`
  Results          []Result         `gorm:"foreignKey:test_id"        json:"results"`
  VMMCores         string           `gorm:"column:vmm_cores"         json:"vmm_cores"`
  KernelCores      string           `gorm:"column:kernel_cores"      json:"kernel_cores"`
  BenchToolCores   string           `gorm:"column:benchtool_cores"   json:"benchtool_cores"`
}

func (u *Test) BeforeCreate(tx *gorm.DB) (err error) {
  u.UUID = uuid.New()
  return nil
}
