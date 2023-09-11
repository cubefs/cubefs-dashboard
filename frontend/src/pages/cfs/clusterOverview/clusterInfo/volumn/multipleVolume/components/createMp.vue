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
    :title="$t('common.create')+ 'Mp'"
    :visible.sync="dialogFormVisible"
    width="700px"
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
        <el-input v-model="forms.name" disabled class="input"></el-input>
      </el-form-item>
      <el-form-item :label="$t('volume.metaparsize') + ':'" prop="start">
        <el-input
          v-model.number="forms.start"
          class="input"
          :placeholder="$t('volume.inputmetaparsize')"
        ></el-input>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button type="primary" @click="doCheck">{{ $t('button.submit') }}</el-button>
      <el-button type="primary" @click="close">{{ $t('button.cancel') }}</el-button>
    </div>
  </el-dialog>
</template>
<script>
import { createMetaPartition } from '@/api/cfs/cluster'
import Mixin from '@/pages/cfs/clusterOverview/mixin'
export default {
  mixins: [Mixin],
  data() {
    return {
      userList: [],
      forms: {
        name: '',
        start: '',
      },
      dialogFormVisible: false,
    }
  },
  computed: {
    rules() {
      return {
        start: [
          {
            required: true,
            message: this.$t('volume.inputmetaparsize'),
            trigger: 'blur',
          },
        ],
        name: [
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
    },
    clearData() {
      this.forms = {
        name: '',
        start: '',
      }
    },
    async doCheck() {
      await this.$refs.form.validate()
      const { name, start } = this.forms
      await createMetaPartition({
        name,
        start: +start,
        cluster_name: this.clusterName,
      })
      this.$message.success( this.$t('volume.mpsuc'))
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
  width: 60%;
}
.dialog-footer {
  text-align: center;
}
</style>
