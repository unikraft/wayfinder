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
	"fmt"
	"math"
	"strconv"
)

type ParamSpec struct {
	Name        string   `json:"Name,omitempty"      yaml:"name,omitempty"`
	Description string   `json:"Description,omitempty" yaml:"description,omitempty"`
	Type        string   `json:"Type,omitempty"      yaml:"type,omitempty"`
	Default     string   `json:"Default,omitempty"   yaml:"default,omitempty"`
	Only        []string `json:"Only,omitempty"      yaml:"only,omitempty"`
	Min         string   `json:"Min,omitempty"       yaml:"min,omitempty"`
	Max         string   `json:"Max,omitempty"       yaml:"max,omitempty"`
	Step        string   `json:"Step,omitempty"      yaml:"step,omitempty"`
	StepMode    string   `json:"Step_mode,omitempty" yaml:"step_mode,omitempty"`

	// Child parameters
	Params *ParamSpec `yaml:"params"`
	When   string     `json:"When,omitempty" yaml:"when,omitempty"`
	If     string     `json:"If,omitempty"   yaml:"if,omitempty"`
}

type ParamPermutation struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
	Value       string `yaml:"value"`
	Cond        string `json:"Cond,omitempty" yaml:"cond,omitempty"`
	When        string `json:"when"`
}

// parseParamInt attends to string parameters and its possible permutations
func parseParamStr(param *ParamSpec) ([]ParamPermutation, error) {
	var params []ParamPermutation

	if len(param.Only) > 0 {
		for _, val := range param.Only {
			params = append(params, ParamPermutation{
				Name:        param.Name,
				Description: param.Description,
				Type:        param.Type,
				Value:       val,
				Cond:        param.If,
				When:        param.When,
			})
		}
	} else if len(param.Default) > 0 {
		params = append(params, ParamPermutation{
			Name:        param.Name,
			Description: param.Description,
			Type:        param.Type,
			Value:       param.Default,
			Cond:        param.If,
			When:        param.When,
		})
	}

	return params, nil
}

// parseParamInt attends to integer parameters and its possible permutations
func parseParamInt(param *ParamSpec) ([]ParamPermutation, error) {
	var params []ParamPermutation

	// Parse values in only
	if len(param.Only) > 0 {
		for _, val := range param.Only {
			params = append(params, ParamPermutation{
				Name:        param.Name,
				Description: param.Description,
				Type:        param.Type,
				Value:       val,
				Cond:        param.If,
				When:        param.When,
			})
		}

		// Parse range between min and max
	} else if len(param.Min) > 0 {
		min, err := strconv.Atoi(param.Min)
		if err != nil {
			return nil, err
		}

		max, err := strconv.Atoi(param.Max)
		if err != nil {
			return nil, err
		}

		if max < min {
			return nil, fmt.Errorf(
				"min can't be greater than max for %s: %d < %d", param.Name, min, max,
			)
		}

		// Figure out the step
		step := 1
		if len(param.Step) > 0 {
			step, err = strconv.Atoi(param.Step)
			if err != nil || step == 0 {
				return nil, fmt.Errorf(
					"invalid step for %s: %s", param.Name, param.Step,
				)
			}
		}

		// Use iterative step
		if len(param.StepMode) == 0 || param.StepMode == "increment" {
			for i := min; i <= max; i += step {
				params = append(params, ParamPermutation{
					Name:  param.Name,
					Type:  param.Type,
					Value: strconv.Itoa(i),
					Cond:  param.If,
					When:  param.When,
				})
			}

			// Use exponential step
		} else if param.StepMode == "power" {
			for i, j := min, min; i <= max; j++ {
				params = append(params, ParamPermutation{
					Name:  param.Name,
					Type:  param.Type,
					Value: strconv.Itoa(i),
					Cond:  param.If,
					When:  param.When,
				})
				i = int(math.Pow(float64(step), float64(j)))
			}

			// Unknown step mode
		} else {
			return nil, fmt.Errorf(
				"unknown step mode for param %s: %s", param.Name, param.StepMode,
			)
		}

	} else if len(param.Default) > 0 {
		params = append(params, ParamPermutation{
			Name:  param.Name,
			Type:  param.Type,
			Value: param.Default,
			Cond:  param.If,
			When:  param.When,
		})
	}

	return params, nil
}

// paramPermutations discovers all the possible variants of a particular
// parameter based on its type and options.
func paramPermutations(param *ParamSpec) ([]ParamPermutation, error) {
	switch t := param.Type; t {
	case "str",
		"string":
		param.Type = "str" // force rename the type
		return parseParamStr(param)
	case "int",
		"integer":
		param.Type = "int" // force rename the type
		return parseParamInt(param)
	}
	return nil, fmt.Errorf(
		"unknown parameter type: \"%s\" for %s", param.Type, param.Name,
	)
}
