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

package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type MapStr map[string]string

func (m MapStr) Value() (driver.Value, error) {
	if len(m) == 0 {
		m = map[string]string{}
	}
	b, e := json.Marshal(m)
	return string(b), e
}

func (m *MapStr) Scan(v interface{}) error {
	b, ok := v.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("invalid type:%T", v))
	}
	return json.Unmarshal(b, m)
}

type Map map[string]interface{}

func (m Map) Value() (driver.Value, error)  {
	if len(m) == 0 {
		m = map[string]interface{}{}
	}
	b, e := json.Marshal(m)
	return string(b), e
}

func (m *Map) Scan(v interface{}) error  {
	b, ok := v.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("invalid type:%T", v))
	}
	return json.Unmarshal(b, m)
}

type Values map[string][]string

func (val Values) Value() (driver.Value, error)  {
	if len(val) == 0 {
		val = map[string][]string{}
	}
	b, e := json.Marshal(val)
	return string(b), e
}

func (val *Values) Scan(v interface{}) error  {
	b, ok := v.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("invalid type:%T", v))
	}
	return json.Unmarshal(b, val)
}