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

package cluster

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs/proto"

	"github.com/cubefs/cubefs-dashboard/backend/helper/httputils"
)

func Get(c *gin.Context, addr string) (*proto.ClusterView, error) {
	reqUrl := "http://" + addr + proto.AdminGetCluster
	req, err := httputils.DoRequestOvertime(c, reqUrl, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	output := httputils.Output{Data: &proto.ClusterView{}}
	_, err = httputils.HandleResponse(c, req, err, &output)
	if err != nil {
		return nil, err
	}
	if output.Code != proto.ErrCodeSuccess {
		return nil, errors.New(output.Msg)
	}
	return output.Data.(*proto.ClusterView), nil
}

func GetNodes(c *gin.Context, addr, node string, isMetaNode bool) ([]string, error) {
	clusterInfo, err := Get(c, addr)
	if err != nil {
		return nil, err
	}
	nodes := make([]string, 0)
	if isMetaNode {
		for _, v := range clusterInfo.MetaNodes {
			nodes = handleNodeAppend(nodes, v.Addr, node)
		}
	} else {
		for _, v := range clusterInfo.DataNodes {
			nodes = handleNodeAppend(nodes, v.Addr, node)
		}
	}
	return nodes, nil
}

func handleNodeAppend(nodes []string, item, node string) []string {
	if node == "" {
		nodes = append(nodes, item)
		return nodes
	}
	if strings.Contains(item, node) {
		nodes = append(nodes, item)
	}
	return nodes
}
