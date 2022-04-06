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

// Host type which extends gorm.Model
type Host struct {
	Base

	Jobs []Job `gorm:"foreignKey:host_id"                      json:"jobs"`

	DmiUUID  string `gorm:"column:dmi_uuid"                         json:"dmi_uuid"`
	Hostname string `gorm:"column:hostname"                         json:"hostname"`

	// contents of lscpu
	Architecture         string  `gorm:"column:architecture"                     json:"architecture"`
	ByteOrder            string  `gorm:"column:byte_order"                       json:"byte_order"`
	AddressSizesPhysical string  `gorm:"column:address_sizes_physical"           json:"address_sizes_physical"`
	AddressSizesVirtual  string  `gorm:"column:address_sizes_virtual"            json:"address_sizes_virtual"`
	CPUs                 uint    `gorm:"column:cpus"                             json:"cpus"`
	ThreadsPerCore       uint    `gorm:"column:threads_per_core"                 json:"threads_per_core"`
	CoresPerSocket       uint    `gorm:"column:cores_per_socket"                 json:"cores_per_socket"`
	Sockets              uint    `gorm:"column:sockets"                          json:"sockets"`
	NUMAnodes            uint    `gorm:"column:numa_nodes"                       json:"numa_nodes"`
	VendorID             string  `gorm:"column:vendor_id"                        json:"vendor_id"`
	CPUFamily            uint    `gorm:"column:cpu_family"                       json:"cpu_family"`
	Model                uint    `gorm:"column:model"                            json:"model"`
	ModelName            string  `gorm:"column:model_name"                       json:"model_name"`
	Stepping             uint    `gorm:"column:stepping"                         json:"stepping"`
	CPUMHz               float32 `gorm:"column:cpu_mhz;      type:decimal(9,4)"  json:"cpu_mhz"`
	CPUMaxMHz            float32 `gorm:"column:cpu_max_mhz;  type:decimal(9,4)"  json:"cpu_max_mhz"`
	CPUMinMHz            float32 `gorm:"column:cpu_min_mhz;  type:decimal(9,4)"  json:"cpu_min_mhz"`
	BogoMIPS             float32 `gorm:"column:bogo_mips;    type:decimal(9,4)"  json:"bogo_mips"`
	Virtualization       string  `gorm:"column:virtualization"                   json:"virtualization"`
	L1dCache             string  `gorm:"column:l1d_cache"                        json:"l1d_cache"`
	L1iCache             string  `gorm:"column:l1i_cache"                        json:"l1i_cache"`
	L2Cache              string  `gorm:"column:l2_cache"                         json:"l2_cache"`
	L3Cache              string  `gorm:"column:l3_cache"                         json:"l3_cache"`
	NUMANode0CPUs        string  `gorm:"column:numa_node0_cpus"                  json:"numa_node0_cpus"`
	NUMANode1CPUs        string  `gorm:"column:numa_node1_cpus"                  json:"numa_node1_cpus"`
	Flags                string  `gorm:"column:flags"                            json:"flags"`
}
