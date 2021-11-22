package scheduler
// SPDX-License-Identifier: BSD-3-Clause
//
// Authors: Alexander Jung <a.jung@lancs.ac.uk>
//
// Copyright (c) 2021, Lancaster University.  All rights reserved.
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
  "fmt"
  "time"
  "context"
  "reflect"

  "github.com/adjust/rmq/v4"
  "github.com/go-redis/redis/v8"
  "github.com/erda-project/erda-infra/base/servicehub"
  
  "github.com/unikraft/wayfinder/pkg/sys"
  "github.com/erda-project/erda-infra/base/logs"
  "github.com/unikraft/wayfinder/modules/tester"
  "github.com/unikraft/wayfinder/modules/builder"
  "github.com/unikraft/wayfinder/modules/postgres"
  "github.com/unikraft/wayfinder/internal/coremap"
)

type Interface interface {
  TaskQueue() rmq.Queue
  CoreMap()  *coremap.CoreMap
}

var (
  queueType = reflect.TypeOf((*rmq.Queue)(nil))
)

type config struct {
  PrefetchLimit   int64         `file:"prefetch_limit" env:"SCHEDULER_PREFETCH_LIMIT" default:"1000"`
  PollDuration    time.Duration `file:"poll_duration"  env:"SCHEDULER_POLL_DURATION"  default:"100ms"`
  GraceTime       time.Duration `file:"grace_time"     env:"SCHEDULER_GRACE_TIME"     default:"5s"`
  MaxRetries      int64         `file:"max_retries"    env:"SCHEDULER_MAX_RETRIES"    default:"3"`
}

type provider struct {
  Cfg        *config
  Log         logs.Logger
  Redis      *redis.Client       `autowired:"redis-client"`
  DB          postgres.Interface `autowired:"postgres"`
  Builder    *builder.Service    `autowired:"builder"`
  Tester     *tester.Service     `autowired:"tester"`

  // Internal
  taskQueue   rmq.Queue
  redisErr    chan error
  coreMap    *coremap.CoreMap
}

// Init this is optional
func (p *provider) Init(ctx servicehub.Context) error {
  // Initialize global task queue
  p.redisErr = make(chan error)
  conn, err := rmq.OpenConnectionWithRedisClient("tasks", p.Redis, p.redisErr)
  if err != nil {
    return fmt.Errorf("could not create connect to redis: %s", err)
  }

  p.taskQueue, err = conn.OpenQueue("task-queue")
  if err != nil {
    return fmt.Errorf("could not create job queue: %s", err)
  }

  if err := p.taskQueue.StartConsuming(p.Cfg.PrefetchLimit, p.Cfg.PollDuration); err != nil {
    return fmt.Errorf("cannot start consuming jobs: %s", err)
  }

  cpuInfo, err := sys.GetCpuInfo()
  if err != nil {
    return fmt.Errorf("could not get host CPU information: %s", err)
  }

  // Initialize the coremap
  p.coreMap, err = coremap.NewFromStr([]string{
    cpuInfo.NUMANode0CPUs,
    cpuInfo.NUMANode1CPUs,
  })
  if err != nil {
    return fmt.Errorf("could not initialize core map: %s", err)
  }

  return nil
}

func (p *provider) Run(ctx context.Context) error {

  // TODO: Implement a *real* scheduler.  For now this consumer accepts any task
  // it receives and attempts to complete it.

  if _, err := p.taskQueue.AddConsumer("task-consumer", NewTaskConsumer(p)); err != nil {
    return fmt.Errorf("cannot initialize task consumer: %s", err)
  }

  // TODO: Catch redis channel errors?

  return nil
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
  switch {
    case ctx.Service() == "scheduler-queue",
         ctx.Type() == queueType:
    return p.taskQueue
  }
  return p
}

func (p *provider) CoreMap() *coremap.CoreMap {
  return p.coreMap
}

func init() {
  servicehub.Register("scheduler", &servicehub.Spec{
    Services:             []string{
      "scheduler",
      "scheduler-queue",
    },
    Types:                []reflect.Type{
      queueType,
    },
    Dependencies:         []string{
      "postgres",
      "builder",
      "tester",
    },
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
