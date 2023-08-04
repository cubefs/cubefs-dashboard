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
	"time"

	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs/blobstore/util/log"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log.Infof(
			"[REQ_BEG] %s %s %15s",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
		)

		c.Next()
		log.Infof(
			"[REQ_END] %v, size: %d",
			time.Now().Sub(start),
			c.Writer.Size(),
		)
	}
}
