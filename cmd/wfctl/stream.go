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
	streamCmd = &cobra.Command{
		Use:                   "stream [OPTIONS...] last_nr_of_lines",
		Aliases:               []string{"sr"},
		Short:                 `Prints the logs.`,
		Run:                   doStreamCmd,
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
	}

	streamCfg = &streamConfig{}
)

type streamConfig struct {
	Level         int32
	PermutationId string
	JobId         string
}

func init() {
	streamCmd.PersistentFlags().Int32VarP(
		&streamCfg.Level,
		"level",
		"l",
		0,
		"Limit the output shown. 0-6, higher is more restrictive.",
	)

	streamCmd.PersistentFlags().StringVarP(
		&streamCfg.JobId,
		"jobId",
		"j",
		"",
		"Filter by a specific job using it's id.",
	)

	streamCmd.PersistentFlags().StringVarP(
		&streamCfg.PermutationId,
		"permutationId",
		"p",
		"",
		"Filter by a specific permutation using it's id.",
	)
}

// doPpCmd
func doStreamCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
	}

	lastNrOfLines, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("invalid number of lines:", err)
		os.Exit(1)
	}

	response, err := Wayfinder.JobService.RedirectDebug(context.TODO(), &proto.RedirectDebugRequest{
		DebugLevel:    proto.DebugLevel(streamCfg.Level),
		JobId:         streamCfg.JobId,
		PermutationId: streamCfg.PermutationId,
		LastNrOfLines: int64(lastNrOfLines),
	})
	if err != nil {
		fmt.Printf("could not retrieve the debug: %s\n", err)
		os.Exit(1)
	}

	fmt.Print(response.Output)
}
