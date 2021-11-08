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
  "os"
  "fmt"
  "net"
  "time"
  "path"
  // "strconv"
  "strings"

	libvirt "github.com/libvirt/libvirt-go"
  // "github.com/unikraft/wayfinder/pkg/sys"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
  // "github.com/opencontainers/runtime-spec/specs-go"

  "github.com/unikraft/wayfinder/internal/bridge"
  "github.com/unikraft/wayfinder/internal/metrics"
  // "github.com/milosgajdos/tenus"
  // "github.com/unikraft/wayfinder/internal/coremap"
)

const (
	x86_64   = "x86_64"
	pci440fx = "pc-i440fx-3.1"

  defaultMemoryValue = 64
  defaultMemoryUnit  = "MiB"
)

var (
	archToMachine = map[string]string{
		x86_64: pci440fx,
	}
)

type Service struct {
  p *provider
}

type Domain struct {
	*metrics.Measurable

  p       *provider
  pid      int
  args     string
	uuid     string
  timer    time.Time
  runtime *time.Duration
  ip       net.IP
  subnet   string
  bridge  *bridge.Bridge
  config  *libvirtxml.Domain
  domain  *libvirt.Domain
}

func (s *Service) NewDomain(pid int, uuid, kernel, initrd, args string, cores []uint64) (*Domain, error) {
  // This maintains an open door for debug purpose
	console := libvirtxml.DomainConsole{
		TTY: "/dev/pts/4",
	}

  iface := libvirtxml.DomainInterface{
    // TODO: Do we have to generate a new unique mac address?
    // MAC: &DomainInterfaceMAC{
    //   Address: "06:39:b4:00:00:46",
    // },
    Model: &libvirtxml.DomainInterfaceModel{
      Type: "virtio",
    },
    Source: &libvirtxml.DomainInterfaceSource{
      Bridge: &libvirtxml.DomainInterfaceSourceBridge{
        Bridge: s.p.Cfg.Bridge,
      },
    },
  }

  logDir := path.Join(s.p.Cfg.LogDir, uuid)
  logFile := path.Join(logDir, "domain.log") // TODO: configuration opt?
  if _, err := os.Stat(logDir); os.IsNotExist(err) {
    os.MkdirAll(logDir, os.ModePerm)
  }

  serial := libvirtxml.DomainSerial{
    // Source: &libvirtxml.DomainChardevSource{
    //   Pty: &libvirtxml.DomainChardevSourcePty{},
    // },
    // Target: &libvirtxml.DomainSerialTarget{
    //   Type: "isa-serial",
    //   Port: &port,
    //   Model: &libvirtxml.DomainSerialTargetModel{
    //     Name: "isa-serial",
    //   },
    // },
    Log: &libvirtxml.DomainChardevLog{
      File:   logFile,
      Append: "on",
    },
  }

  config := &libvirtxml.Domain{
		Type:        "kvm",
		Name:         uuid,
		UUID:         uuid,
    Title:        uuid,
    Description:  uuid,
		Devices:     &libvirtxml.DomainDeviceList{
      Emulator:   s.p.Cfg.Emulator,
			Consoles: []libvirtxml.DomainConsole{console},
      Serials:  []libvirtxml.DomainSerial{serial},
      Interfaces: []libvirtxml.DomainInterface{iface},
		},
    OS:          &libvirtxml.DomainOS{
      Kernel:     kernel,
      Initrd:     initrd,
      Type:      &libvirtxml.DomainOSType{
        Type:     "hvm",
        Arch:     x86_64,
        Machine:  archToMachine[x86_64],
      },
    },
    VCPU:         &libvirtxml.DomainVCPU{
      Placement: "static",
      CPUSet:    strings.Trim(strings.Join(strings.Fields(fmt.Sprint(cores)), ","), "[]"),
      Value:     uint(len(cores)),
    },
    Memory:      &libvirtxml.DomainMemory{
      Value:      defaultMemoryValue,
      Unit:       defaultMemoryUnit,
    },
    OnCrash:      "destroy",
    OnPoweroff:   "destroy",
    OnReboot:     "destroy",
	}

  domain := &Domain{
    p:      s.p,
    pid:    pid,
    args:   args,
    config: config,
  }

  return domain, nil
}

