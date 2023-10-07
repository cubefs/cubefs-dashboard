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
    <div class="mg-bt-s">
      <div>
        <el-radio-group v-model="activeName" style="margin-bottom: 10px;">
          <el-radio-button label="partition">{{ $t('common.copies') }}</el-radio-button>
          <el-radio-button label="blobstoreVolume" :disabled="!ebsClusterList || !ebsClusterList.length">{{ $t('common.ec') }}</el-radio-button>
        </el-radio-group>
      </div>
      <div v-if="activeName === 'partition'" class="flex">
        <span class="fontType"><span>{{ $t('common.total') }}{{ $t('common.nodes') }}:</span> <span class="mg-lf-m"></span>{{ info.node }}</span>
        <span class="fontType mg-lf-m"><span>{{ $t('common.total') }}{{ $t('common.partitions') }}:</span> <span class="mg-lf-m"></span>{{ info.partition }}</span>
        <span class="fontType mg-lf-m"><span>{{ $t('common.broken') }}{{ $t('common.partitions') }}:</span> <span class="mg-lf-m"></span><span class="bad_partition" @click="showDialog('DP' + this.$t('common.status'))">{{ badDataPartitionNum }}</span>/{{ (badDataPartitionNum / info.partition * 100).toFixed()+'%' || '0%' }}</span>
        <span class="fontType mg-lf-m"><span>{{ $t('common.total') }}{{ $t('common.size') }}:</span> <span class="mg-lf-m"></span>{{ info.total |renderSize }}</span>
        <div class="mg-lf-m progress">
          <span>{{ info.used |renderSize }}/{{ (info.used/info.total*100).toFixed(0)+'%' }}</span>
          <el-progress
            v-if="info.node!==0"
            :stroke-width="10"
            :show-text="false"
            :percentage="info.used/info.total*100"
            :color="[
              { color: '#f56c6c', percentage: 100 },
              { color: '#e6a23c', percentage: 80 },
              { color: '#5cb87a', percentage: 60 },
              { color: '#1989fa', percentage: 40 },
              { color: '#6f7ad3', percentage: 20 },
            ]"
          >
          </el-progress>
        </div>
      </div>
    </div>
    <component
      :is="items[activeName].component"
      v-if="items[activeName].name === activeName"
      :info.sync="info"
    />
    <el-dialog
      v-if="DataPartitionDetailDialogVisible"
      :title="$t('common.broken') + 'DP' + $t('common.details')"
      width="65%"
      :visible.sync="DataPartitionDetailDialogVisible"
      center
    >
      <el-table
        :data="PartitionTableData"
        style="width: 100%"
      >
        <el-table-column
          :label="$t('common.partitionid')"
          prop="PartitionID"
          :width="100"
        ></el-table-column>
        <el-table-column :label="$t('common.volumename')" prop="VolName"></el-table-column>
        <el-table-column
          :label="$t('common.copies')"
          prop="ReplicaNum"
        ></el-table-column>
        <el-table-column label="isRecovering" :width="100">
          <template slot-scope="scope">
            <span>{{ scope.row.IsRecover }}</span>
          </template>
        </el-table-column>
        <el-table-column label="Leader" prop="Leader" width="180"></el-table-column>
        <el-table-column label="Members" width="180">
          <template slot-scope="scope">
            <div v-for="item in scope.row.Members" :key="item">{{ item }}</div>
          </template>
        </el-table-column>
        <el-table-column
          :label="$t('common.status')"
          prop="Status"
          :width="150"
        ></el-table-column>
      </el-table>
    </el-dialog>
    <el-dialog
      v-if="DataPartitionDialogVisible"
      :title="$t('common.broken') + 'DP'"
      width="65%"
      :visible.sync="DataPartitionDialogVisible"
      center
      top="5vh"
    >
      <div>{{ $t('volume.copymiss') }}</div>
      <el-table
        max-height="350"
        :data="LackReplicaDataPartitionIDs"
        style="width: 100%"
      >
        <el-table-column
          :label="$t('common.id')"
          type="index"
        >
        </el-table-column>
        <el-table-column
          :label="$t('common.partitionid')"
          prop="id"
        >
          <template slot-scope="scope">
            <div>{{ scope.row }}</div>
          </template></el-table-column>
        <el-table-column
          :label="$t('common.copies')"
          prop="ReplicaNum"
        ></el-table-column>
        <el-table-column
          :label="$t('common.action')"
        >
          <template slot-scope="scope">
            <el-button
              size="medium"
              type="text"
              @click="showDetail(scope.row, 2)"
            >{{ $t('common.detail') }}</el-button>
          </template>
          ></el-table-column>
      </el-table>
      <div>{{ $t('volume.leadermiss') }}</div>
      <el-table
        max-height="300"
        :data="CorruptDataPartitionIDs"
        style="margin-top:5px"
      >
        <el-table-column
          :label="$t('common.id')"
          type="index"
        >
        </el-table-column>
        <el-table-column
          :label="$t('common.partitionid')"
          prop="id"
        >
          <template slot-scope="scope">
            <div>{{ scope.row }}</div>
          </template></el-table-column>
        <el-table-column
          :label="$t('common.copies')"
          prop="ReplicaNum"
        ></el-table-column>
        <el-table-column
          :label="$t('common.action')"
        >
          <template slot-scope="scope">
            <el-button
              size="medium"
              type="text"
              @click="showDetail(scope.row, 2)"
            >{{ $t('common.detail') }}</el-button>
          </template>
          ></el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>
<script>
import Partition from './partition.vue'
import BlobstoreVolume from './blobStoreVolume'
import { renderSize } from '@/utils'
import mixin from '@/pages/cfs/clusterOverview/mixin'
import graphics from '../components/graphics'
export default {
  name: '',
  components: { Partition, BlobstoreVolume },
  filters: {
    renderSize(val) {
      const data = renderSize(val, 1)
      return data
    },
  },
  mixins: [mixin, graphics],
  props: [''],
  data () {
    return {
      activeName: 'partition',
      items: {
        partition: {
          name: 'partition',
          component: 'Partition',
        },
        blobstoreVolume: {
          name: 'blobstoreVolume',
          component: 'BlobstoreVolume',
        },
      },
      info: {
        node: 0,
        partition: 0,
        total: 0,
        used: 0,
      },
    }
  },
  computed: {},
  watch: {},
  created() {},
  beforeMount() {},
  mounted() {},
  methods: {},
}
</script>
<style lang='scss' scoped>
.fontType{
font-family: 'Microsoft YaHei';
font-style: normal;
font-weight: 400;
font-size: 14px;
line-height: 14px;
/* identical to box height, or 167% */
color: #000000;
}
.progress{
  width: 100px;
  position: relative;
  top: -5px;
  left: 10px;
}
.bad_partition {
  cursor: pointer;
  color: #38c59f;
}
</style>
