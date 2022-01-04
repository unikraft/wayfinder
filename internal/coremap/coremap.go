package coremap
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

  "github.com/unikraft/wayfinder/internal/parsecpusets"
  "github.com/unikraft/wayfinder/pkg/sys"
)

type Core struct {
  id           uint64
  numaNodeId   uint64
  socketId     uint64
  cacheGroupId uint64
  busy         bool
  activity     interface{}
}

type CPUCacheGroup struct {
  id         uint64
  cores      map[uint64]*Core
  totalCores int
}

type NumaNode struct {
  id         uint64
  cores      map[uint64]*Core
  totalCores int
}

type CPUSocket struct {
  id          uint64
  numaNodes   map[uint64]*NumaNode
  cacheGroups map[uint64]*CPUCacheGroup
}

type CoreMap struct {
  sync.RWMutex
  sockets map[uint64]*CPUSocket
}

// Restriction level:
const (
  // Any core will do
  CoreOptionNoRestriction = iota
  // Only cores from the same socket will do
  CoreOptionSameSocket
  // Only cores from the same NUMA node and the same socket will do
  CoreOptionSameNUMA
  // Only cores sharing the same cache/socket will do
  CoreOptionSameCache
  // Only cores sharing the same cache/NUMA node/socket will do
  CoreOptionSameCacheAndNUMA
)

type CoreRestriction int

// Helper function
func contains(array []uint64, element uint64) bool {
  for _, a := range array {
      if a == element {
          return true
      }
  }
  return false
}

// Create a new CoreMap with pre-populated socket/numa/cache groups
func New(availableCpuSet string, cpuLayoutInfo *[]sys.CPULayoutInfo) (*CoreMap, error) {
  if cpuLayoutInfo == nil {
    return nil, fmt.Errorf("no cpu layout available")
  }

  availableCores, _ := parsecpusets.ParseCpuSets(availableCpuSet)

  // Create and populate the sockets with their numa nodes and cache groups
  sockets := make(map[uint64]*CPUSocket)
  for _, cpu := range *cpuLayoutInfo {
    if _, ok := sockets[cpu.Socket]; !ok {
      sockets[cpu.Socket] = &CPUSocket{
        id:          cpu.Socket,
        numaNodes:   make(map[uint64]*NumaNode),
        cacheGroups: make(map[uint64]*CPUCacheGroup),
      }
    }
  }

  // Create and populate the numa nodes with their cores
  for _, cpu := range *cpuLayoutInfo {
    if _, ok := sockets[cpu.Socket].numaNodes[cpu.Node]; !ok {
      sockets[cpu.Socket].numaNodes[cpu.Node] = &NumaNode{
        id:         cpu.Node,
        cores:      make(map[uint64]*Core),
        totalCores: 0,
      }
    } else {
      if contains(availableCores, cpu.CPU) {
        sockets[cpu.Socket].numaNodes[cpu.Node].totalCores++
        sockets[cpu.Socket].numaNodes[cpu.Node].cores[cpu.CPU] = &Core{
          id:           cpu.CPU,
          numaNodeId:   cpu.Node,
          socketId:     cpu.Socket,
          cacheGroupId: cpu.L3Cache,
          busy:         false,
          activity:     nil,
        }
      }
    }
  }

  // Create and populate the cache groups with their cores (reuse cores from numa nodes)
  for _, cpu := range *cpuLayoutInfo {
    if _, ok := sockets[cpu.Socket].cacheGroups[cpu.L3Cache]; !ok {
      sockets[cpu.Socket].cacheGroups[cpu.L3Cache] = &CPUCacheGroup{
        id:         cpu.L3Cache,
        cores:      make(map[uint64]*Core),
        totalCores: 0,
      }
    } else {
      if contains(availableCores, cpu.CPU) {
        sockets[cpu.Socket].cacheGroups[cpu.L3Cache].totalCores++
        sockets[cpu.Socket].cacheGroups[cpu.L3Cache].cores[cpu.CPU] =
          sockets[cpu.Socket].numaNodes[cpu.Node].cores[cpu.CPU]
      }
    }
  }

  return &CoreMap{
    sockets: sockets,
  }, nil
}

