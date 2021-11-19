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
  "strconv"
  "io/ioutil"
)

// ProcPIDIO defines the fields of a /proc/[pid]/io file
// cf. https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/Documentation/filesystems/proc.txt?id=HEAD
type ProcPIDIO struct {
  // The process ID.
  PID                   int
  //  number of bytes the process read, using any read-like system call (from
  //  files, pipes, tty...).
  Rchar                 uint64
  // number of bytes the process wrote using any write-like system call.
  Wchar                 uint64
  // number of read-like system call invocations that the process performed.
  Syscr                 uint64
  // number of write-like system call invocations that the process performed.
  Syscw                 uint64
  // number of bytes the process directly read from disk.
  Read_bytes            uint64
  // number of bytes the process originally dirtied in the page-cache (assuming
  // they will go to disk later).
  Write_bytes           uint64
  // number of bytes the process "un-dirtied" - e.g. using an "ftruncate" call
  // that truncated pages from the page-cache.
  Cancelled_write_bytes uint64
}

// GetProcPIDIO reads and returns the io for a process from the proc fs
func GetProcPIDIO(procfs string, pid int) ProcPIDIO {
  stats := ProcPIDIO{PID: pid}
  filepath := fmt.Sprint(procfs, "/", strconv.Itoa(pid), "/io")
  filecontent, err := ioutil.ReadFile(filepath)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Cannot read proc io: %s\n", err)
    return ProcPIDIO{}
  }

  ioFormat := "rchar: %d\n" +
              "wchar: %d\n" +
              "syscr: %d\n" +
              "syscw: %d\n" +
              "read_bytes: %d\n" +
              "write_bytes: %d\n" +
              "cancelled_write_bytes: %d\n"

  _, err = fmt.Sscanf(
    string(filecontent), ioFormat,
    &stats.Rchar,
    &stats.Wchar,
    &stats.Syscr,
    &stats.Syscw,
    &stats.Read_bytes,
    &stats.Write_bytes,
    &stats.Cancelled_write_bytes,
  )

  if err != nil {
    fmt.Fprintf(os.Stderr, "Cannot parse proc io: %s\n", err)
    return ProcPIDIO{}
  }

  return stats
}
