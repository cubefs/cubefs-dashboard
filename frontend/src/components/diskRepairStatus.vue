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
  <el-progress
    v-if="filterStatus"
    :percentage="percentageInner"
    :width="width"
    :color="[
      { color: '#f56c6c', percentage: 100 },
      { color: '#e6a23c', percentage: 80 },
      { color: '#5cb87a', percentage: 60 },
      { color: '#1989fa', percentage: 40 },
      { color: '#6f7ad3', percentage: 20 },
    ]"
  >
  </el-progress>
</template>
<script>
export default {
  props: {
    percentage: {
      default: 0,
      type: Number,
    },
    status: {
      default: 1,
      type: Number,
    },
    isOffLine: {
      default: false,
      type: Boolean,
    },
    width: {
      default: 126,
      type: Number,
    },
  },
  computed: {
    percentageInner() {
      if (this.isOffLine) {
        return this.percentage || 0
      }
      if (this.status === 2) {
        return 0
      }
      if (this.status === 4) {
        return 100
      }
      if (this.status === 3 || this.percentage) {
        return this.percentage || 0
      }
      return 0
    },
    filterStatus() {
      if (!([2, 3].includes(this.status) || this.isOffLine)) {
        return false
      }
      return true
    },
  },
}
</script>
<style lang="scss" scoped>
::v-deep .el-progress__text {
  font-size: 12px !important;
  display: block;
}
</style>
