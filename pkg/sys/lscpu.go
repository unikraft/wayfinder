package sys

// SPDX-License-Identifier: BSD-3-Clause
//
// Adapted from: https://stackoverflow.com/a/35593102
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
	"os/exec"
	"strconv"
	"strings"
)

type CpuInfo struct {
	Architecture         string  `json:"architecture"`
	ByteOrder            string  `json:"byte_order"`
	AddressSizesPhysical string  `json:"address_sizes_physical"`
	AddressSizesVirtual  string  `json:"address_sizes_virtual"`
	CPUs                 uint    `json:"cpus"`
	ThreadsPerCore       uint    `json:"threads_per_core"`
	CoresPerSocket       uint    `json:"cores_per_socket"`
	Sockets              uint    `json:"sockets"`
	NUMAnodes            uint    `json:"numa_nodes"`
	VendorID             string  `json:"vendor_id"`
	CPUFamily            uint    `json:"cpu_family"`
	Model                uint    `json:"model"`
	ModelName            string  `json:"model_name"`
	Stepping             uint    `json:"stepping"`
	CPUMHz               float32 `json:"cpu_mhz"`
	CPUMaxMHz            float32 `json:"cpu_max_mhz"`
	CPUMinMHz            float32 `json:"cpu_min_mhz"`
	BogoMIPS             float32 `json:"bogo_mips"`
	Virtualization       string  `json:"virtualization"`
	L1dCache             string  `json:"l1d_cache"`
	L1iCache             string  `json:"l1i_cache"`
	L2Cache              string  `json:"l2_cache"`
	L3Cache              string  `json:"l3_cache"`
	NUMANode0CPUs        string  `json:"numa_node0_cpus"`
	NUMANode1CPUs        string  `json:"numa_node1_cpus"`
	Flags                string  `json:"flags"`
}

type CPULayoutInfo struct {
	CPU     uint64 `json:"cpu"`
	Core    uint64 `json:"core"`
	Socket  uint64 `json:"socket"`
	Node    uint64 `json:"node"`
	L3Cache uint64 `json:"l3_cache"`
}

func GetCpuInfo() (*CpuInfo, error) {
	out, err := exec.Command("lscpu").Output()
	if err != nil {
		return nil, fmt.Errorf("could not execute lscpu: %s", err)
	}

	outstring := strings.TrimSpace(string(out))
	lines := strings.Split(outstring, "\n")
	c := &CpuInfo{}

	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])

		switch key {
		case "Architecture":
			c.Architecture = value

		case "Byte Order":
			c.ByteOrder = value

		case "Address sizes":
			sizes := strings.Split(value, ",")
			if len(sizes) < 2 {
				continue
			}

			c.AddressSizesPhysical = strings.TrimSpace(sizes[0])
			c.AddressSizesVirtual = strings.TrimSpace(sizes[1])

		case "CPU(s)":
			t, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse number of cores: %s", err)
			}

			c.CPUs = uint(t)

		case "Thread(s) per core":
			t, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse number of threads per core: %s", err)
			}

			c.ThreadsPerCore = uint(t)

		case "Core(s) per socket":
			t, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse number of cores per socket: %s", err)
			}

			c.CoresPerSocket = uint(t)

		case "Socket(s)":
			t, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse number of sockets: %s", err)
			}

			c.Sockets = uint(t)

		case "NUMA node(s)":
			t, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse NUMA nodes: %s", err)
			}

			c.NUMAnodes = uint(t)

		case "Vendor ID":
			c.VendorID = value

		case "CPU family":
			t, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse CPU family: %s", err)
			}

			c.CPUFamily = uint(t)

		case "Model":
			t, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse model: %s", err)
			}

			c.Model = uint(t)

		case "Model name":
			c.ModelName = value

		case "Stepping":
			t, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("could not parse stepping: %s", err)
			}

			c.Stepping = uint(t)

		case "CPU MHz":
			t, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("could not parse CPU MHz: %s", err)
			}

			c.CPUMHz = float32(t)

		case "CPU max MHz":
			if strings.Contains(value, ",") {
				value = strings.ReplaceAll(value, ",", ".")
			}

			t, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("could not parse CPU max MHz: %s", err)
			}

			c.CPUMaxMHz = float32(t)

		case "CPU min MHz":
			if strings.Contains(value, ",") {
				value = strings.ReplaceAll(value, ",", ".")
			}

			t, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("could not parse CPU min MHz: %s", err)
			}

			c.CPUMinMHz = float32(t)

		case "BogoMIPS":
			t, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("could not parse bogoMIPS: %s", err)
			}

			c.BogoMIPS = float32(t)

		case "Virtualization":
			c.Virtualization = value

		case "L1d cache":
			c.L1dCache = value

		case "L1i cache":
			c.L1iCache = value

		case "L2 cache":
			c.L2Cache = value

		case "L3 cache":
			c.L3Cache = value

		case "NUMA node0 CPU(s)":
			c.NUMANode0CPUs = value

		case "NUMA node1 CPU(s)":
			c.NUMANode1CPUs = value

		case "Flags":
			c.Flags = value

		}
	}

	return c, nil
}

func GetCpuLayoutInfo() (*[]CPULayoutInfo, error) {
	out, err := exec.Command("lscpu", "--all", "--parse").Output()
	if err != nil {
		return nil, fmt.Errorf("could not execute \"lscpu --all --parse\": %s", err)
	}

	outstring := strings.TrimSpace(string(out))
	lines := strings.Split(outstring, "\n")
	c := make([]CPULayoutInfo, 0)

	for i := 4; i < len(lines); i++ {
		fields := strings.Split(lines[i], ",")
		CPUValue, _ := strconv.ParseUint(fields[0], 10, 64)
		CoreValue, _ := strconv.ParseUint(fields[1], 10, 64)
		SocketValue, _ := strconv.ParseUint(fields[2], 10, 64)
		NodeValue, _ := strconv.ParseUint(fields[3], 10, 64)
		L3CacheValue, _ := strconv.ParseUint(fields[8], 10, 64)
		c = append(c, CPULayoutInfo{
			CPU:     CPUValue,
			Core:    CoreValue,
			Socket:  SocketValue,
			Node:    NodeValue,
			L3Cache: L3CacheValue,
		})
	}

	return &c, nil
}
