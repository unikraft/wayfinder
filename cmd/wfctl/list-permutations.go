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

	"github.com/spf13/cobra"

	"github.com/unikraft/wayfinder/api/proto"
)

var (
	listPerCmd = &cobra.Command{
		Use:                   "list-permutations [OPTIONS...] jobID offset limit",
		Aliases:               []string{"lp"},
		Short:                 `List permutations for a job.`,
		Run:                   doLpCmd,
		Args:                  cobra.ExactArgs(3),
		DisableFlagsInUseLine: true,
	}

	// pauseCfg = &pauseConfig{}
)

// type pauseConfig struct {
// }

func init() {
}

// doLpCmd
func doLpCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
	}

	jobId, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid id:", err)
		os.Exit(1)
	}

	offset, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Invalid offset:", err)
		os.Exit(1)
	}

	limit, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Invalid limit:", err)
		os.Exit(1)
	}

	var response *proto.ListPermutationsResponse
	response, err = Wayfinder.JobService.ListPermutations(context.TODO(), &proto.ListPermutationsRequest{
		Id:     int64(jobId),
		Offset: int64(offset),
		Limit:  int64(limit),
	})
	if err != nil {
		fmt.Printf("could not get list: %s\n", err)
		os.Exit(1)
	}

	// TODO Pretty print this
	for _, perm := range response.Permutations {
		fmt.Printf("Permutation with checksum %s\n", perm.Checksum)
		fmt.Printf("%+v\n", perm.Builds)
		fmt.Printf("%+v\n", perm.Params)
		fmt.Printf("%+v\n\n", perm.Tests)
	}
}