func (d *Domain) SetArgs(args string) {
  replaceVars := map[string]string{}

  if d.ip != nil {
    replaceVars["$WAYFINDER_DOMAIN_IP_ADDR"] = d.ip.String()
    replaceVars["$WAYFINDER_DOMAIN_IP_SUBNET"] = d.p.Cfg.Subnet
    replaceVars["$WAYFINDER_DOMAIN_IP_MASK"] = d.subnet
  }

  for k, v := range replaceVars {
    args = strings.Replace(args, k, v, -1)
  }

  d.config.OS.Cmdline = args
}

func (d *Domain) SetCores(cores []int64) error {
  // TODO: Lazy hack: use the static list of cores and re-populate a CoreMap
  // so we can figure out Numa Affinities.  We could have been passed this
  // information by the request.  Actually, thinking about, it's probably good
  // to sanity check the request.

  // cpuInfo, err := sys.GetCpuInfo()
  // if err != nil {
  //   return fmt.Errorf("could not get host CPU information: %s", err)
  // }

  // // Initialize a temporary coremap
  // coreMap, err := coremap.NewFromStr([]string{
  //   cpuInfo.NUMANode0CPUs,
  //   cpuInfo.NUMANode1CPUs,
  // })
  // if err != nil {
  //   return fmt.Errorf("could not get determine core mappings: %s", err)
  // }

  // // kernelCores := make([]*coremap.Core, len(cores))
  // // for _, coreId := range cores {
  // //   core, err := coreMap.FindCore(coreId)
  // //   if err != nil {
  // //     return fmt.Errorf("could not get determine core mappings: %s", err)
  // //   }

  // //   fmt.Printf("NUMA NODE =%d\n", core.NumaNodeId())

  // //   kernelCores = append(kernelCores, core)
  // // }

  // // Create a smaller map of numa nodes to cores we can use
  // numaNodesMap := make(map[int64][]*coremap.Core)

  // for _, coreId := range cores {
  //   core, err := coreMap.FindCore(coreId)
  //   if err != nil {
  //     return fmt.Errorf("could not get determine core mappings: %s", err)
  //   }

  //   numaNodesMap[core.NumaNodeId()] = append(
  //     numaNodesMap[core.NumaNodeId()],
  //     core,
  //   )
  //   // if _, ok := numaNodesMap[core.NumaNodeId()]; ok {
  //   // } else {
  //   //   numaNodesMap[core.NumaNodeId()] = []*coremap.Core{core}
  //   // }
  // }
  
  // cells := []libvirtxml.DomainCell{}

  // for _, numaCores := range numaNodesMap {
  //   coreIds := []string{}
  //   for _, core := range numaCores {
  //     coreIds = append(coreIds, strconv.Itoa(int(core.Id())))
  //   }

  //   coreList := strings.Join(coreIds[:], ",")
  //   cells = append(cells, libvirtxml.DomainCell{
  //     CPUs: coreList,
  //     // TODO: Cache locality awareness
  //     // Caches: []libvirtxml.DomainCellCache
  //   })
  // }

  // d.config.CPU = &libvirtxml.DomainCPU{
  //   Numa: &libvirtxml.DomainNuma{
  //     Cell: cells,
  //   },
  //   Topology: &libvirtxml.DomainCPUTopology{
  //     // TODO: Calculate correct number for Sockets:
  //     Sockets: 1,
  //     Cores:   len(cores),
  //     Threads: 1,
  //   },
  // }

  // fmt.Printf("%#v\n", cells)

  return nil
}

// func (d *Domain) AddBridge(bridge *bridge.Bridge) error {
func (d *Domain) CreateBridge(name, hostIface, subnet, stateDir string) error {
  var err error

  d.bridge = bridge.New(name, hostIface, subnet, stateDir)
  d.ip, err = d.bridge.Create(d.pid, false)
  d.subnet = subnet
  if err != nil {
    return fmt.Errorf("could not allocate IP: %s", err)
  }

  d.p.Log.Debugf("Domain IP: %s", d.ip.String())

  return nil
}

func (d *Domain) CreateDefaultBridge() error {
  bridgeStateDir := path.Join(d.p.Cfg.BridgeStateDir, d.config.Name)
  // TODO: Apparently these config values aren't being passed correctly?
  return d.CreateBridge(
    d.p.Cfg.Bridge,
    d.p.Cfg.HostIface,
    d.p.Cfg.Subnet,
    bridgeStateDir,
  )
}

func (d *Domain) Pid() int {
  return d.pid
}

