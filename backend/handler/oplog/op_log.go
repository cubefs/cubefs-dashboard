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

package oplog

import (
	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs/blobstore/util/log"

	"github.com/cubefs/cubefs-dashboard/backend/helper/codes"
	"github.com/cubefs/cubefs-dashboard/backend/helper/ginutils"
	"github.com/cubefs/cubefs-dashboard/backend/model"
)

type ListOutPut struct {
	Page    int                  `json:"page"`
	PerPage int                  `json:"per_page"`
	Count   int64                `json:"count"`
	Data    []model.OperationLog `json:"data"`
}

func List(c *gin.Context) {
	input := &model.FindOpLogParam{}
	if !ginutils.Check(c, input) {
		return
	}
	oplogs, count, err := new(model.OperationLog).Find(input)
	if err != nil {
		log.Errorf("oplog.Find failed. input:%+v, err:%+v", input, err)
		ginutils.Send(c, codes.DatabaseError.Code(), err.Error(), nil)
		return
	}
	output := &ListOutPut{
		Page:    input.Page,
		PerPage: input.PerPage,
		Count:   count,
		Data:    oplogs,
	}
	ginutils.Send(c, codes.OK.Code(), codes.OK.Msg(), output)
}
