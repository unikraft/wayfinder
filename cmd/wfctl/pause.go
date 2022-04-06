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
	pauseCmd = &cobra.Command{
		Use:                   "pause [OPTIONS...] TimeMS",
		Aliases:               []string{"pp"},
		Short:                 `Pause all Redis queues for given milliseconds.`,
		Run:                   doPpCmd,
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
	}

	pauseCfg = &pauseConfig{}
)

type pauseConfig struct {
}

func init() {
}

// doPpCmd
func doPpCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
	}

	pauseTime, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid pause time:", err)
		os.Exit(1)
	}

	_, err = Wayfinder.JobService.PauseRedisQueues(context.TODO(), &proto.PauseRedisQueuesRequest{
		Time: int64(pauseTime),
	})
	if err != nil {
		fmt.Printf("could not start pausing: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sucessfully paused the queue\n")
}
