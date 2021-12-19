package container
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
  "time"
  "path"
  "strings"
  "path/filepath"

  "golang.org/x/sys/unix"
  "github.com/otiai10/copy"
  "github.com/opencontainers/runc/libcontainer"
  "github.com/erda-project/erda-infra/base/logs"
  "github.com/opencontainers/runtime-spec/specs-go"
  v1 "github.com/google/go-containerregistry/pkg/v1"
  "github.com/opencontainers/runc/libcontainer/configs"
  "github.com/opencontainers/runc/libcontainer/specconv"

  "github.com/unikraft/wayfinder/api/proto"
  "github.com/unikraft/wayfinder/internal/bridge"
  "github.com/unikraft/wayfinder/internal/printer"
)

var (
  entrypointFile = "/entrypoint.sh"
  defaultEnvironment = []string{
    "TERM=xterm",
    "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
  }
  defaultMountFlags = unix.MS_NOEXEC | unix.MS_NOSUID | unix.MS_NODEV
  defaultCapabilities = []string{
    "CAP_CHOWN",
    "CAP_DAC_OVERRIDE",
    "CAP_FSETID",
    "CAP_FOWNER",
    "CAP_MKNOD",
    "CAP_NET_RAW",
    "CAP_SETGID",
    "CAP_SETUID",
    "CAP_SETFCAP",
    "CAP_SETPCAP",
    "CAP_NET_BIND_SERVICE",
    "CAP_SYS_CHROOT",
    "CAP_KILL",
    "CAP_AUDIT_WRITE",
  }
  availableDevices = map[string]*configs.Device{
    "/dev/kvm": &configs.Device{
      Path:       "/dev/kvm",
      FileMode:   0432,
      Uid:        0,
      Gid:        104,
      DeviceRule: configs.DeviceRule{
        Type:        configs.CharDevice,
        Major:       10,
        Minor:       232,
        Permissions: "rwm",
        Allow:       true,
      },
    },
    "/dev/net/tun": &configs.Device{
      Path:       "/dev/net/tun",
      FileMode:   438,
      Uid:        0,
      Gid:        104,
      DeviceRule: configs.DeviceRule{
        Type:        configs.CharDevice,
        Major:       10,
        Minor:       200,
        Permissions: "rwm",
        Allow:       true,
      },
    },
    "/dev/random": &configs.Device{
      Path:       "/dev/random",
      FileMode:   0444,
      Uid:        0,
      Gid:        104,
      DeviceRule: configs.DeviceRule{
        Type:        configs.CharDevice,
        Major:       1,
        Minor:       8,
        Permissions: "rwm",
        Allow:       true,
      },
    },
    "/dev/urandom": &configs.Device{
      Path:       "/dev/urandom",
      FileMode:   0444,
      Uid:        0,
      Gid:        104,
      DeviceRule: configs.DeviceRule{
        Type:        configs.CharDevice,
        Major:       1,
        Minor:       9,
        Permissions: "rwm",
        Allow:       true,
      },
    },
  }
)

type Service struct {
  P *Provider
}

type Container struct {
  p               *Provider
  Log              logs.Logger
  spec             proto.Container
  image            v1.Image
  imageManifestHex string
  bridge          *bridge.Bridge
  config          *configs.Config
  env            []string
  container       libcontainer.Container
  process         *libcontainer.Process
  timer            time.Time
  exitCode         int
}

