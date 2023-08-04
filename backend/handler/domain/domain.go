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

package domain

import (
	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs/blobstore/util/log"

	"github.com/cubefs/cubefs-dashboard/backend/helper/codes"
	"github.com/cubefs/cubefs-dashboard/backend/helper/ginutils"
	"github.com/cubefs/cubefs-dashboard/backend/service/domain"
)

func Status(c *gin.Context) {
	addr, err := ginutils.GetClusterMaster(c)
	if err != nil {
		return
	}
	out, err := domain.Status(c, addr)
	if err != nil {
		log.Errorf("domain.Status failed. addr:%+v,err:%+v", addr, err)
		ginutils.Send(c, codes.ThirdPartyError.Code(), err.Error(), nil)
		return
	}
	ginutils.Send(c, codes.OK.Code(), codes.OK.Msg(), out)
}

func Info(c *gin.Context) {
	addr, err := ginutils.GetClusterMaster(c)
	if err != nil {
		return
	}
	out, err := domain.Info(c, addr)
	if err != nil {
		log.Errorf("domain.Info failed. addr:%+v,err:%+v", addr, err)
		ginutils.Send(c, codes.ThirdPartyError.Code(), err.Error(), nil)
		return
	}
	ginutils.Send(c, codes.OK.Code(), codes.OK.Msg(), out)
}
