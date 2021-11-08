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
  "fmt"
  "bytes"
  "strconv"
  "io/ioutil"

  "github.com/unikraft/wayfinder/config"
)

// ProcPIDSchedStat defines the fields of a /proc/[pid]/schedstat file
// cf. https://www.kernel.org/doc/Documentation/scheduler/sched-stats.txt
type ProcPIDSchedStat struct {
  // The process ID.
  PID int
  // time spent on the cpu
  Cputime uint64
  // time spent waiting on a runqueue
  Runqueue uint64
  // # of timeslices run on this cpu
  Timeslices uint64
}

// GetProcPIDSchedStat reads and returns the schedstat for a process from the proc fs
func GetProcPIDSchedStat(pid int) ProcPIDSchedStat {
  stats := ProcPIDSchedStat{PID: pid}
  filepath := fmt.Sprint(config.RuntimeConfig.ProcFS, "/", strconv.Itoa(pid), "/schedstat")
  filecontent, _ := ioutil.ReadFile(filepath)

  _, err := fmt.Fscan(
    bytes.NewBuffer(filecontent),
    &stats.Cputime,
    &stats.Runqueue,
    &stats.Timeslices,
  )

  if err != nil {
    return ProcPIDSchedStat{}
  }

  return stats
}
