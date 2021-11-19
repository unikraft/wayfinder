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

// ProcMount describes one entry (row) of /proc/mounts
type ProcMount struct {
  Device         string
  Mountpoint     string
  FileSystemType string
  Options        string
  dummy1         int
  dummy2         int
}

// GetProcMounts reads and returns an array of mount points defined in /proc/mounts
func GetProcMounts(procfs string) []ProcMount {
  mounts := []ProcMount{}
  filepath := fmt.Sprint(procfs, "/mounts")

  file, _ := os.Open(filepath)
  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)

  format := "%s %s %s %s %d %d"

  for scanner.Scan() {
    row := scanner.Text()
    mount := ProcMount{}

    _, err := fmt.Sscanf(
      string(row), format,
      &mount.Device,
      &mount.Mountpoint,
      &mount.FileSystemType,
      &mount.Options,
      &mount.dummy1,
      &mount.dummy2,
    )
    if err != nil {
      fmt.Fprintf(os.Stderr, "Cannot parse row in proc mounts: %s\n", err)
      continue
    }
    mounts = append(mounts, mount)
  }

  return mounts
}
