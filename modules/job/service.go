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
  "context"
  "encoding/json"

  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"

  "github.com/unikraft/wayfinder/spec"
  "github.com/unikraft/wayfinder/api/proto"
  "github.com/unikraft/wayfinder/internal/gzip"
  "github.com/unikraft/wayfinder/internal/models"
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

  totalPermutations, err := parsed.TotalPermutations()
  if err != nil {
    return nil, status.Errorf(codes.InvalidArgument, "job specification invalid: %s", err)
  }

  // Save job to database
  job, err := s.p.DB.Repos().Jobs().CreateJob(&models.Job{
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

func (s *service) StartJob(ctx context.Context, req *proto.StartJobRequest) (*proto.StartJobResponse, error) {
  s.p.Log.Infof("requested to start job %d...", req.Id)

  job := &models.Job{}
  if err := s.p.DB.Repos().Jobs().FindJob(uint(req.Id), job); err != nil {
    return nil, status.Errorf(codes.NotFound, "job with Id=%d not found", req.Id)
  }

  parsed, err := spec.ParseJobSpec(job.Config)
  if err != nil {
    return nil, status.Errorf(codes.InvalidArgument, "job specification invalid: %s", err)
  }

  parsed.Id        = job.Id
  parsed.Scheduler = req.Scheduler
  parsed.IsolLevel = req.IsolLevel
  parsed.IsolSplit = req.IsolSplit
  parsed.Repeats   = req.Repeats

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

func (s *service) GetJob(ctx context.Context, req *proto.GetJobRequest) (*proto.GetJobResponse, error) {
  s.p.Log.Infof("requested to get job %d...", req.Id)

  job := &models.Job{}
  if err := s.p.DB.Repos().Jobs().FindJob(uint(req.Id), job); err != nil {
    return nil, status.Errorf(codes.NotFound, "job with Id=%d not found", req.Id)
  }

  jobBytes, err := json.Marshal(job)
  if err != nil {
    return nil, status.Errorf(codes.Internal, "could not marshal job: %s", err)
  }

  return &proto.GetJobResponse{
    Success: true,
    Data:   jobBytes,
  }, nil
}

func (s *service) GetJobResults(ctx context.Context, req *proto.GetJobResultsRequest) (*proto.GetJobResultsResponse, error) {
  s.p.Log.Infof("requested to get job %d results...", req.Id)

  results, err := s.p.DB.Repos().Results().FindResults(uint(req.Id))

  if err != nil {
    return nil, status.Errorf(codes.NotFound, "job with Id=%d not found", req.Id)
  }

  resultsBytes, err := json.Marshal(results)

  if err != nil {
    return nil, status.Errorf(codes.Internal, "could not marshal results: %s", err)
  }

  return &proto.GetJobResultsResponse{
    Success: true,
    Data:    resultsBytes,
  }, nil
}

func (s *service) ListJobs(ctx context.Context, req *proto.ListJobsRequest) (*proto.ListJobsResponse, error) {
  s.p.Log.Infof("requested to list jobs...")

  offset := int(req.Offset)
  limit  := int(req.Limit)

  jobs, err := s.p.DB.Repos().Jobs().ListJobs(offset, limit)

  if err != nil {
    return nil, status.Errorf(codes.Internal, "could not list jobs: %s", err)
  }

  jobsBytes, err := json.Marshal(jobs)

  if err != nil {
    return nil, status.Errorf(codes.Internal, "could not marshal jobs: %s", err)
  }

  return &proto.ListJobsResponse{
    Success: true,
    Data:    jobsBytes,
  }, nil
}
