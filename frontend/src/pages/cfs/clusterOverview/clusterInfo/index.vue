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
    <el-card>
      <el-button
        v-auth="'CFS_USERS_LIST'"
        class="fl-rt"
        type="primary"
        @click.stop="userDialogVisible = true"
      >租户管理</el-button>
      <div class="detail_title label">
        <div>
          <span><span>集群名称:</span> <span class="mg-lf-m fontType">{{ clusterInfo.name }}</span></span>
          <span
            class="mg-lf-m"
          ><span>总容量:</span><span class="mg-lf-m fontType">{{ clusterInfo.data_total }}</span></span>
          <span
            class="mg-lf-m"
          ><span>总文件数:</span><span class="mg-lf-m fontType">{{ dentry_count }}</span></span>
          <span
            class="mg-lf-m"
          ><span>总卷数:</span><span class="mg-lf-m fontType">{{ volTotal }}</span></span>
          <span
            class="mg-lf-m"
          ><span>cli:</span> <span class="mg-lf-m fontType">{{ clusterInfo.cli }}</span></span>
        </div>
        <div class="m-t-5">
          <span>
            <span>管理服务地址:</span>
            <span
              v-for="item in clusterInfo.master_addr"
              :key="item"
              class="mg-lf-m fontType address"
              :class="{
                leader_address: item === clusterInfo.leader_addr
              }"
            >{{ item }}{{ item === clusterInfo.leader_addr ? '(主)' : '(从)' }}</span>
          </span>
        </div>
        <el-dialog
          top="10vh"
          title="租户列表"
          :visible.sync="userDialogVisible"
        >
          <UserManage />
        </el-dialog>
      </div>
      <el-collapse v-model="collapseName">
        <el-collapse-item name="1">
          <span slot="title" class="collapse-title">图形看板</span>
          <GraphicsBoard
            :cluster-info="clusterInfo"
          />
        </el-collapse-item>
      </el-collapse>
    </el-card>
  </div>
</template>
<script>
import { mapGetters } from 'vuex'
import GraphicsBoard from './components/graphicsBoard.vue'
import UserManage from './user'
import { getVolList } from '@/api/cfs/cluster'
export default {
  name: '',
  components: {
    GraphicsBoard,
    UserManage,
  },
  data () {
    return {
      collapseName: '1',
      activeName: 'volumn',
      userDialogVisible: false,
      tabs: [
        {
          label: '卷管理',
          name: 'volumn',
          component: 'Volumn',
        },
        {
          label: '数据管理',
          name: 'dataManage',
          component: 'DataManage',
        },
        {
          label: '元数据管理',
          name: 'metaDataManage',
          component: 'MetaDataManage',
        },
      ],
      volTotal: 0,
      dentry_count: 0,
    }
  },
  computed: {
    ...mapGetters('clusterInfoModule', {
      curClusterInfo: 'clusterInfog',
    }),
    clusterInfo() {
      return this.curClusterInfo.clusterInfo
    },
  },
  watch: {},
  created() {
  },
  beforeMount() {},
  mounted() {
    this.getData()
  },
  methods: {
    async getData() {
      this.dataList = []
      this.volTotal = 0
      const { data: volList } = await getVolList({
        cluster_name: this.clusterInfo.name,
      })
      this.dentry_count = 0
      volList.forEach(item => {
        this.dentry_count += item.dentry_count
      })
      this.volTotal = volList.length || 0
    },
  },
}
</script>
<style lang='scss' scoped>
.fontType {
  font-weight: bold;
}
::v-deep.el-collapse-item__header{
      padding:0px 20px;
      flex:1 0 auto;
      order: -1;
    }
    ::v-deep.collapse-title {
      flex: 1 0 90%; //位于左侧
      order: 1;
    }
.detail_title {
  line-height: 20px;
  padding-left: 8px;
  font-size: 14px;
  margin-bottom: 20px;
}
.label {
  color: #666;
}

.address {
  background: #e0e2e6;
  padding: 3px;
  border-radius: 4px;
}

.leader_address {
  color: #2fc29b;
}

</style>
