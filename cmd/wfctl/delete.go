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
	deleteCmd = &cobra.Command{
		Use:                   "delete [OPTIONS...] jobID",
		Aliases:               []string{"dl"},
		Short:                 `Delete a job from the database.`,
		Run:                   doDlCmd,
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
	}

	deleteCfg = &deleteConfig{}
)

type deleteConfig struct {
	Purge   bool
	Cascade bool
}

func init() {
	deleteCmd.PersistentFlags().BoolVarP(
		&deleteCfg.Purge,
		"purge",
		"p",
		false,
		"Permanently delete entries, instead of marking them as deleted. (WIP)",
	)

	deleteCmd.PersistentFlags().BoolVarP(
		&deleteCfg.Cascade,
		"cascade",
		"c",
		false,
		"Cascade a deletion to all elements that relate to a job.",
	)
}

// doPpCmd
func doDlCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
	}

	jobId, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid job ID:", err)
		os.Exit(1)
	}

	_, err = Wayfinder.JobService.DeleteJob(context.TODO(), &proto.DeleteJobRequest{
		Id:      int64(jobId),
		Purge:   deleteCfg.Purge,
		Cascade: deleteCfg.Cascade,
	})
	if err != nil {
		fmt.Printf("could not start deleting: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sucessfully deleted jobs\n")
}
