package bridge

// SPDX-License-Identifier: BSD-3-Clause
//
// Authors: Alexander Jung <a.jung@lancs.ac.uk>
//
// Copyright (c) 2020, Lancaster University.  All rights reserved.
//               2021, Unikraft UG.  All rights reserved.
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
	"net"

	"github.com/lancs-net/netns/bridge"
	"github.com/lancs-net/netns/network"
	"github.com/opencontainers/runtime-spec/specs-go"
)

type Bridge struct {
	Name      string
	Interface string
	Subnet    string
	StateDir  string
	netOpt    network.Opt
	brOpt     bridge.Opt
	client    *network.Client
}

// Init prepares netns
func New(name, iface, subnet, stateDir string) *Bridge {
	var err error

	b := &Bridge{
		Name:      name,
		Interface: iface,
		Subnet:    subnet,
		StateDir:  stateDir,
	}

	// Create the bridge using netns
	b.netOpt = network.Opt{
		ContainerInterface: b.Interface,
		BridgeName:         b.Name,
		StateDir:           b.StateDir,
	}

	b.brOpt = bridge.Opt{
		Name:   b.Name,
		IPAddr: b.Subnet,
	}

	b.client, err = network.New(b.netOpt)
	if err != nil {
		return nil
	}

	return b
}

// Create a veth pair with the bridge
func (b *Bridge) Create(pid int, configure bool) (net.IP, error) {
	return b.client.Create(&specs.State{
		Pid: pid,
	}, b.brOpt, "", configure)
}

// Delete a veth pair from the bridge
func (b *Bridge) Delete(pid int, ip net.IP) error {
	return b.client.Delete(pid, ip)
}
