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
    <div class="fontTypeSpan">
      <span
        class="mg-lf-m"
      ><span>{{ $t('common.host') }}:</span><span class="mg-lf-m">{{ curHost.host }}</span></span>
      <span
        class="mg-lf-m"
      ><span>{{ $t('common.idc') }}:</span><span class="mg-lf-m">{{ curHost.idc }}</span></span>
      <span
        class="mg-lf-m"
      ><span>{{ $t('common.rack') }}:</span><span class="mg-lf-m">{{ curHost.rack }}</span></span>
      <span
        class="mg-lf-m"
      ><span>{{ $t('common.total') }}{{ $t('common.size') }}:</span><span class="mg-lf-m">{{ curHost.size | readablizeBytes }}</span></span>
      <span
        class="mg-lf-m"
      ><span>{{ $t('common.used') }}:</span><span class="mg-lf-m">{{ curHost.used| readablizeBytes }}</span></span>
      <span
        class="mg-lf-m"
      ><span>{{ $t('common.free') }}:</span><span class="mg-lf-m">{{ curHost.free | readablizeBytes }}</span></span>
    </div>
    <el-row class="">
      <el-button
        type="text"
        size="mini"
        class="back"
        @click="goBack"
      >{{ $t('component.goback') }}
      </el-button>
    </el-row>
    <el-row class="inside">
      <UTablePage :data="data" class="list-table" :has-page="false" sort-by-order default-should-sort-key="disk_id">
        <el-table-column
          prop="disk_id"
          :label="$t('common.disk') + 'id'"
          width="90"
          sortable="custom"
        >
        </el-table-column>
        <!-- <el-table-column prop="host" label="主机" sortable="custom"> </el-table-column>
        <el-table-column prop="idc" label="机房" sortable="custom"></el-table-column> -->
        <el-table-column prop="path" label="path" width="180"></el-table-column>
        <el-table-column :label="$t('resource.inwritingstripeid')" width="150" prop="warittingIdsCount" sortable="custom">
          <template slot-scope="scope">
            <span v-if="scope.row.volume_ids" class="color9">
              <span v-for="(link, index) in scope.row.volume_ids" :key="link">
                <a @click="toDataBlock(link)">{{ link }}</a>
                <span v-if="index !== scope.row.volume_ids.length - 1"> | </span>
              </span>
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="disk_load" :label="$t('resource.inwritingstripecnt')" width="150" sortable="custom"></el-table-column>
        <el-table-column prop="readonly" :label="$t('common.ro')" width="100" sortable="custom">
          <template slot-scope="scope">
            <span
              :class="[
                'rw',
                `rw-${scope.row.readonly ? 'readonly' : 'readwrite'}`,
              ]"
            >
              {{ scope.row.readonly ? $t('common.ro') : $t('common.rw') }}</span>
          </template></el-table-column>
        <el-table-column prop="free_chunk_cnt" :label="$t('common.free') + ' chunk'" width="110" sortable="custom"></el-table-column>
        <el-table-column prop="used_chunk_cnt" :label="$t('common.used') + ' chunk'" width="110" sortable="custom"></el-table-column>

        <el-table-column :label="$t('common.usage')" width="100" prop="percentage" sortable="custom">
          <template slot-scope="scope">
            <el-progress
              :percentage="scope.row.percentage"
              :color="[
                { color: '#f56c6c', percentage: 100 },
                { color: '#e6a23c', percentage: 80 },
                { color: '#5cb87a', percentage: 60 },
                { color: '#1989fa', percentage: 40 },
                { color: '#6f7ad3', percentage: 20 },
              ]"
            >
            </el-progress> </template></el-table-column>
        <el-table-column prop="status" :label="$t('common.disk') + $t('common.status')" width="100" sortable="custom">
          <template slot-scope="scope">
            <div
              class="text-center"
              :class="['ant-tag', `ant-tag-${filterStatus(scope.row.status)}`]"
            >
              {{ filterStatus(scope.row.status) }}</div><br>
            <DiskRepairStatus
              :status="scope.row.status"
              :is-off-line="!isShowoffLine(scope.row.disk_id)"
              :percentage="
                +(scope.row.repair_schedule * 100).toFixed(2)
              "
            >
            </DiskRepairStatus>
          </template>
        </el-table-column>
        <el-table-column prop="size" :label="$t('common.total') + $t('common.size')" width="100" sortable="custom">
          <template slot-scope="scope">
            {{ scope.row.size | readablizeBytes }}
          </template>
        </el-table-column>
        <el-table-column prop="used" :label="$t('common.used')" width="100" sortable="custom">
          <template slot-scope="scope">
            {{ scope.row.used | readablizeBytes }}
          </template></el-table-column>
        <el-table-column prop="free" :label="$t('common.free')" width="100" sortable="custom">
          <template slot-scope="scope">
            {{ scope.row.free | readablizeBytes }}
          </template>
        </el-table-column>
        <el-table-column
          fixed="right"
          prop="cluster"
          :label="$t('common.action')"
          :width="160"
          align="center"
        >
          <template slot-scope="{ row }">
            <el-button
              v-if="row.readonly"
              v-auth="'BLOBSTORE_DISKS_ACCESS'"
              type="text"
              size="mini"
              :disabled="row.status > 1 || !isShowoffLine(row.disk_id)"
              @click="changeRW(row)"
            >{{ $t('common.rw') }}</el-button>
            <el-button
              v-if="!row.readonly"
              v-auth="'BLOBSTORE_DISKS_ACCESS'"
              type="text"
              size="mini"
              :disabled="row.status > 1 || !isShowoffLine(row.disk_id)"
              @click="changeRW(row)"
            >{{ $t('common.ro') }}</el-button>
            <!-- <el-popover
              v-if="haveOrder(row, 'NODE_DROP', 'DISK_DROP')"
              placement="top-start"
              title="磁盘下线工单正在审批中"
              trigger="hover"
              :disabled="!haveOrder(row, 'NODE_DROP', 'DISK_DROP')"
            >
              <span><a style="color:#00c9c9; cursor: pointer;" class="link" @click.stop.prevent="goOrderDetail(haveOrder(row, 'NODE_DROP', 'DISK_DROP'))">工单详情</a></span>
              <el-button
                slot="reference"
                type="text"
                size="mini"
                class="dis-btn ml-8"
              >下线</el-button>
            </el-popover>
            <el-button
              v-if="isShowoffLine(row.disk_id) && !haveOrder(row, 'NODE_DROP', 'DISK_DROP') && row.status <= 1&& !row.readonly"
              type="text"
              class="dis-btn"
              size="mini"
              title="必须切只读"
            >下线</el-button>
            <el-button
              v-if="isShowoffLine(row.disk_id) && !haveOrder(row, 'NODE_DROP', 'DISK_DROP') && row.status <= 1&& row.readonly"
              type="text"
              size="mini"
              :disabled="haveOrder(row, 'DISK_SET')"
              @click="offLine(row)"
            >下线</el-button>
            <el-button
              v-if="!isShowoffLine(row.disk_id)"
              type="text"
              size="mini"
              class="waring-status"
              disabled
            >下线中</el-button> -->
            <el-popover
              v-if="haveOrder(row, 'DISK_SET')"
              placement="top-start"
              :title="$t('resource.setbrokenjobonfly')"
              trigger="hover"
              :disabled="!haveOrder(row, 'DISK_SET')"
            >
              <span><a style="color:#00c9c9; cursor: pointer;" class="link" @click.stop.prevent="goOrderDetail(haveOrder(row, 'DISK_SET'))">{{ $t('common.worksheet') }}{{ $t('common.detail') }}</a></span>
              <el-button
                slot="reference"
                v-auth="'BLOBSTORE_DISKS_SET'"
                type="text"
                size="mini"
                class="dis-btn ml-8"
              > {{ $t('resource.setbroken') }}</el-button>
            </el-popover>
            <el-button
              v-else
              v-auth="'BLOBSTORE_DISKS_SET'"
              type="text"
              size="mini"
              :disabled="row.status > 1 || !isShowoffLine(row.disk_id) || haveOrder(row, 'NODE_DROP', 'DISK_DROP')"
              @click="setBad(row)"
            >{{ $t('resource.setbroken') }}</el-button>
            <el-button
              v-auth="'BLOBSTORE_DISKS_PROBE'"
              type="text"
              size="mini"
              :disabled="+row.status == 1"
              @click="register(row)"
            >{{ $t('resource.diskregist') }}</el-button>
          </template>
        </el-table-column>
      </UTablePage>
    </el-row>
  </div>
