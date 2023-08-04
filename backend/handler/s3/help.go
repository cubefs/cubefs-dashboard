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
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

const DELIMITER = "/"

type FileV2 struct {
	Name     string             `json:"name"`
	Hash     *string            `json:"hash"`
	Path     string             `json:"path"`
	MimeType *string            `json:"mime_type"`
	Endpoint string             `json:"endpoint"`
	Bucket   string             `json:"bucket"`
	FileSize *int64             `json:"file_size"`
	PutTime  int64              `json:"put_time"`
	Meta     map[string]*string `json:"x-amz-meta,omitempty"`
}

func getContents(prefix, vol string, s3client *s3.S3, out *s3.ListObjectsV2Output) (prefixes []string, contents []*FileV2, directory *FileV2) {
	contents = make([]*FileV2, 0)
	prefixes = make([]string, 0)

	if out.CommonPrefixes != nil {
		for _, commonPrefix := range out.CommonPrefixes {
			prefixes = append(prefixes, *commonPrefix.Prefix)
		}
	}

	for _, content := range out.Contents {
		mimeType := "application/octet-stream"
		var meta map[string]*string

		headReq, headOut := s3client.HeadObjectRequest(&s3.HeadObjectInput{
			Bucket: aws.String(vol),
			Key:    content.Key,
		})

		if headReq.Send() == nil {
			mimeType = *headOut.ContentType
			meta = headOut.Metadata
		}

		fileOut := &FileV2{
			Name:     path.Base(*content.Key),
			FileSize: content.Size,
			Hash:     content.ETag,
			Path:     *content.Key,
			PutTime:  content.LastModified.Unix() * 10000000,
			MimeType: &mimeType,
			Endpoint: s3client.Endpoint,
			Bucket:   vol,
			Meta:     meta,
		}
		// directory detail
		if *content.Key == "" {
			continue
		}
		if *content.Key == prefix {
			directory = fileOut
			continue
		}
		contents = append(contents, fileOut)
	}
	return
}
