package proc
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
  "os"
  "fmt"
  "bufio"
  "strconv"
  "strings"
)

// ProcCpuinfo represents one entry of cpuinfo proc file
type ProcCpuinfo struct {
  Processor       int
  VendorID        string
  CPUFamily       int
  Model           int
  ModelName       string
  Stepping        int
  Microcode       string
  CPUMhz          float32
  CacheSize       string
  PhysicalID      int
  Siblings        int
  CoreID          int
  CPUCores        int
  ApicID          int
  InitialApicID   int
  Fpu             string
  FpuException    string
  CpuidLevel      int
  Wp              string
  Flags           string
  Bugs            string
  Bogomips        float32
  ClflushSize     int
  CacheAlignment  int
  AddressSizes    string
  PowerManagement string
}

// GetProcCpuinfo reads the proc cpuinfo file
func GetProcCpuinfo(procfs string) []ProcCpuinfo {
  stats := []ProcCpuinfo{}

  filepath := fmt.Sprint(procfs, "/cpuinfo")

  file, _ := os.Open(filepath)
  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)

  core := 0
  for scanner.Scan() {
    row := scanner.Text()
    rowfields := strings.Split(row, ":")

    if len(stats) <= core {
      // create stats for new core
      stats = append(stats, ProcCpuinfo{})
    }
    if len(rowfields) < 2 {
      // new core detected
      core++
      continue
    }

    key := strings.TrimSpace(rowfields[0])
    value := strings.TrimSpace(rowfields[1])

    switch key {
    case "processor":
      stats[core].Processor, _ = strconv.Atoi(value)
    case "vendor_id":
      stats[core].VendorID = value
    case "cpu family":
      stats[core].CPUFamily, _ = strconv.Atoi(value)
    case "model":
      stats[core].Model, _ = strconv.Atoi(value)
    case "model name":
      stats[core].ModelName = value
    case "stepping":
      stats[core].Stepping, _ = strconv.Atoi(value)
    case "microcode":
      stats[core].Microcode = value
    case "cpu MHz":
      tmp, _ := strconv.ParseFloat(value, 32)
      stats[core].CPUMhz = float32(tmp)
    case "cache size":
      stats[core].CacheSize = value
    case "physical id":
      stats[core].PhysicalID, _ = strconv.Atoi(value)
    case "siblings":
      stats[core].Siblings, _ = strconv.Atoi(value)
    case "core id":
      stats[core].CoreID, _ = strconv.Atoi(value)
    case "cpu cores":
      stats[core].CPUCores, _ = strconv.Atoi(value)
    case "apicid":
      stats[core].ApicID, _ = strconv.Atoi(value)
    case "initial apicid":
      stats[core].InitialApicID, _ = strconv.Atoi(value)
    case "fpu":
      stats[core].Fpu = value
    case "fpu_exception":
      stats[core].FpuException = value
    case "cpuid level":
      stats[core].CpuidLevel, _ = strconv.Atoi(value)
    case "wp":
      stats[core].Wp = value
    case "flags":
      stats[core].Flags = value
    case "bugs":
      stats[core].Bugs = value
    case "bogomips":
      tmp, _ := strconv.ParseFloat(value, 32)
      stats[core].Bogomips = float32(tmp)
    case "clflush size":
      stats[core].ClflushSize, _ = strconv.Atoi(value)
    case "cache_alignment":
      stats[core].CacheAlignment, _ = strconv.Atoi(value)
    case "address sizes":
      stats[core].AddressSizes = value
    case "power management":
      stats[core].PowerManagement = value
    }
  }

  return stats
}
