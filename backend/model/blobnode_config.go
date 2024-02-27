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

package model

import (
	"time"

	"gorm.io/gorm"

	"github.com/cubefs/cubefs-dashboard/backend/helper/types"
	"github.com/cubefs/cubefs-dashboard/backend/model/mysql"
)

type NodeConfig struct {
	Id            int64        `gorm:"primaryKey" json:"id"`
	Node          string       `gorm:"type:varchar(255);index;not null;default:''" json:"node"`
	ClusterId     int64        `gorm:"type:bigint(20);not null;default:0" json:"cluster_id"`
	Configuration types.MapStr `gorm:"type:varchar(2048);not null;default:'{}'" json:"configuration"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

func (e *NodeConfig) Upsert(node, key, value string, clusterId int64) error {
	db := mysql.GetDB()
	nodeConf := &NodeConfig{Node: node, ClusterId: clusterId, Configuration: map[string]string{}}
	err := db.Where("node = ? and cluster_id = ? ", node, clusterId).First(nodeConf).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err != nil && err == gorm.ErrRecordNotFound {
		nodeConf.Configuration = map[string]string{}
	}
	nodeConf.Configuration[key] = value
	return db.Save(nodeConf).Error
}

func (e *NodeConfig) One(node string, clusterId int64) error {
	selector := map[string]interface{}{
		"node":       node,
		"cluster_id": clusterId,
	}
	return mysql.GetDB().Where(selector).Last(e).Error
}

type NodeConfigFailure struct {
	Id           int64     `gorm:"primaryKey" json:"id"`
	Node         string    `gorm:"type:varchar(255);not null;default:''" json:"node"`
	ClusterId    int64     `gorm:"type:bigint(20);not null;default:0" json:"cluster_id"`
	Key          string    `gorm:"type:varchar(255);not null;default:''" json:"key"`
	Value        string    `gorm:"type:varchar(255);not null;default:''" json:"value"`
	FailedReason string    `gorm:"type:varchar(255);not null;default:''"  json:"failed_reason"`
	CreatedAt    time.Time `json:"created_at"`
}

func (e *NodeConfigFailure) Insert() error {
	return mysql.GetDB().Create(e).Error
}

func (e *NodeConfigFailure) One(node string, clusterId int64) error {
	selector := map[string]interface{}{
		"node":       node,
		"cluster_id": clusterId,
	}
	return mysql.GetDB().Where(selector).Last(e).Error
}
