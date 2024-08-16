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
  <el-card>
    <div style="display: flex; justify-content: space-between;">
      <el-button
        v-auth="'CFS_USERS_CREATE'"
        icon="el-icon-circle-plus"
        type="primary"
        @click.stop="createUser"
      >{{ $t('common.create') }}{{ $t('common.tenant') }}</el-button>
      <div style="display: flex;">
        <el-input
          v-model="params.userID"
          class="m-r-15"
          :placeholder="$t('tenant.inputtenantid')"
          clearable
        ></el-input>
        <el-button
          type="primary"
          @click="getData"
        >{{ $t('button.search') }}</el-button>
      </div>
    </div>
    <el-row class="userInfo">
      <u-page-table :data="dataList" :page-size="page.per_page">
        <!-- <el-table-column label="序号" type="index"></el-table-column> -->
        <el-table-column :label="$t('tenant.tenantid')" prop="user_id" :width="100"></el-table-column>
        <el-table-column :label="$t('common.tenant') + $t('common.type')" prop="user_type" :width="100"></el-table-column>
        <el-table-column label="AK" prop="access_key"></el-table-column>
        <el-table-column label="SK" prop="secret_key"></el-table-column>
        <el-table-column :label="$t('common.createtime')" prop="create_time" :width="100"></el-table-column>
        <el-table-column :label="$t('common.action')" :width="120" align="center" fixed="right">
          <template slot-scope="scope">
            <el-button
              v-auth="'CFS_USERS_DELETE'"
              type="text"
              size="small"
              v-if="scope.row.user_type != 'root'"
              @click="deleteUser(scope.row)"
            >{{ $t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </u-page-table>
    </el-row>
    <el-row>
      <el-row>
        <el-col :span="12">
          <div class="detail_title">{{ $t('common.permissions') }} {{ $t('common.detail') }}</div>
        </el-col>
        <el-col :span="12">
          <el-row class="search">
            <el-input
              v-model.trim="params.volName"
              :placeholder="$t('volume.inputvolume')"
              clearable
              class="input"
            ></el-input>
            <el-button
              type="primary"
              class="search-btn"
              @click="getData"
            >{{ $t('button.search') }}</el-button>
          </el-row>
        </el-col>
      </el-row>
      <el-row>
        <u-page-table :data="dataListAuth" :page-size="pageAuth.per_page">
          <!-- <el-table-column label="序号" type="index"></el-table-column> -->
          <el-table-column :label="$t('tenant.tenantid')" prop="user_id"></el-table-column>
          <el-table-column :label="$t('common.volume')" prop="volume"></el-table-column>
          <!-- <el-table-column label="子目录" prop="business_team"></el-table-column> -->
          <el-table-column :label="$t('common.business')" prop="business"></el-table-column>
          <el-table-column :label="$t('common.permissions')" prop="policy">
            <template slot-scope="scope">
              <el-radio
                v-for="item in scope.row.policy"
                :key="item"
                disabled
                :value="getTableItemPolicyVal(item)"
                :label="getTableItemPolicyVal(item)"
              >{{ getTableItemPolicyVal(item) }}</el-radio>
            </template>
          </el-table-column>
          <el-table-column :label="$t('common.action')" :width="120">
            <template slot-scope="scope">
              <el-button
                v-auth="'CFS_USERS_POLICIES_DELETE'"
                type="text"
                size="small"
                @click="deleteUserPolicies(scope.row)"
                v-if="scope.row.user_id != 'root' && scope.row.policy.length >= 1 && scope.row.policy[0] != 'owner'"
              >{{ $t('common.clear') + $t('common.permissions') }}</el-button>
              <el-button
                v-auth="'CFS_USERS_VOLS_TRANSFER'"
                type="text"
                size="small"
                :disabled="true"
                v-if="scope.row.user_id != 'root' && scope.row.policy.length == 1 && scope.row.policy[0] == 'owner'"
              >{{ $t('privileges.CFS_USERS_VOLS_TRANSFER') }}</el-button>
            </template>
          </el-table-column>
        </u-page-table>
      </el-row>
    </el-row>
    <CreateUser ref="createUser" @refresh="refresh" />
  </el-card>
</template>
<script>
import UPageTable from '@/pages/components/uPageTable'
import CreateUser from './components/createUser'
import {deleteUser, deleteUserPolicy, deleteVol, getUserList} from '@/api/cfs/cluster'
import Mixin from '@/pages/cfs/clusterOverview/mixin'
export default {
  components: {
    CreateUser,
    UPageTable,
  },
  mixins: [Mixin],
  data() {
    return {
      userTotal: 0,
      dataList: [],
      dataListAuth: [],
      params: {
        userID: '', // 输入租户id
        volName: '', // 输入卷名
      },
      page: {
        per_page: 5, // 页面大小
      },
      pageAuth: {
        per_page: 5, // 页面大小
      },
    }
  },
  computed: {},
  watch: {},
  created() {
    this.getData()
  },
  methods: {
    getTableItemPolicyVal(val) {
      return val.split(':').slice(-1).join('').toLowerCase()
    },
    refresh() {
      this.params = {
        userID: '',
        volName: '',
      }
      this.getData()
    },
    async getData() {
      this.dataList = []
      this.dataListAuth = []
      this.userTotal = 0
      const { userID, volName } = this.params
      if (this.clusterName) {
        const res = await getUserList({
          keywords: userID,
          cluster_name: this.clusterName,
          vol_name: volName,
        })
        const { users, policy } = res.data
        this.dataList = users || []
        this.userTotal = users.length || 0
        this.dataListAuth = policy || []
      }
    },
    createUser() {
      this.$refs.createUser.open()
    },
    async deleteUser(row) {
      try {
        await this.$confirm(this.$t('volume.confirmdeletetenant') + row.user_id, this.$t('common.notice'), {
          confirmButtonText: this.$t('common.yes'),
          cancelButtonText: this.$t('common.no'),
        })
        const res = await deleteUser({
          id: row.user_id,
          cluster_name: this.clusterName,
        })
        if (res.code === 200) {
          this.$message.success(this.$t('common.delete') + this.$t('common.xxsuc') + res.data)
          await this.getData()
        } else {
          this.$message({
            showClose: true,
            message: this.$t('common.delete') + this.$t('common.failed') + '\n' + res.data,
            type: 'error',
            duration: 10000
          })
        }
      } catch (e) {}
    },
    async deleteUserPolicies(row) {
      try {
        await this.$confirm(this.$t('volume.confirmdeletepolicies') + row.user_id + '  ' + row.volume, this.$t('common.notice'), {
          confirmButtonText: this.$t('common.yes'),
          cancelButtonText: this.$t('common.no'),
        })
        const res = await deleteUserPolicy({
          user_id: row.user_id,
          volume: row.volume,
          cluster_name: this.clusterName,
        })
        if (res.code === 200) {
          this.$message.success(this.$t('common.clear') + this.$t('common.permissions') + this.$t('common.xxsuc') + res.data)
          await this.getData()
        } else {
          this.$message({
            showClose: true,
            message: this.$t('common.clear') + this.$t('common.permissions') + this.$t('common.failed') + '\n' + res.data,
            type: 'error',
            duration: 10000
          })
        }
      } catch (e) {}
    },


  },
}
</script>
<style lang="scss" scoped>
.noborder {
  border: none;
}
.input {
  position: absolute;
  width: 300px;
  right: 70px;
}
.search {
  position: relative;
  text-align: right;
}
.search-btn {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  height: 31px;
}
.userInfo {
  margin-bottom: 40px;
}
.detail_title {
  height: 20px;
  line-height: 20px;
  border-left: 4px solid #66cc99;
  padding-left: 8px;
  font-size: 14px;
  margin-bottom: 20px;
}
</style>
