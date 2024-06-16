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
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cubefs/cubefs/blobstore/util/log"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"github.com/cubefs/cubefs-dashboard/backend/helper"
	"github.com/cubefs/cubefs-dashboard/backend/helper/codes"
	"github.com/cubefs/cubefs-dashboard/backend/helper/enums"
	"github.com/cubefs/cubefs-dashboard/backend/helper/ginutils"
	"github.com/cubefs/cubefs-dashboard/backend/helper/pool"
	"github.com/cubefs/cubefs-dashboard/backend/helper/types"
	"github.com/cubefs/cubefs-dashboard/backend/model"
	"github.com/cubefs/cubefs-dashboard/backend/service/cluster"
)

type CreateInput struct {
	Name       string         `json:"name" binding:"required"`
	MasterAddr types.StrSlice `json:"master_addr" binding:"required"`
	IDC        string         `json:"idc"`
	Cli        string         `json:"cli"`
	Domain     string         `json:"domain"`
	ConsulAddr string         `json:"consul_addr" binding:"omitempty,url"`
	Tag        string         `json:"tag"`
	S3Endpoint string         `json:"s3_endpoint" binding:"omitempty,url"`
	VolType    enums.VolType  `json:"vol_type"`
	Region     string         `json:"region"`
}

func (input *CreateInput) Check() error {
	switch input.VolType {
	case enums.VolTypeLowFrequency:
		if input.Region == "" {
			return errors.New("region is needed")
		}
		return nil
	default:
		return nil
	}
}

func (input *CreateInput) checkAddr(c *gin.Context) error {
	if input.VolType == enums.VolTypeLowFrequency {
		addrs := []string{input.ConsulAddr}
		addrs = append(addrs, input.MasterAddr...)
		return checkAddrIpPort(c, addrs)
	}
	return checkMasterAddr(c, input.MasterAddr)
}

func checkMasterAddr(c *gin.Context, masterAddr types.StrSlice) error {
	for _, v := range masterAddr {
		if err := checkIpPort(c, v); err != nil {
			return err
		}
		if _, err := cluster.Get(c, v); err != nil {
			log.Errorf("[%s] get cluster failed. err:%+v", v, err)
			ginutils.Fail(c, codes.ThirdPartyError.Code(), fmt.Sprintf("ip: %s 無效", v))
			return err
		}
	}
	return nil
}

func checkIpPort(c *gin.Context, addr string) error {
	idx := strings.Index(addr, "//")
	if idx > 0 {
		addr = addr[idx+2:]
	}
	if err := helper.CheckIpPort(addr); err != nil {
		log.Errorf("checkIpPort failed. addr:%s, err:%+v", addr, err)
		ginutils.Fail(c, codes.InvalidArgs.Code(), err.Error())
		return err
	}
	return nil
}

func checkAddrIpPort(c *gin.Context, addrs []string) error {
	for _, addr := range addrs {
		if err := checkIpPort(c, addr); err != nil {
			return err
		}
	}
	return nil
}

func Create(c *gin.Context) {
	input := &CreateInput{}
	if !ginutils.Check(c, input) {
		return
	}
	if input.checkAddr(c) != nil {
		return
	}
	if err := checkNameExists(c, input.Name); err != nil {
		return
	}
	clusterInfo, err := cluster.Get(c, input.MasterAddr[0])
	if err != nil {
		ginutils.Send(c, codes.ThirdPartyError.Error(), err.Error(), nil)
		return
	}
	if err := checkTagExists(c, clusterInfo.Name); err != nil {
		return
	}
	input.Tag = clusterInfo.Name
	m, err := create(input)
	if err != nil {
		ginutils.Send(c, codes.DatabaseError.Code(), err.Error(), nil)
		return
	}
	ginutils.Send(c, codes.OK.Code(), codes.OK.Msg(), m)
}

type UpdateInput struct {
	Id         int64          `json:"id" binding:"required"`
	MasterAddr types.StrSlice `json:"master_addr"`
	IDC        string         `json:"idc"`
	Cli        string         `json:"cli"`
	Domain     string         `json:"domain"`
	ConsulAddr string         `json:"consul_addr" binding:"omitempty,url"`
	Tag        string         `json:"tag"`
	S3Endpoint string         `json:"s3_endpoint" binding:"omitempty,url"`
}