func (s *Service) NewContainer(id string) (*Container, error) {
  if len(id) == 0 {
    return nil, fmt.Errorf("container id cannot be empty")
  }

  container := &Container{
    p:    s.P,
    env:  defaultEnvironment,
    spec: proto.Container{
      Id: id,
    },
    config: &configs.Config{
      Rootfs: path.Join(s.P.Cfg.ContainerRootDir, id),
      Namespaces: configs.Namespaces([]configs.Namespace{
        {Type: configs.NEWNS},
        {Type: configs.NEWUTS},
        {Type: configs.NEWIPC},
        {Type: configs.NEWPID},
        {Type: configs.NEWNET},
      }),
      Cgroups: &configs.Cgroup{
        Name:      id,
        Parent:    "",
        Resources: &configs.Resources{
          MemorySwappiness: nil,
        },
      },
      MaskPaths: []string{
        "/proc/acpi",
        "/proc/asound",
        "/proc/kcore",
        "/proc/keys",
        "/proc/latency_stats",
        "/proc/timer_list",
        "/proc/timer_stats",
        "/proc/sched_debug",
        "/sys/firmware",
        "/proc/scsi",
      },
      ReadonlyPaths: []string{
        "/proc/bus",
        "/proc/fs",
        "/proc/irq",
        "/proc/sys",
        "/proc/sysrq-trigger",
      },
      Hostname: id,
      Mounts: []*configs.Mount{
        {
          Source:      "proc",
          Destination: "/proc",
          Device:      "proc",
          Flags:       defaultMountFlags,
        },
        {
          Source:      "tmpfs",
          Destination: "/dev",
          Device:      "tmpfs",
          Flags:       unix.MS_NOSUID | unix.MS_STRICTATIME,
          Data:        "mode=755",
        },
        {
          Source:      "devpts",
          Destination: "/dev/pts",
          Device:      "devpts",
          Flags:       unix.MS_NOSUID | unix.MS_NOEXEC,
          Data:        "newinstance,ptmxmode=0666,mode=0620,gid=5",
        },
        {
          Device:      "tmpfs",
          Source:      "shm",
          Destination: "/dev/shm",
          Data:        "mode=1777,size=65536k",
          Flags:       defaultMountFlags,
        },
        {
          Source:      "mqueue",
          Destination: "/dev/mqueue",
          Device:      "mqueue",
          Flags:       defaultMountFlags,
        },
        {
          Source:      "sysfs",
          Destination: "/sys",
          Device:      "sysfs",
          Flags:       defaultMountFlags | unix.MS_RDONLY,
        },
        {
          Destination: "/sys/fs/cgroup",
          Device:      "cgroup",
          Source:      "cgroup",
          Flags:       defaultMountFlags | unix.MS_RDONLY,
        },
      },
      Networks: []*configs.Network{
        {
          Type:    "loopback",
          Address: "127.0.0.1/0",
          Gateway: "localhost",
        },
      },
      Rlimits: []configs.Rlimit{
        {
          Type: unix.RLIMIT_NOFILE,
          Hard: uint64(1025),
          Soft: uint64(1025),
        },
      },
    },
    Log: s.P.Log.Sub(id),
  }

  // The `StartContainer` hook is the closest way to run code before the
  // process is executed by libcontainer[0].  However, this configuration is
  // passed via a JSON serialized object which uses a path path to an
  // executable script and not Go code, like that below.  As a result, we
  // must use the earliest placable hook, `CreateRuntime`, we can call a
  // within libcontainer before it runs the code specified by the run.
  // [0]: https://github.com/opencontainers/runc/blob/v1.0.0-rc92/libcontainer/standard_init_linux.go#L214-L220
  container.config.Hooks = configs.Hooks{
    configs.CreateRuntime: configs.HookList{
      configs.NewFunctionHook(func(st *specs.State) error {
        container.timer = time.Now()
        return nil
      }),
    },
  }

  return container, nil
}

func (s *Service) LoadContainer(id string) (*Container, error) {
  return nil, nil  
}

func (c *Container) SetCommands(commands string) error {
  entrypointPath := path.Join(c.config.Rootfs, entrypointFile)

  f, err := os.OpenFile(
    entrypointPath,
    os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
    os.ModePerm,
  )
  if err != nil {
    return fmt.Errorf("could not create temporary cmd file: %s", err)
  }

  // TODO: does sh support bash?
  _, err = f.WriteString("#!/usr/bin/env bash\n")
  if err != nil {
    return fmt.Errorf("could not write to temporary cmd file: %s", err)
  }

  _, err = f.WriteString(commands)
  if err != nil {
    return fmt.Errorf("could not write to temporary cmd file: %s", err)
  }

  f.Close()

  return nil
}

func (c *Container) CreateBridge(name, iface, subnet, stateDir string) {
  c.AddBridge(bridge.New(name, iface, subnet, stateDir))
}

