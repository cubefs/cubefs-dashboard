/**
 * Copyright 2023 The CubeFS Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
 * implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */


import { putPresignature, upload, completemultiFun } from './common'

export const getChunksFiles = ({ file }, { chunkSize }) => {
  const fileList = []
  if (!file.size) {
    const tmp = file.slice(0, file.size)
    const newfile = new window.File([tmp], file.name)
    fileList.push(newfile)
  } else {
    for (let i = 0; i < file.size; i = i + chunkSize) {
      const tmp = file.slice(i, Math.min((i + chunkSize), file.size))
      const newfile = new window.File([tmp], file.name)
      fileList.push(newfile)
    }
  }
  return fileList
}

export const putFile = (params, options) => {
  const presigned = new Promise((resolve, reject) => {
    putPresignature(params).then(async res => {
      const uploadUrl = options.getDataFun({ res, params: params })
      if (!uploadUrl.length) {
        throw new Error('不能从接口中取到预签名url，请检查接口或入参getDataFun方法是否正确')
      }
      const fileList = getChunksFiles(options, params)
      upload(uploadUrl, fileList, options.xhrOptions, params.prefix, params.chunkSize).then(async (data) => {
        const resData = [...data]
        if (resData.length > 1) {
          const uploadId = await options.completemultiData(res)
          await completemultiFun({ upload_id: uploadId, ...params })
          resolve(resData)
        }
        resolve(resData)
      }, err => reject(err))
    }, err => reject(err))
  })
  return presigned
}
