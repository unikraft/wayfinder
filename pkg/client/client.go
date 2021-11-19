package client
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

  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials"

  "github.com/unikraft/wayfinder/api/proto"
)

// ClientProvider provides all service clients.
type ClientProvider interface {
  JobService() proto.JobServiceClient
}

// Config contains the details of the remote gRPC server
type Config struct {
  // The server address in the format of host:port
  Addr               string

  // The server name used to verify the hostname returned by the TLS handshake
  ServerNameOverride string

  // The file containing the CA root cert file
  CAFile             string

  // block until the connection is up
  Block                bool
}

type Client struct {
  Cfg          *Config
  conn         *grpc.ClientConn
  dialopts    []grpc.DialOption
  callopts    []grpc.CallOption

  // All available services exposed as global attributes
  JobService   proto.JobServiceClient
}

// New creates a new client
func New(config *Config) (*Client, error) {
  c := &Client{
    Cfg: config,
  }

  var dialopts []grpc.DialOption

  if len(c.Cfg.CAFile) > 0 {
    creds, err := credentials.NewClientTLSFromFile(c.Cfg.CAFile, c.Cfg.ServerNameOverride)
    if err != nil {
      return nil, fmt.Errorf("fail to create tls credentials %s", err)
    }

    dialopts = append(dialopts, grpc.WithTransportCredentials(creds))
  } else {
    dialopts = append(dialopts, grpc.WithInsecure())
  }

  c.dialopts = dialopts

  dialopts = nil

  if c.Cfg.Block {
    dialopts = append(dialopts, grpc.WithBlock())
  }

  conn, err := c.NewConnect(dialopts...)
  if err != nil {
    return nil, fmt.Errorf("fail to dial: %s", err)
  }

  c.conn = conn

  // Initialize all clients
  c.JobService = proto.NewJobServiceClient(c.conn)

  return c, nil
}

func (c *Client) Get() *grpc.ClientConn {
  return c.conn
}

func (c *Client) NewConnect(opts ...grpc.DialOption) (*grpc.ClientConn, error) {
  return grpc.Dial(c.Cfg.Addr, append(opts, c.dialopts...)...)
}
