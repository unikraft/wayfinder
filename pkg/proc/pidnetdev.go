package proc

// SPDX-License-Identifier: BSD-3-Clause
//
// Authors: Alexander Jung <a.jung@lancs.ac.uk>
//
// Copyright (c) 2020, Lancaster University.  All rights reserved.
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
	"strconv"
	"strings"
)

// ProcPIDNetDev represents entries of /proc/<pid>/net/dev file
type ProcPIDNetDev struct {
	// The process ID.
	PID int

	Dev string

	ReceivedBytes      uint64
	ReceivedPackets    uint64
	ReceivedErrs       uint64
	ReceivedDrop       uint64
	ReceivedFifo       uint64
	ReceivedFrame      uint64
	ReceivedCompressed uint64
	ReceivedMulticast  uint64

	TransmittedBytes      uint64
	TransmittedPackets    uint64
	TransmittedErrs       uint64
	TransmittedDrop       uint64
	TransmittedFifo       uint64
	TransmittedColls      uint64
	TransmittedCarrier    uint64
	TransmittedCompressed uint64
}

// GetProcPIDNetDev reads the net/dev file for given pid and device name from
// procfs
func GetProcPIDNetDev(procfs string, pid int, dev string) ProcPIDNetDev {
	stats := ProcPIDNetDev{
		PID: pid,
		Dev: dev,
	}

	filepath := fmt.Sprint(procfs, "/net/dev")
	if pid != 0 {
		filepath = fmt.Sprint(procfs, "/", strconv.Itoa(pid), "/net/dev")
	}

	file, _ := os.Open(filepath)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	format := "" + dev + ": %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d"

	foundDevStats := false
	for scanner.Scan() {
		row := strings.Trim(scanner.Text(), " ")

		if !strings.Contains(row, dev+":") {
			// only parse line with specified device
			continue
		}

		foundDevStats = true

		_, err := fmt.Sscanf(
			string(row), format,
			&stats.ReceivedBytes,
			&stats.ReceivedPackets,
			&stats.ReceivedErrs,
			&stats.ReceivedDrop,
			&stats.ReceivedFifo,
			&stats.ReceivedFrame,
			&stats.ReceivedCompressed,
			&stats.ReceivedMulticast,

			&stats.TransmittedBytes,
			&stats.TransmittedPackets,
			&stats.TransmittedErrs,
			&stats.TransmittedDrop,
			&stats.TransmittedFifo,
			&stats.TransmittedColls,
			&stats.TransmittedCarrier,
			&stats.TransmittedCompressed,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot parse row in proc net/dev: %s\n", err)
			return ProcPIDNetDev{}
		}
	}

	if !foundDevStats {
		fmt.Fprintf(os.Stderr, "could not find network device %s\n", dev)
		return ProcPIDNetDev{}
	}

	return stats
}
