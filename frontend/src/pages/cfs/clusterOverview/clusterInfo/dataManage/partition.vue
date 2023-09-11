<!--
 Copyright 2023 The CubeFS Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
 implied. See the License for the specific language governing
 permissions and limitations under the License.
-->

<template>
  <el-card class="container">
    <div class="darwerPosition">
      <FilterTableData
        ref="filterTableData"
        :data-list="originDataList"
        :types="['STATUS', 'RECOVER']"
        @filterData="filterData"
        style="flex: 0 0 auto; margin-right: 20px;"
      ></FilterTableData>
      <div class="search">
        <el-radio-group v-model="radio" @change="onChangeRadio">
          <el-radio label="vol">
            <span class="label">{{ $t('common.volumename') }}</span>
            <el-select
              v-model="params.name"
              filterable
              :placeholder="$t('resource.selectvol')"
              class="input"
              :disabled="radio !== 'vol'"
              @change="onSelectChange"
            >
              <el-option
                v-for="item in volNameList"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              >
              </el-option>
            </el-select>
          </el-radio>
          <el-radio label="id">
            <span class="label">{{ $t('common.partitionid') }}</span>
            <el-input
              v-model.trim="params.zoneId"
              :placeholder="$t('volume.inputparid')"
              clearable
              class="input"
              :disabled="radio === 'vol'"
              @keyup.native.enter="getData"
            ></el-input></el-radio>
        </el-radio-group>
        <el-button
          type="primary"
          class="search-btn"
          @click="onSearchClick"
        >{{ $t('button.search') }}</el-button>
        <el-button
          type="primary"
          class="search-btn"
          @click="onExportClick"
        >{{ $t('button.export') }}</el-button>
      </div>
    </div>
    <u-page-table :data="dataList" :page-size="page.per_page">
      <!-- <el-table-column label="序号" type="index"></el-table-column> -->
      <el-table-column
        :label="$t('common.partitionid')"
        prop="PartitionID"
        :width="120"
        sortable
      ></el-table-column>
      <el-table-column :label="$t('common.volumename')" prop="VolName" sortable></el-table-column>
      <el-table-column
        :label="$t('common.copies')"
        prop="ReplicaNum"
        sortable
        :width="100"
      ></el-table-column>
      <el-table-column
        v-if="isID"
        :label="$t('common.files')"
        prop="FileCount"
        sortable
      ></el-table-column>
      <el-table-column label="isRecovering" :width="150">
        <template slot-scope="scope">
          <span>{{ scope.row.IsRecover }}</span>
        </template>
      </el-table-column>
      <el-table-column label="Leader" prop="Leader"></el-table-column>
      <el-table-column label="Members">
        <template slot-scope="scope">
          <div v-for="item in scope.row.Members" :key="item">{{ item }}</div>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('common.status')"
        prop="status"
        :width="150"
      ></el-table-column>
      <el-table-column
        :label="$t('common.action')"
        :width="150"
      >
        <template slot-scope="scope">
          <MoreOPerate :count="2" :i18n="i18n">
            <!-- <el-button
              size="medium"
              type="text"
              @click="handleLoad(scope.row)"
            >load</el-button> -->
            <el-button
              size="medium"
              type="text"
              @click="handleDetail(scope.row)"
            >{{ $t('common.detail') }}</el-button>
          </MoreOPerate>
        </template>
      </el-table-column>
    </u-page-table>
    <Detail ref="detail" />
  </el-card>
