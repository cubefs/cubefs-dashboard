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

package config

type ListOutput struct {
	Configs Config `json:"configs"`
}

type Config struct {
	Balance           string `json:"balance"`
	BlobDelete        string `json:"blob_delete"`
	CodeMode          string `json:"code_mode"`
	DiskDrop          string `json:"disk_drop"`
	DiskRepair        string `json:"disk_repair"`
	DiskRepire        string `json:"disk_repire"`
	ShardRepair       string `json:"shard_repair"`
	VolInspect        string `json:"vol_inspect"`
	VolumeChunkSize   string `json:"volume_chunk_size"`
	VolumeInspect     string `json:"volume_inspect"`
	VolumeReserveSize string `json:"volume_reserve_size"`
}