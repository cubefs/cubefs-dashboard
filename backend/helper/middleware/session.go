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
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs/blobstore/util/log"

	"github.com/cubefs/cubefs-dashboard/backend/helper/enums"
)

func Session() gin.HandlerFunc  {
	return func(c *gin.Context) {
		// check whether the user session is valid
		session := sessions.Default(c)
		sessionId, err := c.Cookie("sessionId")
		if err != nil {
			log.Errorf("c.Cookie err: %+v", err)
			c.Next()
			return
		}
		sessionData := session.Get(sessionId)
		if sessionData == nil {
			log.Error("sessionData is nil")
			c.Next()
			return
		}
		c.Set(enums.LoginUser, sessionData)
		c.Next()
	}
}
