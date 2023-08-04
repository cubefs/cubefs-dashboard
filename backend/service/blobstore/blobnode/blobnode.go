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

package blobnode

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs-dashboard/backend/helper"
	"github.com/cubefs/cubefs-dashboard/backend/helper/httputils"
	"github.com/cubefs/cubefs-dashboard/backend/service/blobstore/api"
)

func DiskProbe(c *gin.Context, nodeAddr, path string) error {
	reqUrl := nodeAddr + api.PathDiskProbe
	b, _ := json.Marshal(map[string]string{"path": path})
	resp, err := httputils.DoRequestBlobstore(c, reqUrl, http.MethodPost, bytes.NewReader(b), nil)
	if err != nil {
		return err
	}
	var output interface{}
	_, err = httputils.HandleResponse(c, resp, err, &output)
	if err != nil {
		return err
	}
	return nil
}

type ConfigReloadInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ConfigReload(c *gin.Context, nodeAddr string, input *ConfigReloadInput) error {
	reqUrl := nodeAddr + api.PathConfigReload + "?" + helper.BuildUrlParams(input)
	resp, err := httputils.DoRequestBlobstore(c, reqUrl, http.MethodPost, nil, nil)
	if err != nil {
		return err
	}
	var output interface{}
	_, err = httputils.HandleResponse(c, resp, err, &output)
	if err != nil {
		return err
	}
	return nil
}
