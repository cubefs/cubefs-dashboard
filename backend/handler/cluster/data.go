package cluster

import (
	"github.com/cubefs/cubefs/proto"

	"github.com/cubefs/cubefs-dashboard/backend/model"
)

type Cluster struct {
	model.Cluster       `json:",inline"`
	OriginalName        string                   `json:"original_name"`
	LeaderAddr          string                   `json:"leader_addr"`
	MetaNodeSum         int                      `json:"meta_node_sum"`
	MetaTotal           string                   `json:"meta_total"`
	MetaUsed            string                   `json:"meta_used"`
	MetaUsedRatio       string                   `json:"meta_used_ratio"`
	DataNodeSum         int                      `json:"data_node_sum"`
	DataTotal           string                   `json:"data_total"`
	DataUsed            string                   `json:"data_used"`
	DataUsedRatio       string                   `json:"data_used_ratio"`
	MetaNodes           []proto.NodeView         `json:"meta_nodes,omitempty" `
	DataNodes           []proto.NodeView         `json:"data_nodes,omitempty"`
	BadPartitionIDs     []proto.BadPartitionView `json:"bad_partition_ids"`
	BadMetaPartitionIDs []proto.BadPartitionView `json:"bad_meta_partition_ids"`
}
