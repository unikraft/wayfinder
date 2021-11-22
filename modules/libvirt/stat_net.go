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
  "strings"

  libvirt "github.com/libvirt/libvirt-go"
  libvirtxml "github.com/libvirt/libvirt-go-xml"
  "github.com/unikraft/wayfinder/pkg/proc"
  "github.com/unikraft/wayfinder/internal/metrics"
)

func (d *Domain) NetLookup() error {
  var ifs []string
  xmldoc, _ := d.domain.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
  domcfg := &libvirtxml.Domain{}
  domcfg.Unmarshal(xmldoc)

  if domcfg.Devices == nil {
    return fmt.Errorf("no devices for domain %s", d.config.UUID)
  } else if domcfg.Devices.Interfaces == nil {
    return fmt.Errorf("no device interfaces for domain %s", d.config.UUID)
  }

  for _, devInterface := range domcfg.Devices.Interfaces {
    if devInterface.Target != nil {
      ifs = append(ifs, devInterface.Target.Dev)
    }
  }

  d.AddMeasurement("net_interfaces", metrics.CreateMeasurement(ifs))

  return nil
}

func (d *Domain) NetMeasure() error {
  pid, err := d.Pid()
  if err != nil {
    return fmt.Errorf("domain has no pid: %s", err)
  }

  // get stats from net/dev for domain interfaces
  ifs := d.GetMetricStringArray("net_interfaces")
  statsSum := proc.ProcPIDNetDev{}

  for _, devname := range ifs {
    devStats := proc.GetProcPIDNetDev(d.p.Cfg.ProcFS, pid, devname)

    statsSum.ReceivedBytes += devStats.ReceivedBytes
    statsSum.ReceivedPackets += devStats.ReceivedPackets
    statsSum.ReceivedErrs += devStats.ReceivedErrs
    statsSum.ReceivedDrop += devStats.ReceivedDrop
    statsSum.ReceivedFifo += devStats.ReceivedFifo
    statsSum.ReceivedFrame += devStats.ReceivedFrame
    statsSum.ReceivedCompressed += devStats.ReceivedCompressed
    statsSum.ReceivedMulticast += devStats.ReceivedMulticast

    statsSum.TransmittedBytes += devStats.TransmittedBytes
    statsSum.TransmittedPackets += devStats.TransmittedPackets
    statsSum.TransmittedErrs += devStats.TransmittedErrs
    statsSum.TransmittedDrop += devStats.TransmittedDrop
    statsSum.TransmittedFifo += devStats.TransmittedFifo
    statsSum.TransmittedColls += devStats.TransmittedColls
    statsSum.TransmittedCarrier += devStats.TransmittedCarrier
    statsSum.TransmittedCompressed += devStats.TransmittedCompressed
  }

  d.AddMeasurement("net_ReceivedBytes", metrics.CreateMeasurement(uint64(statsSum.ReceivedBytes)))
  d.AddMeasurement("net_ReceivedPackets", metrics.CreateMeasurement(uint64(statsSum.ReceivedPackets)))
  d.AddMeasurement("net_ReceivedErrs", metrics.CreateMeasurement(uint64(statsSum.ReceivedErrs)))
  d.AddMeasurement("net_ReceivedDrop", metrics.CreateMeasurement(uint64(statsSum.ReceivedDrop)))
  d.AddMeasurement("net_ReceivedFifo", metrics.CreateMeasurement(uint64(statsSum.ReceivedFifo)))
  d.AddMeasurement("net_ReceivedFrame", metrics.CreateMeasurement(uint64(statsSum.ReceivedFrame)))
  d.AddMeasurement("net_ReceivedCompressed", metrics.CreateMeasurement(uint64(statsSum.ReceivedCompressed)))
  d.AddMeasurement("net_ReceivedMulticast", metrics.CreateMeasurement(uint64(statsSum.ReceivedMulticast)))
  d.AddMeasurement("net_TransmittedBytes", metrics.CreateMeasurement(uint64(statsSum.TransmittedBytes)))
  d.AddMeasurement("net_TransmittedPackets", metrics.CreateMeasurement(uint64(statsSum.TransmittedPackets)))
  d.AddMeasurement("net_TransmittedErrs", metrics.CreateMeasurement(uint64(statsSum.TransmittedErrs)))
  d.AddMeasurement("net_TransmittedDrop", metrics.CreateMeasurement(uint64(statsSum.TransmittedDrop)))
  d.AddMeasurement("net_TransmittedFifo", metrics.CreateMeasurement(uint64(statsSum.TransmittedFifo)))
  d.AddMeasurement("net_TransmittedColls", metrics.CreateMeasurement(uint64(statsSum.TransmittedColls)))
  d.AddMeasurement("net_TransmittedCarrier", metrics.CreateMeasurement(uint64(statsSum.TransmittedCarrier)))
  d.AddMeasurement("net_TransmittedCompressed", metrics.CreateMeasurement(uint64(statsSum.TransmittedCompressed)))

  return nil
}

func (d *Domain) NetPrint() map[string]string {
  ifsRaw := d.GetMetricStringArray("net_interfaces")
  interfaces := strings.Join(ifsRaw, ";")

  return map[string]string{
    "net_receivedBytes": d.GetMetricDiffUint64("net_ReceivedBytes", true),
    "net_transmittedBytes": d.GetMetricDiffUint64("net_TransmittedBytes", true),
    "net_receivedPackets": d.GetMetricDiffUint64("net_ReceivedPackets", true),
    "net_receivedErrs": d.GetMetricDiffUint64("net_ReceivedErrs", true),
    "net_receivedDrop": d.GetMetricDiffUint64("net_ReceivedDrop", true),
    "net_receivedFifo": d.GetMetricDiffUint64("net_ReceivedFifo", true),
    "net_receivedFrame": d.GetMetricDiffUint64("net_ReceivedFrame", true),
    "net_receivedCompressed": d.GetMetricDiffUint64("net_ReceivedCompressed", true),
    "net_receivedMulticast": d.GetMetricDiffUint64("net_ReceivedMulticast", true),
    "net_transmittedPackets": d.GetMetricDiffUint64("net_TransmittedPackets", true),
    "net_transmittedErrs": d.GetMetricDiffUint64("net_TransmittedErrs", true),
    "net_transmittedDrop": d.GetMetricDiffUint64("net_TransmittedDrop", true),
    "net_transmittedFifo": d.GetMetricDiffUint64("net_TransmittedFifo", true),
    "net_transmittedColls": d.GetMetricDiffUint64("net_TransmittedColls", true),
    "net_transmittedCarrier": d.GetMetricDiffUint64("net_TransmittedCarrier", true),
    "net_transmittedCompressed": d.GetMetricDiffUint64("net_TransmittedCompressed", true),
    "net_interfaces": interfaces,
  }
}
