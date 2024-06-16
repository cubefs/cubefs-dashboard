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

package migrate

import (
	"github.com/cubefs/cubefs-dashboard/backend/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	migrations = append(migrations, &gormigrate.Migration{
		ID: "2024020500_cluster_add_fields",
		Migrate: func(tx *gorm.DB) error {
			if !tx.Migrator().HasColumn(new(model.Cluster), "region") {
				err := tx.Migrator().AddColumn(new(model.Cluster), "region")
				if err != nil {
					return err
				}
			}
			if !tx.Migrator().HasColumn(new(model.Cluster), "status") {
				err := tx.Migrator().AddColumn(new(model.Cluster), "status")
				if err != nil {
					return err
				}
			}
			return nil
		},
	}, &gormigrate.Migration{
		ID: "2024020501_operation_logs_add_field",
		Migrate: func(tx *gorm.DB) error {
			if tx.Migrator().HasColumn(new(model.OperationLog), "cluster_id") {
				return nil
			}
			return tx.Migrator().AddColumn(new(model.OperationLog), "cluster_id")
		},
	}, &gormigrate.Migration{
		ID: "2024020502_node_config_add_field",
		Migrate: func(tx *gorm.DB) error {
			if tx.Migrator().HasColumn(new(model.NodeConfig), "cluster_id") {
				return nil
			}
			return tx.Migrator().AddColumn(new(model.NodeConfig), "cluster_id")
		},
	}, &gormigrate.Migration{
		ID: "2024020503_node_config_failures_add_field",
		Migrate: func(tx *gorm.DB) error {
			if tx.Migrator().HasColumn(new(model.NodeConfigFailure), "cluster_id") {
				return nil
			}
			return tx.Migrator().AddColumn(new(model.NodeConfigFailure), "cluster_id")
		},
	}, &gormigrate.Migration{
		ID: "2024020503_vol_add_field",
		Migrate: func(tx *gorm.DB) error {
			if tx.Migrator().HasColumn(new(model.Vol), "cluster_id") {
				return nil
			}
			return tx.Migrator().AddColumn(new(model.Vol), "cluster_id")
		},
	}, &gormigrate.Migration{
		ID: "2024020503_user_add_field",
		Migrate: func(tx *gorm.DB) error {
			if tx.Migrator().HasColumn(new(model.User), "cluster_id") {
				return nil
			}
			return tx.Migrator().AddColumn(new(model.User), "cluster_id")
		},
	})
}
