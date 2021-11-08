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
  "os"
  "fmt"
  "bufio"

  "github.com/unikraft/wayfinder/config"
)

// some avg10=0.00 avg60=0.00 avg300=0.00 total=109155294
// full avg10=0.00 avg60=0.00 avg300=0.00 total=71768841

// ProcPressureMetric describes the type for Metric in a ProcPressure element
type ProcPressureMetric string

const (
  // ProcPressureMetricSome defines the metric type "some" for a ProcPressure element
  ProcPressureMetricSome ProcPressureMetric = "some"
  // ProcPressureMetricFull defines the metric type "full" for a ProcPressure element
  ProcPressureMetricFull ProcPressureMetric = "Full"
)

// ProcPressureResource describes the resource (cpu,io,mem)
type ProcPressureResource string

const (
  // ProcPressureResourceCPU defines the resource type "cpu" for a ProcPressure element
  ProcPressureResourceCPU ProcPressureResource = "cpu"
  // ProcPressureResourceIO defines the metric type "io" for a ProcPressure element
  ProcPressureResourceIO ProcPressureResource = "io"
  // ProcPressureResourceMemory defines the metric type "memory" for a ProcPressure element
  ProcPressureResourceMemory ProcPressureResource = "memory"
)

// ProcPressure describes one row in /proc/pressure/{io,cpu,mem}
type ProcPressure struct {
  Metric ProcPressureMetric // some or full
  Avg10  float64
  Avg60  float64
  Avg300 float64
  Total  uint64
}

// GetProcPressure reads and returns the pressures for the given resource
func GetProcPressure(resource ProcPressureResource) []ProcPressure {
  pressures := []ProcPressure{}
  filepath := fmt.Sprint(config.RuntimeConfig.ProcFS, "/pressure/", resource)

  file, err := os.Open(filepath)
  if err != nil {
    // cannot open file ...
    return pressures
  }
  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)

  format := "%s avg10=%f avg60=%f avg300=%f total=%d"

  for scanner.Scan() {
    row := scanner.Text()
    pressure := ProcPressure{}

    _, err := fmt.Sscanf(
      string(row), format,
      &pressure.Metric,
      &pressure.Avg10,
      &pressure.Avg60,
      &pressure.Avg300,
      &pressure.Total,
    )
    if err != nil {
      fmt.Fprintf(os.Stderr, "Cannot parse row in proc pressure %s: %s\n", resource, err)
      continue
    }
    pressures = append(pressures, pressure)
  }
  return pressures
}
