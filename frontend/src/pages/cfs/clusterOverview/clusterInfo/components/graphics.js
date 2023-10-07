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

import { getDataNodeList, getBadDataP, getBadMetaP, getMetaNodeList, getDataNodeDiskList, getDataPartitionList, getMetaPartitionList, batchOfflineDisk } from '@/api/cfs/cluster'
import { mapGetters } from 'vuex'
import { toByte, renderSize } from '@/utils'
export default {
  name: '',
  data() {
    return {
      chartData: [],
      mpDpchartData: [],
      dataPartition: 0,
      metaPartition: 0,
      badDataPartitionNum: 0,
      badMetaPartitionNum: 0,
      diskNum: 0,
      badDiskNum: 0,
      badDataPartition: [],
      badMetaPartition: [],
      DiskDialogVisible: false,
      DataPartitionDialogVisible: false,
      DataPartitionDetailDialogVisible: false,
      MetaPartitionDialogVisible: false,
      MetaPartitionDetailDialogVisible: false,
      PartitionTableData: [],
      MetaPartitionTableData: [],
      LackReplicaDataPartitionIDs: [],
      CorruptDataPartitionIDs: [],
      LackReplicaMetaPartitionIDs: [],
      CorruptMetaPartitionIDs: [],
      diskTableData: [
        {
          PartitionIDs: [
            7569,
            7567,
          ],
          Path: '10.52.136.208:17310:/home/service/var/data1',
        },
      ],
      seletedDisk: [],
    }
  },
  computed: {
    ...mapGetters('clusterInfoModule', {
      curClusterInfo: 'clusterInfog',
    }),
    tableColumns() {
      return [
        {
          title: 'Path',
          key: 'diskPath',
        },
        {
          title: 'id',
          key: 'id',
        },
      ]
    },
  },
  watch: {},
  async created() {
    this.getMpDpStatus()
  },
  beforeMount() { },
  mounted() { },
  methods: {
    mpDpChartValueFormat(value) {
      return renderSize(value)
    },
    newArrFn(arr) {
      const newArr = []
      for (let i = 0; i < arr.length; i++) {
        newArr.indexOf(arr[i]) === -1 ? newArr.push(arr[i]) : newArr
      }
      return newArr
    },
    async getMpDpStatus() {
      const res = await getBadDataP({ cluster_name: this.curClusterInfo.clusterName })
      this.LackReplicaDataPartitionIDs = res.data.LackReplicaDataPartitionIDs
      this.CorruptDataPartitionIDs = res.data.CorruptDataPartitionIDs
      this.badDataPartition = this.newArrFn(res.data.LackReplicaDataPartitionIDs.concat(res.data.CorruptDataPartitionIDs))
      this.badDataPartitionNum = this.badDataPartition.length
      this.diskTableData = res.data.BadDataPartitionIDs
      const res1 = await getBadMetaP({ cluster_name: this.curClusterInfo.clusterName })
      // this.LackReplicaMetaPartitionIDs = ['265', '266']
      this.LackReplicaMetaPartitionIDs = res1.data.LackReplicaMetaPartitionIDs
      this.CorruptMetaPartitionIDs = res1.data.CorruptMetaPartitionIDs
      this.badMetaPartition = this.newArrFn(res1.data.LackReplicaMetaPartitionIDs.concat(res1.data.CorruptMetaPartitionIDs))
      this.badMetaPartitionNum = this.badMetaPartition.length
      this.chartData.push({ data: [{ name: this.$t('common.health'), value: this.dataPartition - this.badDataPartitionNum }, { name: this.$t('common.broken'), value: this.badDataPartitionNum }], title: 'DP' + this.$t('common.status') })
      this.chartData.push({ data: [{ name: this.$t('common.health'), value: this.metaPartition - this.badMetaPartitionNum }, { name: this.$t('common.broken'), value: this.badMetaPartitionNum }], title: 'MP' + this.$t('common.status') })
    },
    calcMpDpData() {
      const { data_total: dataTotal, data_used: dataUsed, meta_total: metaTotal, meta_used: metaUsed } = this.curClusterInfo.clusterInfo

      const calcData = (total, used) => {
        const unitTotal = total.match(/[a-z|A-Z]+/gi)[0]
        const unitUsed = used.match(/[a-z|A-Z]+/gi)[0]
        const newTotal = toByte(Number(total.match(/\d+(\.\d+)?/gi)[0]), unitTotal)
        const newUsed = toByte(Number(used.match(/\d+(\.\d+)?/gi)[0]), unitUsed)
        return {
          totol: newTotal,
          used: newUsed,
          noUsed: newTotal - newUsed,
        }
      }
      const dataInfo = calcData(dataTotal, dataUsed)
      const metaInfo = calcData(metaTotal, metaUsed)
      this.mpDpchartData.push({ data: [{ name: this.$t('common.used'), value: dataInfo.used }, {name: this.$t('common.free'), value: dataInfo.noUsed}], title: this.$t('common.data') + this.$t('common.usage') })
      this.mpDpchartData.push({ data: [{ name: this.$t('common.used'), value: metaInfo.used }, {name: this.$t('common.free'), value: metaInfo.noUsed}], title: this.$t('common.meta') + this.$t('common.usage') })
    },
    showDialog(name) {
      const disk = this.$t('common.disk') + this.$t('common.status')
      const DP = 'DP' + this.$t('common.status')
      const MP = 'MP' + this.$t('common.status')
      if (name === disk) {
        this.DiskDialogVisible = true
      }
      if (name === DP) {
        this.DataPartitionDialogVisible = true
      }
      if (name === MP) {
        this.MetaPartitionDialogVisible = true
      }
    },
    async showDetail(id, type) {
      if (type === 2) {
        this.DataPartitionDetailDialogVisible = true
        const res = await getDataPartitionList({
          id,
          cluster_name: this.curClusterInfo.clusterName,
        })
        this.PartitionTableData = res.data
      } else {
        this.MetaPartitionDetailDialogVisible = true
        const res = await getMetaPartitionList({
          id,
          cluster_name: this.curClusterInfo.clusterName,
        })
        this.metaPartitionTableData = res.data
      }
    },
    handleSelectionDiskData(val) {
      this.seletedDisk = val
    },
    async batchOfflineDisk() {
      if (!this.seletedDisk.length) {
        this.$message.warning(this.$t('resource.choosedisk'))
      }
      // await batchOfflineDisk(params)
    },
  },
}
