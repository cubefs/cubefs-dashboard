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
    width="950px"
  >
    <el-form ref="form" :model="form" label-width="120px" :rules="rules" style="width: 900px">
      <el-form-item :label="$t('common.role') + $t('common.code')" prop="role_code">
        <el-input v-model="form.role_code" :disabled="type !== 'create'" :placeholder="$t('rolemgt.inputcode') + ', ' + $t('rolemgt.inputrule')" style="width: 790px" />
      </el-form-item>
      <el-form-item :label="$t('common.role') + $t('common.name')" prop="role_name">
        <el-input v-model="form.role_name" :disabled="type !== 'create'" :placeholder="$t('rolemgt.inputname')" style="width: 790px" />
      </el-form-item>
      <!-- <div>角色权限</div>
      <div v-for="item in codeList" :key="item.title">
        <el-form-item :label="item.title">
          <el-checkbox-group v-model="form.permission_ids" :disabled="type === 'delete'">
            <el-checkbox
              v-for="_item in item.children"
              :key="_item"
              class="permission_checkbox"
              :label="_item"
              :title="$t(_item)"
            >{{ $t(_item) }}</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </div> -->
      <el-form-item :label="$t('common.role') + ' ' + $t('common.privilege')" prop="roles">
        <el-table :data="codeList" style="width: 97%" :header-cell-style="{ background: '#f8f9fc', color: '#606266', fontWeight: '550' }">
          <el-table-column prop="title" :label="$t('common.modules')" width="120" />
          <el-table-column :label="$t('common.privilege')">
            <template slot-scope="scope">
              <el-checkbox
                v-for="item in scope.row.children"
                :key="item.id"
                v-model="item.checked"
                :title="$t('privileges.'+item.id)"
                :disabled="type === 'delete'"
                @change="value => {
                  itemSelect(value, item, scope.row)
                }"
              >
                <el-tooltip v-if="item.id === 'CFS_USERS_CREATE' || item.id === 'CFS_VOLS_CREATE'" placement="top" :content="$t('rolemgt.tenantcheckedfirst')">
                  <div>
                    <i class="el-icon-question" /> {{ $t('privileges.'+item.id) }}
                  </div>
                </el-tooltip>
                <span v-else>
                  <!-- {{ $t(item.id) }} -->
                  {{ $t('privileges.'+item.id) }}
                </span>
              </el-checkbox>
            </template>
          </el-table-column>
          <el-table-column :label="$t('common.all')" width="100">
            <template slot-scope="scope">
              <el-checkbox
                v-model="scope.row.all"
                :disabled="type === 'delete'"
                @change="
                  value => {
                    allSelect(value, scope.row)
                  }
                "
              />
            </template>
          </el-table-column>
        </el-table>
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
  getAuthList,
  roleCreate,
  roleUpdate,
  roleDelete,
} from '@/api/auth'
import { backendAuthids, getCodeList } from '@/utils/auth'

