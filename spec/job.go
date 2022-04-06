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
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/evolbioinfo/gotree/tree"
	"gopkg.in/yaml.v2"

	"github.com/unikraft/wayfinder/api/proto"
)

type JobSpec struct {
	// Parsable in the configuration file
	Params     []ParamSpec `json:"params"  yaml:"params"`
	Build      BuildSpec   `json:"build"   yaml:"build"`
	Test       TestSpec    `json:"test"    yaml:"test"`
	tree       *tree.Tree
	validPerms int

	// Desired additional information by the scheduler
	Id               uint               `json:"id"`
	Scheduler        proto.JobScheduler `json:"scheduler"`
	SeqScheduler     bool               `json:"seq_scheduler"`
	IsolLevel        proto.JobIsolLevel `json:"isol_level"`
	IsolSplit        proto.JobIsolSplit `json:"isol_split"`
	PermutationLimit uint64             `json:"permutation_limit"`
	Repeats          uint64             `json:"repeats"`
	CurrentPerm      JobPermutation     `json:"current"`
	DryRun           bool               `json:"dry_run"`
	Name             string             `json:"name"`
}

type JobPermutation struct {
	Id            uint               `json:"id"`
	JobId         uint               `json:"job_id"`
	Params        []ParamPermutation `json:"params"`
	BuildChecksum string             `json:"build_checksum"`
	Checksum      string             `json:"checksum"`
}

const (
	Grid = iota
	Random
	Guided
)

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
			if param.Type == "str" {
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

	expResult, err := expression.Evaluate(nil)
	if err != nil {
		return false, fmt.Errorf("could not evaluate condition: %s", err)
	}

	return expResult.(bool), nil
}

// next recursively iterates across paramters to generate a set of tasks
func (j *JobSpec) next(
	i int,
	permutationsChan chan *JobPermutation,
	errChan chan error,
	treeChan chan *tree.Tree,
	canPublish chan bool,
	allNr uint64,
	limit uint64,
	curr []ParamPermutation,
	parentNode *tree.Node,
) (
	nextPNr uint64,
	err error,
) {
	// List all permutations for this parameter
	params, err := paramPermutations(&j.Params[i-1])
	if err != nil {
		errChan <- err
		return 0, err
	}

	for _, param := range params {
		if len(curr) > 0 {
			last := curr[len(curr)-1]
			if last.Name == param.Name {
				curr = curr[:len(curr)-1]
			}
		}

		curr = append(curr, param)

		// Create a node with the parameter
		node := j.tree.NewNode()
		node.SetDepth(i)
		j.validPerms++
		node.SetId(j.validPerms)
		node.SetName(fmt.Sprintf("%s: %s", param.Name, param.Value))
		j.tree.ConnectNodes(parentNode, node)

		// Break when there are no more parameters to iterate over, thus creating
		// the task.
		if i == len(j.Params) {
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
						return 0, fmt.Errorf("could not evaluate expression for param %v: %s", param, err)
					}
					if !shouldEval {
						perm.Params[i].Value = ""
					}
				}
			}

			perm.Checksum = perm.checksum()
			perm.BuildChecksum = perm.buildChecksum()

			// Block save until task is processed
			if j.SeqScheduler {
				_ = <-canPublish
			}

			// Save permutation
			permutationsChan <- perm
			allNr++

			if allNr >= limit {
				treeChan <- j.tree
				return 0, nil
			}

			// Otherwise, recursively parse parameters in-order
		} else {
			nextPNr, err := j.next(i+1, permutationsChan, errChan, treeChan, canPublish, 0, limit, curr, node)
			if err != nil {
				// If this has occured, we've already sent the error to the channel
				return 0, err
			}

			allNr += nextPNr
		}
	}

	if allNr >= limit {
		treeChan <- j.tree
	}
	return allNr, nil
}

