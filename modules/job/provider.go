package job

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
	"fmt"
	"reflect"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/pkg/transport"
	"github.com/go-redis/redis/v8"

	"github.com/unikraft/wayfinder/api/proto"
	"github.com/unikraft/wayfinder/modules/container"
	"github.com/unikraft/wayfinder/modules/postgres"
)

type config struct {
	PrefetchLimit int64         `file:"prefetch_limit" env:"JOB_PREFETCH_LIMIT" default:"1000"`
	PollDuration  time.Duration `file:"poll_duration"  env:"JOB_POLL_DURATION"  default:"100ms"`
}

type provider struct {
	Cfg       *config
	Log       logs.Logger
	Register  transport.Register
	Redis     *redis.Client      `autowired:"redis-client"`
	DB        postgres.Interface `autowired:"postgres"`
	TaskQueue rmq.Queue          `autowired:"scheduler-queue"`
	Container *container.Service `autowired:"container"` // SavedDir
	JobQueue  rmq.Queue
	RedisErr  chan error
	service   *service
}

func (p *provider) Init(ctx servicehub.Context) error {
	// Initialize global job queue
	p.RedisErr = make(chan error)
	conn, err := rmq.OpenConnectionWithRedisClient("jobs", p.Redis, p.RedisErr)
	if err != nil {
		return fmt.Errorf("could not create connect to redis: %s", err)
	}

	p.JobQueue, err = conn.OpenQueue("job-queue")
	if err != nil {
		return fmt.Errorf("could not create job queue: %s", err)
	}

	if p.Register != nil {
		p.service = &service{p: p}

		proto.RegisterJobServiceImp(p.Register, p.service)
	}

	if err := p.JobQueue.StartConsuming(p.Cfg.PrefetchLimit, p.Cfg.PollDuration); err != nil {
		return fmt.Errorf("cannot start consuming jobs: %s", err)
	}

	if _, err := p.JobQueue.AddConsumer("job-consumer", NewJobConsumer(p)); err != nil {
		return fmt.Errorf("cannot initialize job consumer: %s", err)
	}

	return nil
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	switch {
	case ctx.Service() == "job",
		ctx.Service() == "wayfinder.JobService",
		ctx.Type() == proto.JobServiceServerType(),
		ctx.Type() == proto.JobServiceHandlerType():
		return p.service
	}
	return p
}

func (p *provider) Start() error {
	return nil
}

func (p *provider) Close() error {
	<-p.JobQueue.StopConsuming()
	return nil
}

func init() {
	servicehub.Register("job", &servicehub.Spec{
		Services: []string{
			"job",
			"wayfinder.JobService",
		},
		Types: []reflect.Type{
			proto.JobServiceClientType(),
			proto.JobServiceServerType(),
			proto.JobServiceHandlerType(),
		},
		Dependencies: []string{
			"redis",
			"postgres-repos",
			"scheduler-queue",
		},
		OptionalDependencies: []string{
			"service-register",
		},
		Description: "",
		ConfigFunc: func() interface{} {
			return &config{}
		},
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
