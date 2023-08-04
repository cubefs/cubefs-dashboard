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
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs/blobstore/common/proto"

	"github.com/cubefs/cubefs-dashboard/backend/helper/httputils"
	"github.com/cubefs/cubefs-dashboard/backend/service/blobstore/clustermgr/service"
)

func Request(c *gin.Context, method, consulAddr, path, query string, body io.Reader, header map[string]string) (*http.Response, error) {
	hosts, err := GetSchedulers(c, consulAddr)
	if err != nil {
		return nil, err
	}
	if len(hosts) == 0 {
		return nil, errors.New("scheduler host not found")
	}
	var resp *http.Response
	for _, host := range hosts {
		reqUrl := host + path
		if query != "" {
			reqUrl += "?" + query
		}
		resp, err = httputils.DoRequestBlobstore(c, reqUrl, method, body, header)
		if err == nil {
			return resp, nil
		}
	}
	return resp, err
}

func GetSchedulers(c *gin.Context, consulAddr string) ([]string, error) {
	hosts := make([]string, 0)
	nodes, err := service.List(c, consulAddr)
	if err != nil {
		return hosts, err
	}
	for _, node := range nodes {
		if node.Name == proto.ServiceNameScheduler {
			hosts = append(hosts, node.Host)
		}
	}
	return hosts, nil
}
