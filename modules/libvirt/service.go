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
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	libvirt "github.com/libvirt/libvirt-go"
	libvirtxml "github.com/libvirt/libvirt-go-xml"

	"github.com/unikraft/wayfinder/api/proto"
	"github.com/unikraft/wayfinder/internal/bridge"
	"github.com/unikraft/wayfinder/internal/metrics"
	"github.com/unikraft/wayfinder/internal/strutils"
	"github.com/unikraft/wayfinder/pkg/sys"
)

const (
	x86_64   = "x86_64"
	pci440fx = "pc-i440fx-3.1"
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

	monitors []*proto.TestMonitor
	p        *provider
	pid      int
	fakePid  int
	args     string
	// uuid    string
	timer time.Time
	// runtime *time.Duration
	ip      net.IP
	subnet  string
	bridge  *bridge.Bridge
	config  *libvirtxml.Domain
	domain  *libvirt.Domain
	measure bool
}

func (s *Service) NewDomain(fakePid int, uuid, kernel, initrd, args string, inputDisks []*proto.BuildOutputDiskImage,
	cores []uint64, memoryValue uint, memoryUnit string, monitors []*proto.TestMonitor) (*Domain, error) {
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

	disks := []libvirtxml.DomainDisk{}

	for _, inputDisk := range inputDisks {
		libvirtDisk := libvirtxml.DomainDisk{
			Source: &libvirtxml.DomainDiskSource{
				File: &libvirtxml.DomainDiskSourceFile{
					File: inputDisk.Path,
				},
			},
			Target: &libvirtxml.DomainDiskTarget{
				Dev: inputDisk.Name,
				Bus: "virtio",
			},
		}

		switch inputDisk.Type {
		case proto.BuildOutputDiskImageType_BUILD_OUTPUT_DISK_RAW:
			libvirtDisk.Device = "disk"
			libvirtDisk.Driver = &libvirtxml.DomainDiskDriver{
				Name: "qemu",
				Type: "raw",
			}
		}

		disks = append(disks, libvirtDisk)
	}

	logDir := path.Join(s.p.Cfg.LogDir, uuid)
	logFile := path.Join(logDir, "domain.log") // TODO: configuration opt?
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, os.ModePerm)
	}

	serial := libvirtxml.DomainSerial{
		Log: &libvirtxml.DomainChardevLog{
			File:   logFile,
			Append: "on",
		},
	}

	config := &libvirtxml.Domain{
		Type:        "kvm",
		Name:        uuid,
		UUID:        uuid,
		Title:       uuid,
		Description: uuid,
		Devices: &libvirtxml.DomainDeviceList{
			Emulator:   s.p.Cfg.Emulator,
			Consoles:   []libvirtxml.DomainConsole{console},
			Serials:    []libvirtxml.DomainSerial{serial},
			Interfaces: []libvirtxml.DomainInterface{iface},
			Disks:      disks,
		},
		OS: &libvirtxml.DomainOS{
			Kernel: kernel,
			Initrd: initrd,
			Type: &libvirtxml.DomainOSType{
				Type:    "hvm",
				Arch:    x86_64,
				Machine: archToMachine[x86_64],
			},
		},
		VCPU: &libvirtxml.DomainVCPU{
			Placement: "static",
			CPUSet:    strings.Trim(strings.Join(strings.Fields(fmt.Sprint(cores)), ","), "[]"),
			Value:     uint(len(cores)),
		},
		Memory: &libvirtxml.DomainMemory{
			Value: memoryValue,
			Unit:  memoryUnit,
		},
		OnCrash:    "destroy",
		OnPoweroff: "destroy",
		OnReboot:   "destroy",
	}

	domain := &Domain{
		Measurable: metrics.NewMeasurable(),
		p:          s.p,
		fakePid:    fakePid,
		args:       args,
		config:     config,
		monitors:   monitors,
	}

	return domain, nil
}

func (d *Domain) SetArgs(args string) error {
	replaceVars := map[string]string{}

	if d.ip != nil {
		// TODO: Do this conversion earlier? not in the main test loop?
		gwAddr, gwNet, err := net.ParseCIDR(d.p.Cfg.Subnet)
		if err != nil {
			return fmt.Errorf("could not parse configuration mask: %s", err)
		}

		replaceVars["$WAYFINDER_DOMAIN_IP_ADDR"] = d.ip.String()
		replaceVars["$WAYFINDER_DOMAIN_IP_GW_ADDR"] = gwAddr.String()
		replaceVars["$WAYFINDER_DOMAIN_IP_MASK"] = fmt.Sprintf(
			"%d.%d.%d.%d",
			gwNet.Mask[0],
			gwNet.Mask[1],
			gwNet.Mask[2],
			gwNet.Mask[3],
		)
	}

	for k, v := range replaceVars {
		args = strings.Replace(args, k, v, -1)
	}

	d.config.OS.Cmdline = args

	return nil
}

