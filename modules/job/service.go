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
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v2"

	"github.com/unikraft/wayfinder/api/proto"
	"github.com/unikraft/wayfinder/internal/gzip"
	"github.com/unikraft/wayfinder/internal/models"
	"github.com/unikraft/wayfinder/spec"
)

type service struct {
	p *provider
}

func (s *service) CreateJob(ctx context.Context, req *proto.CreateJobRequest) (*proto.CreateJobResponse, error) {
	var err error
	var data string

	s.p.Log.Infof("received new job create request...")

	if req.Compressed {
		var buf bytes.Buffer
		err = gzip.Decompress(&buf, req.Data)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not decompress data")
		}

		data = buf.String()
	} else {
		data = string(req.Data)
	}

	parsed, err := spec.ParseJobSpec(data)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "job specification invalid: %s", err)
	}

	name := parsed.Name
	if req.Name != "" {
		name = req.Name
	}

	totalPermutations, err := parsed.TotalPermutations()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "job specification invalid: %s", err)
	}

	// Save job to database
	job, err := s.p.DB.Repos().Jobs().CreateJob(&models.Job{
		Name:              name,
		Config:            data,
		TotalPermutations: totalPermutations,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not save job: %s", err)
	}

	return &proto.CreateJobResponse{
		Success: true,
		Id:      int64(job.Id),
	}, nil
}

// Searches for a job with the given id and creates a permutation within it and runs it
func (s *service) CreatePermutationJob(ctx context.Context, req *proto.CreatePermutationJobRequest) (*proto.CreatePermutationJobResponse, error) {
	s.p.Log.Infof("received new permutation job create request...")

	job := models.Job{}

	// Check if job exists
	if err := s.p.DB.Repos().Jobs().FindJob(req.Id, 0, 1, &job); err != nil {
		return nil, status.Errorf(codes.NotFound, "job with Id=%d not found", req.Id)
	}

	newPermutationInJob := spec.JobSpec{}
	yaml.Unmarshal([]byte(job.Config), &newPermutationInJob)

	// Format permutation
	newPermutationInJob.Id = uint(req.Id)
	newPermutationInJob.PermutationLimit = 1

	// Calculate checksum for given permutation
	if len(newPermutationInJob.CurrentPerm.Checksum) == 0 {
		newPermutationInJob.CurrentPerm.Params = make([]spec.ParamPermutation, 0)

		// Save name and type to currentperm
		for _, p := range req.Params {
			newPermutationInJob.CurrentPerm.Params = append(
				newPermutationInJob.CurrentPerm.Params, spec.ParamPermutation{
					Name: p.Name,
					Type: p.Type,
				})
		}

		md5val := md5.New()
		md5valBuild := md5.New()
		for i, param := range newPermutationInJob.CurrentPerm.Params {
			var value string
			switch param.Type {
			case "int", "integer":
				value = fmt.Sprintf("%d", req.Params[i].ValueInt)
			case "str", "string":
				value = req.Params[i].ValueStr
			}
			newPermutationInJob.CurrentPerm.Params[i].Value = value
			io.WriteString(md5val, fmt.Sprintf("%s=%s\n", param.Name, value))

			// Build no matter what
			io.WriteString(md5valBuild, fmt.Sprintf("%s=%s\n", param.Name, param.Value))
		}

		newPermutationInJob.CurrentPerm.Checksum = fmt.Sprintf("%x", md5val.Sum(nil))
		newPermutationInJob.CurrentPerm.BuildChecksum = fmt.Sprintf("%x", md5valBuild.Sum(nil))
		newPermutationInJob.CurrentPerm.JobId = uint(req.Id)
		newPermutationInJob.CurrentPerm.Id = 0
	}

	// Publish permutation to queue
	taskBytes, err := json.Marshal(newPermutationInJob)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not marshal job: %s", err)
	}

	var compressedBytes bytes.Buffer
	w := zlib.NewWriter(&compressedBytes)
	w.Write(taskBytes)
	w.Close()

	s.p.TaskQueue.PublishBytes(compressedBytes.Bytes())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not publish job to queue: %s", err)
	}

	return &proto.CreatePermutationJobResponse{
		Success: true,
		Id:      int64(job.Id),
	}, nil
}

