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

package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs/blobstore/util/log"

	"github.com/cubefs/cubefs-dashboard/backend/helper/codes"
	"github.com/cubefs/cubefs-dashboard/backend/helper/ginutils"
	"github.com/cubefs/cubefs-dashboard/backend/helper/types"
	"github.com/cubefs/cubefs-dashboard/backend/model"
)

func RecordOpLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlPath := c.FullPath()
		method := c.Request.Method
		optype, err := new(model.OpType).FindByUniqKey(urlPath, method)
		if err != nil {
			c.Next()
			return
		}
		if !optype.Record {
			c.Next()
			return
		}

		user, err := ginutils.GetLoginUser(c)
		if err != nil {
			log.Errorf("ginutils.GetLoginUser failed.err:%+v", err)
			ginutils.Send(c, codes.Unauthorized.Code(), err.Error(), nil)
			c.Abort()
			return
		}

		var service string
		if strings.Contains(c.Request.RequestURI, "blobstore") {
			service = "blobstore"
		}
		if strings.Contains(c.Request.RequestURI, "cfs") {
			service = "cubefs"
		}

		var query map[string][]string
		query = c.Request.URL.Query()

		bodyParams := types.Map{}
		if c.Request.Body != nil {
			b, _ := ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			if err = json.Unmarshal(b, &bodyParams); err != nil {
				log.Errorf("parse body params failed.err:%+v", err)
			}
		}

		opLog := &model.OperationLog{
			Service:     service,
			Cluster:     c.Param("cluster"),
			UserId:      user.Id,
			UserName:    user.UserName,
			OpTypeId:    optype.Id,
			OpTypeEN:    optype.NameEN,
			OpTypeCN:    optype.NameCN,
			URI:         urlPath,
			Method:      method,
			QueryParams: query,
			BodyParams:  bodyParams,
			Result:      nil,
			CreateTime:  time.Now(),
		}

		bw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bw
		c.Next()
		c.Writer.Status()
		result := map[string]interface{}{}
		_ = json.Unmarshal(bw.body.Bytes(), &result)
		opLog.Result = types.Map{}
		if val, ok := result["code"]; ok {
			opLog.Result["code"] = val
		}
		if val, ok := result["msg"]; ok {
			opLog.Result["msg"] = val
		}
		if err := opLog.Create(); err != nil {
			log.Errorf("opLog.Create failed. opLog:%+v, err:%+v", opLog, err)
		}
	}
}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
