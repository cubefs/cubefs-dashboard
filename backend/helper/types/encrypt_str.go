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
	"errors"
	"fmt"

	"database/sql/driver"

	"github.com/cubefs/cubefs-dashboard/backend/helper/crypt"
)

type EncryptStr string

func (e EncryptStr) Value() (driver.Value, error)  {
	if e == "" {
		return nil, nil
	}
	return crypt.Encrypt([]byte(e))
}

func (e *EncryptStr) Scan(v interface{}) error {
	b, ok := v.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("invalid type:%T", v))
	}
	s, err := crypt.Decrypt(string(b))
	if err != nil {
		return err
	}
	*e = EncryptStr(s)
	return nil
}