func Update(c *gin.Context) {
	input := &UpdateInput{}
	if !ginutils.Check(c, input) {
		return
	}
	cm := new(model.Cluster)
	if err := cm.FindId(input.Id); err != nil {
		log.Errorf("cm.FindId. id:%s, error:%+v", input.Id, err)
		ginutils.Send(c, codes.DatabaseError.Code(), err.Error(), nil)
		return
	}
	set := bson.M{"update_time": time.Now()}
	if len(input.MasterAddr) > 0 {
		var err error
		if cm.VolType == enums.VolTypeLowFrequency {
			err = checkAddrIpPort(c, input.MasterAddr)
		} else {
			err = checkMasterAddr(c, input.MasterAddr)
		}
		if err != nil {
			return
		}
		set["master_addr"] = input.MasterAddr.String()
		clusterInfo, err := cluster.Get(c, input.MasterAddr[0])
		if err != nil {
			ginutils.Send(c, codes.ThirdPartyError.Error(), err.Error(), nil)
			return
		}
		clusterModel, err := new(model.Cluster).FindTag(clusterInfo.Name)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Errorf("by cluster by tag failed. err:%+v", err)
			ginutils.Send(c, codes.DatabaseError.Code(), err.Error(), nil)
			return
		}
		if clusterModel != nil && clusterModel.Name != cm.Name {
			err = errors.New("cluster tag exists")
			ginutils.Send(c, codes.InvalidArgs.Code(), err.Error(), nil)
			return
		}
		input.Tag = clusterInfo.Name
	}

	if cm.VolType == enums.VolTypeLowFrequency && input.ConsulAddr != "" {
		err := checkIpPort(c, input.ConsulAddr)
		if err != nil {
			return
		}
		set["consul_addr"] = input.ConsulAddr
	}

	handleStrBson(set, "idc", input.IDC)
	handleStrBson(set, "cli", input.Cli)
	handleStrBson(set, "domain", input.Domain)
	handleStrBson(set, "tag", input.Tag)
	handleStrBson(set, "s3_endpoint", input.S3Endpoint)

	if err := new(model.Cluster).Update(input.Id, set); err != nil {
		log.Errorf("update cluster failed. err:%+v", err)
		ginutils.Send(c, codes.DatabaseError.Code(), err.Error(), nil)
		return
	}
	ginutils.Send(c, codes.OK.Code(), codes.OK.Msg(), nil)
}

func handleStrBson(set bson.M, key, val string) {
	if val == "" {
		return
	}
	set[key] = val
}

func checkNameExists(c *gin.Context, name string) error {
	cm, err := new(model.Cluster).FindName(name)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("find cluster by name failed. err:%+v", err)
		ginutils.Send(c, codes.DatabaseError.Code(), err.Error(), nil)
		return err
	}
	if cm != nil && cm.Id != 0 {
		err = errors.New("cluster name exists")
		ginutils.Send(c, codes.InvalidArgs.Code(), err.Error(), nil)
		return err
	}
	return nil
}

func checkTagExists(c *gin.Context, tag string) error {
	cm, err := new(model.Cluster).FindTag(tag)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("by cluster by tag failed. err:%+v", err)
		ginutils.Send(c, codes.DatabaseError.Code(), err.Error(), nil)
		return err
	}
	if cm != nil && cm.Id != 0 {
		err = errors.New("cluster tag exists")
		ginutils.Send(c, codes.InvalidArgs.Code(), err.Error(), nil)
		return err
	}
	return nil
}

func create(input *CreateInput) (*model.Cluster, error) {
	m := &model.Cluster{}
	if err := copier.Copy(m, input); err != nil {
		log.Errorf("copy Cluster failed. err:%+v", err)
		return m, err
	}
	m.Status = enums.ClusterStatusOn
	if err := m.Create(); err != nil {
		log.Errorf("create cluster failed. err:%+v", err)
		return m, err
	}
	return m, nil
}

type ListInput struct {
	*model.FindClusterParam
}

func (input *ListInput) Check() error {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PerPage <= 0 {
		input.PerPage = 15
	}
	return nil
}

type ListOutput struct {
	Page     int       `json:"page"`
	PerPage  int       `json:"per_page"`
	Count    int64     `json:"count"`
	Clusters []Cluster `json:"clusters"`
}

