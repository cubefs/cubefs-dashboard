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

type User struct {
	Id         uint64           `gorm:"primaryKey" json:"id"`
	ClusterId  int64            `gorm:"type:bigint(20);not null;default:0" json:"cluster_id"`
	Name       string           `gorm:"type:varchar(50);not null;default:'';uniqueIndex" json:"name"`
	Role       int              `gorm:"type:tinyint(4);not null;default:3" json:"role"`
	AccessKey  types.EncryptStr `gorm:"type:varchar(500);not null;default:''" json:"access_key"`
	SecretKey  types.EncryptStr `gorm:"type:varchar(500);not null;default:''" json:"secret_key"`
	CreatorId  int              `gorm:"not null;default:0;index" json:"creator_id"`
	CreateTime time.Time        `gorm:"not null;default:CURRENT_TIMESTAMP(3)" json:"create_time"`
}

func (u *User) Create() error {
	return mysql.GetDB().Create(u).Error
}
