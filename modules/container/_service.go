package container
// SPDX-License-Identifier: BSD-3-Clause
//
// Authors: Alexander Jung <a.jung@lancs.ac.uk>
//
// Copyright (c) 2021, Lancaster University.  All rights reserved.
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
  // "os"
  // "fmt"
  // // "time"
  // "path"
  // "strings"
  // "context"

  // "golang.org/x/sys/unix"
  // "google.golang.org/grpc/codes"
  // "google.golang.org/grpc/status"
  // // "github.com/otiai10/copy"
  // // "github.com/novln/docker-parser"
  // "github.com/opencontainers/runc/libcontainer"
  // "github.com/opencontainers/runtime-spec/specs-go"
  // "github.com/opencontainers/runc/libcontainer/specconv"
  // "github.com/opencontainers/runc/libcontainer/configs"

  // "github.com/unikraft/wayfinder/api/proto"
  // "github.com/unikraft/wayfinder/pkg/common/errors"
  // "github.com/unikraft/wayfinder/internal/log"
)


// func (s *Service) CreateContainer(ctx context.Context, req *proto.CreateContainerRequest) (*proto.CreateContainerResponse, error) {
//   if req.Image == "" {
//     return nil, errors.NewMissingParameterError("Image")
//   }

//   if req.Name == "" {
//     return nil, errors.NewMissingParameterError("Name")
//   }

//   if len(req.Cores) == 0 {
//     return nil, errors.NewMissingParameterError("Cores")
//   }

//   if len(req.Commands) == 0 {
//     return nil, errors.NewMissingParameterError("Commands")
//   }

//   image, err := PullImage(req.Image, s.p.Cfg.CacheDir)
//   if err != nil {
//     return nil, status.Errorf(codes.Internal, "could not pull image: %s", err)
//   }
  
//   digest, err := image.Digest()
//   if err != nil {
//     return nil, status.Errorf(codes.Internal, "could not process image digest: %s", err)
//   }

//   s.p.Log.Infof("Pulled: %s", digest)

//   allowedDevices := specconv.AllowedDevices
//   for _, device := range req.Devices {
//     if val, ok := defaultDevices[device]; ok {
//       allowedDevices = append(allowedDevices, val)
//     }
//   }
  
//   var allowedDeviceRules []*configs.DeviceRule
//   for _, device := range allowedDevices {
//     allowedDeviceRules = append(allowedDeviceRules, &device.DeviceRule)
//   }

//   capabilities := defaultCapabilities
//   for _, capability := range req.Capabilities {
//     capabilities = append(capabilities, capability)
//   }
  
//   // Create the root file system on disk
//   rootfs := path.Join(s.p.Cfg.ContainerRootDir, req.Name)
//   if _, err := os.Stat(rootfs); os.IsNotExist(err) {
//     os.MkdirAll(rootfs, os.ModePerm)
//   }

//   bridgeStateDir := path.Join(s.p.Cfg.BridgeStateDir, req.Name)

//   bridge := NewBridge(
//     s.p.Cfg.Bridge,
//     s.p.Cfg.HostIface,
//     s.p.Cfg.Subnet,
//     bridgeStateDir,
//   )

//   config := &configs.Config{
//     Rootfs: rootfs,
//     Capabilities: &configs.Capabilities{
//       Bounding: capabilities,
//       Effective: capabilities,
//       Inheritable: capabilities,
//       Permitted: capabilities,
//       Ambient: capabilities,
//     },
//     Namespaces: configs.Namespaces([]configs.Namespace{
//       {Type: configs.NEWNS},
//       {Type: configs.NEWUTS},
//       {Type: configs.NEWIPC},
//       {Type: configs.NEWPID},
//       {Type: configs.NEWNET},
//     }),
//     Cgroups: &configs.Cgroup{
//       Name:      req.Name,
//       Parent:    "",
//       Resources: &configs.Resources{
//         MemorySwappiness: nil,
//         Devices:          allowedDeviceRules,
//         // Join the core ids together in a comma separated listed
//         CpusetCpus:       strings.Trim(
//           strings.Join(strings.Fields(fmt.Sprint(req.Cores)), ","), "[]",
//         ),
//         // Set the share to 100 so that the container has the whole CPU share
//         CpuShares:        100,
//       },
//     },
//     MaskPaths: []string{
//       "/proc/acpi",
//       "/proc/asound",
//       "/proc/kcore",
//       "/proc/keys",
//       "/proc/latency_stats",
//       "/proc/timer_list",
//       "/proc/timer_stats",
//       "/proc/sched_debug",
//       "/sys/firmware",
//       "/proc/scsi",
//     },
//     ReadonlyPaths: []string{
//       "/proc/bus",
//       "/proc/fs",
//       "/proc/irq",
//       "/proc/sys",
//       "/proc/sysrq-trigger",
//     },
//     Devices:  allowedDevices,
//     Hostname: req.Name,
//     Mounts: []*configs.Mount{
//       {
//         Source:      "proc",
//         Destination: "/proc",
//         Device:      "proc",
//         Flags:       defaultMountFlags,
//       },
//       {
//         Source:      "tmpfs",
//         Destination: "/dev",
//         Device:      "tmpfs",
//         Flags:       unix.MS_NOSUID | unix.MS_STRICTATIME,
//         Data:        "mode=755",
//       },
//       {
//         Source:      "devpts",
//         Destination: "/dev/pts",
//         Device:      "devpts",
//         Flags:       unix.MS_NOSUID | unix.MS_NOEXEC,
//         Data:        "newinstance,ptmxmode=0666,mode=0620,gid=5",
//       },
//       {
//         Device:      "tmpfs",
//         Source:      "shm",
//         Destination: "/dev/shm",
//         Data:        "mode=1777,size=65536k",
//         Flags:       defaultMountFlags,
//       },
//       {
//         Source:      "mqueue",
//         Destination: "/dev/mqueue",
//         Device:      "mqueue",
//         Flags:       defaultMountFlags,
//       },
//       {
//         Source:      "sysfs",
//         Destination: "/sys",
//         Device:      "sysfs",
//         Flags:       defaultMountFlags | unix.MS_RDONLY,
//       },
//       {
//         Destination: "/sys/fs/cgroup",
//         Device:      "cgroup",
//         Source:      "cgroup",
//         Flags:       defaultMountFlags | unix.MS_RDONLY,
//       },
//     },
//     Networks: []*configs.Network{
//       {
//         Type:    "loopback",
//         Address: "127.0.0.1/0",
//         Gateway: "localhost",
//       },
//     },
//     Rlimits: []configs.Rlimit{
//       {
//         Type: unix.RLIMIT_NOFILE,
//         Hard: uint64(1025),
//         Soft: uint64(1025),
//       },
//     },
//     Hooks: configs.Hooks{
//       configs.Prestart: configs.HookList{
//         configs.NewFunctionHook(func(st *specs.State) error {
//           ip, err := bridge.Create(st)
//           if err != nil {
//             return fmt.Errorf("could not create bridge: %s", err)
//           }

