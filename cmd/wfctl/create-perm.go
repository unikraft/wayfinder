package main

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
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/unikraft/wayfinder/api/proto"
)
  
var (
  cpCmd = &cobra.Command{
    Use:                   "create-permutation [OPTIONS...] ID",
    Aliases:               []string{"cp"},
    Short:                 `Create a single unique permutation linked to a job.`,
    Run:                   doCpCmd,
    Args:                  cobra.ExactArgs(1),
    DisableFlagsInUseLine: true,
  }

  cpCfg = &createPermutationConfig{}
)

type createPermutationConfig struct {
  jobId    int
  params []string
}

func init() {
   cpCmd.PersistentFlags().StringArrayVarP(
    &cpCfg.params,
    "set",
    "s",
    []string{},
    "Parameters set within the permutation.",
  )
}

func parseSet(args []string) []*proto.Param {
  params := make([]*proto.Param, 0)
  
  for _, arg := range args {
    elems := strings.Split(arg, "=")
    param := proto.Param{}

    param.Name = elems[0]
    value, err := strconv.ParseInt(elems[1], 10, 64)
    if err == nil {
      param.Type = "int"
      param.ValueInt = value
    } else {
      param.Type = "str"
      param.ValueStr = elems[1]
    }


    params = append(params, &param)
  }

  return params
}

// doStartCmd
func doCpCmd(cmd *cobra.Command, args []string) {
  if len(args) == 0 {
    cmd.Help()
  }
  
  jobId, err := strconv.Atoi(args[0])
  if err != nil || jobId == 0 {
    fmt.Printf("invalid job ID: %d", jobId)
    os.Exit(1)
  }

  _, err = Wayfinder.JobService.CreatePermutationJob(context.TODO(), &proto.CreatePermutationJobRequest{
    Id:     int64(jobId),
    Params: parseSet(cpCfg.params),
  })
  if err != nil {
    fmt.Printf("could not start permutation: %s\n", err)
    os.Exit(1)
  }

  fmt.Printf("Successfully added permutation to job with ID=%d\n", jobId)
}
