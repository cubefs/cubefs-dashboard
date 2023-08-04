// Copyright 2023 The CubeFS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package scheduler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs/blobstore/api/scheduler"

	"github.com/cubefs/cubefs-dashboard/backend/helper"
	"github.com/cubefs/cubefs-dashboard/backend/helper/httputils"
)

func LeaderStats(c *gin.Context, consulAddr string) (*scheduler.TasksStat, error) {
	resp, err := Request(c, http.MethodGet, consulAddr, scheduler.PathStatsLeader, "", nil, nil)
	if err != nil {
		return nil, err
	}
	output := &scheduler.TasksStat{}
	_, err = httputils.HandleResponse(c, resp, err, output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func Stats(c *gin.Context, nodeAddr string) (*scheduler.TasksStat, error) {
	reqUrl := nodeAddr + scheduler.PathStats
	resp, err := httputils.DoRequestBlobstore(c, reqUrl, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	output := &scheduler.TasksStat{}
	_, err = httputils.HandleResponse(c, resp, err, output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

type DiskMigratingStatsInput struct {
	DiskId   uint32 `json:"disk_id"`
	TaskType string `json:"task_type"`
}

func DiskMigratingStats(c *gin.Context, consulAddr string, input *DiskMigratingStatsInput) (*scheduler.DiskMigratingStats, error) {
	resp, err := Request(c, http.MethodGet, consulAddr, scheduler.PathStatsDiskMigrating, helper.BuildUrlParams(input), nil, nil)
	if err != nil {
		return nil, err
	}
	output := &scheduler.DiskMigratingStats{}
	_, err = httputils.HandleResponse(c, resp, err, output)
	if err != nil {
		return nil, err
	}
	return output, nil
}
