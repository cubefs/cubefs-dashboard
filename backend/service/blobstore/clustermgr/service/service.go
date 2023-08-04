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

package service

import (
	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/cubefs/cubefs/blobstore/api/clustermgr"

	"github.com/cubefs/cubefs-dashboard/backend/helper/httputils"
	"github.com/cubefs/cubefs-dashboard/backend/service/blobstore/api"
)

func List(c *gin.Context, consulAddr string) ([]clustermgr.ServiceNode, error) {
	reqUrl := api.PathServiceList
	resp, err := api.DoRequestBlobstore(c, consulAddr, reqUrl, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	output := &clustermgr.ServiceInfo{Nodes: make([]clustermgr.ServiceNode, 0)}
	_, err = httputils.HandleResponse(c, resp, err, output)
	if err != nil {
		return nil, err
	}
	return output.Nodes, nil
}

func Get(c *gin.Context, consulAddr, name string) ([]clustermgr.ServiceNode, error) {
	reqUrl := api.PathServiceGet + "?name=" + name
	resp, err := api.DoRequestBlobstore(c, consulAddr, reqUrl, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	output := &clustermgr.ServiceInfo{Nodes: make([]clustermgr.ServiceNode, 0)}
	_, err = httputils.HandleResponse(c, resp, err, output)
	if err != nil {
		return nil, err
	}
	return output.Nodes, nil
}

type UnregisterInput struct {
	Name string `json:"name"`
	Host string `json:"host"`
}

func Unregister(c *gin.Context, consulAddr string, input *UnregisterInput) error {
	reqUrl := api.PathServiceList
	b, _ := json.Marshal(input)
	resp, err := api.DoRequestBlobstore(c, consulAddr, reqUrl, http.MethodPost, bytes.NewReader(b), nil)
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
