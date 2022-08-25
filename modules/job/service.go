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
	"os"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v2"

	"github.com/acarl005/stripansi"
	"github.com/unikraft/wayfinder/api/proto"
	"github.com/unikraft/wayfinder/internal/gzip"
	"github.com/unikraft/wayfinder/internal/models"
	"github.com/unikraft/wayfinder/spec"
)

type service struct {
	p *provider
}

func (s *service) CheckJobExists(ctx context.Context, req *proto.CheckJobExistsRequest) (*proto.CheckJobExistsResponse, error) {
	s.p.Log.Infof("checking if job already exists")
	exists := false

	if err := s.p.DB.Repos().Jobs().CheckJobExists(req.Checksum, &exists); err != nil {
		return &proto.CheckJobExistsResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "could not check for existing jobs")
	}

	return &proto.CheckJobExistsResponse{
		Success: true,
		Exists:  exists,
	}, nil
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
		TotalPermutations: totalPermutations.String(),
		Checksum:          req.Checksum,
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
	newPermutationInJob.PermutationLimit = "1"

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
			io.WriteString(md5valBuild, fmt.Sprintf("%s=%s\n", param.Name, value))
		}

		newPermutationInJob.CurrentPerm.Checksum = fmt.Sprintf("%x", md5val.Sum(nil))
		newPermutationInJob.CurrentPerm.BuildChecksum = fmt.Sprintf("%x", md5valBuild.Sum(nil))
		newPermutationInJob.CurrentPerm.JobId = uint(req.Id)
		newPermutationInJob.CurrentPerm.Id = 0
	}

	perm, err := s.p.DB.Repos().Permutations().FindOrCreateFromJobSpec(&newPermutationInJob)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not save job to DB: %s", err)
	}
	newPermutationInJob.CurrentPerm.Id = uint(perm.Id)

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
		Id:      int64(perm.Id),
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
	parsed.LaxMode = req.LaxMode

	// Set the permutation limit to the maximum if set to 0
	if req.PermutationLimit == "0" {
		parsed.PermutationLimit = job.TotalPermutations
	} else {
		parsed.PermutationLimit = fmt.Sprint(req.PermutationLimit)
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
				Uuid:             build.UUID.String(),
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
				Uuid:             test.UUID.String(),
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
			Id:       uint64(permutation.Id),
		})
	}

	return permutations
}

