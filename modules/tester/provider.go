package tester

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

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/pkg/transport"

	"github.com/erda-project/erda-infra/base/logs"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/unikraft/wayfinder/api/proto"
	"github.com/unikraft/wayfinder/modules/container"
	"github.com/unikraft/wayfinder/modules/libvirt"
	"github.com/unikraft/wayfinder/modules/postgres"
)

type config struct {
	DefaultStartDelay     time.Duration `file:"default_start_delay"      env:"TESTER_DEFAULT_START_DELAY"      default:"5s"`
	MetricsFreq           time.Duration `file:"metrics_freq"             env:"TESTER_METRICS_FREQ"             default:"1s"`
	InfluxReceiverType    string        `file:"metrics_receiver_type"    env:"TESTER_METRICS_RECEIVER_TYPE"    default:"http"`
	InfluxUsername        string        `file:"metrics_username"         env:"TESTER_METRICS_USERNAME"         default:""`
	InfluxPassword        string        `file:"metrics_password"         env:"TESTER_METRICS_PASSWORD"         default:""`
	InfluxToken           string        `file:"metrics_token"            env:"TESTER_METRICS_TOKEN"            default:""`
	InfluxBucket          string        `file:"metrics_bucket"           env:"TESTER_METRICS_BUCKET"           default:"wayfinder"`
	InfluxOrg             string        `file:"metrics_org"              env:"TESTER_METRICS_ORG"              default:"wayfinder"`
	InfluxReceiver        string        `file:"metrics_receiver"         env:"TESTER_METRICS_RECEIVER"         default:"127.0.0.1:8086"`
	InfluxReceiverTimeout time.Duration `file:"metrics_receiver_timeout" env:"TESTER_METRICS_RECEIVER_TIMEOUT" default:"5s"`
}

type provider struct {
	Cfg           *config
	Log           logs.Logger
	Register      transport.Register
	service       *Service
	DB            postgres.Interface `autowired:"postgres"`
	Container     *container.Service `autowired:"container"`
	Libvirt       *libvirt.Service   `autowired:"libvirt"`
	metricsClient influxdb2.Client
}

func (p *provider) Init(ctx servicehub.Context) error {
	p.service = &Service{
		p:     p,
		tests: make(map[string]*test),
	}

	if p.Register != nil {
		proto.RegisterTesterServiceImp(p.Register, p.service)
	}

	p.metricsClient = influxdb2.NewClient(p.Cfg.InfluxReceiverType+"://"+p.Cfg.InfluxReceiver, p.Cfg.InfluxToken)

	_, err := p.metricsClient.Health(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to metrics server: %s", err)
	}

	return nil
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	switch {
	case ctx.Service() == "tester",
		ctx.Service() == "wayfinder.TesterService",
		ctx.Type() == proto.TesterServiceServerType(),
		ctx.Type() == proto.TesterServiceHandlerType():
		return p.service
	}

	return p
}

func (p *provider) PushMetrics(testUuid string, metrics map[string]interface{}) error {
	tags := make(map[string]string)
	tags["test_uuid"] = testUuid

	point := influxdb2.NewPoint("domain", tags, metrics, time.Now())

	p.metricsClient.WriteAPI(p.Cfg.InfluxOrg, p.Cfg.InfluxBucket).WritePoint(point)

	return nil
}

func init() {
	servicehub.Register("tester", &servicehub.Spec{
		Services: []string{
			"tester",
			"wayfinder.TesterService",
		},
		Types: []reflect.Type{},
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
