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

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs-dashboard/backend/config"
	"github.com/cubefs/cubefs-dashboard/backend/handler/auth"
)

type authRouter struct{}

func (r *authRouter) Register(engine *gin.Engine) {
	group := engine.Group(config.Conf.Prefix.Api + "/console/auth")

	group.POST("/login", auth.LoginHandler)
	group.POST("/logout", auth.LogoutHandler)

	user := group.Group("/user")
	{
		user.GET("/list", auth.GetUserHandler)
		user.POST("/create", auth.CreateUserHandler)
		user.PUT("/update", auth.UpdateUserHandler)
		user.PUT("/self/update", auth.UpdateSelfUserHandler)
		user.DELETE("/delete", auth.DeleteUserHandler)
		user.PUT("/password/update", auth.UpdateUserPasswordHandler)
		user.PUT("/password/self/update/", auth.UpdateSelfUserPasswordHandler)
		user.GET("/permission", auth.GetUserPermissionHandler)
	}
	role := group.Group("/role")
	{
		role.GET("/list", auth.GetRoleHandler)
		role.POST("/create", auth.CreateRoleHandler)
		role.PUT("/update", auth.UpdateRoleHandler)
		role.DELETE("/delete", auth.DeleteRoleHandler)
	}
	permission := group.Group("/permission")
	{
		permission.GET("/list", auth.GetPermissionHandler)
		permission.POST("/create", auth.CreatePermissionHandler)
		permission.PUT("/update", auth.UpdatePermissionHandler)
		permission.DELETE("/delete", auth.DeletePermissionHandler)
	}
}
