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

	"github.com/spf13/cobra"

	"github.com/unikraft/wayfinder/api/proto"
)

var (
  flushCmd = &cobra.Command{
    Use:                   "flush [OPTIONS...] 'task'/'job'/'all'",
    Aliases:               []string{"fq"},
    Short:                 `Flush elements from the Redis queues.`,
    Run:                   doFqCmd,
    Args:                  cobra.ExactArgs(1),
    DisableFlagsInUseLine: true,
  }

  flushCfg = &flushConfig{}
)
  
type flushConfig struct {
}

func init() {
}

// doFqCmd
func doFqCmd(cmd *cobra.Command, args []string) {
  if len(args) == 0 {
      cmd.Help()
  }

  flushType := args[0]
  var flushTypeId proto.RedisQueueType

  switch flushType {
    case "task":
      flushTypeId = proto.RedisQueueType_REDIS_QUEUE_TYPE_PERMUTATION
    case "job":
      flushTypeId = proto.RedisQueueType_REDIS_QUEUE_TYPE_JOB
    case "all":
      flushTypeId = proto.RedisQueueType_REDIS_QUEUE_TYPE_ALL
    default:
      fmt.Println("Invalid flush type:", flushType)
      os.Exit(1)
  }

  _, err := Wayfinder.JobService.FlushRedisQueue(context.TODO(), &proto.FlushRedisQueueRequest{
      QueueType: flushTypeId,
  })
  if err != nil {
      fmt.Printf("could not start flushing: %s\n", err)
      os.Exit(1)
  }

  fmt.Printf("Sucessfully started flushing queue %s\n", flushType)
}
  