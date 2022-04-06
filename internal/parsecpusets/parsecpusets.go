package parsecpusets
// SPDX-License-Identifier: BSD-3-Clause
//
// Authors: Alexander Jung <a.jung@lancs.ac.uk>
//
// Copyright (c) 2020, Lancaster University.  All rights reserved.
//               2021, Unikraft UG.  All rights reserved.
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
  "strconv"
  "strings"
)

func ParseCpuSets(cpuSets string) ([]uint64, error) {
  var cpus []uint64

  if strings.Contains(cpuSets, ",") {
    c := strings.Split(cpuSets, ",")

    for i := range c {
      moreCpus, err := parseRange(c[i])
      if err != nil {
        return nil, err
      }

      // If we have found more cpus, the format is A-C,E-G
      if len(moreCpus) > 0 {
        cpus = append(cpus, moreCpus...)
      
      // If not, then the format is A,B,C,E,F,G
      } else {
        j, err := strconv.ParseUint(c[i], 10, 64)
        if err != nil {
          return nil, fmt.Errorf("invalid syntax for CPU set: %s", c[i])
        }

        cpus = append(cpus, j)
      }
    }

  // Maybe the range is simply A-C
  } else {
    moreCpus, err := parseRange(cpuSets)
    if err != nil {
      return nil, err
    }

    cpus = append(cpus, moreCpus...)
  }

  return cpus, nil
}

// func parseDelim(cpuSets string) ([]uint64, error) {
//   var cpus []uint64
//   if strings.Contains(cpuSets, ",") {
//     c := strings.Split(cpuSets, ",")

//     for i := range c {
//       j, err := strconv.ParseUint(c[i], 10, 64)
//       if err != nil {
//         return nil, fmt.Errorf("invalid syntax for CPU set: %s", c[i])
//       }

//       cpus = append(cpus, j)
//     }
//   }

//   return cpus, nil
// }

func parseRange(cpuSets string) ([]uint64, error) {
  var cpus []uint64

  if res := strings.Contains(cpuSets, "-"); res {
    c := strings.Split(cpuSets, "-")
    if len(c) > 2 {
      return nil, fmt.Errorf("invalid syntax for CPU sets: %s", cpuSets)
    }

    start, err := strconv.ParseUint(c[0], 10, 64)
    if err != nil {
      return nil, fmt.Errorf("invalid syntax for CPU sets: %s", cpuSets)
    }

    end, err := strconv.ParseUint(c[1], 10, 64)
    if err != nil {
      return nil, fmt.Errorf("invalid syntax for CPU sets: %s", cpuSets)
    }
    
    for i := start; i <= end; i++ {
      cpus = append(cpus, i)
    }
  }

  return cpus, nil
}
