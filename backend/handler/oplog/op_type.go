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
	"github.com/jinzhu/copier"

	"github.com/cubefs/cubefs/blobstore/util/log"

	"github.com/cubefs/cubefs-dashboard/backend/helper/codes"
	"github.com/cubefs/cubefs-dashboard/backend/helper/ginutils"
	"github.com/cubefs/cubefs-dashboard/backend/model"
)

type CreateOpTypeInput struct {
	NameEN string `json:"name_en" binding:"required"`
	NameCN string `json:"name_cn" binding:"required"`
	URI    string `json:"uri" binding:"required"`
	Method string `json:"method" binding:"required"`
	Record bool   `json:"record"`
}

func CreateOpType(c *gin.Context) {
	input := &CreateOpTypeInput{}
	if !ginutils.Check(c, input) {
		return
	}
	optype := new(model.OpType)
	_ = copier.Copy(optype, input)
	if err := optype.Create(); err != nil {
		log.Errorf("optype.Create failed. input:%+v, err:%+v", input, err)
		ginutils.Send(c, codes.DatabaseError.Code(), err.Error(), nil)
		return
	}
	ginutils.Send(c, codes.OK.Code(), codes.OK.Msg(), nil)
}

func UpdateOpType(c *gin.Context) {
	input := &model.OpTypeUpdateParam{}
	if !ginutils.Check(c, input) {
		return
	}
	optype := new(model.OpType)
	if err := optype.UpdateId(input); err != nil {
		log.Errorf("optype.update failed. input:%+v, err:%+v", input, err)
		ginutils.Send(c, codes.DatabaseError.Code(), err.Error(), nil)
		return
	}
	ginutils.Send(c, codes.OK.Code(), codes.OK.Msg(), nil)
}

type ListOpTypeInput struct {
	Record *bool `form:"record"`
}

func ListOpType(c *gin.Context) {
	input := &ListOpTypeInput{}
	if !ginutils.Check(c, input) {
		return
	}
	optypes, err := new(model.OpType).Find(input.Record)
	if err != nil {
		log.Errorf("optype.List failed. input:%+v, err:%+v", input, err)
		ginutils.Send(c, codes.DatabaseError.Code(), err.Error(), nil)
		return
	}
	ginutils.Send(c, codes.OK.Code(), codes.OK.Msg(), optypes)
}
