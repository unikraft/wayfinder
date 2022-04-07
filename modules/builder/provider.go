package builder

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
	"reflect"
	"time"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/pkg/transport"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/unikraft/wayfinder/api/proto"
	"github.com/unikraft/wayfinder/modules/container"
	"github.com/unikraft/wayfinder/modules/postgres"
)

type config struct {
	OutputDir   string        `file:"outputdir"    env:"BUILDER_OUTPUTDIR"    default:"/tmp/wayfinder/builds"`
	KillTimeout time.Duration `file:"kill_timeout" env:"BUILDER_KILL_TIMEOUT" default:"0s"`
	Retries     int64         `file:"retries"      env:"BUILDER_RETRIES"      default:"1"`
}

type provider struct {
	// ctx       servicehub.Context
	Cfg       *config
	Log       logs.Logger
	Register  transport.Register
	service   *Service
	DB        postgres.Interface `autowired:"postgres"`
	Container *container.Service `autowired:"container"`
}

func (p *provider) Init(ctx servicehub.Context) error {
	p.service = &Service{
		p:      p,
		builds: make(map[string]*build),
	}

	if p.Register != nil {
		proto.RegisterBuilderServiceImp(p.Register, p.service)
	}

	return nil
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	switch {
	case ctx.Service() == "builder",
		ctx.Service() == "wayfinder.BuilderService",
		ctx.Type() == proto.BuilderServiceServerType(),
		ctx.Type() == proto.BuilderServiceHandlerType():
		return p.service
	}

	return p
}

func init() {
	servicehub.Register("builder", &servicehub.Spec{
		Services: []string{
			"builder",
			"wayfinder.BuilderService",
		},
		Types: []reflect.Type{
			proto.BuilderServiceClientType(),
			proto.BuilderServiceServerType(),
			proto.BuilderServiceHandlerType(),
		},
		Dependencies: []string{
			"postgres",
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
