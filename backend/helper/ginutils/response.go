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
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs/blobstore/util/log"

	"github.com/cubefs/cubefs-dashboard/backend/helper/codes"
)

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"msg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Send(ctx *gin.Context, code int, msg string, data interface{}) {
	log.Infof("[REQ_CODE] %d", code)
	sendResponse(ctx, code, msg, data)
}

func Success(ctx *gin.Context, data interface{}) {
	Send(ctx, codes.OK.Code(), codes.OK.Msg(), data)
}

func sendResponse(ctx *gin.Context, code int, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, &Result{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}
