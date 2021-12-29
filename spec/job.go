package spec
// SPDX-License-Identifier: BSD-3-Clause
//
// Authors: Alexander Jung <a.jung@lancs.ac.uk>
//                         <alex@unikraft.io>
//
// Copyright (c) 2021, Lancaster University.  All rights reserved.
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
  "io"
  "fmt"
  "crypto/md5"
  "strings"

  "github.com/Knetic/govaluate"
  "gopkg.in/yaml.v2"
  "github.com/evolbioinfo/gotree/tree"

  "github.com/unikraft/wayfinder/api/proto"
)

type JobSpec struct {
  // Parsable in the configuration file
  Params         []ParamSpec          `json:"params"  yaml:"params"`
  Build            BuildSpec          `json:"build"   yaml:"build"`
  Test             TestSpec           `json:"test"    yaml:"test"`
  tree            *tree.Tree
  
  // Desired additional information by the scheduler
  Id               uint               `json:"id"`
  Scheduler        proto.JobScheduler `json:"scheduler"`
  IsolLevel        proto.JobIsolLevel `json:"isol_level"`
  IsolSplit        proto.JobIsolSplit `json:"isol_split"`
  PermutationLimit int64              `json:"permutation_limit"`
  CurrentPerm      JobPermutation     `json:"current"`
}

type JobPermutation struct {
  Id         uint             `json:"id"`
  JobId      uint             `json:"job_id"`
  Params   []ParamPermutation `json:"params"`
  Checksum   string           `json:"checksum"`
}

// ParseJobSpec accepts a YAML string input and returns a parsed object
func ParseJobSpec(spec string) (*JobSpec, error) {  
  if len(spec) == 0 {
    return nil, fmt.Errorf("spec is empty")
  }

  job := JobSpec{}

  err := yaml.Unmarshal([]byte(spec), &job)
  if err != nil {
    return nil, fmt.Errorf("could not parse YAML: %s", err)
  }

  if len(job.Params) == 0 {
    return nil, fmt.Errorf("job spec has no parameters")
  }

  job.tree = tree.NewTree()

  return &job, nil
}

// next recursively iterates across paramters to generate a set of tasks
// func (j *JobSpec) next(i int, limit int64, perms []*JobPermutation, curr []ParamPermutation) ([]*JobPermutation, error) {
//   // List all permutations for this parameter
//   params, err := paramPermutations(&j.Params[i])
//   if err != nil {
//     return nil, err
//   }
  
//   // func (t *Tree) tree.NewNode() *Node {
//   // func (t *Tree) tree.ConnectNodes(parent *Node, child *Node) *Edge {
  
//   for _, param := range params {
//     if len(curr) > 0 {
//       last := curr[len(curr)-1]
//       if last.Name == param.Name {
//         curr = curr[:len(curr)-1]
//       }
//     }

//     curr = append(curr, param)

//     // Break when there are no more parameters to iterate over, thus creating
//     // the task.
//     if i + 1 == len(j.Params) {
//       var p = make([]ParamPermutation, len(j.Params))
//       copy(p, curr)
//       perm := &JobPermutation{
//         Params:   p,
//         Outputs: &j.Outputs,
//       }
//       perms = append(perms, perm)

//     // Otherwise, recursively parse parameters in-order    
//     } else {
//       nextPerms, err := j.next(i + 1, limit, nil, curr)
//       if err != nil {
//         return nil, err
//       }

//       perms = append(perms, nextPerms...)
//     }
//   }

//   return perms, nil
// }

// Extracts all params and formats them in an evaluable way
func createParamMapForEval(params []ParamPermutation) map[string]string {
  paramsForCond := make(map[string]string)
  for _, param := range params {
    var value string
    if param.Value == "y" {
      value = "true"
    } else if param.Value == "n" {
      value = "false"
    } else {
      if (param.Type == "str") {
        value = "'" + param.Value + "'"
      } else {
        value = param.Value
      }
    }
    paramsForCond[param.Name] = value
  }
  return paramsForCond
}

// Replace all the params in the condition with their values
func replaceSymbols(unformattedCond string, permMap map[string]string) string {
  cond := unformattedCond

  for k, v := range permMap {
    cond = strings.Replace(cond, k, v, -1)
  }

  return cond
}

