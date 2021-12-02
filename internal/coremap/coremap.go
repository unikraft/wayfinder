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
)

type Core struct {
  id          uint64
  numaNodeId  uint64
  busy        bool
  activity    interface{}
}

type NumaNode struct {
  id         uint64
  cores      map[uint64]*Core
  totalCores int
}

type CoreMap struct {
  sync.RWMutex
  numaNodes map[uint64]*NumaNode
}

// Helper function
func contains(array []uint64, element uint64) bool {
  for _, a := range array {
      if a == element {
          return true
      }
  }
  return false
}

// NumaNode creates a fixed-length map of cores with their ID as index.
func NewNumaNode(id uint64, cores []uint64, availableCores []uint64) *NumaNode {
  numaNode := &NumaNode{
    id:         id,
    cores:      make(map[uint64]*Core, len(cores)),
    totalCores: len(cores),
  }

  // Add the core ID as index to the map
  for i := 0; i < len(cores); i++ {
    if contains(availableCores, cores[i]) {
      numaNode.cores[cores[i]] = &Core{
        id:         cores[i],
        numaNodeId: id,
        busy:       false,
        activity:   nil,
      }
    }
  }

  return numaNode
}

// Create a new CoreMap with pre-populated NumaNode array
func New(numaNodes map[uint64]*NumaNode) *CoreMap {
  return &CoreMap{
    numaNodes: numaNodes,
  }
}

// Create a new CoreMap based on string array of CPU set notation
func NewFromStr(numaNodesStr []string, availableCpuSet string) (*CoreMap, error) {
  if len(numaNodesStr) == 0 {
    return nil, fmt.Errorf("no NUMA node strings provided")
  }

  availableCores, _ := parsecpusets.ParseCpuSets(availableCpuSet)
  numaNodes := make(map[uint64]*NumaNode, len(numaNodesStr))

  for i, numaNodeStr := range numaNodesStr {
    cores, err := parsecpusets.ParseCpuSets(numaNodeStr)
    if err != nil {
      return nil, fmt.Errorf("could not parse NUMA node string: %s", err)
    }

    numaNodes[uint64(i)] = NewNumaNode(uint64(i), cores, availableCores)
  }

  return New(numaNodes), nil
}

// Get the core's numa node map
func (cm *CoreMap) NumaNodes() map[uint64]*NumaNode {
  return cm.numaNodes
}

// Retrieve a list of cores which are free
func (cm *CoreMap) FindAllFreeCoresOnNumaNode(numaNodeId uint64) ([]*Core, error) {
  if _, ok := cm.numaNodes[numaNodeId]; !ok {
    return nil, fmt.Errorf("could not find NUMA node with id=%d", numaNodeId)
  }

  cm.RLock()
  defer cm.RUnlock()

  var freeCores []*Core
  for _, core := range cm.numaNodes[numaNodeId].cores {
    if !core.busy {
      freeCores = append(freeCores, core)
    }
  }

  return freeCores, nil
}

// Retrieve a list of cores which are free
func (cm *CoreMap) FindAllFreeCoresAcrossAllNumaNodes() []*Core {
  cm.RLock()
  defer cm.RUnlock()

  var freeCores []*Core
  for _, numaNode := range cm.numaNodes {
    for _, core := range numaNode.cores {
      if !core.busy {
        freeCores = append(freeCores, core)
      }
    }
  }
  
  return freeCores
}

