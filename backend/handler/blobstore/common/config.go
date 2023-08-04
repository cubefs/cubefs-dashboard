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

package common

import (
	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs/blobstore/common/codemode"
	"github.com/cubefs/cubefs-dashboard/backend/service/blobstore/clustermgr/config"
)

type CodeMode struct {
	ModeName codemode.CodeModeName `json:"mode_name"`
	CodeMode codemode.CodeMode     `json:"code_mode"`
	Enable   bool                  `json:"enable"`
	Tactic   codemode.Tactic       `json:"tactic"`
}

func GetCodeMode(c *gin.Context, consulAddr string) ([]CodeMode, error) {
	policies, err := config.Get(c, consulAddr, "code_mode")
	if err != nil {
		return nil, err
	}

	modes := make([]CodeMode, 0)
	for _, policy := range policies {
		item := CodeMode{
			ModeName: policy.ModeName,
			CodeMode: policy.ModeName.GetCodeMode(),
			Enable:   policy.Enable,
			Tactic:   policy.ModeName.Tactic(),
		}
		modes = append(modes, item)
	}
	return modes, nil
}
