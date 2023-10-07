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
  <div>
    <div v-auth="'CFS_DISKS_DECOMMISSION'" class="operate-wrap">
      <span>{{ $t('common.batch') }}</span>
      <el-select v-model="operateType" style="width: 100px;margin: 0 20px 0 10px;">
        <el-option :value="1" :label="$t('common.offline')"></el-option>
      </el-select>
      <el-button type="primary" @click="batchOperate">{{ $t('common.run') }}</el-button>
    </div>
    <u-page-table :data="dataList" :page-size="page.per_page" @selection-change="handleSelectionChange">
      <!-- <el-table-column label="序号" type="index"></el-table-column> -->
      <el-table-column type="selection" width="80"></el-table-column>
      <el-table-column :label="$t('resource.diskpath')" prop="path" sortable></el-table-column>
      <el-table-column :label="$t('common.partitions')" prop="partitions" sortable>
        <template slot-scope="scope">
          <a @click="toPath(scope.row.path)">{{ scope.row.partitions }}</a>
          <!-- <router-link
            :to="{
              name: 'dataNodePartiInfo',
              query: {
                diskPath: scope.row.path,
                ...$route.query,
              },
            }"
            tag="div"
            class="link"
          >{{ scope.row.partitions }}</router-link> -->
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('common.total') + $t('common.size')"
        prop="total"
        sortable
        :sort-method="sortMethodTotal"
      ></el-table-column>
      <el-table-column
        :label="$t('common.allocated')"
        prop="allocated"
        sortable
        :sort-method="sortMethodAvai"
      ></el-table-column>
      <el-table-column
        :label="$t('common.used')"
        prop="used"
        sortable
        :sort-method="sortMethodUsed"
      ></el-table-column>
      <el-table-column
        :label="$t('common.usage')"
        prop="usage_ratio"
        sortable
        :sort-method="sortMethodUsedRatio"
      ></el-table-column>
      <el-table-column :label="$t('common.status')" prop="status"></el-table-column>
      <el-table-column :label="$t('common.action')">
        <template slot-scope="scope">
          <MoreOPerate :count="2" :i18n="i18n">
            <el-button
              v-auth="'CFS_DISKS_DECOMMISSION'"
              size="medium"
              type="text"
              :disabled="!scope.row.partitions"
              @click="handleOffLine(scope.row)"
            >{{ $t('common.offline') }}</el-button>
          </MoreOPerate>
        </template>
      </el-table-column>
    </u-page-table>
  </div>
</template>
<script>
import { getDataNodeDiskList, offLineDisks } from '@/api/cfs/cluster'
import MoreOPerate from '@/pages/components/moreOPerate'
import UPageTable from '@/pages/components/uPageTable'
import { sortSizeWithUnit } from '@/utils'
import Mixin from '@/pages/cfs/clusterOverview/mixin'
export default {
  components: {
    MoreOPerate,
    UPageTable,
  },
  mixins: [Mixin],
  props: {
    curNode: {
      type: Object,
      default() {
        return {}
      },
    },
    path: {
      type: String,
      default() {
        return ''
      },
    },
  },
  data() {
    return {
      dataList: [],
      page: {
        per_page: 5, // 页面大小
      },
      selectedData: [],
      operateType: 1,
      i18n: this.$i18n,
    }
  },
  computed: {
    addr() {
      return this.curNode.addr || ''
    },
  },
  watch: {},
  created() {
    this.refresh()
  },
  methods: {
    sortMethodUsedRatio(a, b) {
      const ausage_ratio = parseFloat(a.usage_ratio)
      const busage_ratio = parseFloat(b.usage_ratio)
      if (ausage_ratio < busage_ratio) return -1
      if (ausage_ratio > busage_ratio) return 1
      return 0
    },
    sortMethodTotal(a, b) {
      return sortSizeWithUnit(a.total, b.total)
    },
    sortMethodAvai(a, b) {
      return sortSizeWithUnit(a.allocated, b.allocated)
    },
    sortMethodUsed(a, b) {
      return sortSizeWithUnit(a.used, b.used)
    },
    refresh() {
      this.getData()
    },
    async getData() {
      this.dataList = []
      const res = await getDataNodeDiskList({
        cluster_name: this.clusterName,
        addr: this.addr,
      })
      this.dataList = (res.data || []).sort((a, b) => {
        if (a.path < b.path) return -1
        if (a.path > b.path) return 1
        if ((a.path = b.path)) return 0
      })
    },
    async handleOffLine({ path }) {
      try {
        await this.$confirm( this.$t('resource.offlinenodedisk') + `(${this.addr} : ${path})`, this.$t('common.notice'), {
          confirmButtonText: this.$t('common.yes'),
          cancelButtonText: this.$t('common.no'),
          type: 'warning',
        })
        await offLineDisks({
          disks: [path],
          addr: this.addr,
          cluster_name: this.clusterName,
        })
        this.$message.success(this.$t('common.offline') + this.$t('common.xxsuc'))
        this.refresh()
      } catch (e) { }
    },
    toPath(path) {
      this.$emit('pathChange', path)
    },
    handleSelectionChange(val) {
      this.selectedData = val
    },
    async batchOperate() {
      if (!this.operateType) {
        this.$message.warning(this.$t('resource.chooseaction'))
        return
      }
      if (!this.selectedData.length) {
        this.$message.warning(this.$t('resource.least1'))
        return
      }
      try {
        await this.$confirm(this.$t('resource.confirmbatchoffline'), this.$t('common.notice'), { type: 'warning' })
        const disks = this.selectedData.map(item => item.path)
        await offLineDisks({
          cluster_name: this.clusterName,
          addr: this.addr || '',
          disks,
        })
        this.$message.success(this.$t('common.offline') + this.$t('common.xxsuc'))
        this.refresh()
      } catch (e) { }
    },
  },
}
</script>
<style lang="scss" scoped>
.operate-wrap {
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.noborder {
  border: none;
}

.input {
  width: 300px;
}

.search {
  position: relative;
  text-align: right;
}

.link {
  color: #66cc99;
  cursor: pointer;
}

.mr-l {
  color: #66cc99;
  cursor: pointer;
}
</style>