func (c *Container) AddBridge(bridge *bridge.Bridge) {
  c.bridge = bridge

  c.config.Hooks[configs.Prestart] = configs.HookList{
    configs.NewFunctionHook(func(st *specs.State) error {
      ip, err := c.bridge.Create(st.Pid, true)
      if err != nil {
        return fmt.Errorf("could not create bridge: %s", err)
      }

      c.Log.Infof("Container IP: %s", ip)

      return nil
    }),
  }
}

func (c *Container) CreateDefaultBridge() {
  bridgeStateDir := path.Join(c.p.Cfg.BridgeStateDir, c.spec.Id)
  c.CreateBridge(
    c.p.Cfg.Bridge,
    c.p.Cfg.HostIface,
    c.p.Cfg.Subnet,
    bridgeStateDir,
  )
}

func (c *Container) AddEnvVars(envvars []string) {
  for _, envvar := range envvars {
    c.env = append(c.env, envvar)
  }
}

func (c *Container) SetDevices(devices []string) {
  allDevices := specconv.AllowedDevices
  for _, device := range devices {
    if val, ok := availableDevices[device]; ok {
      allDevices = append(allDevices, val)
    }
  }
  
  var allowedDeviceRules []*configs.DeviceRule
  for _, device := range allDevices {
    allowedDeviceRules = append(allowedDeviceRules, &device.DeviceRule)
  }

  c.config.Devices = allDevices
  c.config.Cgroups.Resources.Devices = allowedDeviceRules
}

func (c *Container) SetCapabilities(capabilities []string) {
  allCapabilities := defaultCapabilities
  for _, capability := range capabilities {
    allCapabilities = append(allCapabilities, capability)
  }

  c.config.Capabilities = &configs.Capabilities{
    Bounding: allCapabilities,
    Effective: allCapabilities,
    Inheritable: allCapabilities,
    Permitted: allCapabilities,
    Ambient: allCapabilities,
  }
}

func (c *Container) SetImage(image string) {
  c.spec.Image = image
}

func (c *Container) PullImage() error {
  if c.image != nil {
    return nil
  }

  var err error = nil
  c.image, c.imageManifestHex, err = PullImage(c.spec.Image, c.p.Cfg.CacheDir, c.p.Cfg.SavedDir)
  if err != nil {
    return fmt.Errorf("could not pull image: %s", err)
  }

  // If the image is nil it means that the image was not found locally and copied
  if c.image == nil {
    return nil
  }

  _, err = c.image.Digest()
  if err != nil {
    return fmt.Errorf("could not process image digest: %s", err)
  }

  return nil
}

func (c *Container) AttachImage() error {
  err := UnpackImage(c.image, c.imageManifestHex, c.p.Cfg.CacheDir, c.p.Cfg.SavedDir, c.config.Rootfs, true) // r.Config.AllowOverride
  if err != nil {
    return fmt.Errorf("could not attach image: %s", err)
  }

  return nil
}

func (c *Container) PullAndAttachImage(image string) error {
  c.SetImage(image)

  c.Log.Infof("Pulling image %s...", image)
  if err := c.PullImage(); err != nil {
    return err
  }

  return c.AttachImage()
}

func (c *Container) SetCores(cores []uint64) {
  // Set the share to 100 so that the container has the whole CPU share
  c.config.Cgroups.Resources.CpuShares = 100
  c.config.Cgroups.Resources.CpusetCpus = strings.Trim(
    strings.Join(strings.Fields(fmt.Sprint(cores)), ","), "[]",
  )
}

func (c *Container) Init() error {
  var err error

  c.container, err = c.p.factory.Create(c.spec.Id, c.config)
  if err != nil {
    return fmt.Errorf("could not create container: %s", err)
  }

  return nil
}

func (c *Container) Start() error {
  if c.bridge == nil {
    c.CreateDefaultBridge()
  }

  // TODO: We need to create a logger so we can save the output
  p := printer.New(c.Log, path.Join(c.p.Cfg.LogDir, c.spec.Id))

  c.process = &libcontainer.Process{
    Cwd:      "/",
    Env:      c.env,
    User:     "root",
    Stdout:   p,
    Stderr:   p,
    Init:     true,
    Args:   []string{entrypointFile},
  }

  if err := c.container.Run(c.process); err != nil {
    return fmt.Errorf("could not run task process: %s", err)
  }

  return nil
}

