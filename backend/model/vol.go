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

	"github.com/cubefs/cubefs-dashboard/backend/model/mysql"
)

type Vol struct {
	Id              uint64    `gorm:"primaryKey" json:"id"`
	Name            string    `gorm:"type:varchar(100);not null;default:'';index" json:"name"`
	Owner           string    `gorm:"type:varchar(50);not null;default:'';index" json:"owner"`
	Capacity        uint64    `gorm:"not null;default:0" json:"capacity"`
	CacheCap        int       `gorm:"not null;default:0" json:"cache_cap"`
	CrossZone       bool      `gorm:"not null;default:0" json:"cross_zone"`
	Business        string    `gorm:"type:varchar(200);not null;default:''" json:"business"`
	DefaultPriority bool      `gorm:"not null;default:0" json:"default_priority"`
	ReplicaNumber   int       `gorm:"type:tinyint(4);not null;default:0" json:"replica_number"`
	VolType         int       `gorm:"type:tinyint(1);not null;default:0" json:"vol_type"`
	CreatorId       int       `gorm:"not null;default:0;index" json:"creator_id"`
	CreateTime      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP(3)" json:"create_time"`
}

func (v *Vol) Create() error {
	return mysql.GetDB().Create(v).Error
}

type FindVolsParam struct {
	Owner   string `form:"owner" binding:"required"`
	Page    int    `form:"page"`
	PerPage int    `form:"per_page"`
}

func (p *FindVolsParam) Check() error {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PerPage <= 0 {
		p.PerPage = 10
	}
	return nil
}

func (v *Vol) Find(param *FindVolsParam) ([]Vol, int64, error) {
	db := mysql.GetDB().Model(&Vol{})
	if param.Owner != "" {
		db = db.Where("owner = ?", param.Owner)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	vols := make([]Vol, 0)
	err := db.Scopes(mysql.Paginate(param.PerPage, param.Page)).Find(&vols).Error
	return vols, count, err
}