// Create an expression from the condition and evaluate it to true/false
// If the expression evaluation fails the param is ignored
func evalExpression(cond string) (bool, error) {
  expression, err := govaluate.NewEvaluableExpression(cond)
  if err != nil {
      return false, fmt.Errorf("could not parse condition: %s", err)
  }

  expResult, err := expression.Evaluate(nil);
  if err != nil {
      return false, fmt.Errorf("could not evaluate condition: %s", err)
  }

  return expResult.(bool), nil
}

// next recursively iterates across paramters to generate a set of tasks
func (j *JobSpec) next(i int, permutations chan *JobPermutation,
                       errors chan error, done chan bool, all []*JobPermutation,
                       limit int64, curr []ParamPermutation) ([]*JobPermutation, error) {
  // List all permutations for this parameter
  params, err := paramPermutations(&j.Params[i])
  if err != nil {
    errors <- err
    return nil, err
  }

  // TODO: Use gotree
  // func (t *Tree) tree.NewNode() *Node
  // func (t *Tree) tree.ConnectNodes(parent *Node, child *Node) *Edge
  
  for _, param := range params {
    if len(curr) > 0 {
      last := curr[len(curr)-1]
      if last.Name == param.Name {
        curr = curr[:len(curr)-1]
      }
    }

    curr = append(curr, param)

    // Break when there are no more parameters to iterate over, thus creating
    // the task.
    if i + 1 == len(j.Params) {
      var p = make([]ParamPermutation, len(j.Params))

      copy(p, curr)

      perm := &JobPermutation{
        JobId:  j.Id,
        Params: p,
      }
      
      // Remove value from false conditions. This will generate tasks with duplicate checksums.
      // They will need to be uniqued after the whole generation is done.
      var paramMap map[string]string
      shouldBuildMap := true
      for i, param := range perm.Params {
        if param.Cond != "" {
          // Creates a key-value map for evaluating conditionals when the first conditional is detected
          if shouldBuildMap {
            paramMap = createParamMapForEval(perm.Params)
            shouldBuildMap = false
          }
          shouldEval, err := evalExpression(replaceSymbols(param.Cond, paramMap))
          if err != nil {
              return nil, fmt.Errorf("could not evaluate expression for param %v: %s", param, err)
          }
          if !shouldEval {
            perm.Params[i].Value = ""
          }
        }
      }

      perm.Checksum = perm.checksum()

      // Save permutation
      permutations <- perm
      all = append(all, perm)

      if len(all) >= int(limit) {
        done <- true
        return all, nil
      }

    // Otherwise, recursively parse parameters in-order    
    } else {
      nextPerms, err := j.next(i + 1, permutations, errors, done, nil, limit, curr)
      if err != nil {
        // If this has occured, we've already sent the error to the channel
        return nil, err
      }

      all = append(all, nextPerms...)
    }
  }

  if len(all) >= int(limit) {
    done <- true
  }
  return all, nil
}

// func (j *JobSpec) Permutations(ch chan []*JobPermutation, limit int64) (error) {
//   if j.permutations != nil {
//     return j.permutations, nil
//   }

//   var perm []*JobPermutation

//   perm, err := j.next(0, limit, perm, nil)
//   if err != nil {
//     return nil, err
//   }

//   j.permutations = perm

//   return perm, nil
// }

// Permutations returns a list of all possible tasks based on parameterisation
func (j *JobSpec) Permutations(limit int64) (chan *JobPermutation, chan error, chan bool) {
  done := make(chan bool)
  errors := make(chan error)
  permutations := make(chan *JobPermutation)

  var all []*JobPermutation

  go func() {
    j.next(0, permutations, errors, done, all, limit, nil)
  }()

  return permutations, errors, done
}

// TotalPermutations calculates the total number of permutations for the job
func (j *JobSpec) TotalPermutations() (uint, error) {
  
  if len(j.Params) == 0 {
    return 0, fmt.Errorf("no parameters")
  }

  var total uint = 1

  for _, param := range j.Params {
    params, err := paramPermutations(&param)
    if err != nil {
      return 0, fmt.Errorf("could not parse parameter: %s", err)
    }

    total *= uint(len(params))
  }

  return total, nil
}

// checksum ...
func (jp *JobPermutation) checksum() string {
  if len(jp.Checksum) == 0 {

    // Calculate the UUID based on a reproducible md5 seed
    md5val := md5.New()
    for _, param := range jp.Params {
      io.WriteString(md5val, fmt.Sprintf("%s=%s\n", param.Name, param.Value))
    }

    jp.Checksum = fmt.Sprintf("%x", md5val.Sum(nil))
  }

  return jp.Checksum
}
