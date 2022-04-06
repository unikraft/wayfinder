package job

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
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/erda-project/erda-infra/base/logs"

	"github.com/unikraft/wayfinder/internal/models"
	"github.com/unikraft/wayfinder/spec"
)

type JobConsumer struct {
	p   *provider
	Log logs.Logger
}

func NewJobConsumer(p *provider) *JobConsumer {
	return &JobConsumer{p: p}
}

func (c *JobConsumer) Consume(delivery rmq.Delivery) {
	specBytes := delivery.Payload()
	job := spec.JobSpec{}

	// Check if we received the full job specification
	err := json.Unmarshal([]byte(specBytes), &job)
	if err != nil {
		if err := delivery.Reject(); err != nil {
			c.p.Log.Errorf("failed to reject job: %s", err)
		} else {
			c.p.Log.Errorf("could not unmarshal job %s", err)
		}
	}

	// TODO: check validity of job specification
	// e.g.: if err = job.Validity(); err != nil {}

	// Create a new logger for this job
	c.Log = c.p.Log.Sub(fmt.Sprintf("%d", job.Id))

	// Start the job
	err = c.StartJob(&job)
	if err != nil {
		if err := delivery.Reject(); err != nil {
			c.p.Log.Errorf("failed to reject job: %s", err)
		} else {
			c.p.Log.Errorf("could not complete job %s", err)
		}
	} else {
		delivery.Ack()
	}
}

func (c *JobConsumer) StartJob(jobSpec *spec.JobSpec) error {
	// Grab the job from the database
	job := &models.Job{}
	if err := c.p.DB.DB().Where("id = ?", jobSpec.Id).First(&job).Error; err != nil {
		return fmt.Errorf("job with Id=%d not found", jobSpec.Id)
	}

	// Update the state of the job
	if err := c.p.DB.Repos().Jobs().SetStatusRunningByJobId(job.Id); err != nil {
		return fmt.Errorf("could not update job state: %s", err)
	}

	// Calculate max possible number of permutations (unconditioned)
	maxPerm, _ := jobSpec.TotalPermutations()
	c.Log.Infof("Calculated max number of permutations: %d", maxPerm)

	c.Log.Info("calculating all permutation values...")

	permutations, errors, done, remaining, canPublish, err := jobSpec.Permutations(
		uint(jobSpec.Scheduler), jobSpec.PermutationLimit, maxPerm)
	if err != nil {
		return fmt.Errorf("could not calculate permutations: %s", err)
	}

	// See if the generator is sequential or not
	isSequential := <-canPublish

	// Send first signal to the generator
	canPublish <- true

	var totalPermutations uint64 = 0
	uniquePermutations := make(map[string]bool)
	for {
		select {
		// Error occured during permutation calculation
		case err := <-errors:
			return fmt.Errorf("could not calculate permutation: %s", err)

		// All permutations calculated
		case generatedTree := <-done:
			if totalPermutations < jobSpec.PermutationLimit {
				remaining <- jobSpec.PermutationLimit - totalPermutations
				continue
			}
			c.Log.Infof("calculated number of permutations: %d", totalPermutations)
			if generatedTree != nil {
				// TODO export generated tree
			}
			return nil

		// New permutation
		case perm := <-permutations:

			// Add the Job ID to the permutation
			perm.JobId = jobSpec.Id
			jobSpec.CurrentPerm = *perm

			// Check if the permutation is unique, if it isn't skip it
			if _, notUnqiue := uniquePermutations[perm.Checksum]; notUnqiue {
				continue
			} else {
				uniquePermutations[perm.Checksum] = true
				taskBytes, err := json.Marshal(*jobSpec)
				if err != nil {
					// Or should we just drop the permutation?
					errors <- fmt.Errorf("could not marshal permutation: %s", err)
				}

				totalPermutations++
				var compressedBytes bytes.Buffer
				w := zlib.NewWriter(&compressedBytes)
				w.Write(taskBytes)
				w.Close()

				if isSequential {
					// Busy wait to free up space for the task executor
					for {
						if c.p.Redis.LLen(c.p.Redis.Context(), "rmq::queue::[task-queue]::ready").Val() == 0 {
							break
						}
						time.Sleep(time.Second)
					}

					canPublish <- true
				}

				// Publish permutation to the task queue.  This will be picked up by the
				// scheduler.
				for i := 0; i < 1+int(jobSpec.Repeats); i++ {
					c.p.TaskQueue.PublishBytes(compressedBytes.Bytes())
				}
			}
		}
	}
}