// Returns a list of cores based on the given level of restriction
// This should be used when requesting free cores from the coremap
func (cm *CoreMap) FindFreeCores(level CoreRestriction) ([]*Core) {
  switch(level) {
    case CoreOptionNoRestriction:
      return cm.findAllFreeCoresAcrossAllNumaNodes()
    case CoreOptionSameSocket:
      return cm.findFreeCoresOnSameSocket()
    case CoreOptionSameNUMA:
      return cm.findFreeCoresOnSameNumaNode()
    case CoreOptionSameCache:
      return cm.findFreeCoresOnSameCache()
    case CoreOptionSameCacheAndNUMA:
      return nil // TODO
    default:
      return nil
  }
}

// Returns a list of cores sharing the same socket
func (cm *CoreMap) findFreeCoresOnSameSocket() ([]*Core) {
  var freeCores [][]*Core

  cm.RLock()
  defer cm.RUnlock()

  // Create list of all free cores for each socket
  for socketCrt, socket := range cm.sockets {
    for _, cacheGroup := range socket.cacheGroups {
      for _, core := range cacheGroup.cores {
        if !core.busy {
          freeCores[socketCrt] = append(freeCores[socketCrt], core)
        }
      }
    }
  }

  // Find the socket with the most free cores and return it
  maxCores := 0
  var maxCoresList []*Core
  for _, cores := range freeCores {
    if len(cores) > maxCores {
      maxCores = len(cores)
      maxCoresList = cores
    }
  }

  return maxCoresList
}

// Returns a list of cores sharing the same cache
func (cm *CoreMap) findFreeCoresOnSameCache() ([]*Core) {
  var freeCores [][]*Core

  cm.RLock()
  defer cm.RUnlock()

  coreCrt := 0
  // Create list of all free cores for each cache group
  for _, socket := range cm.sockets {
    for _, cacheGroup := range socket.cacheGroups {
      for _, core := range cacheGroup.cores {
        if !core.busy {
          freeCores[coreCrt] = append(freeCores[coreCrt], core)
        }
      }
      coreCrt++
    }
  }

  // Find the cache group with the most free cores and return it
  maxCores := 0
  var maxCoresList []*Core
  for _, cores := range freeCores {
    if len(cores) > maxCores {
      maxCores = len(cores)
      maxCoresList = cores
    }
  }

  return maxCoresList
}

// Returns a list of cores on the same NUMA node
func (cm *CoreMap) findFreeCoresOnSameNumaNode() ([]*Core) {
  var freeCores [][]*Core

  cm.RLock()
  defer cm.RUnlock()

  coreCrt := 0
  // Create list of all free cores for each NUMA mode
  for _, socket := range cm.sockets {
    for _, numaNode := range socket.numaNodes {
      for _, core := range numaNode.cores {
        if !core.busy {
          freeCores[coreCrt] = append(freeCores[coreCrt], core)
        }
      }
      coreCrt++
    }
  }

  // Find the NUMA node with the most free cores and return it
  maxCores := 0
  var maxCoresList []*Core
  for _, cores := range freeCores {
    if len(cores) > maxCores {
      maxCores = len(cores)
      maxCoresList = cores
    }
  }

  return maxCoresList
}

// Retrieve a list of cores which are free
func (cm *CoreMap) FindAllFreeCoresOnNumaNode(numaNodeId uint64, socketId uint64) ([]*Core, error) {
  if _, ok := cm.sockets[socketId].numaNodes[numaNodeId]; !ok {
    return nil, fmt.Errorf("could not find NUMA node with id=%d", numaNodeId)
  }

  cm.RLock()
  defer cm.RUnlock()

  var freeCores []*Core
  for _, core := range cm.sockets[socketId].numaNodes[numaNodeId].cores {
    if !core.busy {
      freeCores = append(freeCores, core)
    }
  }

  return freeCores, nil
}

// Retrieve a list of cores which are free
func (cm *CoreMap) findAllFreeCoresAcrossAllNumaNodes() []*Core {
  cm.RLock()
  defer cm.RUnlock()

  var freeCores []*Core
  for _, socket := range cm.sockets {
    for _, numaNode := range socket.numaNodes {
      for _, core := range numaNode.cores {
        if !core.busy {
          freeCores = append(freeCores, core)
        }
      }
    }
  }
  
  return freeCores
}

