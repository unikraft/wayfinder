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
  "github.com/spf13/cobra"

  "github.com/erda-project/erda-infra/base/servicehub"

  "github.com/unikraft/wayfinder/internal/logs"
  v "github.com/unikraft/wayfinder/internal/version"

  // import all providers
  _ "github.com/erda-project/erda-infra/providers/health"
  _ "github.com/erda-project/erda-infra/providers/httpserver"
  _ "github.com/erda-project/erda-infra/providers/grpcserver"
  _ "github.com/erda-project/erda-infra/providers/grpcclient"
  _ "github.com/erda-project/erda-infra/providers/serviceregister"

  _ "github.com/unikraft/wayfinder/modules/job"
  _ "github.com/unikraft/wayfinder/modules/postgres"
  _ "github.com/unikraft/wayfinder/modules/hostconfig"

  _ "github.com/unikraft/wayfinder/api/client"
)

var (
  version   = "No version provided"
  commit    = "No commit provided"
  buildTime = "No build timestamp provided"

  configFile string
)

// Build the cobra command that handles our command line tool.
func NewRootCommand() *cobra.Command {
  rootCmd := &cobra.Command{
    Use:                   "wayfinderd -c wayfinderd.yaml",
    Short:                 `wayfinder: OS Configuration Micro-Benchmarking Framework`,
    Long:                  `
Wayfinder is a generic OS performance evaluation platform.  Wayfinder is fully
automated and ensures both the accuracy and reproducibility of results, all the
while speeding up how fast tests are run on a system. Wayfinder is easily
extensible and offers convenient APIs to:

  - Implement custom configuration space exploration techniques,
  - Add new benchmarks; and,
  - Support additional OS projects.
`,
    Run:                    doRootCmd,
    DisableFlagsInUseLine: true,
    PersistentPreRunE:     func(cmd *cobra.Command, args []string) error {
      showVer, err := cmd.Flags().GetBool("version")
      if err != nil {
        fmt.Printf("%s\n", err)
        os.Exit(0)
      }
      if showVer {
        fmt.Printf(
          "wayfinder %s (%s) built %s\n",
          version,
          commit,
          buildTime,
        )
        os.Exit(0)
      }

      return nil
    },
  }

  rootCmd.PersistentFlags().BoolP(
    "version",
    "V",
    false,
    "Show version and quit",
  )

  rootCmd.PersistentFlags().StringVarP(
    &configFile,
    "config",
    "c",
    "wayfinderd.yaml",
    "config file",
  )

  // Subcommands
  rootCmd.AddCommand(runcInitCmd)

  return rootCmd
}

// doRootCmd starts the main system
func doRootCmd(cmd *cobra.Command, args []string) {
  fmt.Printf(" _       __            _____           __           \n")
  fmt.Printf("| |     / /___ ___  __/ __(_)___  ____/ /__  _____  \n")
  fmt.Printf("| | /| / / __ `/ / / / /_/ / __ \\/ __  / _ \\/ ___/\n")
  fmt.Printf("| |/ |/ / /_/ / /_/ / __/ / / / / /_/ /  __/ /      \n")
  fmt.Printf("|__/|__/\\__,_/\\__, /_/ /_/_/ /_/\\__,_/\\___/_/   \n")
  fmt.Printf("             /____/                                 \n")
  fmt.Printf(" %s\n", v.String())

  hub := servicehub.New(servicehub.WithLogger(logs.Logger{}))

  hub.Run("wayfinderd", configFile, args...)
}

func main() {
  v.SetVersion(&v.Version{
    Version:   version,
    Commit:    commit,
    BuildTime: buildTime,
  })

  cmd := NewRootCommand()
  if err := cmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
