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

package node

import (
	"math"
	"strconv"
)

const (
	Active      = "Active"
	InActive    = "Inactive"
	UnAvailable = "Unavailable"
)

func FormatNodeStatus(isActive bool, badDisks ...string) string {
	if len(badDisks) > 0 {
		return UnAvailable
	}
	if isActive {
		return Active
	}
	return InActive
}

func FormatAvailable(available, used, total uint64) uint64 {
	if used > total {
		return  0
	}
	return available
}

func FormatWritableStr(isWritable bool) string {
	if isWritable {
		return "writable"
	}
	return "readonly"
}

func FormatDiskAndPartitionStatus(status int8) string {
	switch status {
	case 1:
		return "ReadOnly"
	case 2:
		return "ReadWrite"
	case -1:
		return "Unavailable"
	default:
		return "Unknown"
	}
}

func FormatUint64(num uint64) string {
	if num >= math.MaxInt64 {
		return "unlimited"
	}
	return strconv.FormatUint(num, 10)
}