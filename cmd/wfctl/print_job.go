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
  "context"
  "strconv"

  "gopkg.in/yaml.v2"
  "github.com/spf13/cobra"

  "github.com/unikraft/wayfinder/api/proto"
)

var (
  printJobCmd = &cobra.Command{
    Use:                   "print-job [OPTIONS...] ID",
    Aliases:               []string{"pj"},
    Short:                 `Outputs information about a particular job`,
    Run:                   doPrintJobCmd,
    Args:                  cobra.ExactArgs(1),
    DisableFlagsInUseLine: true,
  }
  showConfig bool
)

func init() {
  printJobCmd.PersistentFlags().BoolVarP(
    &showConfig,
    "config",
    "c",
    false,
    "Only show input configuration",
  )
}

// doPrintJobCmd 
func doPrintJobCmd(cmd *cobra.Command, args []string) {
  if len(args) == 0 || args[0] == "" {
    cmd.Help()
  }
  
  jobId, err := strconv.ParseInt(args[0], 10, 64)
  if err != nil {
    fmt.Printf("invalid job ID: %s", err)
    os.Exit(1)
  }

  resp, err := Wayfinder.JobService.GetJob(context.TODO(), &proto.GetJobRequest{
    Id: jobId,
  })
  if err != nil {
    fmt.Printf("could not retrieve job: %s\n", err)
    os.Exit(1)
  }

  var out []byte

  if showConfig {
    out = []byte(resp.Job.Config)
  } else {
    out, err = yaml.Marshal(&resp.Job)
    if err != nil {
      fmt.Printf("could not serialise YAML: %s\n", err)
      os.Exit(1)
    }
  }

  fmt.Printf("%s", out)
}
