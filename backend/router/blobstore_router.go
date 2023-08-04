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

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs-dashboard/backend/config"
	"github.com/cubefs/cubefs-dashboard/backend/handler/blobstore"
	"github.com/cubefs/cubefs-dashboard/backend/handler/blobstore/conf"
	"github.com/cubefs/cubefs-dashboard/backend/handler/blobstore/disk"
	"github.com/cubefs/cubefs-dashboard/backend/handler/blobstore/node"
	"github.com/cubefs/cubefs-dashboard/backend/handler/blobstore/service"
	"github.com/cubefs/cubefs-dashboard/backend/handler/blobstore/volume"
)

type blobRouter struct{}

func (b *blobRouter) Register(engine *gin.Engine) {
	group := engine.Group(config.Conf.Prefix.Api + "/console/blobstore")

	region := group.Group("/:cluster")
	{
		region.GET("/clusters/list", blobstore.ListClusters)
	}

	group = group.Group("/:cluster/:id")
	{
		// overview
		group.GET("/stat", blobstore.Overview)

		// raft
		group.POST("/leadership/transfer", blobstore.LeadershipTransfer)
		group.POST("/member/remove", blobstore.MemberRemove)
	}

	nodes := group.Group("/nodes")
	{
		nodes.GET("/list", node.List)
		nodes.POST("/access", node.Access)
		nodes.POST("/drop", node.Drop)
		nodes.POST("/offline", node.Offline)
		nodes.POST("/config/reload", node.ConfigReload)
		nodes.GET("/config/info", node.ConfigInfo)
		nodes.GET("/config/failures", node.ConfigFailures)
	}

	volumes := group.Group("/volumes")
	{
		volumes.GET("/list", volume.List)
		volumes.GET("/writing/list", volume.WritingList)
		volumes.GET("/v2/list", volume.V2List)
		volumes.GET("/allocated/list", volume.AllocatedList)
		volumes.GET("/get", volume.Get)
	}

	disks := group.Group("/disks")
	{
		disks.GET("/list", disk.List)
		disks.GET("/info", disk.Info)
		disks.GET("/dropping/list", disk.DroppingList)
		disks.POST("/access", disk.Access)
		disks.POST("/set", disk.SetBroken)
		disks.POST("/drop", disk.Drop)
		disks.POST("/probe", disk.Probe)
		disks.GET("/stats/migrating", disk.StatMigrating)
	}

	configs := group.Group("/config")
	{
		configs.GET("/list", conf.List)
		configs.POST("/set", conf.Set)
	}

	services := group.Group("/services")
	{
		services.GET("/list", service.List)
		services.GET("/get", service.Get)
		services.POST("/offline", service.Offline)
	}
}