func (s *service) listPermutations(permutations []models.Permutation) []*proto.Permutation {
	var permutationsResult []*proto.Permutation

	for _, permutation := range permutations {
		var params []*proto.Param
		var status proto.JobPermutationStatus
		var builds []*proto.Build
		var tests []*proto.Test

		for _, param := range permutation.Params {
			params = append(params, &proto.Param{
				Name:     param.Name,
				ValueInt: param.ValueInt,
				ValueStr: param.ValueStr,
			})
		}

		status = permutation.Status

		for _, build := range permutation.Builds {
			builds = append(builds, &proto.Build{
				PermutationId: uint64(build.PermutationId),
				Status:        int32(build.Status),
				KernelPath:    build.KernelPath,
				InitRdPath:    build.InitRdPath,
				LogPath:       build.LogPath,
				Cores:         build.Cores,
			})
		}

		for _, test := range permutation.Tests {
			var results []*proto.Result

			for _, result := range test.Results {
				results = append(results, &proto.Result{
					Name:       result.Name,
					Type:       int32(result.Type),
					ValueStr:   result.ValueStr,
					ValueInt:   result.ValueInt,
					ValueFloat: result.ValueFloat,
				})
			}

			tests = append(tests, &proto.Test{
				PermutationId:  uint64(test.PermutationId),
				Status:         int32(test.Status),
				Results:        results,
				KernelCores:    test.KernelCores,
				BenchToolCores: test.BenchToolCores,
			})
		}

		permutationsResult = append(permutationsResult, &proto.Permutation{
			Uuid:     permutation.UUID[:],
			JobId:    int64(permutation.JobId),
			Checksum: permutation.Checksum,
			Params:   params,
			Status:   int64(status),
			Builds:   builds,
			Tests:    tests,
		})
	}

	return permutationsResult
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

func (s *service) GetPermutationJob(ctx context.Context, req *proto.GetPermutationJobRequest) (*proto.GetPermutationJobResponse, error) {
	s.p.Log.Infof("requested to get permutation %d from job %d...", req.PermId, req.Id)

	job := &models.Job{}
	if err := s.p.DB.Repos().Jobs().FindJobPermutation(req.Id, req.PermId, job); err != nil {
		return nil, status.Errorf(codes.NotFound, "job with Id=%d not found", req.Id)
	}

	found := false
	response := &proto.Permutation{}

	for _, permutation := range s.createPermutation(job) {
		if permutation.Id == uint64(req.PermId) {
			response = permutation
			found = true
			break
		}
	}

	return &proto.GetPermutationJobResponse{
		Success:     found,
		Permutation: response,
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
		}, status.Errorf(codes.Internal, "could not delete job with id=%d: %s", req.Id, err)
	}

	if req.Cascade {
		for _, permutation := range job.Permutations {
			// Delete results
			if err := s.p.DB.Repos().Results().DeleteResultsByPermutationId(int64(permutation.Id), req.Purge); err != nil {
				return &proto.DeleteJobResponse{
					Success: false,
				}, status.Errorf(codes.Internal, "could not delete results for permutation: %s", err)
			}

			// Delete Builds
			if err := s.p.DB.Repos().Builds().DeleteBuildsByPermutationId(int64(permutation.Id), req.Purge); err != nil {
				return &proto.DeleteJobResponse{
					Success: false,
				}, status.Errorf(codes.Internal, "could not delete builds for permutation: %s", err)
			}

			// Delete Tests
			if err := s.p.DB.Repos().Tests().DeleteTestsByPermutationId(int64(permutation.Id), req.Purge); err != nil {
				return &proto.DeleteJobResponse{
					Success: false,
				}, status.Errorf(codes.Internal, "could not delete tests for permutation: %s", err)
			}

			// Delete Permutations
			if err := s.p.DB.Repos().Permutations().DeleteById(int64(permutation.Id), req.Purge); err != nil {
				return &proto.DeleteJobResponse{
					Success: false,
				}, status.Errorf(codes.Internal, "could not delete permutation: %s", err)
			}
		}
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

func (s *service) ListPermutations(ctx context.Context, req *proto.ListPermutationsRequest) (*proto.ListPermutationsResponse, error) {
	s.p.Log.Infof("requested to list permutations for a job...")

	offset := int(req.Offset)
	limit := int(req.Limit)
	id := int(req.Id)
	var permutations []models.Permutation

	err := s.p.DB.Repos().Jobs().ListPermutationsForJob(id, offset, limit, &permutations)
	if err != nil {
		return &proto.ListPermutationsResponse{
			Success:      false,
			Total:        0,
			Permutations: nil,
		}, status.Errorf(codes.Internal, "could not list permutations: %s", err)
	}

	return &proto.ListPermutationsResponse{
		Success:      true,
		Total:        int64(len(permutations)),
		Permutations: s.listPermutations(permutations),
	}, nil
}

func (s *service) RedirectDebug(ctx context.Context, req *proto.RedirectDebugRequest) (*proto.RedirectDebugResponse, error) {
	s.p.Log.Infof("requested to redirect debug...")
	output := ""

	// Skip if incorrect request
	if req.LastNrOfLines <= 0 && req.LastNrOfLines > 100000 {
		return nil, status.Errorf(codes.InvalidArgument, "lastNrOfLines must be >= 0 and <= 100000")
	}

	// Retrieve last X lines from log file
	fName := s.p.Cfg.GeneralLogPath
	// Read an approximate number of lines from the end of the file
	lastBytes := req.LastNrOfLines * 100

	// Open log file
	logFile, err := os.OpenFile(fName, os.O_RDONLY, 0666)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not open log file: %s", err)
	}
	defer logFile.Close()

	// Seek to almost the end of the log
	stat, err := os.Stat(fName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not stat log file: %s", err)
	}
	fileSize := stat.Size()
	start := fileSize - lastBytes
	if start < 0 {
		start = 0
		lastBytes = fileSize
	}

	// Read last lastBytes bytes
	buf := make([]byte, lastBytes)
	_, err = logFile.ReadAt(buf, start)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not read log file: %s", err)
	}

	// Filter output content based on log level/jobId/permutationId
	// Skip first incomplete line
	firstToSkip := strings.Index(string(buf), "\n[")
	if firstToSkip == -1 {
		firstToSkip = 0
	}

	// The list of possible prefixes
	logLinePrefixes := make([]string, 7)
	logLinePrefixes[0] = ""
	logLinePrefixes[1] = "[:)]"
	logLinePrefixes[2] = "[DEBU]"
	logLinePrefixes[3] = "[INFO]"
	logLinePrefixes[4] = "WARN"
	logLinePrefixes[5] = "[ERRO]"
	logLinePrefixes[6] = "[FATA]"

	req.JobId = fmt.Sprintf("[%s]", req.JobId)
	req.PermutationId = fmt.Sprintf("[%s]", req.PermutationId)
	lines := strings.SplitAfterN(string(buf), "\n", firstToSkip)
	nrLinesOmitted := 0
	for i, line := range lines {
		// Strip
		line = stripansi.Strip(line)
		lines[i] = line
		// Check if the line contains the jobId
		if req.JobId != "[]" && !strings.Contains(line, req.JobId) {
			lines[i] = ""
			nrLinesOmitted++
			continue
		}
		// Check if the line contains the permutationId
		if req.PermutationId != "[]" && !strings.Contains(line, req.PermutationId) {
			lines[i] = ""
			nrLinesOmitted++
			continue
		}
		// Check if the line contains any prefix (including the ones more restrictive)
		shouldRemoveLine := true
		for i := req.DebugLevel; i <= proto.DebugLevel_DEBUG_LEVEL_FATAL; i++ {
			if strings.Contains(line, logLinePrefixes[i]) {
				shouldRemoveLine = false
				break
			}
		}

		if shouldRemoveLine {
			lines[i] = ""
			nrLinesOmitted++
		}
	}

	// If there are more lines than asked, remove from the first ones
	if int64(len(lines)-nrLinesOmitted) >= req.LastNrOfLines {
		linesToRemove := int64(len(lines)-nrLinesOmitted) - req.LastNrOfLines - 1
		for i, line := range lines {
			if line != "" {
				lines[i] = ""
				linesToRemove--
				if linesToRemove == 0 {
					break
				}
			}
		}
	}

	// Return the output content
	output = strings.Join(lines, "")

	return &proto.RedirectDebugResponse{
		Success: true,
		Output:  output,
	}, nil
}

// Returns the permutation status for a permutation Id
func (s *service) PermutationStatus(ctx context.Context, req *proto.PermutationStatusRequest) (*proto.PermutationStatusResponse, error) {
	s.p.Log.Infof("requested to get permutation status...")

	var permStatus proto.JobPermutationStatus

	err := s.p.DB.Repos().Permutations().GetPermutationStatus(int64(req.Id), &permStatus)
	if err != nil {
		return &proto.PermutationStatusResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "could not get permutation status: %s", err)
	}

	return &proto.PermutationStatusResponse{
		Success: true,
		Status:  permStatus,
	}, nil
}
