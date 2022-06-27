package spec

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

type TestMonitorSpec struct {
	Name     string `yaml:"name"`
	Commands string `yaml:"commands"`
}

type TestKernelSpec struct {
	Memory string `yaml:"memory"`
	Args   string `yaml:"args"`
	Cores  uint64 `yaml:"cores"`
}

type TestBenchToolSpec struct {
	Image        string            `yaml:"image"`
	Devices      []string          `yaml:"devices"`
	Capabilities []string          `yaml:"capabilities"`
	Commands     string            `yaml:"commands"`
	Monitors     []TestMonitorSpec `yaml:"monitors"`
	Cores        uint64            `yaml:"cores"`
	StartDelay   uint64            `yaml:"start_delay"` // in seconds
	Duration     uint64            `yaml:"duration"`    // in seconds
	Environment  map[string]string `yaml:"environment"`
}

type TestResultSpec struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
	Type string `yaml:"type"` // only "int", "integer", "str", "string", "float" or "bool"
}

type TestSpec struct {
	Name      string            `yaml:"name"`
	Kernel    TestKernelSpec    `yaml:"kernel"`
	BenchTool TestBenchToolSpec `yaml:"benchtool"`
	Results   []TestResultSpec  `yaml:"results"`
}
