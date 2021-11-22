package main
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
  "bytes"
  "context"
  "io/ioutil"

  "github.com/spf13/cobra"

  "github.com/unikraft/wayfinder/api/proto"
  "github.com/unikraft/wayfinder/internal/gzip"
)

var (
  createCmd = &cobra.Command{
    Use:                   "create [OPTIONS...] FILE",
    Aliases:               []string{"cj"},
    Short:                 `Create a new job to be executed by the wayfinder server.`,
    Run:                   doCreateCmd,
    Args:                  cobra.ExactArgs(1),
    DisableFlagsInUseLine: true,
  }
)

// doCreateCmd 
func doCreateCmd(cmd *cobra.Command, args []string) {
  if len(args) == 0 {
    cmd.Help()
  }
  
  filePath := args[0]
  if _, err := os.Stat(filePath); err != nil {
    fmt.Printf("file does not exist: %s\n", filePath)
    os.Exit(1)
  }

  // Slurp the file contents into memory
  dat, err := ioutil.ReadFile(filePath)
  if err != nil {
    fmt.Printf("could not read file: %s: %s\n", err, filePath)
    os.Exit(1)
  }

  if len(dat) == 0 {
    fmt.Printf("file is empty: %s\n", filePath)
    os.Exit(1)
  }

  var buf bytes.Buffer
  err = gzip.Compress(&buf, []byte(dat))
  if err != nil {
    fmt.Printf("could not compress file: %s\n", err)
    os.Exit(1)
  }

  resp, err := Wayfinder.JobService.CreateJob(context.TODO(), &proto.CreateJobRequest{
    Data:       buf.Bytes(),
    Compressed: true,
  })
  if err != nil {
    fmt.Printf("could not create job: %s\n", err)
    os.Exit(1)
  }

  fmt.Printf("Successfully created new job with ID=%d\n", resp.Id)
}
