package container

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
	"os"
	"reflect"

	"github.com/opencontainers/runc/libcontainer"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/pkg/transport"
)

type config struct {
	ContainerRootDir string   `file:"rootdir"          env:"CONTAINER_ROOTDIR"         default:"/var/lib/wayfinder/containers"`
	CacheDir         string   `file:"cachedir"         env:"CONTAINER_CACHEDIR"        default:"/var/lib/wayfinder/cache"`
	SavedDir         string   `file:"saveddir"         env:"CONTAINER_SAVEDDIR"        default:"/var/lib/wayfinder/saved"`
	LogDir           string   `file:"logdir"           env:"CONTAINER_LOGDIR"          default:"/var/lib/wayfinder/logs"`
	RegistryAddr     string   `yaml:"registry"         env:"CONTAINER_REGISTRY"        default:"localhost:5000"`
	HostIface        string   `yaml:"host_iface"       env:"CONTAINER_HOST_IFACE"      default:"eth0"`
	Bridge           string   `yaml:"bridge"           env:"CONTAINER_BRIDGE"          default:"wayfinder0"`
	BridgeStateDir   string   `yaml:"bridge_statedir"  env:"CONTAINER_BRIDGE_STATEDIR" default:"/var/lib/wayfinder/bridges"`
	Subnet           string   `yaml:"subnet"           env:"CONTAINER_SUBNET"          default:"172.88.0.1/16"`
	AuthType         string   `yaml:"auth_type"        env:"CONTAINER_AUTH_TYPE"       default:"anonymous"`
	AuthUsername     string   `yaml:"auth_username"    env:"CONTAINER_AUTH_USERNAME"   default:"wayfinder"`
	AuthPassword     string   `yaml:"auth_password"    env:"CONTAINER_AUTH_PASSWORD"   default:"wayfinder"`
	AuthToken        string   `yaml:"auth_token"       env:"CONTAINER_AUTH_TOKEN"      default:""`
	Environment      []string `yaml:"environment"`
}

type Provider struct {
	Cfg      *config
	Log      logs.Logger
	Register transport.Register
	service  *Service
	image    *Image
	factory  libcontainer.Factory
}

var (
	factoryType = reflect.TypeOf((libcontainer.Factory)(nil))
)

func (p *Provider) Init(ctx servicehub.Context) error {
	var err error

	if p.Register != nil {
		p.service = &Service{P: p}
		p.image = &Image{P: p}
		// proto.RegisterBuilderServiceImp(p.Register, p.Service)
	}

	p.factory, err = libcontainer.New(
		p.Cfg.CacheDir,
		libcontainer.Cgroupfs,
		// There's a special "hidden" entry method in wayfinderd
		// see: cmd/wayfinderd/runc.go
		libcontainer.InitArgs(os.Args[0], "runc-init"),
	)
	if err != nil {
		return fmt.Errorf("could not initialize container factory: %s", err)
	}

	return nil
}

func (p *Provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	switch {
	case ctx.Service() == "container",
		ctx.Service() == "wayfinder.ContainerService":
		// ctx.Type() == factoryType:
		return p.service
	case ctx.Service() == "container-factory",
		ctx.Type() == factoryType:
		return p.factory
	}

	return p
}

func init() {
	servicehub.Register("container", &servicehub.Spec{
		Services: []string{
			"container",
			"container-factory",
		},
		Types:        []reflect.Type{},
		Dependencies: []string{},
		OptionalDependencies: []string{
			"service-register",
		},
		Description: "",
		ConfigFunc: func() interface{} {
			return &config{}
		},
		Creator: func() servicehub.Provider {
			return &Provider{}
		},
	})
}