func (c *Container) Wait() (time.Duration, error) {
  state, err := c.process.Wait()
  if err != nil {
    return 0, fmt.Errorf("could not wait for container to finish: %s", err)
  }

  c.exitCode = state.ExitCode()

  return time.Since(c.timer), nil
}

func (c *Container) StartAndWait() (time.Duration, error) {
  if err := c.Start(); err != nil {
    return time.Duration(0), err
  }

  return c.Wait()
}

func (c *Container) SaveOutput(dest, output string) (string, error) {
  newPath := path.Join(dest, c.spec.Id)
  if _, err := os.Stat(newPath); os.IsNotExist(err) {
    os.MkdirAll(newPath, os.ModePerm)
  }

  newFile := path.Join(newPath, filepath.Base(output))
  err := copy.Copy(path.Join(c.config.Rootfs, output), newFile)
  if err != nil {
    return "", fmt.Errorf("could not copy output: %s: %s", output, err)
  }

  return newFile, nil
}

func (c *Container) GetResultInt(loc string) (int64, error) {
  var result int64

  intPath := path.Join(c.config.Rootfs, loc)
  if _, err := os.Stat(intPath); os.IsNotExist(err) {
    return 0, fmt.Errorf("result path does not exist: %s", intPath)
  }

  file, err := os.Open(intPath)
  if err != nil {
    return 0, fmt.Errorf("could not open file: %s: %s", intPath, err)
  }

  _, err = fmt.Fscanf(file, "%d", &result)
  if err != nil {
    return 0, fmt.Errorf("could not read integer from path: %s: %s", intPath, err)
  }

  return result, nil
}

func (c *Container) GetResultFloat(loc string) (float64, error) {
  var result float64

  floatPath := path.Join(c.config.Rootfs, loc)
  if _, err := os.Stat(floatPath); os.IsNotExist(err) {
    return 0, fmt.Errorf("result path does not exist: %s", floatPath)
  }

  file, err := os.Open(floatPath)
  if err != nil {
    return 0, fmt.Errorf("could not open file: %s: %s", floatPath, err)
  }

  _, err = fmt.Fscanf(file, "%f", &result)
  if err != nil {
    return 0, fmt.Errorf("could not read integer from path: %s: %s", floatPath, err)
  }

  return result, nil
}

func (c *Container) GetResultStr(loc string) (string, error) {
  var result string

  strPath := path.Join(c.config.Rootfs, loc)
  if _, err := os.Stat(strPath); os.IsNotExist(err) {
    return "", fmt.Errorf("result path does not exist: %s", strPath)
  }

  file, err := os.Open(strPath)
  if err != nil {
    return "", fmt.Errorf("could not open file: %s: %s", strPath, err)
  }

  _, err = fmt.Fscanf(file, "%s", &result)
  if err != nil {
    return "", fmt.Errorf("could not read integer from path: %s: %s", strPath, err)
  }

  return result, nil
}

func (c *Container) GetResultBool(loc string) (bool, error) {
  var result string

  boolPath := path.Join(c.config.Rootfs, loc)
  if _, err := os.Stat(boolPath); os.IsNotExist(err) {
    return false, fmt.Errorf("result path does not exist: %s", boolPath)
  }

  file, err := os.Open(boolPath)
  if err != nil {
    return false, fmt.Errorf("could not open file: %s: %s", boolPath, err)
  }

  _, err = fmt.Fscanf(file, "%s", &result)
  if err != nil {
    return false, fmt.Errorf("could not read integer from path: %s: %s", boolPath, err)
  }

  result = strings.ToLower(result)

  switch result {
    case "1",
         "y",
         "yes",
         "true":
      return true, nil
  }

  return false, nil
}

func (c *Container) Destroy() error {
  cacheDir := path.Join(c.p.Cfg.CacheDir, c.spec.Id)
  rootFs := c.config.Rootfs

  c.container.Destroy()
  c.container = nil
  
  if err := os.RemoveAll(cacheDir); err != nil {
    return fmt.Errorf("could not destroy container: %s", err)
  }

  if err := os.RemoveAll(rootFs); err != nil {
    return fmt.Errorf("could not destroy container: %s", err)
  }

  return nil
}