// Retrieve a list of numa nodes which are free
func (cm *CoreMap) FindAllFreeNumaNodes() ([]*NumaNode, error) {
  var freeNumaNodes []*NumaNode
  
  for _, socket := range cm.sockets {
    for _, numaNode := range socket.numaNodes {
      freeCores := []*Core{}
  
      for _, core := range numaNode.cores {
        if !core.busy {
          freeCores = append(freeCores, core)
        }
      }
  
      if len(freeCores) == numaNode.totalCores {
        freeNumaNodes = append(freeNumaNodes, numaNode)
      }
    }
  }

  return freeNumaNodes, nil
}

// Find a core on a NUMA node
func (cm *CoreMap) FindCoreOnNumaNode(coreId, numaNodeId uint64, socketId uint64) (*Core, error) {
  if _, ok := cm.sockets[socketId].numaNodes[numaNodeId]; !ok {
    return nil, fmt.Errorf("could not find NUMA node with id=%d", numaNodeId)
  }

  for _, core := range cm.sockets[socketId].numaNodes[numaNodeId].cores {
    if core.id == coreId {
      return core, nil
    }
  }

  return nil, fmt.Errorf("could not find core with id=%d on NUMA node with id=%d", coreId, numaNodeId)
}

// Find a core based on its ID
func (cm *CoreMap) FindCore(coreId uint64) (*Core, error) {
  for _, socket := range cm.sockets {
    for _, numaNode := range socket.numaNodes {
      for _, core := range numaNode.cores {
        if core.id == coreId {
          return core, nil
        }
      }
    }
  }

  return nil, fmt.Errorf("could not find core with id=%d", coreId)
}

// Set the activity of a core on a NUMA node
func (cm *CoreMap) SetCoreActivityOnNumaNode(coreId, numaNodeId uint64, socketId uint64, activity interface{}) error {
  core, err := cm.FindCoreOnNumaNode(coreId, numaNodeId, socketId)
  if err != nil {
    return fmt.Errorf("could not set activity: %s", err)
  }

  if core.busy {
    return fmt.Errorf("could not set activity on busy core id=%d", core.id)
  }

  cm.RLock()

  cm.sockets[socketId].numaNodes[numaNodeId].cores[coreId].activity = activity
  cm.sockets[socketId].numaNodes[numaNodeId].cores[coreId].busy = true

  cm.RUnlock()

  return nil
}

// Set the activity of a core
func (cm *CoreMap) SetCoreActivity(coreId uint64, activity interface{}) error {  
  core, err := cm.FindCore(coreId)
  if err != nil {
    return fmt.Errorf("could not set core activity for id=%d: %s", coreId, err)
  }

  return cm.SetCoreActivityOnNumaNode(coreId, core.numaNodeId, core.socketId, activity)
}

// Release a core from its activity
func (cm *CoreMap) ReleaseCoreOnNumaNode(coreId, numaNodeId uint64, socketId uint64) error {
  if _, ok := cm.sockets[socketId].numaNodes[numaNodeId]; !ok {
    return fmt.Errorf("invalid NUMA node id=%d", numaNodeId)
  }

  if _, ok := cm.sockets[socketId].numaNodes[numaNodeId].cores[coreId]; !ok {
    return fmt.Errorf("invalid core id=%d for NUMA node with id=%d", coreId, numaNodeId)
  }

  cm.RLock()
  
  cm.sockets[socketId].numaNodes[numaNodeId].cores[coreId].activity = nil
  cm.sockets[socketId].numaNodes[numaNodeId].cores[coreId].busy = false

  cm.RUnlock()

  return nil
}

// Release a core from its activity
func (cm *CoreMap) ReleaseCore(coreId uint64) error {
  core, err := cm.FindCore(coreId)
  if err != nil {
    return fmt.Errorf("could not find core with id=%d: %s", coreId, err)
  }

  return cm.ReleaseCoreOnNumaNode(coreId, core.numaNodeId, core.socketId)
}

