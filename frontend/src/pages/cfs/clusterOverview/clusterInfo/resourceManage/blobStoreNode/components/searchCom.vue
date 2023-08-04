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
    <el-form :inline="true" :model="forms">
      <el-form-item v-if="hasStatus" label="状态：">
        <el-select v-model="forms.status" placeholder="请选择" clearable style="width: 180px">
          <el-option
            v-for="item in statusList"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          ></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="机房：">
        <el-select v-model="forms.idc" placeholder="请选择" clearable style="width: 180px">
          <el-option
            v-for="item in ebsIdcList"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          ></el-option>
        </el-select>
      </el-form-item>
      <el-form-item v-if="hasHost" label="主机：">
        <el-input v-model="forms.host" placeholder="请输入" clearable style="width: 180px" @clear="clearHost">
        </el-input>
      </el-form-item>
      <el-button
        type="primary"
        icon="el-icon-search"
        @click="searchClick"
      >搜索</el-button>
    </el-form>
  </div>
</template>
<script>
import { mapGetters } from 'vuex'
import { nodeStatusList } from '@/pages/cfs/status.conf'
export default {
  components: {
  },
  props: {
    hasStatus: {
      type: Boolean,
      default: false,
    },
    hasHost: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      forms: {
        cluster: '',
        idc: '',
        status: '',
        host: '',
      },
    }
  },
  computed: {
    clusterInfo() {
      return JSON.parse(sessionStorage.getItem('clusterInfo'))
    },
    ...mapGetters('clusterInfoModule', ['ebsIdcList']),
  },
  created() {
    this.initData()
    this.searchClick()
  },
  mounted() {
    // this.searchClick()
  },
  methods: {
    clearHost() {
      this.$router.push({
        name: 'nodeManage',
      })
    },
    initData() {
      this.statusList = nodeStatusList
    },
    searchClick(flag = false) {
      this.forms.cluster = this.clusterInfo.name
      if (!this.hasStatus) {
        delete this.forms.status
      }
      if (!this.hasHost) {
        delete this.forms.host
      }
      this.$emit('set-forms', this.forms, flag)
    },
  },
}
</script>
<style lang="scss" scoped>
</style>
