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

package s3

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cubefs/cubefs/blobstore/util/log"
	"github.com/gin-gonic/gin"

	"github.com/cubefs/cubefs-dashboard/backend/helper/codes"
	"github.com/cubefs/cubefs-dashboard/backend/helper/ginutils"
	"github.com/cubefs/cubefs-dashboard/backend/model"
	"github.com/cubefs/cubefs-dashboard/backend/service/user"
	"github.com/cubefs/cubefs-dashboard/backend/service/vol"
)

func GetS3Client(c *gin.Context, owner, volName string) (*s3.S3, bool) {
	cluster, err := ginutils.GetCluster(c)
	if err != nil {
		log.Errorf("get cluster by id failed.cluster:%s,err:%+v", c.Param(ginutils.Cluster), err)
		return nil, false
	}
	if len(cluster.MasterAddr) == 0 {
		log.Errorf("no master addr. cluster:%s", cluster.Name)
		ginutils.Send(c, codes.NotFound.Code(), "no master addr", nil)
		return nil, false
	}
	if cluster.Tag == "" {
		log.Errorf("no tag. cluster:%s", cluster.Name)
		ginutils.Send(c, codes.NotFound.Code(), "no region", nil)
		return nil, false
	}
	if cluster.S3Endpoint == "" {
		log.Errorf("no s3_endpoint. cluster:%s", cluster.Name)
		ginutils.Send(c, codes.NotFound.Code(), "no s3_endpoint", nil)
		return nil, false
	}
	addr := cluster.MasterAddr[0]
	if owner == "" {
		volInfo, err := vol.GetByName(c, addr, volName)
		if err != nil {
			log.Errorf("vol.ClientGet failed. addr:%s,owner:%s,vol:%s,err:%+v", addr, owner, volName, err)
			ginutils.Send(c, codes.ThirdPartyError.Code(), err.Error(), nil)
			return nil, false
		}
		owner = volInfo.Owner
	}
	ownerUser, err := user.Info(c, addr, owner)
	return NewClient(cluster, ownerUser), true
}

func NewClient(cluster *model.Cluster, owner *user.InfoOutput) *s3.S3 {
	conf := aws.Config{
		Region:                         aws.String(cluster.Tag),
		Endpoint:                       aws.String(cluster.S3Endpoint),
		S3ForcePathStyle:               aws.Bool(true),
		DisableRestProtocolURICleaning: aws.Bool(true),
		Credentials:                    credentials.NewStaticCredentials(owner.AccessKey, owner.SecretKey, ""),
	}
	return s3.New(session.Must(session.NewSessionWithOptions(session.Options{Config: conf})))
}

func GetExpire(expire int64) time.Duration {
	return time.Duration(expire) * time.Minute
}