func (s *service) StartJob(ctx context.Context, req *proto.StartJobRequest) (*proto.StartJobResponse, error) {
	s.p.Log.Infof("requested to start job %d...", req.Id)

	job := &models.Job{}
	if err := s.p.DB.Repos().Jobs().FindJob(req.Id, 0, 1, job); err != nil {
		return nil, status.Errorf(codes.NotFound, "job with Id=%d not found", req.Id)
	}

	parsed, err := spec.ParseJobSpec(job.Config)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "job specification invalid: %s", err)
	}

	parsed.Id = job.Id
	parsed.Scheduler = req.Scheduler
	parsed.IsolLevel = req.IsolLevel
	parsed.IsolSplit = req.IsolSplit
	parsed.Repeats = req.Repeats
	parsed.DryRun = req.DryRun
	parsed.SeqScheduler = req.SeqScheduler

	// Set the permutation limit to the maximum if set to 0
	if req.PermutationLimit == 0 {
		parsed.PermutationLimit = uint64(job.TotalPermutations)
	} else {
		parsed.PermutationLimit = uint64(req.PermutationLimit)
	}

	specBytes, err := json.Marshal(parsed)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not marshal job: %s", err)
	}

	err = s.p.JobQueue.PublishBytes(specBytes)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not publish job to queue: %s", err)
	}

	return &proto.StartJobResponse{
		Success: true,
	}, nil
}

func (s *service) createPermutation(job *models.Job) []*proto.Permutation {
	var permutations []*proto.Permutation

	for _, permutation := range job.Permutations {
		var results []*proto.Result
		var params []*proto.Param
		var builds []*proto.Build
		var tests []*proto.Test

		for _, result := range permutation.Results {
			results = append(results, &proto.Result{
				Name:       result.Name,
				Type:       int32(result.Type),
				ValueStr:   result.ValueStr,
				ValueInt:   int64(result.ValueInt),
				ValueBool:  result.ValueBool,
				ValueFloat: result.ValueFloat,
			})
		}

		for _, param := range permutation.Params {
			params = append(params, &proto.Param{
				Name:     param.Name,
				Type:     param.Type,
				ValueStr: param.ValueStr,
				ValueInt: int64(param.ValueInt),
			})
		}

		for _, build := range permutation.Builds {
			builds = append(builds, &proto.Build{
				PermutationId:    uint64(build.PermutationId),
				Status:           int32(build.Status),
				Runtime:          build.Runtime.String(),
				WayfinderVersion: build.WayfinderVersion,
				KernelPath:       build.KernelPath,
				InitRdPath:       build.InitRdPath,
				LogPath:          build.LogPath,
				Cores:            build.Cores,
			})
		}

		for _, test := range permutation.Tests {
			tests = append(tests, &proto.Test{
				PermutationId:    uint64(test.PermutationId),
				Status:           int32(test.Status),
				Runtime:          test.Runtime.String(),
				WayfinderVersion: test.WayfinderVersion,
				Results:          results,
				VmmCores:         test.VMMCores,
				KernelCores:      test.KernelCores,
				BenchToolCores:   test.BenchToolCores,
			})
		}

		permutations = append(permutations, &proto.Permutation{
			Uuid:     permutation.UUID[:],
			JobId:    int64(permutation.JobId),
			Checksum: permutation.Checksum,
			Status:   int64(permutation.Status),
			Params:   params,
			Builds:   builds,
			Tests:    tests,
		})
	}

	return permutations
}

func (s *service) GetJob(ctx context.Context, req *proto.GetJobRequest) (*proto.GetJobResponse, error) {
	s.p.Log.Infof("requested to get job %d...", req.Id)

	job := &models.Job{}
	if err := s.p.DB.Repos().Jobs().FindJob(req.Id, int(req.Offset), int(req.Limit), job); err != nil {
		return nil, status.Errorf(codes.NotFound, "job with Id=%d not found", req.Id)
	}

	response := &proto.Job{
		Id:                int64(job.Id),
		Status:            job.Status,
		HostId:            uint64(job.HostId),
		Config:            job.Config,
		CompletedAt:       job.CompletedAt.String(),
		TotalPermutations: job.TotalPermutations,
		Permutations:      s.createPermutation(job),
	}

	return &proto.GetJobResponse{
		Success: true,
		Job:     response,
	}, nil
}

func (s *service) DeleteJob(ctx context.Context, req *proto.DeleteJobRequest) (*proto.DeleteJobResponse, error) {
	s.p.Log.Infof("deleting job with id=%d...", req.Id)

	job := &models.Job{}
	if err := s.p.DB.Repos().Jobs().FindJob(req.Id, 0, 1, job); err != nil {
		return &proto.DeleteJobResponse{
			Success: false,
		}, status.Errorf(codes.NotFound, "job with id=%d not found", req.Id)
	}

	if err := s.p.DB.Repos().Jobs().DeleteJob(req.Id, req.Purge); err != nil {
		return &proto.DeleteJobResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "could not delete job with id=%d: %s", err)
	}

	return &proto.DeleteJobResponse{
		Success: true,
	}, nil
}

