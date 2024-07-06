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

package vol

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs/proto"

	"github.com/cubefs/cubefs-dashboard/backend/config"
	"github.com/cubefs/cubefs-dashboard/backend/helper"
	"github.com/cubefs/cubefs-dashboard/backend/helper/crypt"
	"github.com/cubefs/cubefs-dashboard/backend/helper/httputils"
)

type CreateInput struct {
	Name            string `json:"name"`
	Owner           string `json:"owner"`
	Capacity        uint64 `json:"capacity"`
	Description     string `json:"description"`
	VolType         int    `json:"volType"`
	CrossZone       bool   `json:"crossZone"`
	DefaultPriority bool   `json:"defaultPriority"`
	ReplicaNumber   int    `json:"replicaNum"`
	CacheCap        int    `json:"cacheCap"`
	FollowerRead    bool   `json:"followerRead"`
}

func (c *CreateInput) QueryParams() string {
	if c == nil {
		return ""
	}
	c.DefaultPriority = !c.DefaultPriority
	return helper.BuildUrlParams(c)
}

func Create(c *gin.Context, clusterAddr string, input *CreateInput) (interface{}, error) {
	reqUrl := "http://" + clusterAddr + proto.AdminCreateVol + "?" + input.QueryParams() + "&clientIDKey=" + config.Conf.Server.ClientIDKey
	resp, err := httputils.DoRequestNoCookie(c, reqUrl, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	output := httputils.Output{}
	_, err = httputils.HandleResponse(c, resp, err, &output)
	if err != nil {
		return nil, err
	}
	if output.Code != proto.ErrCodeSuccess {
		return nil, errors.New(output.Msg)
	}
	return output.Data, nil
}

type UpdateInput struct {
	Name           string `json:"name"`
	AuthKey        string `json:"authKey"`
	CacheCap       uint64 `json:"cacheCap"`
	CacheThreshold int    `json:"cacheThreshold"`
	CacheTTL       int    `json:"cacheTTL"`
	ReplicaNumber  *uint8 `json:"replicaNum,omitempty"`
	FollowerRead   *bool  `json:"followerRead,omitempty"`
}

func Update(c *gin.Context, clusterAddr string, input *UpdateInput) (interface{}, error) {
	reqUrl := "http://" + clusterAddr + proto.AdminUpdateVol + "?" + helper.BuildUrlParams(input) + "&clientIDKey=" + config.Conf.Server.ClientIDKey
	resp, err := httputils.DoRequestNoCookie(c, reqUrl, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	output := httputils.Output{}
	_, err = httputils.HandleResponse(c, resp, err, &output)
	if err != nil {
		return nil, err
	}
	if output.Code != proto.ErrCodeSuccess {
		return nil, errors.New(output.Msg)
	}
	return output.Data, nil
}

type ExpandInput struct {
	Name     string `json:"name"`
	Capacity uint64 `json:"capacity"`
	AuthKey  string `json:"authKey"`
}

func Expand(c *gin.Context, clusterAddr string, input *ExpandInput) error {
	reqUrl := "http://" + clusterAddr + proto.AdminVolExpand + "?" + helper.BuildUrlParams(input) + "&clientIDKey=" + config.Conf.Server.ClientIDKey
	resp, err := httputils.DoRequestNoCookie(c, reqUrl, http.MethodGet, nil, nil)
	if err != nil {
		return err
	}
	output := httputils.Output{}
	_, err = httputils.HandleResponse(c, resp, err, &output)
	if err != nil {
		return err
	}
	if output.Code != proto.ErrCodeSuccess {
		return errors.New(output.Msg)
	}
	return nil
}

func Shrink(c *gin.Context, clusterAddr string, input *ExpandInput) error {
	reqUrl := "http://" + clusterAddr + proto.AdminVolShrink + "?" + helper.BuildUrlParams(input) + "&clientIDKey=" + config.Conf.Server.ClientIDKey
	resp, err := httputils.DoRequestNoCookie(c, reqUrl, http.MethodGet, nil, nil)
	if err != nil {
		return err
	}
	output := httputils.Output{}
	_, err = httputils.HandleResponse(c, resp, err, &output)
	if err != nil {
		return err
	}
	if output.Code != proto.ErrCodeSuccess {
		return errors.New(output.Msg)
	}
	return nil
}

func Get(c *gin.Context, clusterAddr, keywords string) (*[]proto.VolInfo, error) {
	reqUrl := "http://" + clusterAddr + proto.AdminListVols + "?keywords=" + keywords + "&clientIDKey=" + config.Conf.Server.ClientIDKey
	resp, err := httputils.DoRequestNoCookie(c, reqUrl, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	output := httputils.Output{Data: &[]proto.VolInfo{}}
	_, err = httputils.HandleResponse(c, resp, err, &output)
	if err != nil {
		return nil, err
	}
	if output.Code != proto.ErrCodeSuccess {
		return nil, errors.New(output.Msg)
	}
	return output.Data.(*[]proto.VolInfo), nil
}

func GetByName(c *gin.Context, clusterAddr, name string) (*proto.SimpleVolView, error) {
	reqUrl := "http://" + clusterAddr + proto.AdminGetVol + "?name=" + name
	resp, err := httputils.DoRequestNoCookie(c, reqUrl, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	output := httputils.Output{Data: &proto.SimpleVolView{}}
	_, err = httputils.HandleResponse(c, resp, err, &output)
	if err != nil {
		return nil, err
	}
	if output.Code != proto.ErrCodeSuccess {
		return nil, errors.New(output.Msg)
	}
	return output.Data.(*proto.SimpleVolView), nil
}

func ClientGet(c *gin.Context, clusterAddr, owner, name string) (*ClientVol, error) {
	authkey := crypt.Md5Encryption(owner)
	reqUrl := "http://" + clusterAddr + proto.ClientVol + "?name=" + name + "&authKey=" + authkey
	resp, err := httputils.DoRequestNoCookie(c, reqUrl, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	output := httputils.Output{Data: &ClientVol{}}
	_, err = httputils.HandleResponse(c, resp, err, &output)
	if err != nil {
		return nil, err
	}
	if output.Code != proto.ErrCodeSuccess {
		return nil, errors.New(output.Msg)
	}
	return output.Data.(*ClientVol), nil
}
