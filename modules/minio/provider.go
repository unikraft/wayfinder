package minio
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

  
  "github.com/minio/minio-go/v7"
  "github.com/minio/minio-go/v7/pkg/credentials"

  "github.com/erda-project/erda-infra/base/logs"
  "github.com/erda-project/erda-infra/base/servicehub"
)

// Interface .
type Interface interface {
  Client()     (*minio.Client)
}

var (
  interfaceType = reflect.TypeOf((*Interface)(nil)).Elem()
  clientType    = reflect.TypeOf((*minio.Client)(nil))
)

type config struct {
  Endpoint        string `file:"endpoint"          env:"MINIO_ENDPOINT"`
  AccessKeyId     string `file:"access_key_id"     env:"MINIO_ACCESS_KEY_ID"`
  SecretAccessKey string `file:"secret_access_key" env:"MINIO_SECRET_ACCESS_KEY"`
  UseSSL          bool   `file:"use_ssl"           env:"MINIO_USE_SSL"            default:"false"`
  DefaultBucket   string `file:"default_bucket"    env:"MINIO_DEFAULT_BUCKET"     default:"wayfinder"`
  DefaultLocation string `file:"default_location"  env:"MINIO_DEFAULT_LOCATION"   default:"eu-west-1"`
}

type provider struct {
  Cfg     *config
  Log      logs.Logger
  client  *minio.Client
  clients  map[string]*minio.Client
  lock     sync.Mutex
}

func (p *provider) Init(ctx servicehub.Context) error {
  if len(p.Cfg.DefaultBucket) <= 0 {
    return nil
  }

  c, err := p.OpenBucket(ctx, p.Cfg.DefaultBucket)

  if err != nil {
    return err

  }
  p.client = c
  return nil
}

func (p *provider) Client() *minio.Client {
  if p.client != nil {
    return p.client
  }
  c, _ := p.OpenBucket(context.TODO(), p.Cfg.DefaultBucket)
  return c
}

func (p *provider) OpenBucketLocation(ctx context.Context, bucket, location string) (*minio.Client, error) {
  p.lock.Lock()
  defer p.lock.Unlock()

  if c, ok := p.clients[bucket]; ok {
    return c, nil
  }

  c, err := minio.New(p.Cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(p.Cfg.AccessKeyId, p.Cfg.SecretAccessKey, ""),
		Secure: p.Cfg.UseSSL,
	})

	err = c.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		if _, err := c.BucketExists(ctx, bucket); err != nil {
      return nil, fmt.Errorf("could not make bucket: %s", err)
    }
	}

  p.clients[bucket] = c
  return c, nil
}

func (p *provider) OpenBucket(ctx context.Context, bucket string) (*minio.Client, error) {
  return p.OpenBucketLocation(ctx, bucket, p.Cfg.DefaultLocation)
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
  switch {
    case ctx.Service() == "minio-client",
         ctx.Type() == clientType:
      return p.Client()
  }

  return p
}

func init() {
  servicehub.Register("minio", &servicehub.Spec{
    Services:             []string{
      "minio",
      "minio-client",
    },
    Types:                []reflect.Type{
      interfaceType,
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
        clients:            make(map[string]*minio.Client),
      }
    },
  })
}
