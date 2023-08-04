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
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs-dashboard/backend/helper/enums"
)

type LoginUser struct {
	Id       int
	UserName string
	Email    string
	Phone    string
}

func GetLoginUser(c *gin.Context) (*LoginUser, error) {
	sessionData, exists := c.Get(enums.LoginUser)
	if !exists {
		return nil, errors.New("session not exists")
	}
	user, ok := sessionData.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid session data")
	}
	return &LoginUser{
		Id:       user["Id"].(int),
		UserName: user["UserName"].(string),
		Email:    user["Email"].(string),
		Phone:    user["Phone"].(string),
	}, nil
}