// Generate random permutations of the parameters until limit is reached,
// then wait to see if more are needed (in the case when there are colisions)
func (j *JobSpec) random(
	permutationsChan chan *JobPermutation,
	errChan chan error,
	treeChan chan *tree.Tree,
	canPublish chan bool,
	remainingChan chan uint64,
	limit uint64,
) (
	uint64,
	error,
) {
	var params [][]ParamPermutation = make([][]ParamPermutation, len(j.Params))
	var err error

	// Extract all parameters and prepare them for permutations
	for i, param := range j.Params {
		params[i], err = paramPermutations(&param)
		if err != nil {
			errChan <- err
			return 0, err
		}
	}

	var i uint64 = 0
	for {
		// Generate the given number of random permutations
		for ; i < limit; i++ {
			var p []ParamPermutation = make([]ParamPermutation, len(j.Params))

			// First generate a random permutation for each parameter
			for j, param := range params {
				p[j] = param[rand.Intn(len(param))]
			}

			// Fill the permutation to be returned
			permFinal := &JobPermutation{
				JobId:  j.Id,
				Params: p,
			}

			// Evaluate all conditions
			var paramMap map[string]string
			shouldBuildMap := true
			for j, param := range permFinal.Params {
				if param.Cond != "" {
					if shouldBuildMap {
						paramMap = createParamMapForEval(permFinal.Params)
						shouldBuildMap = false
					}
					shouldEval, err := evalExpression(replaceSymbols(permFinal.Params[j].Cond, paramMap))
					if err != nil {
						err = fmt.Errorf("could not evaluate expression for param %v: %s", param, err)
						errChan <- err
						return 0, err
					}
					if !shouldEval {
						permFinal.Params[j].Value = ""
					}
				}
			}

			// Calculate checksum
			permFinal.Checksum = permFinal.checksum()
			permFinal.BuildChecksum = permFinal.buildChecksum()

			// Block save until task is processed
			if j.SeqScheduler {
				_ = <-canPublish
			}

			// Send the permutation
			permutationsChan <- permFinal
		}

		// Inform the receiver that the generation is done
		treeChan <- j.tree

		// Wait for the receiver to confirm that the generation is done
		var remainingToGenerate uint64 = <-remainingChan
		if remainingToGenerate == 0 {
			break
		} else {
			i = limit - remainingToGenerate
		}
	}

	return 0, nil
}

// Permutations returns a list of all possible tasks based on parameterisation
func (j *JobSpec) Permutations(
	schedulerType uint,
	limit, maxPerm uint64,
) (
	permutationsChan chan *JobPermutation,
	errChan chan error,
	treeChan chan *tree.Tree,
	remainingChan chan uint64,
	canPublishChan chan bool,
	err error,
) {
	done := make(chan *tree.Tree)
	errors := make(chan error)
	permutations := make(chan *JobPermutation)
	remaining := make(chan uint64)
	canPublish := make(chan bool, 2)

	var allNr uint64

	// Reject generation if random and the job wants more than 0.75 of permutations generated
	if schedulerType == Random && limit > 0 && maxPerm > 0 && (float64(limit)/float64(maxPerm) > 0.75) {
		return nil, nil, nil, nil, nil, fmt.Errorf("too many permutations requested")
	}

	canPublish <- j.SeqScheduler
	go func() {
		// If not sequential, ignore first block and notfy the consumer
		if !j.SeqScheduler {
			_ = <-canPublish
		}
		switch schedulerType {
		case Grid:
			j.tree = tree.NewTree()
			node := j.tree.NewNode()
			node.SetName("Root")
			node.SetDepth(0)
			node.SetId(0)
			j.tree.SetRoot(node)

			j.next(1, permutations, errors, done, canPublish, allNr, limit, nil, node)
		case Random:
			j.random(permutations, errors, done, canPublish, remaining, limit)
		case Guided:
			// TODO implement Bayesian guided scheduling
		default:
			errors <- fmt.Errorf("unknown scheduler type: %d", schedulerType)
		}
	}()

	return permutations, errors, done, remaining, canPublish, nil
}

// TotalPermutations calculates the total number of permutations for the job
func (j *JobSpec) TotalPermutations() (uint64, error) {

	if len(j.Params) == 0 {
		return 0, fmt.Errorf("no parameters")
	}

	var total uint64 = 1

	for _, param := range j.Params {
		params, err := paramPermutations(&param)
		if err != nil {
			return 0, fmt.Errorf("could not parse parameter: %s", err)
		}

		total *= uint64(len(params))
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

// Checksum of all `build` and `both` parameters
func (jp *JobPermutation) buildChecksum() string {
	if len(jp.BuildChecksum) == 0 {

		// Calculate the UUID based on a reproducible md5 seed
		md5val := md5.New()
		for _, param := range jp.Params {
			if param.When == "test" {
				continue
			}
			io.WriteString(md5val, fmt.Sprintf("%s=%s\n", param.Name, param.Value))
		}

		jp.BuildChecksum = fmt.Sprintf("%x", md5val.Sum(nil))
	}

	return jp.BuildChecksum
}
