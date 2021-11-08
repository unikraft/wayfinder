package libvirt
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
  "sync"
  "reflect"
  "context"

	libvirt "github.com/libvirt/libvirt-go"

  "github.com/erda-project/erda-infra/base/logs"
  "github.com/erda-project/erda-infra/pkg/transport"
  "github.com/erda-project/erda-infra/base/servicehub"
)

// Interface .
// type Interface interface {
//   Client() (*libvirt.Connect)
// }

var (
  // interfaceType = reflect.TypeOf((*Interface)(nil)).Elem()
  clientType    = reflect.TypeOf((*libvirt.Connect)(nil))
)

type config struct {
  Endpoint       string `file:"endpoint"        env:"LIBVIRT_ENDPOINT"        default:"qemu:///system"`
  Emulator       string `file:"emulator"        env:"LIBVIRT_EMULATOR"        default:"/usr/bin/qemu-system-x86_64"`
  SockDir        string `file:"sockdir"         env:"LIBVIRT_SOCKDIR"         default:"/run/libvirt/qemu/"`
  HostIface      string `yaml:"host_iface"      env:"LIBVIRT_HOST_IFACE"      default:"eth0"`
  Bridge         string `yaml:"bridge"          env:"LIBVIRT_BRIDGE"          default:"wayfinder0"`
  BridgeStateDir string `yaml:"bridge_statedir" env:"LIBVIRT_BRIDGE_STATEDIR" default:"/var/lib/wayfinder/bridges"`
  LogDir         string `yaml:"logdir"          env:"LIBVIRT_LOGDIR"          default:"/var/lib/wayfinder/logs"`
  Subnet         string `yaml:"subnet"          env:"LIBVIRT_SUBNET"          default:"172.88.0.1/16"`
}

type provider struct {
  Cfg       *config
  Log        logs.Logger
  client    *libvirt.Connect
  Register   transport.Register
  clients    map[string]*libvirt.Connect
  service   *Service
  lock       sync.Mutex
}

// eventloop keeps the connection to libvirt daemon active, removing this
// results in client timeouts
func eventloop() {
	for {
		libvirt.EventRunDefaultImpl()
	}
}

func (p *provider) Init(ctx servicehub.Context) error {
  var err error

  if p.Register != nil {
    p.service = &Service{p:p}
    // proto.RegisterBuilderServiceImp(p.Register, p.Service)
  }

  c, err := p.Connect(context.TODO(), p.Cfg.Endpoint)
  if err != nil {
    return fmt.Errorf("could not connect to libvirt: %s", err)
  }

  go eventloop()

  p.client = c
  
  return nil
}
func (p *provider) Client() *libvirt.Connect {
  return p.client
}

func (p *provider) Connect(ctx context.Context, endpoint string) (*libvirt.Connect, error) {
  p.lock.Lock()
  defer p.lock.Unlock()

  if c, ok := p.clients[endpoint]; ok {
    return c, nil
  }

	err := libvirt.EventRegisterDefaultImpl()
	if err != nil {
		return nil, fmt.Errorf("failed to register event loop: %s", err)
	}

  c, err := libvirt.NewConnect(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to libvirt: %s", err)
	}

  p.clients[endpoint] = c
  return c, nil
}

func (p *provider) Close(ctx context.Context, endpoint string) error {
  p.lock.Lock()
  defer p.lock.Unlock()

  c, ok := p.clients[endpoint]
  if !ok {
    return nil
  }

  _, err := c.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection: %s", err)
	}

  delete(p.clients, endpoint)

  return nil  
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
  switch {
    case ctx.Service() == "libvirt",
         ctx.Service() == "wayfinder.LibvirtService":
        //  ctx.Type() == clientType:
      return p.service
    
    case ctx.Service() == "libvirt-client",
         ctx.Type() == clientType:
      return p.Client()
  }

  return p
}

func init() {
  servicehub.Register("libvirt", &servicehub.Spec{
    Services:             []string{
      "libvirt",
      "libvirt-client",
    },
    Types:                []reflect.Type{
      // interfaceType,
      clientType,
    },
    Dependencies:         []string{},
    OptionalDependencies: []string{
      "service-register",
    },
    Description:            "",
    ConfigFunc:             func() interface{} {
      return &config{}
    },
    Creator:                func() servicehub.Provider {
      return &provider{
        clients:            make(map[string]*libvirt.Connect),
      }
    },
  })
}