// func (d *Domain) Pid() (int, error) {
//   if d.pid > 0 {
//     return d.pid, nil
//   }
// 
// 	var pid int
//   pidFile := path.Join(d.p.Cfg.SockDir, fmt.Sprintf("%s.pid", d.config.Name))
// 
//   file, err := os.Open(pidFile)
// 	if err != nil {
// 		return 0, fmt.Errorf("error accessing domain's pid: %s", err)
// 	}
// 
//   _, err = fmt.Fscanf(file, "%d", &pid)
// 	if err != nil {
// 		return 0, fmt.Errorf("could not read pid file: %s: %s", pidFile, err)
// 	}
// 
//   d.pid = pid
// 
// 	return pid, nil
// }

func (d *Domain) Wait() (time.Duration, error) {
  getStateRetries := 3
  
  // Continiously poll the status of the domain to check whether it has exited
  for {
    state, reason, err := d.domain.GetState()
    if err != nil {
      getStateRetries--
      if getStateRetries == 0 {
        return time.Duration(0), fmt.Errorf("could not get state after reties: %s", err)
      }
      // TODO: Measure impact of temporal buffer on test runtime and libvirt
      // daemon connectivity.  E.g. tcp libvirt daemon may need it, unix may not
      time.Sleep(1 * time.Second)
      continue
    }

    switch state {
      // Domain has shutdown as normally
      case libvirt.DOMAIN_SHUTDOWN,
           libvirt.DOMAIN_SHUTOFF:
        return time.Since(d.timer), nil

      // Domain lifecycle has ended abruptly with an non-running state.
      // TODO: The `reason` variable are well-defined enums, these should be
      // encoded to provide a better error message by splitting up this switch
      // into more cases and interpreting the reason.
      case libvirt.DOMAIN_NOSTATE,
           libvirt.DOMAIN_BLOCKED,
           libvirt.DOMAIN_CRASHED,
           libvirt.DOMAIN_PMSUSPENDED:
        return time.Duration(0), fmt.Errorf("domain exited with state: %d, and code: %d", state, reason)
    }
  }
}

func (d *Domain) StartAndWait() (time.Duration, error) {
  if err := d.Start(); err != nil {
    return time.Duration(0), err
  }

  return d.Wait()
}

// Initialise the libvirt domain
func (d *Domain) Init() error {
  err := d.CreateDefaultBridge()
  if err != nil {
    return fmt.Errorf("could not attach to default bridge: %s", err)
  }

  d.SetArgs(d.args)

	doc, err := d.config.Marshal()
	if err != nil {
		return fmt.Errorf("could not marshal domain to xml: %s", err)
	}

  // This prints out the XML:
  // fmt.Printf("%s\n", doc)

	d.domain, err = d.p.Client().DomainDefineXML(doc)
	if err != nil {
		return fmt.Errorf("could not initialise domain: %s", err)
	}

	return d.domain.CreateWithFlags(libvirt.DOMAIN_START_PAUSED)
}

// Start the libvirt domain
func (d *Domain) Start() error {
  d.timer = time.Now()

	return d.domain.Resume()
}

// Pause the libvirt domain
func (d *Domain) Pause() error {
	return d.domain.Suspend()
}

// Resume the libvirt domain
func (d *Domain) Resume() error {
	return d.domain.Resume()
}

// Kill the libvirt domain
func (d *Domain) Destroy() error {
	if err := d.domain.Destroy(); err != nil {
		return fmt.Errorf("unable to destroy domain: %s", err)
	}

	if err := d.domain.Undefine(); err != nil {
		return fmt.Errorf("could not undefine domain: %s", err)
	}

  // Clean up the network
  if err := d.bridge.Delete(d.pid, d.ip); err != nil {
    return fmt.Errorf("could not delete veth pair: %s", err)
  }

	return nil
}

func (d *Domain) IP() net.IP {
  return d.ip
}

func (d *Domain) CollectStatsOneTime() error {
  return nil
}

// TODO: https://gist.github.com/hodgesds/7d8354b51bea65c833817a067e45fd8d
// EventHandleLoop starts an event loop
func EventHandleLoop() error {
	err := libvirt.EventRegisterDefaultImpl()
	if err != nil {
		return fmt.Errorf("failed to register event loop: %s", err)
	}

	for {
		err := libvirt.EventRunDefaultImpl()
		if err != nil {
			return fmt.Errorf("failed to run event loop: %s", err)
		}
	}
}
