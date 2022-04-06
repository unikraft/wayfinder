package proc

// SPDX-License-Identifier: BSD-3-Clause
//
// Authors: Alexander Jung <a.jung@lancs.ac.uk>
//
// Copyright (c) 2021, Lancaster University.  All rights reserved.
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
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ProcStatCPU describes one CPU row in /proc/stat
type ProcStatCPU struct {
	Name string

	// The amount of time, measured in units of USER_HZ
	// (1/100ths of a second on most architectures)
	User      uint64
	Nice      uint64
	System    uint64
	Idle      uint64
	IOWait    uint64
	IRQ       uint64
	SoftIRQ   uint64
	Steal     uint64
	Guest     uint64
	GuestNice uint64
}

// GetProcStatCPU reads and returns the cpu related rows in /proc/stat
func GetProcStatCPU(procfs string) []ProcStatCPU {
	stats := []ProcStatCPU{}
	filepath := fmt.Sprint(procfs, "/stat")

	file, err := os.Open(filepath)
	if err != nil {
		// cannot open file ...
		return stats
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	format := "%s %d %d %d %d %d %d %d %d %d %d"

	for scanner.Scan() {
		row := scanner.Text()
		stat := ProcStatCPU{}

		// filter rows, only consider cpu rows
		if !strings.HasPrefix(row, "cpu") {
			continue
		}

		_, err := fmt.Sscanf(
			string(row), format,
			&stat.Name,
			&stat.User,
			&stat.Nice,
			&stat.System,
			&stat.Idle,
			&stat.IOWait,
			&stat.IRQ,
			&stat.SoftIRQ,
			&stat.Steal,
			&stat.Guest,
			&stat.GuestNice,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot parse row in proc stat: %s\n", err)
			continue
		}
		stats = append(stats, stat)
	}
	return stats
}
