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

package mysql

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"github.com/cubefs/cubefs-dashboard/backend/config"
	"github.com/cubefs/cubefs-dashboard/backend/helper/enums"
)

var _db *gorm.DB

type DB struct {
	Error        error
	DB           *gorm.DB
	RowsAffected int64
}

func Init() error {
	username := config.Conf.Mysql.User
	password := config.Conf.Mysql.Password
	host := config.Conf.Mysql.Host
	port := config.Conf.Mysql.Port
	database := config.Conf.Mysql.Database
	slaveHost := config.Conf.Mysql.SlaveHost
	slavePort := config.Conf.Mysql.SlavePort

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	slaveDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, slaveHost, slavePort, database)

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		CreateBatchSize:                          800,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("connect mysql error: " + err.Error())
		return err
	}

	err = conn.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(dsn)},
		Replicas: []gorm.Dialector{mysql.Open(slaveDsn)},
		Policy:   dbresolver.RandomPolicy{},
	}))
	sqlDB, err := conn.DB()
	if err != nil {
		panic("connect mysql error: " + err.Error())
	}
	sqlDB.SetMaxIdleConns(config.Conf.Mysql.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.Conf.Mysql.MaxOpenConn)
	if enums.GetGinMode(config.Conf.Server.Mode) == gin.DebugMode {
		_db = conn.Debug()
	} else {
		_db = conn
	}
	return nil
}

func GetDB() *gorm.DB {
	return _db
}

func Paginate(limit, page int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if limit <= 0 {
			limit = 10
		}
		offset := (page - 1) * limit
		return db.Limit(limit).Offset(offset)
	}
}
