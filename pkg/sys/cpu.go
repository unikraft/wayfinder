package sys

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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/sys/unix"
)

// SysCPU reflects cpu system information from /sys/devices/system/cpu/cpu*
type SysCPU struct {
	MaxFreq float32
	MinFreq float32
	CurFreq float32
}

// GetSysCPU returns the system CPU information for available cores
func GetSysCPU(cpuSets string) []SysCPU {
	stats := []SysCPU{}

	files, err := filepath.Glob(fmt.Sprintf("/sys/devices/system/cpu/cpu%s", cpuSets))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read cpu infos from sys fs: %s\n", err)
		return stats
	}

	var filepath string
	var filecontent []byte
	for _, f := range files {
		cpuStat := SysCPU{}

		filepath = fmt.Sprint(f + "/cpufreq/cpuinfo_max_freq")
		filecontent, _ = ioutil.ReadFile(filepath)
		fmt.Fscan(
			bytes.NewBuffer(filecontent),
			&cpuStat.MaxFreq,
		)

		filepath = fmt.Sprint(f + "/cpufreq/cpuinfo_min_freq")
		filecontent, _ = ioutil.ReadFile(filepath)
		fmt.Fscan(
			bytes.NewBuffer(filecontent),
			&cpuStat.MinFreq,
		)

		filepath = fmt.Sprint(f + "/cpufreq/scaling_cur_freq")
		filecontent, _ = ioutil.ReadFile(filepath)
		fmt.Fscan(
			bytes.NewBuffer(filecontent),
			&cpuStat.CurFreq,
		)

		stats = append(stats, cpuStat)

	}

	return stats
}

// SetAffinity pins a set of CPUs to a process ID
func SetAffinity(set []uint64, pid int) ([]uint64, error) {
	if len(set) == 0 {
		return set, nil
	}

	var filteredSet []uint64
	num := runtime.NumCPU()
	for _, index := range set {
		if index == 0 || int(index) >= num {
			continue
		}
		filteredSet = append(filteredSet, index)
	}

	if len(filteredSet) == 0 {
		return filteredSet, fmt.Errorf("unable to set affinity: no valid cpu ids specified")
	}

	cpuset := unix.CPUSet{}
	for _, index := range filteredSet {
		cpuset.Set(int(index))
	}

	err := unix.SchedSetaffinity(0, &cpuset)
	if err != nil {
		return filteredSet, err
	}

	return filteredSet, nil
}
