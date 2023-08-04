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
  <el-form-item label="集群">
    <el-select
      v-model="currentClusterId"
      style="width: 180px"
      placeholder="请选择"
      :clearable="clearable"
      @change="handleChange"
    >
      <el-option
        v-for="item in ebsClusterList"
        :key="item.cluster_id"
        :label="item.cluster_id"
        :value="item.cluster_id"
      ></el-option>
    </el-select>
  </el-form-item>
</template>

<script>
import { createDataPartition } from '@/api/cfs/cluster'
import mixin from '@/pages/cfs/clusterOverview/mixin'
export default {
  mixins: [mixin],
  props: {
    // eslint-disable-next-line vue/require-default-prop
    value: {
      // eslint-disable-next-line vue/require-prop-type-constructor
      type: Number | String,
    },
    clearable: {
      type: Boolean,
      default: true,
    },
  },
  data() {
    return {
      currentClusterId: this.value,
    }
  },
  watch: {
    value(val) {
      this.currentClusterId = val
    },
  },
  created() {
    this.currentClusterId = this.ebsClusterList[0]?.cluster_id
    this.$emit('input', this.currentClusterId)
  },
  methods: {
    handleChange(val) {
      this.$emit('input', val)
    },
  },
}
</script>