func (s *service) PauseRedisQueues(ctx context.Context, req *proto.PauseRedisQueuesRequest) (*proto.PauseRedisQueuesResponse, error) {
	s.p.Log.Infof("pausing redis queues...")

	if s.p.Redis.ClientPause(s.p.Redis.Context(), time.Millisecond*time.Duration(req.Time)).Err() != nil {
		return &proto.PauseRedisQueuesResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "could not pause redis queues")
	}

	fmt.Printf("Paused redis queues for %dms\n", req.Time)

	return &proto.PauseRedisQueuesResponse{
		Success: true,
	}, nil
}

func (s *service) FlushRedisQueue(ctx context.Context, req *proto.FlushRedisQueueRequest) (*proto.FlushRedisQueueResponse, error) {
	s.p.Log.Infof("flushing redis queue...")

	switch req.QueueType {
	case proto.RedisQueueType_REDIS_QUEUE_TYPE_JOB:
		s.p.Log.Infof("flushing job queue...")

		nrPurged, err := s.p.JobQueue.PurgeReady()
		if err != nil {
			s.p.Log.Errorf("could not purge job queue: %s", err)
		}

		s.p.Log.Infof("Flushed job queue, %d jobs purged", nrPurged)

		return &proto.FlushRedisQueueResponse{
			Success: true,
		}, nil
	case proto.RedisQueueType_REDIS_QUEUE_TYPE_PERMUTATION:
		s.p.Log.Infof("flushing permutation queue...")

		nrPurged, err := s.p.TaskQueue.PurgeReady()
		if err != nil {
			s.p.Log.Errorf("could not purge task queue: %s", err)
		}

		s.p.Log.Infof("Flushed task queue, %d permutations purged", nrPurged)

		return &proto.FlushRedisQueueResponse{
			Success: true,
		}, nil
	case proto.RedisQueueType_REDIS_QUEUE_TYPE_ALL:
		var totalPurged int64
		s.p.Log.Infof("flushing all queues...")

		nrPurged, err := s.p.JobQueue.PurgeReady()
		if err != nil {
			s.p.Log.Errorf("could not purge job queue: %s", err)
		}
		totalPurged = nrPurged

		nrPurged, err = s.p.TaskQueue.PurgeReady()
		if err != nil {
			s.p.Log.Errorf("could not purge task queue: %s", err)
		}
		totalPurged += nrPurged

		s.p.Log.Infof("Flushed all queues, %d elements purged", totalPurged)

		return &proto.FlushRedisQueueResponse{
			Success: true,
		}, nil
	default:
		return &proto.FlushRedisQueueResponse{
			Success: false,
		}, status.Errorf(codes.InvalidArgument, "invalid queue type: %s", req.QueueType)
	}
}

func (s *service) GetJobResults(ctx context.Context, req *proto.GetJobResultsRequest) (*proto.GetJobResultsResponse, error) {
	s.p.Log.Infof("requested to get job %d results...", req.Id)

	offset := int(req.Offset)
	limit := int(req.Limit)

	results, err := s.p.DB.Repos().Results().FindResults(uint(req.Id), offset, limit)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "job with Id=%d not found", req.Id)
	}

	var resultsResponse []*proto.Result

	for _, result := range results {
		resultsResponse = append(resultsResponse, &proto.Result{
			Name:       result.Name,
			Type:       int32(result.Type),
			ValueStr:   result.ValueStr,
			ValueInt:   result.ValueInt,
			ValueFloat: result.ValueFloat,
			ValueBool:  result.ValueBool,
		})
	}

	return &proto.GetJobResultsResponse{
		Success: true,
		Total:   int64(len(results)),
		Results: resultsResponse,
	}, nil
}

func (s *service) ListJobs(ctx context.Context, req *proto.ListJobsRequest) (*proto.ListJobsResponse, error) {
	s.p.Log.Infof("requested to list jobs...")

	offset := int(req.Offset)
	limit := int(req.Limit)

	jobs, err := s.p.DB.Repos().Jobs().ListJobs(offset, limit)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not list jobs: %s", err)
	}

	var jobsResponse []*proto.Job

	for _, job := range jobs {
		jobsResponse = append(jobsResponse, &proto.Job{
			Id:                int64(job.Id),
			Status:            job.Status,
			HostId:            uint64(job.HostId),
			Config:            job.Config,
			CompletedAt:       job.CompletedAt.String(),
			TotalPermutations: job.TotalPermutations,
			Permutations:      s.createPermutation(job),
		})
	}

	return &proto.ListJobsResponse{
		Success: true,
		Total:   int64(len(jobs)),
		Jobs:    jobsResponse,
	}, nil
}
