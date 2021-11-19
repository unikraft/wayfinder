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
)

// ProcDiskstat defines the fields of one row (one block device) of a
// /proc/diskstats file cf.
// https://www.kernel.org/doc/Documentation/ABI/testing/procfs-diskstats
type ProcDiskstat struct {
  Majornumber        int    //  1 - major number
  Minornumber        int    //  2 - minor mumber
  Devicename         string //  3 - device name
  Reads              uint64 //  4 - reads completed successfully f1
  ReadsMerged        uint64 //  5 - reads merged f2
  SectorsRead        uint64 //  6 - sectors read f3
  TimeReading        uint64 //  7 - time spent reading (ms) f4
  Writes             uint64 //  8 - writes completed f5
  WritesMerged       uint64 //  9 - writes merged f6
  SectorsWritten     uint64 // 10 - sectors written f7
  TimeWriting        uint64 // 11 - time spent writing (ms) f8
  CurrentOps         uint64 // 12 - I/Os currently in progress f9
  TimeForOps         uint64 // 13 - time spent doing I/Os (ms) f10
  WeightedTimeForOps uint64 // 14 - weighted time spent doing I/Os (ms) f11
}

// GetProcDiskstats reads and returns the diskstats from the proc fs
func GetProcDiskstats(procfs string) map[string]ProcDiskstat {
  stats := make(map[string]ProcDiskstat)

  filepath := fmt.Sprint(procfs, "/diskstats")
  file, _ := os.Open(filepath)
  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)

  format := "%d %d %s %d %d %d %d %d %d %d %d %d %d %d"

  for scanner.Scan() {
    row := scanner.Text()
    diskstat := ProcDiskstat{}

    _, err := fmt.Sscanf(
      string(row), format,
      &diskstat.Majornumber,
      &diskstat.Minornumber,
      &diskstat.Devicename,
      &diskstat.Reads,
      &diskstat.ReadsMerged,
      &diskstat.SectorsRead,
      &diskstat.TimeReading,
      &diskstat.Writes,
      &diskstat.WritesMerged,
      &diskstat.SectorsWritten,
      &diskstat.TimeWriting,
      &diskstat.CurrentOps,
      &diskstat.TimeForOps,
      &diskstat.WeightedTimeForOps,
    )

    if err != nil {
      fmt.Fprintf(os.Stderr, "Cannot parse row in proc diskstats: %s\n", err)
      continue
    }

    stats[diskstat.Devicename] = diskstat
  }

  return stats
}