// Release an entire NUMA node of all core activities
func (cm *CoreMap) ReleaseNumaNode(numaNodeId uint64, socketId uint64) error {
  if _, ok := cm.sockets[socketId].numaNodes[numaNodeId]; !ok {
    return fmt.Errorf("could not find NUMA node with id=%d", numaNodeId)
  }

  for coreId, _ := range cm.sockets[socketId].numaNodes[numaNodeId].cores {
    err := cm.ReleaseCoreOnNumaNode(coreId, numaNodeId, socketId)
    if err != nil {
      return fmt.Errorf("could not release numa node with id=%d: %s", numaNodeId, err)
    }
  }
  
  return nil
}

// Release an entire socket of all core activities
func (cm *CoreMap) ReleaseSocket(socketId uint64) error {
  if _, ok := cm.sockets[socketId]; !ok {
    return fmt.Errorf("could not find socket with id=%d", socketId)
  }

  for numaNodeId, _ := range cm.sockets[socketId].numaNodes {
    err := cm.ReleaseNumaNode(numaNodeId, socketId)
    if err != nil {
      return fmt.Errorf("could not release socket with id=%d: %s", socketId, err)
    }
  }

  return nil
}

// Set a core by overwritting it
func (cm *CoreMap) SetCore(newCore *Core) error {
  for socketId, socket := range cm.sockets {
    for numaNodeId, numaNode := range socket.numaNodes {
      for coreId, oldCore := range numaNode.cores {
        if coreId == newCore.id {
          if oldCore.busy {
            return fmt.Errorf("cannot overwrite busy core with id=%d", coreId)
          }
  
          if newCore.numaNodeId != oldCore.numaNodeId {
            return fmt.Errorf("cannot overwrite core, mismatch NUMA node ids %d and %d", newCore.numaNodeId, oldCore.numaNodeId)
          }
          
          cm.RLock()
  
          cm.sockets[socketId].numaNodes[numaNodeId].cores[coreId] = newCore
  
          cm.RUnlock()
  
          return nil
        }
      }
    }
  }


  return fmt.Errorf("could not set core")
}

// Get a NUMA node
func (cm *CoreMap) SetCoreOnNumaNode(newCore *Core, numaNodeId uint64, socketId uint64) error {
  if _, ok := cm.sockets[socketId].numaNodes[numaNodeId]; !ok {
    return fmt.Errorf("could not find NUMA node with id=%d", numaNodeId)
  }

  for coreId, oldCore := range cm.sockets[socketId].numaNodes[numaNodeId].cores {
    if coreId == newCore.id {
      if oldCore.busy {
        return fmt.Errorf("cannot overwrite busy core with id=%d", coreId)
      }

      if newCore.numaNodeId != oldCore.numaNodeId {
        return fmt.Errorf("cannot overwrite core, mismatch NUMA node ids %d and %d", newCore.numaNodeId, oldCore.numaNodeId)
      }
      
      cm.RLock()

      cm.sockets[socketId].numaNodes[numaNodeId].cores[coreId] = newCore

      cm.RUnlock()

      return nil
    }
  }

  return fmt.Errorf("could not set core")
}

// Get a core
func (cm *CoreMap) GetCoreActivity(coreId uint64) (interface{}, error) {
  core, err := cm.FindCore(coreId)
  if err != nil {
    return nil, fmt.Errorf("could not get core with id=%d activity: %s", coreId, err)
  }
  
  return core.activity, nil
}

func (cm *CoreMap) Print() {
  i := 1
  cols := 3
  fmt.Printf("Core map:\n")
  for _, socket := range cm.sockets {
    for _, numaNode := range socket.numaNodes {
      for _, core := range numaNode.cores {
        busy := ""
        if core.busy {
          busy = fmt.Sprintf("%p", &core.activity)
        }
        fmt.Printf(" %2d|%2d|%2d: [%12s] ", socket.id, numaNode.id, core.Id(), busy)
        if int(i) % cols == 0 {
          fmt.Printf("\n")
        }
        i++
      }
    }
  }
  fmt.Printf("\n\n")
}

func (c *Core) Id() uint64 {
  return c.id
}

func (c *Core) NumaNodeId() uint64 {
  return c.numaNodeId
}

func (c *Core) IsBusy() bool {
  return c.busy
}

func (c *Core) Activity() interface{} {
  return c.activity
}