export default {
  name: 'UserDialog',
  data() {
    return {
      dialogVisible: false,
      id: '',
      type: '',
      form: {
        role_code: '',
        role_name: '',
        permission_ids: [],
      },
      authList: [],
      rules: {
        role_code: [
          { required: true, message: this.$t('rolemgt.inputcode'), trigger: 'blur' },
          { pattern: /^[a-zA-Z0-9_]+$/, message: this.$t('rolemgt.inputrule'), trigger: 'blur' },
        ],
        role_name: [{ required: true, message: this.$t('rolemgt.inputname'), trigger: 'blur' }],
      },
      codeList: [],
      codeMap: {},
    }
  },
  computed: {
    title() {
      let title = ''
      switch (this.type) {
        case 'create':
          title = this.$t('common.add') + this.$t('common.role')
          break
        case 'edit':
          title = this.$t('common.edit') + this.$t('common.privilege')
          break
        case 'delete':
          title = this.$t('common.delete') + this.$t('common.role')
          break
      }
      return title
    },
  },
  mounted() {
    this.getAuthList()
  },
  methods: {
    initDialog() {
      this.type = 'create'
      this.form = {
        role_code: '',
        role_name: '',
        permission_ids: [],
      }
      this.codeList = getCodeList(this).map(item => ({
        title: item.title,
        children: item.children.map(_item => ({
          checked: false,
          id: _item,
        })),
        all: false,
      }))
      this.dialogVisible = true
    },
    async editDialog(roleInfo) {
      // eslint-disable-next-line camelcase
      const { id, role_code, role_name, origin_permissions } = roleInfo
      this.form = {
        role_code,
        role_name,
        permission_ids: origin_permissions.map(item => item.id),
      }
      this.echoCode()
      this.id = id
      this.type = 'edit'
      this.dialogVisible = true
    },
    async deleteDialog(roleInfo) {
      // eslint-disable-next-line camelcase
      const { id, role_code, role_name, origin_permissions } = roleInfo
      this.form = {
        role_code,
        role_name,
        permission_ids: origin_permissions.map(item => item.id),
      }
      this.echoCode()
      this.id = id
      this.type = 'delete'
      this.dialogVisible = true
    },
    async getAuthList() {
      const { data: { permissions } } = await getAuthList({ page_size: 500 })
      this.codeMap = permissions.reduce((acc, cur) => {
        acc[cur.auth_code] = cur.id
        return acc
      }, {})
    },
    echoCode() {
      this.codeList = getCodeList(this).map(item => ({
        title: item.title,
        children: item.children.map(_item => ({
          checked: false,
          id: _item,
        })),
        all: false,
      }))
      this.codeList.forEach(item => {
        let flag = true
        item.children.forEach(_item => {
          if (this.form.permission_ids.findIndex(id => id === this.codeMap[_item.id]) !== -1) {
            _item.checked = true
          } else {
            flag = false
          }
        })
        if (flag) {
          item.all = true
        }
      })
    },
    deleteCode(code) {
      const index = this.form.permission_ids.findIndex(item => item === this.codeMap[code])
      if (index !== -1) {
        this.form.permission_ids.splice(index, 1)
      }
    },
    addCode(code) {
      const index = this.form.permission_ids.findIndex(item => item === this.codeMap[code])
      if (index === -1) {
        this.form.permission_ids.push(this.codeMap[code])
      }
    },
    itemSelect(value, item, row) {
      if (value) {
        this.addCode(item.id)
      } else {
        this.deleteCode(item.id)
      }
      row.all = !row.children.some(item => !item.checked)
    },
    allSelect(value, { children }) {
      children.forEach(item => {
        if (value) {
          this.addCode(item.id)
        } else {
          this.deleteCode(item.id)
        }
        item.checked = value
      })
    },
    async submit() {
      await this.$refs.form.validate()
      if (this.type === 'create') {
        await roleCreate({
          ...this.form,
          permission_ids: [...this.form.permission_ids, ...backendAuthids],
          password: 'abcd1234',
        })
        this.$message.success(this.$t('common.created') + this.$t('common.xxsuc'))
        this.$emit('submit')
      } else if (this.type === 'edit') {
        await roleUpdate({
          ...this.form,
          id: this.id,
        })
        this.$message.success(this.$t('common.edit') + this.$t('common.xxsuc'))
        this.$emit('submit')
      } else if (this.type === 'delete') {
        await roleDelete({
          ids: [this.id],
        })
        this.$message.success(this.$t('common.deleted') + this.$t('common.xxsuc'))
        this.$emit('submit')
      }
      this.dialogVisible = false
    },
  },
}

</script>
<style scoped lang='scss'>
::v-deep .el-checkbox {
  margin: 0;
  display: inline-block;
  width: 175px;
  box-sizing: border-box;
  .el-checkbox__label {
    max-width: 145px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    vertical-align: middle;
  }
}
</style>