</template>
<script>
import {
  getNodeList,
  offLineNodeInfo,
  offLineDropNodeList,
  changeNodeInfoRW,
  setNodeInfoBad,
  registerDisk,
} from '@/api/ebs/ebs'
import { readablizeBytes } from '@/utils'
import { nodeStatusMap } from '@/pages/cfs/status.conf'
import UTablePage from '@/components/uPageTable.vue'
import DiskRepairStatus from '@/components/diskRepairStatus'
import mixin from '@/pages/cfs/clusterOverview/mixin'
export default {
  components: {
    UTablePage,
    DiskRepairStatus,
  },
  filters: {
    readablizeBytes(value) {
      return readablizeBytes(value)
    },
  },
  mixins: [mixin],
  inject: ['app'],
  props: {
    detail: {
      type: Object,
      default() {
        return {}
      },
    },
  },
  data() {
    return {
      data: [],
      droppingList: [],
      curHost: {},
    }
  },
  computed: {
    user() {
      return this.$store.state.userInfo
    },
    clusterInfo() {
      return JSON.parse(sessionStorage.getItem('clusterInfo'))
    },
    showHost() {
      return this.$route.query.host
    },
  },
  created() {
    this.getData()
  },
  methods: {
    async register(row) {
      try {
        await this.$confirm(this.$t('resource.confirmregistdisk'), this.$t('common.notice'), {
          confirmButtonText: this.$t('common.yes'),
          cancelButtonText: this.$t('common.no'),
          type: 'warning',
        })
        await registerDisk({ region: this.clusterName, clusterId: this.app.clusterId, path: row.path, host: row.host, disk_id: row.disk_id })
        this.$message.success(this.$t('resource.registsuc'))
        await this.getData()
      } catch (e) {}
    },
    haveOrder(row, type, type2) {
      const temp = (row?.docking_info || []).find(i => i.OrderType === type)
      if (temp) {
        return temp?.OrderUrl
      }
      if (type2) {
        const temp2 = (row?.docking_info || []).find(i => i.OrderType === type2)
        return temp2?.OrderUrl
      }
      return false
    },
    resh() {
      this.refresh = 1
      this.getData()
    },
    filterStatus(v) {
      const temp = Object.entries(nodeStatusMap).filter(
        (temp) => temp[1] === v,
      )?.[0]
      return temp?.[0]
    },
    isShowoffLine(id) {
      return !(
        this.droppingList.filter((item) => {
          return item.disk_id === id
        })?.length || 0
      )
    },
    goBack() {
      // this.$router.replace({query:{...}})
      this.$router.go(-1)
    },
    getForms(forms) {
      this.forms = forms
      this.getData()
    },
    goOrderDetail(url) {
      window.open(url)
    },
    async getData() {
      this.data = []
      await this.getDroppingList()
      const host = this.$route.query.host
      const res = await getNodeList({
        region: this.clusterName,
        clusterId: this.app.clusterId,
        host,
      })
      const nodes = res.data.nodes
      // eslint-disable-next-line camelcase
      const { docking_info } = nodes[host]
      this.curHost = nodes[host].disks.reduce(
        (pre, next) => {
          return {
            ...next,
            host_name: nodes[host]?.host_name,
            sn: nodes[host]?.sn,
            size: pre.size + next.size,
            used: pre.used + next.used,
            free: pre.free + next.free,
          }
        },
      )
      this.data = (nodes[host]?.disks || []).map(item => {
        return {
          ...item,
          docking_info: docking_info,
          percentage: +((item.used / item.size) * 100).toFixed(2),
          warittingIdsCount: (item.volume_ids || []).length,
        }
      })
    },
    // eslint-disable-next-line camelcase
    async setBad({ disk_id }) {
      try {
        await this.$confirm(this.$t('resource.confirmsetbroken'), this.$t('common.notice'), {
          confirmButtonText: this.$t('common.yes'),
          cancelButtonText: this.$t('common.no'),
        })
        const res = await setNodeInfoBad({ region: this.clusterName, clusterId: this.app.clusterId, disk_id })
        if (res.data?.orderDetailUrl) {
          try {
            await this.$confirm(this.$t('resource.brokenjob'), this.$t('common.notice'), {
              confirmButtonText: this.$t('common.view'),
              type: 'success',
            })
            window.open(res.data?.orderDetailUrl)
          } catch (error) {
            this.getData()
          }
        } else {
          this.$message.success(this.$t('resource.setsuc'))
        }
        await this.getData()
      } catch (e) {
        console.log(e)
      }
    },
    async getDroppingList() {
      this.droppingList = []
      const res = await offLineDropNodeList({
        region: this.clusterName,
        clusterId: this.app.clusterId,
      })
      this.droppingList = res.data.disks || []
    },
    // eslint-disable-next-line camelcase
    async offLine({ disk_id, host, path }) {
      try {
        const remark = await this.$prompt(this.$t('resource.autooffline'), this.$t('common.notice'), {
          confirmButtonText: this.$t('common.yes'),
          cancelButtonText: this.$t('common.no'),
          inputValue: '',
        })
        const { user_id: userId, user_name: userName } = this.user
        const res = await offLineNodeInfo({ region: this.clusterName, clusterId: this.app.clusterId, disk_id, user_id: userId, host_name: '', user_name: userName, node: host, disk_path: path, remark: `${remark?.value}`, docking_remark: '审批通过,将下线' })
        if (res.data?.orderDetailUrl) {
          try {
            await this.$confirm(this.$t('resource.cloudjobs'), this.$t('common.notice'), {
              confirmButtonText: this.$t('common.view'),
              type: 'success',
            })
            window.open(res.data?.orderDetailUrl)
          } catch (error) {
            this.getData()
          }
        } else {
          this.$message.success(this.$t('common.offline') + this.$t('common.xxsuc'))
        }
        await this.getData()
      } catch (e) {}
    },
    // eslint-disable-next-line camelcase
    async changeRW({ disk_id, readonly }) {
      try {
        await this.$confirm(
          this.$t('resource.setdiskstatus')+ `${readonly ?  this.$t('common.rw'): this.$t('common.ro')}?`,
          this.$t('common.notice'),
          {
            confirmButtonText: this.$t('common.yes'),
            cancelButtonText: this.$t('common.no'),
            type: 'warning',
          },
        )
        await changeNodeInfoRW({
          region: this.clusterName,
          clusterId: this.app.clusterId,
          disk_id,
          readonly: !readonly,
        })
        this.$message.success(this.$t('resource.setsuc'))
        await this.getData()
      } catch (e) {}
    },
    toDataBlock(id, type) {
      this.$router.push({ query: { vid: id } })
    },
  },
}
</script>
<style lang="scss" scoped>
.dis-btn {
  color: #bfbfbf;
}
.ml-8 {
  margin-left: 8px;
}
.w100 {
  width: 100%;
}
.mb10{
  margin-bottom: 10px;
}
.p-t-26 {
  padding-top: 26px;
}
.back {
  // color: #00c9c9;
  cursor: pointer;
  text-align: right;
  margin-right: 20px;
  float: right;
}
.color9 {
  color: #999;
}
.color6 {
  color: #666;
}
.mr-rt {
  margin-right: 20px;
}
::v-deep .el-progress__text {
  display: block !important;
  font-size: 12px !important;
}
.text-center {
  text-align: center;
  width: 80px;
}
.inside{
  margin: 10px;
}
.fontTypeSpan {
font-family: 'OPPOSans M';
font-style: normal;
font-weight: 400;
font-size: 16px;
line-height: 20px;
}

::v-deep .el-drawer__body {
  overflow: auto;
}
/*2.隐藏滚动条，太丑了*/
::v-deep .el-drawer__container ::-webkit-scrollbar{
    width: 10px;
}
</style>
