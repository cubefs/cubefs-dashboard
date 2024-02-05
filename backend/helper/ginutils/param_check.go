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

package ginutils

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs/blobstore/util/log"

	"github.com/cubefs/cubefs-dashboard/backend/helper/codes"
	"github.com/cubefs/cubefs-dashboard/backend/helper/enums"
	"github.com/cubefs/cubefs-dashboard/backend/model"
)

type Checker interface {
	Check() error
}

func Check(c *gin.Context, v interface{}) bool {
	if v == nil {
		return true
	}
	if err := c.ShouldBind(v); err != nil {
		log.Errorf("parse param error:%+v", err)
		Fail(c, codes.InvalidArgs.Code(), err.Error())
		return false
	}

	checker, ok := v.(Checker)
	if !ok {
		return true
	}

	if err := checker.Check(); err != nil {
		log.Errorf("param check failed. err:%+v", err)
		Fail(c, codes.InvalidArgs.Code(), err.Error())
		return false
	}

	return true
}

const (
	Cluster     = "cluster" // Cluster is the value of cluster.id
	ClusterName = "cluster_name"
	Region      = "region"
)

func GetClusterMaster(c *gin.Context) (string, error) {
	cluster, err := GetCluster(c)
	if err != nil {
		return "", err
	}
	if len(cluster.MasterAddr) == 0 {
		err = ClusterMasterAddrErr(cluster.Name, errors.New("no master addr"))
		log.Errorf("cluster vol_type error or no master_addr.master_addr:%+v,vol_type:%+v,name:%s", cluster.MasterAddr, cluster.VolType, cluster.Name)
		Fail(c, codes.DatabaseError.Error(), err.Error())
		return "", err
	}
	return cluster.MasterAddr[0], nil
}

func GetConsulAddr(c *gin.Context) (string, error) {
	cluster, err := GetCluster(c)
	if err != nil {
		return "", err
	}
	if len(cluster.ConsulAddr) == 0 {
		err = ClusterConsulAddrErr(cluster.Name, errors.New("no consul addr"))
		log.Errorf("GetConsulAddr failed.err:%+v", err)
		Fail(c, codes.DatabaseError.Error(), err.Error())
		return "", err
	}
	if cluster.VolType != enums.VolTypeLowFrequency {
		log.Errorf("cluster vol_type error vol_type:%+v,id:%d,name:%s", cluster.VolType, cluster.Id, cluster.Name)
		Fail(c, codes.DatabaseError.Error(), "vol_type error")
		return "", err
	}
	return cluster.ConsulAddr, nil
}

func CheckAndGetConsul(c *gin.Context, v interface{}) (string, error) {
	if !Check(c, v) {
		return "", errors.New("check param failed")
	}
	return GetConsulAddr(c)
}

func CheckAndGetMaster(c *gin.Context, v interface{}) (string, error) {
	if !Check(c, v) {
		return "", errors.New("check param failed")
	}
	return GetClusterMaster(c)
}

func GetCluster(c *gin.Context) (*model.Cluster, error) {
	idStr := c.Param(Cluster)
	id, err := ParseClusterId(idStr)
	if err != nil {
		log.Errorf("GetClusterId failed, clusterId:%s,err:%+v", idStr, err)
		Fail(c, codes.InvalidArgs.Code(), err.Error())
		return nil, err
	}
	cluster := new(model.Cluster)
	if err = cluster.FindId(id); err != nil {
		err = fmt.Errorf("get cluster(%d) error:%+v", id, err)
		log.Errorf("cluster.FindId failed.id:%s,err:%+v", idStr, err)
		Fail(c, codes.DatabaseError.Error(), err.Error())
		return nil, err
	}
	c.Set(Cluster, cluster.Id)
	c.Set(ClusterName, cluster.Name)
	c.Set(Region, cluster.Region)
	return cluster, nil
}

func ParseClusterId(idStr string) (int64, error) {
	if idStr == "" {
		return 0, errors.New("no cluster id")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse cluster id(%s) failed", idStr)
	}
	return id, nil
}
