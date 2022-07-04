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

	libvirt "github.com/libvirt/libvirt-go"
	"github.com/unikraft/wayfinder/internal/metrics"
	"github.com/unikraft/wayfinder/pkg/proc"
)

const pagesize = 4096

func (d *Domain) MemLookup() error {
	memStats, err := d.domain.MemoryStats(uint32(libvirt.DOMAIN_MEMORY_STAT_NR), 0)
	if err != nil {
		return fmt.Errorf("could not get memory stats: %s", err)
	}

	var total, unused, used uint64

	for _, stat := range memStats {
		if stat.Tag == int32(libvirt.DOMAIN_MEMORY_STAT_UNUSED) {
			unused = stat.Val
		}
		if stat.Tag == int32(libvirt.DOMAIN_MEMORY_STAT_AVAILABLE) ||
			stat.Tag == int32(libvirt.DOMAIN_MEMORY_STAT_ACTUAL_BALLOON) {
			total = stat.Val
		}
	}

	used = total - unused
	d.AddMeasurement("ram_total", metrics.CreateMeasurement(total))
	d.AddMeasurement("ram_used", metrics.CreateMeasurement(used))

	return nil
}

func (d *Domain) MemMeasure() error {
	pid, err := d.Pid()
	if err != nil {
		return fmt.Errorf("domain has no pid: %s", err)
	}

	stats := proc.GetProcPIDStat(d.p.Cfg.ProcFS, pid)

	d.AddMeasurement("ram_vsize", metrics.CreateMeasurement(uint64(stats.VSize)))
	d.AddMeasurement("ram_rss", metrics.CreateMeasurement(uint64(stats.RSS*pagesize)))
	d.AddMeasurement("ram_minflt", metrics.CreateMeasurement(uint64(stats.MinFlt)))
	d.AddMeasurement("ram_cminflt", metrics.CreateMeasurement(uint64(stats.CMinFlt)))
	d.AddMeasurement("ram_majflt", metrics.CreateMeasurement(uint64(stats.MajFlt)))
	d.AddMeasurement("ram_cmajflt", metrics.CreateMeasurement(uint64(stats.CMajFlt)))

	return nil
}

func (d *Domain) MemPrint() map[string]string {
	total, _ := d.GetMetricUint64("ram_total", 0)
	used, _ := d.GetMetricUint64("ram_used", 0)
	vsize, _ := d.GetMetricUint64("ram_vsize", 0)
	rss, _ := d.GetMetricUint64("ram_rss", 0)

	return map[string]string{
		"ram_total":   total,
		"ram_used":    used,
		"ram_vsize":   vsize,
		"ram_rss":     rss,
		"ram_minflt":  d.GetMetricDiffUint64("ram_minflt", false),
		"ram_cminflt": d.GetMetricDiffUint64("ram_cminflt", false),
		"ram_majflt":  d.GetMetricDiffUint64("ram_majflt", false),
		"ram_cmajflt": d.GetMetricDiffUint64("ram_cmajflt", false),
	}
}