func List(c *gin.Context) {
	input := &ListInput{&model.FindClusterParam{}}
	if !ginutils.Check(c, input) {
		return
	}
	clusters, count, err := new(model.Cluster).Find(input.FindClusterParam)
	if err != nil {
		log.Errorf("find clusters failed. input:%+v, err:%+v", input, err)
		ginutils.Fail(c, codes.DatabaseError.Code(), err.Error())
		return
	}
	output := ListOutput{
		Page:     input.Page,
		PerPage:  input.PerPage,
		Count:    count,
		Clusters: make([]Cluster, 0),
	}
	if len(clusters) == 0 {
		ginutils.Success(c, output)
		return
	}

	poolSize := pool.GetPoolSize(30, len(clusters))
	tp := pool.New(poolSize, poolSize)
	clusterChan := make(chan Cluster, len(clusters))
	handleClusters(c, clusters, tp, clusterChan)
	for cc := range clusterChan {
		output.Clusters = append(output.Clusters, cc)
	}

	ginutils.Success(c, output)
}

func handleClusters(c *gin.Context, clusters []model.Cluster, tp pool.TaskPool, clusterChan chan Cluster) {
	wg := sync.WaitGroup{}
	for i := range clusters {
		m := clusters[i]
		wg.Add(1)
		tp.Run(func() {
			defer wg.Done()
			if len(m.MasterAddr) == 0 || m.Status != enums.ClusterStatusOn {
				clusterChan <- Cluster{Cluster: m}
				return
			}
			clusterInfo, err := cluster.Get(c, m.MasterAddr[0])
			if err != nil {
				log.Errorf("get addr(%v) cluster failed: %v", m.MasterAddr[0], err)
			}
			if clusterInfo == nil {
				clusterChan <- Cluster{Cluster: m}
				return
			}
			data := Cluster{
				Cluster:             m,
				OriginalName:        clusterInfo.Name,
				LeaderAddr:          clusterInfo.LeaderAddr,
				MetaNodeSum:         len(clusterInfo.MetaNodes),
				MetaTotal:           helper.GBByteConversion(clusterInfo.MetaNodeStatInfo.TotalGB),
				MetaUsed:            helper.GBByteConversion(clusterInfo.MetaNodeStatInfo.UsedGB),
				MetaUsedRatio:       helper.Percentage(clusterInfo.MetaNodeStatInfo.UsedGB, clusterInfo.MetaNodeStatInfo.TotalGB),
				DataNodeSum:         len(clusterInfo.DataNodes),
				DataTotal:           helper.GBByteConversion(clusterInfo.DataNodeStatInfo.TotalGB),
				DataUsed:            helper.GBByteConversion(clusterInfo.DataNodeStatInfo.UsedGB),
				DataUsedRatio:       helper.Percentage(clusterInfo.DataNodeStatInfo.UsedGB, clusterInfo.DataNodeStatInfo.TotalGB),
				MetaNodes:           clusterInfo.MetaNodes,
				DataNodes:           clusterInfo.DataNodes,
				BadPartitionIDs:     clusterInfo.BadPartitionIDs,
				BadMetaPartitionIDs: clusterInfo.BadMetaPartitionIDs,
			}
			clusterChan <- data
		})
	}
	wg.Wait()
	close(clusterChan)
	tp.Close()
}

type ChangeStatusInput struct {
	ClusterId int64               `json:"cluster_id" binding:"required"`
	Status    enums.ClusterStatus `json:"status"`
}

func ChangeStatus(c *gin.Context) {
	input := &ChangeStatusInput{}
	if !ginutils.Check(c, input) {
		return
	}
	err := new(model.Cluster).Update(input.ClusterId, map[string]interface{}{
		"status": input.Status,
	})
	if err != nil {
		log.Errorf("change cluster status failed,args:%+v,err:%+v", input, err)
		ginutils.Fail(c, codes.DatabaseError.Code(), err.Error())
		return
	}
	ginutils.Success(c, nil)
}

type RemoveInput struct {
	ClusterId int64 `form:"cluster_id" binding:"required"`
}

func Remove(c *gin.Context) {
	input := &RemoveInput{}
	if err := c.BindQuery(input); err != nil {
		log.Errorf("parse param failed,err:%+v", err)
		ginutils.Fail(c, codes.InvalidArgs.Code(), err.Error())
		return
	}
	if err := new(model.Cluster).DeleteId(input.ClusterId); err != nil {
		log.Errorf("delete cluster failed.id:%d,err:%+v", input.ClusterId, err)
		ginutils.Fail(c, codes.DatabaseError.Code(), err.Error())
		return
	}
	ginutils.Success(c, nil)
}
