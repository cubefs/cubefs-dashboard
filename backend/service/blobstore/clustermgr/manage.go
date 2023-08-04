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

package clustermgr

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs/blobstore/api/clustermgr"

	"github.com/cubefs/cubefs-dashboard/backend/helper/httputils"
	"github.com/cubefs/cubefs-dashboard/backend/service/blobstore/api"
)

func Stat(c *gin.Context, addr string, isConsulAddr bool) (*clustermgr.StatInfo, error) {
	var (
		resp *http.Response
		err  error
	)
	if isConsulAddr {
		reqUrl := api.PathStat
		resp, err = api.DoRequestBlobstore(c, addr, reqUrl, http.MethodGet, nil, nil)
	} else {
		reqUrl := addr + api.PathStat
		resp, err = httputils.DoRequestBlobstore(c, reqUrl, http.MethodGet, nil, nil)
	}
	if err != nil {
		return nil, err
	}
	output := &clustermgr.StatInfo{}
	_, err = httputils.HandleResponse(c, resp, err, output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func LeadershipTransfer(c *gin.Context, consulAddr string, peerId uint64) (interface{}, error) {
	reqUrl := api.PathLeadershipTransfer
	b, _ := json.Marshal(map[string]uint64{"peer_id": peerId})
	resp, err := api.DoRequestBlobstore(c, consulAddr, reqUrl, http.MethodPost, bytes.NewReader(b), nil)
	if err != nil {
		return nil, err
	}
	var output interface{}
	_, err = httputils.HandleResponse(c, resp, err, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func MemberRemove(c *gin.Context, consulAddr string, peerId uint64) (interface{}, error) {
	reqUrl := api.PathMemberRemove
	b, _ := json.Marshal(map[string]uint64{"peer_id": peerId})
	resp, err := api.DoRequestBlobstore(c, consulAddr, reqUrl, http.MethodPost, bytes.NewReader(b), nil)
	if err != nil {
		return nil, err
	}
	var output interface{}
	_, err = httputils.HandleResponse(c, resp, err, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}
