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

package config

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs/blobstore/common/codemode"

	"github.com/cubefs/cubefs-dashboard/backend/helper/httputils"
	"github.com/cubefs/cubefs-dashboard/backend/service/blobstore/api"
)

func Get(c *gin.Context, consulAddr string, key string) ([]codemode.Policy, error) {
	reqUrl := api.PathConfigGet
	if key != "" {
		reqUrl += "?key=" + key
	}
	resp, err := api.DoRequestBlobstore(c, consulAddr, reqUrl, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	var output string
	_, err = httputils.HandleResponse(c, resp, err, &output)
	if err != nil {
		return nil, err
	}
	conf := make([]codemode.Policy, 0)
	err = json.Unmarshal([]byte(output), &conf)
	return conf, err
}

type SetInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func Set(c *gin.Context, consulAddr string, input *SetInput) error {
	reqUrl := api.PathConfigSet
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
