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
      <el-input v-model="role_name" prefix-icon="el-icon-search" :placeholder="$t('rolemgt.inputname')" style="width: 240px" clearable />
      <div>
        <el-button v-auth="'AUTH_ROLE_DELETE'" type="text" icon="el-icon-delete" style="color: #ed4014">{{ $t('common.delete') }}</el-button>
        <el-button v-auth="'AUTH_ROLE_CREATE'" type="primary" icon="el-icon-plus" @click="addRole">{{ $t('common.add') + " " + $t('common.role') }}</el-button>
      </div>
    </div>
    <o-page-table
      ref="table"
      :url="url"
      :columns="tableColumns"
      :form-data="{
        role_name_like: role_name
      }"
      data-key="roles"
      method="get"
      page-key="page"
      page-size-key="page_size"
      total-key="count"
      :after-send-hook="afterSendHook"
    />
    <RoleDialog ref="dialog" @submit="onSubmit" />
  </el-card>
</template>

<script>
import RoleDialog from './roleDialog.vue'
export default {
  components: {
    RoleDialog,
  },
  data() {
    return {
      url: '/api/cubefs/console/auth/role/list',
      role_name: '',
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
          title: 'ID',
          key: 'id',
          width: 120,
        },
        {
          title: this.$t('common.role') + " " + this.$t('common.name'),
          key: 'role_name',
          width: 200,
        },
        {
          title: this.$t('common.privilege'),
          key: 'permissions',
          render: (h, { row }) => {
            return <div>
              {
                row.permissions.map((item, index) => {
                  return <span>{ this.$t('privileges.' + item.auth_code) }{ index < row.permissions.length - 1 ? '、' : '' }</span>
                })
              }
            </div>
          },
        },
        {
          title: this.$t('common.action'),
          key: '',
          align: 'center',
          width: 200,
          render: (h, { row }) => {
            return <div>
              <el-button v-auth="AUTH_ROLE_UPDATE" type='text' onClick={() => this.editRole(row)}>{ this.$t('common.edit') }</el-button>
              <el-button v-auth="AUTH_ROLE_DELETE" type='text' style={ [1, 2, 3].includes(row.id) ? {} : { color: '#ed4014' } } disabled={[1, 2, 3].includes(row.id)} onClick={() => this.deleteRole(row)}>{ this.$t('common.delete') }</el-button>
            </div>
          },
        },
      ]
    },
  },
  methods: {
    afterSendHook(data) {
      const codeList = Object.keys(this.$i18n.messages[this.$i18n.locale].privileges)
      data.forEach(item => {
        item.origin_permissions = item.permissions
        item.permissions = item.permissions.filter(_item => _item.is_check).filter(_item => codeList.includes(_item.auth_code))
      })
    },
    addRole() {
      this.$refs.dialog.initDialog()
    },
    editRole(row) {
      this.$refs.dialog.editDialog(row)
    },
    deleteRole(row) {
      this.$refs.dialog.deleteDialog(row)
    },
    async onSubmit() {
      this.$refs.table.search()
      await this.$store.dispatch('moduleUser/setAuth')
    },
    selectInit(row, index) {
      if ([1, 2, 3].includes(row.id)) {
        return false
      } else {
        return true
      }
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