//           s.p.Log.Infof("Container IP: %s", ip)

//           return nil
//         }),
//       },

//       // The `StartContainer` hook is the closest way to run code before the
//       // process is executed by libcontainer[0].  However, this configuration is
//       // passed via a JSON serialized object which uses a path path to an
//       // executable script and not Go code, like that below.  As a result, we
//       // must use the earliest placable hook, `CreateRuntime`, we can call a
//       // within libcontainer before it runs the code specified by the run.
//       // [0]: https://github.com/opencontainers/runc/blob/v1.0.0-rc92/libcontainer/standard_init_linux.go#L214-L220
//       configs.CreateRuntime: configs.HookList{
//         configs.NewFunctionHook(func(st *specs.State) error {
//           // r.timer = time.Now()
//           return nil
//         }),
//       },
//     },
//   }

//   container, err := s.p.factory.Create(req.Name, config)
//   if err != nil {
//     return nil, status.Errorf(codes.Internal, "could not create container: %s", err)
//   }

  
//   // Extract the image to the desired location
//   s.p.Log.Infof("Extracting image to: %s", rootfs)
//   // TODO: SLOW!
//   err = UnpackImage(image, s.p.Cfg.CacheDir, rootfs, true) // r.Config.AllowOverride
//   if err != nil {
//     return nil, status.Errorf(codes.Internal, "could not extract image: %s", err)
//   }

//   // Set the argument as either the path or the cmd of the run
//   entrypointFile := "/entrypoint.sh"
//   entrypointPath := path.Join(rootfs, entrypointFile)

//   f, err := os.OpenFile(
//     entrypointPath,
//     os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
//     os.ModePerm,
//   )
//   if err != nil {
//     return nil, status.Errorf(codes.Internal, "could not create temporary cmd file: %s", err)
//   }

//   // TODO: does sh support bash?
//   _, err = f.WriteString("#!/usr/bin/env sh\n")
//   if err != nil {
//     return nil, status.Errorf(codes.Internal, "could not write to temporary cmd file: %s", err)
//   }

//   _, err = f.WriteString(req.Commands)
//   if err != nil {
//     return nil, status.Errorf(codes.Internal, "could not write to temporary cmd file: %s", err)
//   }

//   f.Close()

//   taskProcess := &libcontainer.Process{
//     Cwd:      "/",
//     // Env:      append(defaultEnvironment, r.Config.Env...),
//     Env:      defaultEnvironment,
//     User:     "root",
//     Stdout:   &log.Logger{},
//     Stderr:   &log.Logger{},
//     Init:     true,
//     Args:   []string{"/entrypoint.sh"},
//   }

//   err = container.Run(taskProcess)
//   if err != nil {
//     return nil, status.Errorf(codes.Unimplemented, "could not run task process: %s", err)
//   }

//   // Wait for the process to finish
//   state, err := taskProcess.Wait()
//   if err != nil {
//     return nil, status.Errorf(codes.Unimplemented, "could not wait for container to finish: %s", err)
//   }

//   fmt.Printf("%#v\n", state.ExitCode()) // , time.Since(r.timer), nil

//   return nil, status.Errorf(codes.Unimplemented, "method CreateContainer not implemented")
// }
