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
    v-if="dialogVisible"
    :title="title"
    :visible.sync="dialogVisible"
    width="540px"
  >
    <el-form ref="form" :model="form" label-width="100px" :rules="rules" style="width: 450px">
      <el-form-item :label="$t('common.username')" prop="user_name">
        <el-input v-model="form.user_name" :disabled="type !== 'create'" :placeholder="$t('usermgt.inputname2')" />
      </el-form-item>
      <el-form-item :label="$t('common.phone')" prop="phone">
        <el-input v-model="form.phone" :disabled="type === 'delete'" :placeholder="$t('common.phone')" />
      </el-form-item>
      <el-form-item :label="$t('common.email')" prop="email">
        <el-input v-model="form.email" :disabled="type === 'delete'" :placeholder="$t('common.email')" />
      </el-form-item>
      <el-form-item :label="$t('common.role')" prop="roles">
        <el-select
          v-model="form.role_ids"
          multiple
          filterable
          :disabled="type === 'delete'"
          :placeholder = "$t('common.select')"
        >
          <el-option
            v-for="item in roleList"
            :key="item.id"
            :label="item.role_name"
            :value="item.id"
          />
        </el-select>
      </el-form-item>
    </el-form>
    <span slot="footer" class="dialog-footer">
      <el-button @click="dialogVisible = false">{{ $t('button.cancel') }}</el-button>
      <el-button v-if="type !== 'delete'" type="primary" @click="submit">{{ $t('button.submit') }}</el-button>
      <el-button v-else type="danger" @click="submit">{{ $t('button.delete') }}</el-button>
    </span>
  </el-dialog>
</template>

<script>
import {
  getRoleList,
  userCreate,
  userUpdate,
  userDelete,
} from '@/api/auth'
export default {
  name: 'UserDialog',
  data() {
    return {
      dialogVisible: false,
      id: '',
      type: '',
      form: {
        user_name: '',
        phone: '',
        email: '',
        role_ids: [],
      },
      roleList: [],
      rules: {
        user_name: [{ required: true, message: this.$t('usermgt.inputname2'), trigger: 'blur' }],
        phone: [{ required: true, message: this.$t('usermgt.phone'), trigger: 'blur' }],
        email: [{ required: true, message: this.$t('usermgt.email'), trigger: 'blur' }],
      },
    }
  },
  computed: {
    title() {
      let title = ''
      switch (this.type) {
        case 'create':
          title = this.$t('common.add') + this.$t('common.user')
          break
        case 'edit':
          title = this.$t('common.edit')
          break
        case 'delete':
          title = this.$t('common.delete') + this.$t('common.user')
          break
      }
      return title
    },
  },
  mounted() {
    this.getRoleList()
  },
  methods: {
    initDialog() {
      this.type = 'create'
      this.form = {
        user_name: '',
        phone: '',
        email: '',
        role_ids: [],
      }
      this.dialogVisible = true
    },
    async editDialog(userInfo) {
      // eslint-disable-next-line camelcase
      const { id, user_name, phone, email, roles } = userInfo
      this.form = {
        user_name,
        phone,
        email,
        role_ids: roles.map(item => item.id),
      }
      this.id = id
      this.type = 'edit'
      this.dialogVisible = true
    },
    async deleteDialog(userInfo) {
      // eslint-disable-next-line camelcase
      const { id, user_name, phone, email, roles } = userInfo
      this.form = {
        user_name,
        phone,
        email,
        role_ids: roles.map(item => item.id),
      }
      this.id = id
      this.type = 'delete'
      this.dialogVisible = true
    },
    async getRoleList() {
      const { data: { roles } } = await getRoleList({ page_size: 500 })
      this.roleList = roles
    },
    async submit() {
      await this.$refs.form.validate()
      if (this.type === 'create') {
        await userCreate({
          ...this.form,
          password: 'abcd1234',
        })
        this.$message.success(this.$t('common.create') + this.$t('common.xxsuc'))
        this.$emit('submit')
      } else if (this.type === 'edit') {
        await userUpdate({
          ...this.form,
          id: this.id,
        })
        this.$message.success(this.$t('common.edit') + this.$t('common.xxsuc'))
        this.$emit('submit')
      } else if (this.type === 'delete') {
        await userDelete({
          ids: [this.id],
        })
        this.$message.success(this.$t('common.delete') + this.$t('common.xxsuc'))
        this.$emit('submit')
      }
      this.dialogVisible = false
    },
  },
}

</script>
<style scoped lang='scss'>
</style>