// Retrieve a list of numa nodes which are free
func (cm *CoreMap) FindAllFreeNumaNodes() ([]*NumaNode, error) {
  var freeNumaNodes []*NumaNode
  
  for _, numaNode := range cm.numaNodes {
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

  return freeNumaNodes, nil
}

// Find a core on a NUMA node
func (cm *CoreMap) FindCoreOnNumaNode(coreId, numaNodeId uint64) (*Core, error) {
  if _, ok := cm.numaNodes[numaNodeId]; !ok {
    return nil, fmt.Errorf("could not find NUMA node with id=%d", numaNodeId)
  }

  for i, core := range cm.numaNodes[numaNodeId].cores {
    if i == coreId {
      return core, nil
    }
  }

  return nil, fmt.Errorf("could not find core with id=%d", coreId)
}

// Find a core based on its ID
func (cm *CoreMap) FindCore(coreId uint64) (*Core, error) {
  for numaId, _ := range cm.numaNodes {
    return cm.FindCoreOnNumaNode(coreId, numaId)
  }

  return nil, fmt.Errorf("could not find core with id=%d", coreId)
}

// Set the activity of a core on a NUMA node
func (cm *CoreMap) SetCoreActivityOnNumaNode(coreId, numaNodeId uint64, activity interface{}) error {
  core, err := cm.FindCoreOnNumaNode(coreId, numaNodeId)
  if err != nil {
    return fmt.Errorf("could not set activity: %s", err)
  }

  if core.busy {
    return fmt.Errorf("could not set activity on busy core id=%d", core.id)
  }

  cm.RLock()

  cm.numaNodes[numaNodeId].cores[coreId].activity = activity
  cm.numaNodes[numaNodeId].cores[coreId].busy = true

  cm.RUnlock()

  return nil
}

// Set the activity of a core
func (cm *CoreMap) SetCoreActivity(coreId uint64, activity interface{}) error {  
  core, err := cm.FindCore(coreId)
  if err != nil {
    return fmt.Errorf("could not set core activity for id=%d: %s", coreId, err)
  }

  return cm.SetCoreActivityOnNumaNode(coreId, core.numaNodeId, activity)
}

// Release a core from its activity
func (cm *CoreMap) ReleaseCoreOnNumaNode(coreId, numaNodeId uint64) error {
  if _, err := cm.FindCore(coreId); err != nil {
    return fmt.Errorf("could not set core with id=%d on NUMA node with id: %d, %s", coreId, numaNodeId, err)
  }

  cm.RLock()
  
  cm.numaNodes[numaNodeId].cores[coreId].activity = nil
  cm.numaNodes[numaNodeId].cores[coreId].busy = false

  cm.RUnlock()

  return nil
}

// Release a core from its activity
func (cm *CoreMap) ReleaseCore(coreId uint64) error {
  core, err := cm.FindCore(coreId)
  if err != nil {
    return fmt.Errorf("could not release core with id=%d: %s", coreId, err)
  }

  return cm.ReleaseCoreOnNumaNode(coreId, core.numaNodeId)
}

// Release an entire NUMA node of all core activities
func (cm *CoreMap) ReleaseNumaNode(numaNodeId uint64) error {
  if _, ok := cm.numaNodes[numaNodeId]; !ok {
    return fmt.Errorf("could not find NUMA node with id=%d", numaNodeId)
  }

  for coreId, _ := range cm.numaNodes[numaNodeId].cores {
    err := cm.ReleaseCoreOnNumaNode(coreId, numaNodeId)
    if err != nil {
      return fmt.Errorf("could not release numa node with id=%d: %s", numaNodeId, err)
    }
  }
  
  return nil
}

// Set a core by overwritting it
func (cm *CoreMap) SetCore(newCore *Core) error {
  for numaNodeId, numaNode := range cm.numaNodes {
    for coreId, oldCore := range numaNode.cores {
      if coreId == newCore.id {
        if oldCore.busy {
          return fmt.Errorf("cannot overwrite busy core with id=%d", coreId)
        }

        if newCore.numaNodeId != oldCore.numaNodeId {
          return fmt.Errorf("cannot overwrite core, mismatch NUMA node ids %d and %d", newCore.numaNodeId, oldCore.numaNodeId)
        }
        
        cm.RLock()

        cm.numaNodes[numaNodeId].cores[coreId] = newCore

        cm.RUnlock()

        return nil
      }
    }
  }

  return fmt.Errorf("could not set core")
}

// Get a NUMA node
func (cm *CoreMap) SetCoreOnNumaNode(newCore *Core, numaNodeId uint64) error {
  if _, ok := cm.numaNodes[numaNodeId]; !ok {
    return fmt.Errorf("could not find NUMA node with id=%d", numaNodeId)
  }

  for coreId, oldCore := range cm.numaNodes[numaNodeId].cores {
    if coreId == newCore.id {
      if oldCore.busy {
        return fmt.Errorf("cannot overwrite busy core with id=%d", coreId)
      }

      if newCore.numaNodeId != oldCore.numaNodeId {
        return fmt.Errorf("cannot overwrite core, mismatch NUMA node ids %d and %d", newCore.numaNodeId, oldCore.numaNodeId)
      }
      
      cm.RLock()

      cm.numaNodes[numaNodeId].cores[coreId] = newCore

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
