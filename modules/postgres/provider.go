package postgres
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
  "time"
  "reflect"

  "gorm.io/gorm"
  gormlog "gorm.io/gorm/logger"
  "gorm.io/driver/postgres"

  "github.com/erda-project/erda-infra/base/logs"
  "github.com/erda-project/erda-infra/base/servicehub"

  "github.com/unikraft/wayfinder/internal/models"
  "github.com/unikraft/wayfinder/internal/repositories"
)

// Interface .
type Interface interface {
  DB()    *gorm.DB
  Repos() *repositories.Repositories
}

var (
  interfaceType = reflect.TypeOf((*Interface)(nil)).Elem()
  gormType      = reflect.TypeOf((*gorm.DB)(nil))
  repoType      = reflect.TypeOf((*repositories.Repositories)(nil))
)

type config struct {
  BaseDSN       string        `file:"base_dsn"       env:"POSTGRES_BASE_DSN"`
  Host          string        `file:"host"           env:"POSTGRES_HOST"         default:"localhost"`
  Port          string        `file:"port"           env:"POSTGRES_PORT"         default:"3306"`
  Username      string        `file:"username"       env:"POSTGRES_USER"         default:"root"`
  LogQueries    bool          `file:"log_queries"    env:"POSTGRES_LOG_QUERIES"  default:"false"`
  Password      string        `file:"password"       env:"POSTGRES_PASSWORD"     default:""`
  Database      string        `file:"database"       env:"POSTGRES_DATABASE"`
  MaxIdleConns  uint64        `file:"max_idle_conns" env:"POSTGRES_MAXIDLECONNS" default:"1"`
  MaxOpenConns  uint64        `file:"max_open_conns" env:"POSTGRES_MAXOPENCONNS" default:"2"`
  MaxLifeTime   time.Duration `file:"max_lifetime"   env:"POSTGRES_MAXLIFETIME"  default:"30m"`
  SSLEnable     bool          `file:"ssl_enable"     env:"POSTGRES_SSL_ENABLE"   default:"false"`
  SSLRootCert   string        `file:"ssl_rootcert"   env:"POSTGRES_SSL_ROOTCERT"`
  SSLCert       string        `file:"ssl_cert"       env:"POSTGRES_SSL_CERT"`
  SSLKey        string        `file:"ssl_key"        env:"POSTGRES_SSL_KEY"`
  EncryptionKey string        `file:"encryption_key" env:"POSTGRES_ENCRYPTION_KEY"`
  MaxRetries    int           `file:"max_retries"    env:"POSTGRES_MAX_RETRIES"  default:"3"`
}

type provider struct {
  Cfg   *config
  Log    logs.Logger
  db    *gorm.DB
  repos *repositories.Repositories
}

func (p *provider) DB() *gorm.DB {
  return p.db
}

func (p *provider) Repos() *repositories.Repositories {
  return p.repos
}

func (c *config) baseDSN() (string, error) {
  if c.BaseDSN != "" {
    return c.BaseDSN, nil
  }

  // connect to default postgres instance first
  baseDSN := fmt.Sprintf(
    "user=%s password=%s port=%s host=%s",
    c.Username,
    c.Password,
    c.Port,
    c.Host,
  )

  baseDSN = baseDSN + " sslmode=disable"

  return baseDSN, nil
}

func (p *provider) Init(ctx servicehub.Context) error {
  baseDSN, err := p.Cfg.baseDSN()
  if err != nil {
    return fmt.Errorf("could not create base DSN: %s", err)
  }

  postgresDSN := baseDSN + " database=postgres"
  targetDSN := baseDSN + " database=" + p.Cfg.Database

  cfg := &gorm.Config{
    FullSaveAssociations: true,
  }
  if p.Cfg.LogQueries {
    cfg.Logger = gormlog.Default.LogMode(gormlog.Info)
  }

  defaultDB, err := gorm.Open(postgres.Open(postgresDSN), cfg)
  if err != nil {
    return fmt.Errorf("could not open default database: %s", err)
  }

  // attempt to create the database
  if p.Cfg.Database != "" {
    defaultDB.Exec(fmt.Sprintf("CREATE DATABASE %s;", p.Cfg.Database))
  }

  // open the database connection
  res, err := gorm.Open(postgres.Open(targetDSN), cfg)

  // retry the connection
  retryCount := 0
  timeout, _ := time.ParseDuration("5s")

  if err != nil {
    for {
      time.Sleep(timeout)
      res, err = gorm.Open(postgres.Open(targetDSN), cfg)

      if retryCount > p.Cfg.MaxRetries {
        return err
      }

      if err == nil {
        goto migrate
      }

      retryCount++
    }
  }

migrate:
  p.Log.Infof("migrating database...")
  p.db = res

  err = models.AutoMigrate(p.db)
  if err != nil {
    return fmt.Errorf("could not migrate database: %s", err)
  }

  var key [32]byte

  for i, b := range []byte(p.Cfg.EncryptionKey) {
    key[i] = b
  }

  p.repos = repositories.NewRepositories(p.db, &key)

  return nil
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
  switch {
    case ctx.Service() == "postgres-client",
         ctx.Type() == gormType:
    return p.db
  
    case ctx.Service() == "postgres-repos",
         ctx.Type() == repoType:
    return p.repos
  }
  return p
}

func init() {
  servicehub.Register("postgres", &servicehub.Spec{
    Services:             []string{
      "postgres",
      "postgres-client",
      "postgres-repos",
    },
    Types:                []reflect.Type{
      interfaceType,
      gormType,
      repoType,
    },
    Dependencies:         []string{},
    OptionalDependencies: []string{
      "service-register",
    },
    Description:            "postgres",
    ConfigFunc:             func() interface{} {
      return &config{}
    },
    Creator:                func() servicehub.Provider {
      return &provider{}
    },
  })
}
