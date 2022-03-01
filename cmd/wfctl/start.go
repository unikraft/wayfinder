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

  "github.com/spf13/cobra"
  "github.com/thediveo/enumflag"

  "github.com/unikraft/wayfinder/api/proto"
)

// TODO iterate over proto.JobScheduler to generate enumflag.Flag from IDs
type SchedulerType enumflag.Flag

const (
  Grid SchedulerType = iota
  Random
  Bayesian
)

type IsolLevelType enumflag.Flag

const (
  // Tasks can be scheduled on any core
  NoIsol IsolLevelType = iota

  // Tasks will be scheduled on cores belonging to the same socket
  DiffSocketIsol

  // Tasks will only be scheduled on the same numa node
  DiffNumaIsol

  // Tasks will be scheduled on cores sharing the same L3 cache
  DiffCacheIsol

  // Tasks will be scheduled on cores on the same Socket/NUMA/L3 cache
  FullIsol
)

type IsolSplitType enumflag.Flag

const (
  SplitBoth IsolSplitType = iota
  SplitBuilds
  SplitTests
)

type StartJobConfig struct {
  Scheduler        proto.JobScheduler
  IsolLevel        proto.JobIsolLevel
  IsolSplit        proto.JobIsolSplit
  PermutationLimit int
  Repeats          int
  DryRun           bool
}

var (
  startCmd = &cobra.Command{
    Use:                   "start [OPTIONS...] ID",
    Aliases:               []string{"sj"},
    Short:                 `Start a job on the wayfinder server.`,
    Run:                   doStartCmd,
    Args:                  cobra.ExactArgs(1),
    DisableFlagsInUseLine: true,
  }

  // Evaluated additional configuration for the job
  jobCfg = &StartJobConfig{}
  
  // Map scheduler enumeration values to their textual representations
  // (value identifiers).
  SchedulerTypeIds = map[SchedulerType][]string{
    Grid:     {"grid"},
    Random:   {"random"},
    Bayesian: {"bayesian"},
  }
  
  // Map isolation level enumeration values to their textual representations
  // (value identifiers).
  IsolLevelTypeIds = map[IsolLevelType][]string{
    NoIsol:         {"none"},
    DiffNumaIsol:   {"numa"},
    DiffCacheIsol:  {"cache"},
    DiffSocketIsol: {"socket"},
    FullIsol:       {"full"},
  }
  
  // Map isolation split enumeration values to their textual representations
  // (value identifiers).
  IsolSplitTypeIds = map[IsolSplitType][]string{
    SplitBoth:   {"both"},
    SplitBuilds: {"builds"},
    SplitTests:  {"tests"},
  }
)

func init() {
  startCmd.PersistentFlags().VarP(
    enumflag.New(
      &jobCfg.Scheduler,
      "grid",
      SchedulerTypeIds,
      enumflag.EnumCaseInsensitive,
    ),
    "scheduler",
    "s",
    "Specify the scheduler for job permutations.",
  )

  startCmd.PersistentFlags().IntVarP(
    &jobCfg.Repeats,
    "repeats",
    "r",
    0,
    "Number of times to repeat a permutation. Useful for random search. (default 0)",
  )

  startCmd.PersistentFlags().BoolVarP(
    &jobCfg.DryRun,
    "dry-run",
    "D",
    false,
    "Specify whether to save output to the database or not.",
  )

  startCmd.PersistentFlags().VarP(
    enumflag.New(
      &jobCfg.IsolLevel,
      "none",
      IsolLevelTypeIds,
      enumflag.EnumCaseInsensitive,
    ),
    "isol-level",
    "i",
    "Specify the level of isolation for job permutations.",
  )

  startCmd.PersistentFlags().VarP(
    enumflag.New(
      &jobCfg.IsolSplit,
      "both",
      IsolSplitTypeIds,
      enumflag.EnumCaseInsensitive,
    ),
    "isol-split",
    "x",
    "Specify the split of isolation for job permutations.",
  )

  startCmd.PersistentFlags().IntVarP(
    &jobCfg.PermutationLimit,
    "permutation-limit",
    "l",
    0,
    "Number of permutations to iterate over (grid - powers of 2; random - exact).  Zero means all.",
  )

  // TODO: Flag to skip existing permutations of this job seen in the database
  // startCmd.PersistentFlags().BoolVarP(
  //   &jobCfg.SkipExisting,
  //   "skip-exisiting",
  //   "S",
  //   false,
  //   "",
  // )
}

// doStartCmd
func doStartCmd(cmd *cobra.Command, args []string) {
  if len(args) == 0 {
    cmd.Help()
  }
  
  jobId, err := strconv.Atoi(args[0])
  if err != nil || jobId == 0 {
    fmt.Printf("invalid job ID: %d", jobId)
    os.Exit(1)
  }

  _, err = Wayfinder.JobService.StartJob(context.TODO(), &proto.StartJobRequest{
    Id:               int64(jobId),
    Scheduler:        jobCfg.Scheduler,
    IsolLevel:        jobCfg.IsolLevel,
    IsolSplit:        jobCfg.IsolSplit,
    PermutationLimit: int64(jobCfg.PermutationLimit),
    Repeats:          uint64(jobCfg.Repeats),
    DryRun:           jobCfg.DryRun,
  })
  if err != nil {
    fmt.Printf("could not start job: %s\n", err)
    os.Exit(1)
  }

  fmt.Printf("Successfully started job with ID=%d\n", jobId)
}
