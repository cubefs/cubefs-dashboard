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
    <div class="bar">
      <el-input v-model="user_name" prefix-icon="el-icon-search" :placeholder="$t('usermgt.inputname')" style="width: 240px" clearable />
      <div>
        <el-button v-auth="'AUTH_USER_DELETE'" type="text" icon="el-icon-delete" style="color: #ed4014">{{ $t('button.delete') }}</el-button>
        <el-button type="primary" icon="el-icon-plus" @click="addUser">{{ $t('common.add') }}{{ $t('common.user') }}</el-button>
      </div>
    </div>
    <o-page-table
      ref="table"
      :url="url"
      :columns="tableColumns"
      :form-data="{
        user_name_like: user_name
      }"
      data-key="users"
      method="get"
      page-key="page"
      page-size-key="page_size"
      total-key="count"
    />
    <UserDialog ref="dialog" @submit="onSubmit" />
  </el-card>
</template>

<script>
import UserDialog from './userDialog.vue'
import { passwordUpdate } from '@/api/auth'
export default {
  components: {
    UserDialog,
  },
  data() {
    return {
      url: '/api/cubefs/console/auth/user/list',
      user_name: '',
    }
  },
  computed: {
    tableColumns() {
      return [
        {
          type: 'selection',
          align: 'center',
          width: 80,
          selectable: this.selectInit,
        },
        {
          title: this.$t('common.name'),
          key: 'user_name',
          width: 200,
        },
        {
          title: this.$t('common.role'),
          key: 'roles',
          render: (h, { row }) => {
            return <div>
              {
                row.roles.map((item, index) => {
                  return <span>{ item.role_name }{ index < row.roles.length - 1 ? '、' : '' }</span>
                })
              }
            </div>
          },
        },
        {
          title: this.$t('common.phone'),
          key: 'phone',
        },
        {
          title: this.$t('common.email'),
          key: 'email',
        },
        {
          title: this.$t('common.action'),
          key: '',
          width: 200,
          render: (h, { row }) => {
            return <div>
              <el-button v-auth="AUTH_USER_UPDATE" type='text' onClick={() => this.editUser(row)}>{ this.$t('common.edit') }</el-button>
              <el-button v-auth="AUTH_USER_PASSWORD_UPDATE" type='text' onClick={() => this.updatePassword(row)}>{this.$t('common.passwd')}</el-button>
              <el-button v-auth="AUTH_USER_DELETE" type='text' style={ row.id === 1 ? {} : { color: '#ed4014' } } disabled={row.id === 1} onClick={() => this.deleteUser(row)}>{ this.$t('common.delete') }</el-button>
            </div>
          },
        },
      ]
    },
  },
  methods: {
    addUser() {
      this.$refs.dialog.initDialog()
    },
    editUser(row) {
      this.$refs.dialog.editDialog(row)
    },
    deleteUser(row) {
      this.$refs.dialog.deleteDialog(row)
    },
    selectInit(row, index) {
      if ([1].includes(row.id)) {
        return false
      } else {
        return true
      }
    },
    updatePassword(row) {
      this.$prompt('<div>' + this.$t('usermgt.inputpasswd') + '<div style="color: red; font-size: 12px;">' + this.$t('usermgt.passwdlen') + '<br />' + this.$t('usermgt.passwdrule') + '</div><div>', this.$t('usermgt.chgpasswd'), {
        confirmButtonText: this.$t('common.submit'),
        cancelButtonText: this.$t('common.cancel'),
        inputValidator: this.checkPassword,
        inputErrorMessage: this.$t('usermgt.wrongpasswd'),
        dangerouslyUseHTMLString: true,
      }).then(async ({ value }) => {
        await passwordUpdate({
          user_name: row.user_name,
          password: value,
        })
        this.$message.success(this.$t('usermgt.updatesuc'))
      }).catch(() => {})
    },
    checkPassword(str) {
      // 密码只能包含数字、大写字母、小写字母、特殊字符（~!@#$%^&*_.?），且至少两种类型以上
      const regList = [/[0-9]/, /[A-Z]/, /[a-z]/, /[~!@#$%^&*_.?]/]
      let num = 0
      regList.forEach(item => {
        if (item.test(str)) {
          num++
        }
      })
      if (str.length > 7 && str.length < 17 && num > 1) {
        return true
      }
      return false
    },
    async onSubmit() {
      this.$refs.table.search()
      await this.$store.dispatch('moduleUser/setAuth')
    },
  },
}

</script>
<style scoped lang='scss'>
.bar {
  display: flex;
  justify-content: space-between;
}
</style>
