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

import { getDataNodeList, getMetaNodeList } from '@/api/cfs/cluster'
import { toByte } from '@/utils'
export default {
  methods: {
    async getDataNodeList(params) {
      this.dataList = []
      if (this.clusterName) {
        const res = await getDataNodeList(params)
        const data = (res.data || []).sort((a, b) => {
          if (a.id < b.id) return -1
          if (a.id > b.id) return 1
          return 0
        })
        return data
      }
    },
    countDataNodeInfo(dataNodeListData) {
      const info = {
        node: 0,
        partition: 0,
        total: 0,
        used: 0,
      }
      dataNodeListData.forEach(item => {
        info.node += 1
        info.partition += item.partition_count
        const unitTotal = item.total.match(/[a-z|A-Z]+/gi)[0]
        const unitUsed = item.used.match(/[a-z|A-Z]+/gi)[0]
        info.total += toByte(Number(item.total.match(/\d+/gi)[0]), unitTotal)
        info.used += toByte(Number(item.used.match(/\d+/gi)[0]), unitUsed)
      })
      return info
    },
    async getMetaNodeList(params) {
      const res = await getMetaNodeList(params)
      const data = (res.data || []).sort((a, b) => {
        if (a.id < b.id) return -1
        if (a.id > b.id) return 1
        return 0
      })
      return data
    },
    countMetaNodeInfo(metaNodeList) {
      const info = {
        node: 0,
        partition: 0,
        total: 0,
        used: 0,
      }
      metaNodeList.forEach(item => {
        info.node += 1
        info.partition += item.partition_count
        const unitTotal = item.total.slice(-2)
        const unitUsed = item.used.slice(-2)
        info.total += toByte(Number(item.total.slice(0, -2)), unitTotal)
        info.used += toByte(Number(item.used.slice(0, -2)), unitUsed)
      })
      return info
    },
  },
}
