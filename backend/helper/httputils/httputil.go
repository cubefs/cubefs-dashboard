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

package httputils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs-dashboard/backend/helper/codes"
)

func HandleResponse(c *gin.Context, resp *http.Response, err error, data interface{}) (int, error) {
	if err != nil {
		return codes.ResultError.Code(), err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("failed :%s", string(s))
		return resp.StatusCode, err
	}
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("read body error :%v", err)
		return codes.ResultError.Code(), err
	}
	if len(s) == 0 {
		return codes.OK.Code(), nil
	}
	err = json.Unmarshal(s, data)
	if err != nil {
		err = fmt.Errorf("resp json decode error:%v (body:%s)", err, string(s))
		return codes.ResultError.Code(), err
	}
	return codes.OK.Code(), err
}

func DoRequestNoCookie(c *gin.Context, url, method string, body io.Reader, extraHeaders map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		err = fmt.Errorf("new http request error :%v, url:%s", err, url)
		return nil, err
	}
	header := c.Request.Header
	for k, v := range extraHeaders {
		header.Set(k, v)
	}

	req.Header = header
	client := &http.Client{
		Transport: http.DefaultTransport,
	}
	return client.Do(req)
}

// DoRequestBlobstore TODO blobstore token header
func DoRequestBlobstore(c *gin.Context, url, method string, body io.Reader, extraHeaders map[string]string) (*http.Response, error) {
	switch method {
	case http.MethodPost:
		if extraHeaders == nil {
			extraHeaders = map[string]string{}
		}
		extraHeaders["Content-Type"] = "application/json"
	}

	return DoRequestNoCookie(c, url, method, body, extraHeaders)
}

func DoRequestOvertime(c *gin.Context, url, method string, body io.Reader, extraHeaders map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		err = fmt.Errorf("new http request error :%v", err)
		return nil, err
	}
	header := c.Request.Header
	for k, v := range extraHeaders {
		header.Set(k, v)
	}
	req.Header = header
	client := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   time.Second,
	}
	return client.Do(req)
}