// func (d *Domain) AddBridge(bridge *bridge.Bridge) error {
func (d *Domain) CreateBridge(name, hostIface, subnet, stateDir string) error {
	var err error

	// TODO sometimes fails
	retries := 0
	for retries < 5 {
		d.bridge = bridge.New(name, hostIface, subnet, stateDir)
		d.ip, err = d.bridge.Create(d.fakePid, false)
		d.subnet = subnet
		if err != nil {
			retries++
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		// fmt.Printf("Allocating bridge: bridge=%#v d.ip=%#v d.subnet=%#v d.fakePid=%#v\n", d.bridge, d.ip, d.subnet, d.fakePid)
		return fmt.Errorf("could not allocate IP: %s", err)
	}

	d.p.Log.Debugf("Domain IP: %s", d.ip.String())

	return nil
}

func (d *Domain) CreateDefaultBridge() error {
	bridgeStateDir := path.Join(d.p.Cfg.BridgeStateDir, d.config.Name)
	return d.CreateBridge(
		d.p.Cfg.Bridge,
		d.p.Cfg.HostIface,
		d.p.Cfg.Subnet,
		bridgeStateDir,
	)
}

func (d *Domain) FakePid() int {
	return d.fakePid
}

func (d *Domain) Pid() (int, error) {
	if d.pid > 0 {
		return d.pid, nil
	}

	var pid int
	pidFile := path.Join(d.p.Cfg.SockDir, fmt.Sprintf("%s.pid", d.config.Name))

	file, err := os.Open(pidFile)
	if err != nil {
		return 0, fmt.Errorf("error accessing domain's pid: %s", err)
	}

	_, err = fmt.Fscanf(file, "%d", &pid)
	if err != nil {
		return 0, fmt.Errorf("could not read pid file: %s: %s", pidFile, err)
	}

	d.pid = pid

	return pid, nil
}

func (d *Domain) Wait() (time.Duration, error) {
	getStateRetries := 3

	// Continiously poll the status of the domain to check whether it has exited
	for {
		state, reason, err := d.domain.GetState()
		if err != nil {
			getStateRetries--
			if getStateRetries == 0 {
				return time.Duration(0), fmt.Errorf("could not get state after retries: %s", err)
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

		// Domain finished unsuccessfully (e.g. crashed)
		// Most of these states do not have a reason, so it remains unknown
		case libvirt.DOMAIN_NOSTATE:
			switch libvirt.DomainNostateReason(reason) {
			case libvirt.DOMAIN_NOSTATE_UNKNOWN:
				return time.Duration(0), fmt.Errorf("domain is in no state for an unknown reason; state: %d, code: %d", state, reason)
			default:
				return time.Duration(0), fmt.Errorf("domain exited with state: %d, code: %d", state, reason)
			}
		case libvirt.DOMAIN_BLOCKED:
			switch libvirt.DomainBlockedReason(reason) {
			case libvirt.DOMAIN_BLOCKED_UNKNOWN:
				return time.Duration(0), fmt.Errorf("domain blocked with an unknown reason; state: %d, code: %d", state, reason)
			default:
				return time.Duration(0), fmt.Errorf("domain exited with state: %d, code: %d", state, reason)
			}
		case libvirt.DOMAIN_CRASHED:
			switch libvirt.DomainCrashedReason(reason) {
			case libvirt.DOMAIN_CRASHED_UNKNOWN:
				return time.Duration(0), fmt.Errorf("domain crashed with an unknown reason; state: %d, code: %d", state, reason)
			case libvirt.DOMAIN_CRASHED_PANICKED:
				return time.Duration(0), fmt.Errorf("domain crashed with a panic; state: %d, code: %d", state, reason)
			default:
				return time.Duration(0), fmt.Errorf("domain crashed and exited with state: %d, code: %d", state, reason)
			}
		case libvirt.DOMAIN_PMSUSPENDED:
			switch libvirt.DomainPMSuspendedReason(reason) {
			case libvirt.DOMAIN_PMSUSPENDED_UNKNOWN:
				return time.Duration(0), fmt.Errorf("domain suspended with an unknown reason; state: %d, code: %d", state, reason)
			default:
				return time.Duration(0), fmt.Errorf("domain suspended and exited with state: %d, code: %d", state, reason)
			}
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
	if err := d.CreateDefaultBridge(); err != nil {
		return fmt.Errorf("could not attach to default bridge: %s", err)
	}

	if err := d.SetArgs(d.args); err != nil {
		return fmt.Errorf("could not set args: %s", err)
	}

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

// Pin the VMM process to defined set of cores
func (d *Domain) PinVMMToCores(cores []uint64) error {
	pid, err := d.Pid()
	if err != nil || pid <= 0 {
		return fmt.Errorf("could not set VMM affinity, no PID: %s", err)
	}

	corelist := strutils.JoinUint64(cores, ",")
	d.p.Log.Debugf("Pinning VMM (pid=%d) to cores [%s]", pid, corelist)

	if _, err := sys.SetAffinity(cores, pid); err != nil {
		return fmt.Errorf(
			"could not set VMM affinity (pid=%d) to cores [%s]: %s",
			pid,
			corelist,
			err,
		)
	}

	return nil
}

// Start the libvirt domain
func (d *Domain) Start() error {
	d.timer = time.Now()

	if err := d.domain.Resume(); err != nil {
		return fmt.Errorf("could not start domain: %s", err)
	}

	timedout, err := d.WaitForStateChange(libvirt.DOMAIN_RUNNING, d.p.Cfg.Timeout)
	if timedout {
		return fmt.Errorf("timed out waiting for domain to start running")
	} else if err != nil {
		return fmt.Errorf("domain is not running: %s", err)
	}

	if err := d.InitMeasurements(); err != nil {
		return fmt.Errorf("could initialize measurements: %s", err)
	}

	go func() {
		for d.measure {
			time.Sleep(d.p.Cfg.MeasureFreq)
			d.MeasureResources()
		}
	}()

	return nil
}

// Pause the libvirt domain
func (d *Domain) Pause() error {
	d.measure = false

	return d.domain.Suspend()
}

// Resume the libvirt domain
func (d *Domain) Resume() error {
	d.measure = true

	return d.domain.Resume()
}

func (d *Domain) WaitForStateChange(expected libvirt.DomainState, timeout time.Duration) (bool, error) {
	tout := time.After(timeout)

	for {
		select {
		case <-tout:
			return true, nil
		default:
			state, _, err := d.domain.GetState()
			if err != nil {
				return false, err
			}

			if state == expected {
				return false, nil
			}
		}
	}
}

// Kill the libvirt domain
func (d *Domain) Destroy() error {
	d.measure = false

	if err := d.domain.Destroy(); err != nil {
		return fmt.Errorf("unable to destroy domain: %s", err)
	}

	if err := d.domain.Undefine(); err != nil {
		return fmt.Errorf("could not undefine domain: %s", err)
	}

	// TODO sometimes fails
	retries := 0
	var err error
	for retries < 5 {
		// Clean up the network
		if err = d.bridge.Delete(d.fakePid, d.ip); err != nil {
			retries++
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	if retries >= 5 {
		return fmt.Errorf("could not delete veth pair: %s", err)
	}

	return nil
}

func (d *Domain) IP() net.IP {
	return d.ip
}

func (d *Domain) InitMeasurements() error {
	d.measure = true

	if err := d.CpuLookup(); err != nil {
		return fmt.Errorf("could not look up CPU cores: %s", err)
	}

	if err := d.MemLookup(); err != nil {
		return fmt.Errorf("could not look up memory: %s", err)
	}

	if err := d.NetLookup(); err != nil {
		return fmt.Errorf("could not look up networks: %s", err)
	}

	if err := d.MetricLookup(); err != nil {
		return fmt.Errorf("could not look up custom metrics: %s", err)
	}

	return nil
}

func (d *Domain) MeasureResources() []error {
	var errs []error

	if err := d.CpuMeasure(); err != nil {
		errs = append(errs, fmt.Errorf("could not measure CPU: %s", err))
	}

	if err := d.MemMeasure(); err != nil {
		errs = append(errs, fmt.Errorf("could not measure memory: %s", err))
	}

	// Throws "could not find network device" error
	if err := d.NetMeasure(); err != nil {
		errs = append(errs, fmt.Errorf("could not measure network: %s", err))
	}

	for _, monitor := range d.monitors {
		if err := d.MetricMeasure(monitor.Name, monitor.Commands); err != nil {
			errs = append(errs, fmt.Errorf("could not measure custom metrics: %s", err))
		}
	}

	return errs
}

func (d *Domain) GetResourceMeasurements() map[string]interface{} {
	res := make(map[string]interface{})

	for k, v := range d.CpuPrint() {
		valFloat, err := strconv.ParseFloat(v, 64)
		if err == nil {
			res[k] = valFloat
		} else {
			res[k] = v
		}
	}

	for k, v := range d.MemPrint() {
		valFloat, err := strconv.ParseFloat(v, 64)
		if err == nil {
			res[k] = valFloat
		} else {
			res[k] = v
		}
	}

	for k, v := range d.NetPrint() {
		valFloat, err := strconv.ParseFloat(v, 64)
		if err == nil {
			res[k] = valFloat
		} else {
			res[k] = v
		}
	}

	for _, monitor := range d.monitors {
		for k, v := range d.MetricPrint(monitor.Name) {
			valFloat, err := strconv.ParseFloat(v, 64)
			if err == nil {
				res[k] = valFloat
			} else {
				res[k] = v
			}
		}
	}

	return res
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