</template>
<script>
import MoreOPerate from '@/pages/components/moreOPerate'
import UPageTable from '@/pages/components/uPageTable'
import FilterTableData from '@/pages/components/filter'
import Detail from './components/Detail'
import {
  getDataPartitionList,
  getVolList,
  loadDataPartition,
} from '@/api/cfs/cluster'
import Mixin from '@/pages/cfs/clusterOverview/mixin'
import NodeMixin from '@/pages/cfs/clusterOverview/clusterInfo/mixin'
export default {
  components: {
    MoreOPerate,
    UPageTable,
    FilterTableData,
    Detail,
  },
  mixins: [Mixin, NodeMixin],
  props: {
    showType: {
      type: String,
      // 对象或数组默认值必须从一个工厂函数获取
      default: function () {
        return 'normalPosition'
      },
    },
    curVol: {
      type: Object,
      default: function () {
        return {}
      },
    },
  },
  data() {
    return {
      radio: 'vol',
      isID: false,
      checkList: [],
      volNameList: [],
      dataList: [],
      originDataList: [],
      params: {
        name: '',
        zoneId: '',
      }, // 输入查询
      page: {
        per_page: 5, // 页面大小
      },
      i18n: this.$i18n,
    }
  },
  computed: {},
  watch: {},
  async mounted() {
    await this.getVolList()
    this.initVolName()
    this.getData()
    this.getNodeData()
  },
  methods: {
    async handleDetail({ PartitionID }) {
      const res = await getDataPartitionList({
        id: PartitionID,
        cluster_name: this.clusterName,
      })
      this.$refs.detail.init(res?.data?.[0] || {})
    },
    onSearchClick() {
      this.$refs.filterTableData.clear()
      this.getData()
    },
    refresh() {
      // this.resetParam()
      this.getData()
    },
    resetParam() {
      this.params = {
        name: '',
        zoneId: '',
      }
    },
    onChangeRadio(v) {
      if (v === 'vol') {
        this.params.name = this.curVol.name || this.volNameList[0]?.value
        this.params.zoneId = ''
        this.getData()
      } else {
        this.params.name = ''
        // this.getData()
        // this.dataList = []
        // this.originDataList = []
      }
    },
    async getVolList() {
      const volList = await getVolList({
        keywords: '',
        cluster_name: this.clusterName,
      })
      this.volNameList = volList.data.map((item) => {
        return { label: item.name, value: item.name }
      })
    },
    async getNodeData() {
      const data = await this.getDataNodeList({
        cluster_name: this.clusterName,
      })
      const info = this.countDataNodeInfo(data)
      this.$emit('update:info', info)
    },
    async getData() {
      this.dataList = []
      const { name, zoneId } = this.params
      const flag = this.radio === 'vol'
      this.isID = !flag && zoneId
      const res = await getDataPartitionList({
        vol_name: flag ? name : '',
        id: !flag ? zoneId : '',
        cluster_name: this.clusterName,
      })
      const temp = (res.data || [])
        .map((item) => {
          return {
            ...item,
            status: item.Status,
          }
        })
        .sort((a, b) => {
          if (a.PartitionID < b.PartitionID) return -1
          if (a.PartitionID > b.PartitionID) return 1
          return 0
        })
      this.dataList = temp
      this.originDataList = [...temp]
    },
    onExportClick() {
      const XLSX = require('xlsx')
      let exportData = ''
      const zhexportData = this.dataList.map(item => ({
        分区ID: item.PartitionID,
        卷名: item.VolName,
        Start: item.Start,
        End: item.End,
        DentryCount: item.DentryCount,
        InodeCount: item.InodeCount,
        MaxInodeID: item.MaxInodeID,
        isRecovering: item.IsRecover,
        Leader: item.Leader,
        Members: item.Members.join(),
        状态: item.Status,
      }))
      const enexportData = this.dataList.map(item => ({
        PartitionID: item.PartitionID,
        VolName: item.VolName,
        Start: item.Start,
        End: item.End,
        DentryCount: item.DentryCount,
        InodeCount: item.InodeCount,
        MaxInodeID: item.MaxInodeID,
        isRecovering: item.IsRecover,
        Leader: item.Leader,
        Members: item.Members.join(),
        Status: item.Status,
      }))
      if (localStorage.getItem('language') == 'zh') {
        exportData = zhexportData
      } else {
        exportData = enexportData
      }
      const wb = XLSX.utils.book_new()
      const ws = XLSX.utils.json_to_sheet(exportData)
      XLSX.utils.book_append_sheet(wb, ws, 'Sheet1')
      XLSX.writeFile(wb, 'data.xlsx')
    },
    initVolName() {
      const volName = this.$route.query.volName || this.curVol.name
      this.params.name = volName || this.volNameList[0]?.value || ''
    },
    onSelectChange() {
      this.getData()
    },
    filterData(data) {
      this.dataList = [...data]
    },
    async handleLoad({ PartitionID }) {
      await loadDataPartition({
        id: PartitionID,
        cluster_name: this.clusterName,
      })
      this.$message.success(this.$t('common.success'))
    },
  },
}
</script>
<style lang="scss" scoped>
.noborder {
  border: none;
}
.input {
  width: 150px;
}
.search {
  position: relative;
  text-align: right;
}
.search-btn {
  margin-left: 10px;
}
.normalPosition{
  display: flex;
  justify-content: space-between;
}
.container{
  position: relative;
}
.darwerPosition{
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  justify-content: flex-end;
}
</style>
