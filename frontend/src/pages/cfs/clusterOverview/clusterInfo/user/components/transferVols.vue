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
  <el-dialog
    :title="$t('privileges.CFS_USERS_VOLS_TRANSFER')"
    :visible.sync="dialogFormVisible"
    width="700px"
    :append-to-body="true"
    @closed="clearData"
  >
    <el-form
      ref="form"
      :model="forms"
      :rules="rules"
      label-width="25%"
      class="mid-block"
    >
      <el-form-item :label="$t('common.volumename') + ':'" prop="volName">
        <el-input v-model="forms.volName" disabled class="input"></el-input>
      </el-form-item>

      <el-form-item :label="$t('common.volumebelongsto') + ':'" prop="userSrc">
        <el-input v-model="forms.userSrc" disabled class="input"></el-input>
      </el-form-item>

      <el-form-item :label="$t('common.transferownershipto') + ':'" prop="userDst">
        <el-select
          v-model="forms.userDst"
          class="input"
          filterable
          reserve-keyword
          remote
          :placeholder="$t('volume.selectowner')"
          :remote-method="remoteMethod"
        >
          <el-option
            v-for="item in userList"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          ></el-option>
        </el-select>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button ref="pol" type="primary" @click="doCheck">{{ $t('button.submit') }}</el-button>
      <el-button ref="pol" type="primary" @click="close">{{ $t('button.cancel') }}</el-button>
    </div>
  </el-dialog>
</template>
<script>
import {getUserNameList, transferVol} from '@/api/cfs/cluster'
import Mixin from '@/pages/cfs/clusterOverview/mixin'
export default {
  mixins: [Mixin],
  data() {
    return {
      userList: [],
      forms: {
        userSrc: '',
        userDst: '',
        volName: '',
      },
      dialogFormVisible: false,
    }
  },
  computed: {
    rules() {
      return {
        userSrc: [
          {
            required: true,
            message: this.$t('usermgt.inputname'),
            trigger: 'blur',
          },
        ],
        userDst: [
          {
            required: true,
            message: this.$t('usermgt.inputname'),
            trigger: 'blur',
          },
        ],
        volName: [
          {
            required: true,
            message: this.$t('volume.inputvolume'),
            trigger: 'blur',
          },
        ],
      }
    },
  },
  created() {},
  methods: {
    initForm(val) {
      this.forms = { ...this.forms, ...val }
    },
    open() {
      this.dialogFormVisible = true
      // 默认显示所有的用户
      this.getUserNameList()
    },
    clearData() {
      this.forms = {
        userSrc: '',
        userDst: '',
        volName: '',
      }
    },
    async doCheck() {
      await this.$refs.form.validate()
      const res = await transferVol({
        volume: this.forms.volName,
        user_src: this.forms.userSrc,
        user_dst: this.forms.userDst,
        cluster_name: this.clusterName,
      })
      if (res.code === 200) {
        this.$message.success(this.$t('privileges.CFS_USERS_VOLS_TRANSFER') + this.$t('common.permissions') + this.$t('common.xxsuc') + res.data)
        this.close()
      } else {
        this.$message({
          showClose: true,
          message: this.$t('privileges.CFS_USERS_VOLS_TRANSFER') + this.$t('common.permissions') + this.$t('common.failed') + '\n' + res.data,
          type: 'error',
          duration: 10000
        })
      }
      this.$emit('refresh')
      this.close()
    },
    close() {
      this.clearData()
      this.dialogFormVisible = false
    },
    async remoteMethod(query) {
      this.getUserNameList(query)
    },
    async getUserNameList(query) {
      const res = await getUserNameList({
        keywords: query,
        cluster_name: this.clusterName,
      })
      this.userList = (res.data || []).map((item) => {
        return {
          label: item,
          value: item,
        }
      })
    },
  },
}
</script>
<style lang="scss" scoped>
.input {
  width: 60%;
}
.dialog-footer {
  text-align: center;
}
</style>
