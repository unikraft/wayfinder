package libvirt

// SPDX-License-Identifier: BSD-3-Clause
//
// Authors: Cezar Craciunoiu <cezar.craciunoiu@gmail.com>
//
// Copyright (c) 2022, Universitatea POLITEHNICA of Bucharest.  All rights reserved.
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
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/unikraft/wayfinder/internal/metrics"
)

// Same as the other metrics, but format is general
func (d *Domain) MetricLookup() error {
	return nil
}

// Runs every measurement script and adds the values obtain
func (d *Domain) MetricMeasure(name string, pid int, command string) error {
	cmd := exec.Command("bash", "-c", command)
	cmd.Env = append(os.Environ(), fmt.Sprintf("DOMAIN_PID=%d", pid))
	// TODO: variablize 'sh' such that user can specify entrypoint program
	result, err := cmd.Output()
	if err != nil {
		return err
	}

	measurement, _ := strconv.ParseFloat(strings.TrimSuffix(string(result), "\n"), 64)

	d.AddMeasurement(name, metrics.CreateMeasurement(measurement))

	return nil
}

// Prints the name-value pairs for each monitor
func (d *Domain) MetricPrint(name string) map[string]string {
	return map[string]string{
		name: d.GetMetricFloat64(name, 0),
	}
}
