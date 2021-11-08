package hostconfig
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
  "strings"
  "reflect"
  "context"

  "github.com/erda-project/erda-infra/base/logs"
  "github.com/erda-project/erda-infra/base/servicehub"

  "github.com/unikraft/wayfinder/pkg/sys"
  "github.com/unikraft/wayfinder/modules/postgres"

  "github.com/unikraft/wayfinder/internal/models"
  "github.com/unikraft/wayfinder/internal/parsecpusets"
)

type config struct {
  CpuSets           string `file:"cpu_sets"         env:"HOSTCONFIG_CPU_SETS"`
  ScalingGovernor   string `file:"scaling_governor" env:"HOSTCONFIG_SCALING_GOVERNOR" default:"performance"`
  Procfs          []string `file:"procfs"           env:"HOSTCONFIG_PROCFS"`
}

type provider struct {
  Cfg    *config
  Log     logs.Logger
  DB      postgres.Interface `autowired:"postgres"`
  procfs *ProcFs
}

func (p *provider) Init(ctx servicehub.Context) error {
  // TODO: Clean up when the servce exits, restoring the host back to its
  // original configuration.
  // p.exited = ctx.Hub().Events().Exited()

  ctx.AddTask(p.PrepareEnvironmentTask)
  ctx.AddTask(p.SaveHostDetailsTask)

  return nil
}

func (p *provider) PrepareEnvironmentTask(ctx context.Context) error {
  var err error
  cpuSets, err := parsecpusets.ParseCpuSets(p.Cfg.CpuSets)
  if err != nil {
    return fmt.Errorf("could not parse CPU sets: %s", err)
  }

  procVals := make(map[string]string)

  for _, val := range p.Cfg.Procfs {
    fields := strings.Split(val, ":")
    if len(fields) < 2 {
      p.Log.Warnf("could not parse: %s: invalid format", val)
      continue
    }

    key := strings.TrimSpace(fields[0])
    value := strings.TrimSpace(fields[1])

    procVals[key] = value
  }

  p.procfs = &ProcFs{
    Log: p.Log.Sub("procfs"),
  }

  // Setup the host environment
  err = p.procfs.Prepare(p.Cfg.ScalingGovernor, cpuSets, procVals)
  if err != nil {
    return fmt.Errorf("could not prepare host environment: %s", err)
  }

  return nil
}

func (p *provider) SaveHostDetailsTask(ctx context.Context) error {
  dmiUuid, err := sys.GetSysDmiUUID()
  if err != nil {
    return err
  }

  // Check if the host already exists in the database
  host := &models.Host{}
  if err := p.DB.DB().Where("dmi_uuid = ?", dmiUuid).First(&host).Error; err != nil {
    cpuInfo, err := sys.GetCpuInfo()
    if err != nil {
      return fmt.Errorf("could not retrieve contents of lscpu: %s", err)
    }

    // Otherwise, save it
    host.Hostname, err        = os.Hostname()
    if err != nil {
      return fmt.Errorf("could not get hostname: %s", err)
    }

    host.DmiUUID              = dmiUuid
    host.Architecture         = cpuInfo.Architecture
    host.ByteOrder            = cpuInfo.ByteOrder
    host.AddressSizesPhysical = cpuInfo.AddressSizesPhysical
    host.AddressSizesVirtual  = cpuInfo.AddressSizesVirtual
    host.CPUs                 = cpuInfo.CPUs
    host.ThreadsPerCore       = cpuInfo.ThreadsPerCore
    host.CoresPerSocket       = cpuInfo.CoresPerSocket
    host.Sockets              = cpuInfo.Sockets
    host.NUMAnodes            = cpuInfo.NUMAnodes
    host.VendorID             = cpuInfo.VendorID
    host.CPUFamily            = cpuInfo.CPUFamily
    host.Model                = cpuInfo.Model
    host.ModelName            = cpuInfo.ModelName
    host.Stepping             = cpuInfo.Stepping
    host.CPUMHz               = cpuInfo.CPUMHz
    host.CPUMaxMHz            = cpuInfo.CPUMaxMHz
    host.CPUMinMHz            = cpuInfo.CPUMinMHz
    host.BogoMIPS             = cpuInfo.BogoMIPS
    host.Virtualization       = cpuInfo.Virtualization
    host.L1dCache             = cpuInfo.L1dCache
    host.L1iCache             = cpuInfo.L1iCache
    host.L2Cache              = cpuInfo.L2Cache
    host.L3Cache              = cpuInfo.L3Cache
    host.NUMANode0CPUs        = cpuInfo.NUMANode0CPUs
    host.NUMANode1CPUs        = cpuInfo.NUMANode1CPUs
    host.Flags                = cpuInfo.Flags

    host, err = p.DB.Repos().Hosts().CreateHost(host)
    if err != nil {
      return fmt.Errorf("could not save host: %s", err)
    }
  }

  // TODO: What happens if the machine is upgraded?

  return nil
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
  return p
}

func init() {
  servicehub.Register("hostconfig", &servicehub.Spec{
    Services:             []string{
      "hostconfig",
    },
    Types:                []reflect.Type{  },
    Dependencies:         []string{},
    OptionalDependencies: []string{
      "service-register",
    },
    Description:            "",
    ConfigFunc:             func() interface{} {
      return &config{}
    },
    Creator:                func() servicehub.Provider {
      return &provider{}
    },
  })
}
