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
    :title="$t('common.create') + $t('common.tenant')"
    :visible.sync="dialogFormVisible"
    width="800px"
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
      <el-form-item :label="$t('tenant.tenantid')+ ':'" prop="user">
        <el-input
          v-model="forms.user"
          class="input"
          :placeholder="$t('tenant.inputtenantid')"
        ></el-input>
        <el-tooltip
          class="item"
          effect="dark"
          :content="$t('tenant.idrule')"
          placement="top"
        >
          <i class="el-icon-question fontS16"></i>
        </el-tooltip>
      </el-form-item>
      <el-form-item :label="$t('tenant.tenantdescrp')" prop="description">
        <el-input
          v-model="forms.description"
          type="textarea"
          :rows="2"
          :placeholder="$t('tenant.inputdescrp')"
          class="input"
        ></el-input>
      </el-form-item>
      <el-form-item :label="$t('common.tenant') + $t('common.type') + ':'" prop="accountType">
        <el-radio-group v-model="forms.accountType">
          <!-- <el-radio :label="1">root</el-radio>
          <el-radio :label="2">admin</el-radio> -->
          <el-radio :label="3">normal</el-radio>
        </el-radio-group>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button ref="pol" type="primary" @click="doCheck">{{ $t('button.submit') }}</el-button>
      <el-button ref="pol" type="primary" @click="close">{{ $t('button.cancel') }}</el-button>
    </div>
  </el-dialog>
</template>
<script>
import { createUser } from '@/api/cfs/cluster'
import Mixin from '@/pages/cfs/clusterOverview/mixin'
export default {
  mixins: [Mixin],
  data() {
    return {
      userList: [],
      forms: {
        user: '',
        description: '',
        accountType: 3,
      },
      dialogFormVisible: false,
    }
  },
  computed: {
    rules() {
      return {
        user: [
          {
            required: true,
            trigger: 'blur',
            validator: (rule, value, cb) => {
              const reg = /^[a-zA-Z][a-zA-Z_0-9]*$/
              if (!value) {
                cb(new Error(this.$t('tenant.inputtenantid')))
              } else if (!reg.test(value)) {
                cb(new Error(this.$t('tenant.idrule')))
              }
              cb()
            },
          },
        ],
      }
    },
  },
  methods: {
    initForm(val) {
      this.forms = { ...val }
    },
    open() {
      this.dialogFormVisible = true
    },
    clearData() {
      this.forms = {
        user: '',
        description: '',
        accountType: 3,
      }
    },
    async doCheck() {
      await this.$refs.form.validate()
      const { user, description, accountType } = this.forms
      await createUser({
        id: user,
        type: accountType,
        description,
        cluster_name: this.clusterName,
      })
      this.$message.success(this.$t('tenant.createtenantsuc'))
      this.$emit('refresh')
      this.close()
    },
    close() {
      this.clearData()
      this.dialogFormVisible = false
    },
  },
}
</script>
<style lang="scss" scoped>
.input {
  width: 65%;
}
.dialog-footer {
  text-align: center;
}
</style>
