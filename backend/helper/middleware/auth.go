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
	"errors"
	"strings"

	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs/blobstore/util/log"

	"github.com/cubefs/cubefs-dashboard/backend/helper/codes"
	"github.com/cubefs/cubefs-dashboard/backend/helper/ginutils"
	"github.com/cubefs/cubefs-dashboard/backend/model"
	"github.com/cubefs/cubefs-dashboard/backend/model/mysql"
	"github.com/cubefs/cubefs-dashboard/backend/service/auth"
)

const (
	KeyPairs    = "cubefs-key"
	SessionName = "cubefs-session"
)

func generateRouterKey(method, fullPath string) string {
	return strings.ToUpper(method) + strings.ReplaceAll(fullPath, "/", "_")
}

func InitSession() gin.HandlerFunc {
	db := mysql.GetDB()
	store := gormsessions.NewStore(db, true, []byte(KeyPairs))
	store.Options(sessions.Options{MaxAge: 8 * 3600, Path: "/"})
	return sessions.Sessions(SessionName, store)
}

func Authorization(c *gin.Context) {
	urlPath := c.FullPath()
	method := c.Request.Method
	db := mysql.GetDB()
	permissions := make([]model.AuthPermission, 0)
	query := db.Model(model.AuthPermission{}).Where("uri = ? and method = ?", urlPath, method).Find(&permissions)
	if query.Error != nil {
		log.Errorf("get permission err: %v", query.Error)
		c.Abort()
		return
	}
	if len(permissions) == 0 {
		log.Infof("%s: %s not require auth", method, urlPath)
		c.Next()
		return
	}

	permission := permissions[0]
	// check whether the request requires login
	if !permission.IsLogin {
		log.Infof("Users do not need to log in")
		c.Next()
		return
	}

	// check whether the request requires check auth
	if !permission.IsCheck {
		log.Infof("%s: %s not require auth", method, urlPath)
		c.Next()
		return
	}

	// check whether the user session is valid
	session := sessions.Default(c)
	sessionId, err := c.Cookie("sessionId")
	if err != nil {
		log.Errorf("get session err: %+v", err)
		ginutils.Send(c, codes.Forbidden.Code(), err.Error(), nil)
		c.Abort()
		return
	}
	sessionData := session.Get(sessionId)
	if sessionData == nil {
		log.Error("get session err")
		err = errors.New("get session err")
		ginutils.Send(c, codes.Forbidden.Code(), err.Error(), nil)
		c.Abort()
		return
	}

	userId := sessionData.(map[string]interface{})["Id"].(int)

	// check whether the user has permission
	userPermission, err := auth.GetUserPermission(userId, nil)
	if err != nil {
		log.Errorf("get the user auth fail")
		ginutils.Send(c, codes.Forbidden.Code(), err.Error(), nil)
		c.Abort()
		return
	}

	for _, item := range userPermission {
		if item.AuthCode == permission.AuthCode {
			c.Next()
			return
		}
	}
	log.Errorf("the user does not have this permission")
	ginutils.Send(c, codes.Forbidden.Code(), "the user does not have this permission", nil)
	c.Abort()
	return
}
