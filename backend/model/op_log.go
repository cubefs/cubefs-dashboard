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

	"github.com/cubefs/cubefs-dashboard/backend/helper/types"
	"github.com/cubefs/cubefs-dashboard/backend/model/mysql"
)

type OperationLog struct {
	Id          uint64       `gorm:"primaryKey" json:"id"`
	Service     string       `gorm:"type:varchar(20);not null;default:''" json:"service"` // blobstore/cubefs
	Cluster     string       `gorm:"type:varchar(255);index;not null;default:''" json:"cluster"`
	ClusterId   int64        `gorm:"type:bigint(20);not null;default:0" json:"cluster_id"`
	UserId      int          `gorm:"index;not null;default:0" json:"user_id"`
	UserName    string       `gorm:"type:varchar(50);not null;default:''" json:"user_name"`
	OpTypeId    int          `gorm:"index;not null;default:0" json:"op_type_id"`
	OpTypeEN    string       `gorm:"type:varchar(100);not null;default:''" json:"op_type_en"`
	OpTypeCN    string       `gorm:"type:varchar(100);not null;default:''" json:"op_type_cn"`
	URI         string       `gorm:"type:varchar(200);not null;default:''" json:"uri"`
	Method      string       `gorm:"type:varchar(10);not null;default:''" json:"method"`
	QueryParams types.Values `gorm:"type:varchar(1024);not null;default:''" json:"query_params"`
	BodyParams  types.Map    `gorm:"type:varchar(1024);not null;default:''" json:"body_params"`
	Result      types.Map    `gorm:"type:varchar(500);not null;default:''" json:"result"`
	CreateTime  time.Time    `gorm:"primaryKey" json:"create_time"`
}

func (o *OperationLog) Create() error {
	return mysql.GetDB().Create(o).Error
}

type FindOpLogParam struct {
	Page     int `form:"page"`
	PerPage  int `form:"per_page"`
	OpTypeId int `form:"op_type_id"`
	UserId   int `form:"user_id"`
}

func (p *FindOpLogParam) Check() error {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PerPage <= 0 {
		p.PerPage = 10
	}
	return nil
}

func (o *OperationLog) Find(param *FindOpLogParam) ([]OperationLog, int64, error) {
	db := mysql.GetDB().Model(&OperationLog{})
	if param.OpTypeId > 0 {
		db = db.Where("op_type_id = ?", param.OpTypeId)
	}
	if param.UserId > 0 {
		db = db.Where("user_id = ?", param.UserId)
	}
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	oplogs := make([]OperationLog, 0)
	err = db.Scopes(mysql.Paginate(param.PerPage, param.Page)).Order("id desc").Find(&oplogs).Error
	return oplogs, count, err
}
