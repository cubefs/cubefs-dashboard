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

type StrSlice []string

func (s StrSlice) Value() (driver.Value, error) {
	if s == nil {
		s = []string{}
	}
	b, err := json.Marshal(s)
	return string(b), err
}

func (s *StrSlice) Scan(v interface{}) error {
	if b, ok := v.([]byte); ok {
		return json.Unmarshal(b, s)
	}
	if str, ok := v.(string); ok {
		return json.Unmarshal([]byte(str), s)
	}
	return errors.New(fmt.Sprintf("invalid data type:%T", v))
}

func (s StrSlice) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}
