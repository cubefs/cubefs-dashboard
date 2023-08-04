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

package conf

import (
	"github.com/cubefs/cubefs/blobstore/common/proto"
	"github.com/cubefs/cubefs/blobstore/util/log"
	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs-dashboard/backend/helper/codes"
	"github.com/cubefs/cubefs-dashboard/backend/helper/ginutils"
	"github.com/cubefs/cubefs-dashboard/backend/service/blobstore/clustermgr/config"
	"github.com/cubefs/cubefs-dashboard/backend/service/blobstore/clustermgr/service"
	"github.com/cubefs/cubefs-dashboard/backend/service/blobstore/scheduler"
)

type ListOutput struct {
	Repair        []interface{} `json:"repair"`
	Drop          []interface{} `json:"drop"`
	Balance       []interface{} `json:"balance"`
	ManualMigrate []interface{} `json:"manual_migrate"`
	Inspect       []interface{} `json:"inspect"`
	ShardRepair   []*TaskData   `json:"shard_repair"`
	BlobDelete    []*TaskData   `json:"blob_delete"`
}

type TaskData struct {
	Host        string      `json:"host"`
	HostName    string      `json:"host_name"`
	ShardRepair interface{} `json:"shard_repair"`
	BlobDelete  interface{} `json:"blob_delete"`
}

func List(c *gin.Context) {
	consulAddr, err := ginutils.GetConsulAddr(c)
	if err != nil {
		return
	}
	taskStat, err := scheduler.LeaderStats(c, consulAddr)
	if err != nil {
		log.Errorf("scheduler.LeaderStats failed.consul_addr:%s,err:%+v", consulAddr, err)
		ginutils.Send(c, codes.ThirdPartyError.Code(), err.Error(), nil)
		return
	}
	output := ListOutput{
		Repair:        []interface{}{taskStat.DiskRepair},
		Drop:          []interface{}{taskStat.DiskDrop},
		Balance:       []interface{}{taskStat.Balance},
		ManualMigrate: []interface{}{taskStat.ManualMigrate},
		Inspect:       []interface{}{taskStat.VolumeInspect},
		ShardRepair:   make([]*TaskData, 0),
		BlobDelete:    make([]*TaskData, 0),
	}
	serviceNodes, err := service.List(c, consulAddr)
	if err != nil {
		log.Errorf("service.List failed.consul_addr:%s,err:%+v", consulAddr, err)
		ginutils.Send(c, codes.ThirdPartyError.Code(), err.Error(), nil)
		return
	}
	for _, node := range serviceNodes {
		if node.Name != proto.ServiceNameScheduler {
			continue
		}
		stats, err := scheduler.Stats(c, node.Host)
		if err != nil {
			log.Errorf("scheduler.Stats failed.node_addr:%s,err:%+v", node.Host, err)
			continue
		}

		dataRepair := &TaskData{Host: node.Host + "_down"}
		if stats.ShardRepair != nil {
			dataRepair = &TaskData{Host: node.Host, ShardRepair: stats.ShardRepair}
		}
		output.ShardRepair = append(output.ShardRepair, dataRepair)

		dataDelete := &TaskData{Host: node.Host + "_down"}
		if stats.BlobDelete != nil {
			dataDelete = &TaskData{Host: node.Host, BlobDelete: stats.BlobDelete}
		}
		output.BlobDelete = append(output.BlobDelete, dataDelete)
	}

	ginutils.Send(c, codes.OK.Code(), codes.OK.Msg(), output)
}

type SetInput struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

func Set(c *gin.Context) {
	input := &SetInput{}
	consulAddr, err := ginutils.CheckAndGetConsul(c, input)
	if err != nil {
		return
	}
	err = config.Set(c, consulAddr, &config.SetInput{Key: input.Key, Value: input.Value})
	if err != nil {
		log.Errorf("config.Set failed.consulAddr:%s,args:%+v,err:%+v", consulAddr, input, err)
		ginutils.Send(c, codes.ThirdPartyError.Code(), err.Error(), nil)
		return
	}
	ginutils.Send(c, codes.OK.Code(), codes.OK.Msg(), nil)
}
