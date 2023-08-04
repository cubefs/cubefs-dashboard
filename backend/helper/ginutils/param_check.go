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
		Send(c, codes.InvalidArgs.Code(), err.Error(), nil)
		return false
	}

	checker, ok := v.(Checker)
	if !ok {
		return true
	}

	if err := checker.Check(); err != nil {
		log.Errorf("param check failed. err:%+v", err)
		Send(c, codes.InvalidArgs.Code(), err.Error(), nil)
		return false
	}

	return true
}

const Cluster = "cluster"

func GetClusterMaster(c *gin.Context) (string, error) {
	name := c.Param(Cluster)
	cluster, err := new(model.Cluster).FindName(name)
	if err != nil {
		err = ClusterMasterAddrErr(name, err)
		log.Errorf("cluster.FindName failed.name:%s,err:%+v", name, err)
		Send(c, codes.DatabaseError.Error(), err.Error(), nil)
		return "", err
	}
	if len(cluster.MasterAddr) == 0 {
		err = ClusterMasterAddrErr(name, errors.New("no master addr"))
		log.Errorf("cluster vol_type error or no master_addr.master_addr:%+v,vol_type:%+v,name:%s", cluster.MasterAddr, cluster.VolType, name)
		Send(c, codes.DatabaseError.Error(), err.Error(), nil)
		return "", err
	}
	return cluster.MasterAddr[0], nil
}

func GetConsulAddr(c *gin.Context) (string, error) {
	name := c.Param(Cluster)
	cluster, err := new(model.Cluster).FindName(name)
	if err != nil {
		err = ClusterConsulAddrErr(name, err)
		log.Errorf("cluster.FindName failed. cluster:%s,err:%+v", name, err)
		Send(c, codes.DatabaseError.Code(), err.Error(), nil)
		return "", err
	}
	if len(cluster.ConsulAddr) == 0 || cluster.VolType != enums.VolTypeLowFrequency {
		err = ClusterConsulAddrErr(name, errors.New("no consul addr"))
		log.Errorf("cluster vol_type error or no consulAddr:%+v,vol_type:%+v,name:%s", cluster.ConsulAddr, cluster.VolType, c.Param(Cluster))
		Send(c, codes.DatabaseError.Error(), err.Error(), nil)
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