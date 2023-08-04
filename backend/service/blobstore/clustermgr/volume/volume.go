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

package volume

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs/blobstore/api/clustermgr"

	"github.com/cubefs/cubefs-dashboard/backend/helper"
	"github.com/cubefs/cubefs-dashboard/backend/helper/httputils"
	"github.com/cubefs/cubefs-dashboard/backend/service/blobstore/api"
)

func Get(c *gin.Context, consulAddr string, vid uint32) (*clustermgr.VolumeInfo, error) {
	reqUrl := api.PathVolumeGet + "?" + fmt.Sprintf("vid=%d", vid)
	resp, err := api.DoRequestBlobstore(c, consulAddr, reqUrl, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	output := &clustermgr.VolumeInfo{}
	_, err = httputils.HandleResponse(c, resp, err, output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

type ListInput struct {
	Marker int `json:"marker"`
	Count  int `json:"count"`
}

func List(c *gin.Context, consulAddr string, input *ListInput) (*clustermgr.ListVolumes, error) {
	reqUrl := api.PathVolumeList + "?" + helper.BuildUrlParams(input)
	resp, err := api.DoRequestBlobstore(c, consulAddr, reqUrl, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	output := &clustermgr.ListVolumes{Volumes: make([]*clustermgr.VolumeInfo, 0)}
	_, err = httputils.HandleResponse(c, resp, err, output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func V2List(c *gin.Context, consulAddr string, status int) (*clustermgr.ListVolumes, error) {
	reqUrl := api.PathV2VolumeList + "?" + fmt.Sprintf("status=%d", status)
	resp, err := api.DoRequestBlobstore(c, consulAddr, reqUrl, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	output := &clustermgr.ListVolumes{Volumes: make([]*clustermgr.VolumeInfo, 0)}
	_, err = httputils.HandleResponse(c, resp, err, output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

type AllocatedListInput struct {
	Host     string `json:"host"`
	CodeMode uint8  `json:"code_mode"`
}

func AllocatedList(c *gin.Context, consulAddr string, input *AllocatedListInput) ([]clustermgr.AllocVolumeInfo, error) {
	reqUrl := api.PathVolumeAllocatedList + "?" + helper.BuildUrlParams(input)
	resp, err := api.DoRequestBlobstore(c, consulAddr, reqUrl, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	output := &clustermgr.AllocatedVolumeInfos{AllocVolumeInfos: make([]clustermgr.AllocVolumeInfo, 0)}
	_, err = httputils.HandleResponse(c, resp, err, output)
	if err != nil {
		return nil, err
	}
	return output.AllocVolumeInfos, nil
}
