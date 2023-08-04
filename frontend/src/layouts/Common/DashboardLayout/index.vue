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
  <el-container class="main-container">
    <CloudContent v-bind="$props"></CloudContent>
  </el-container>
</template>
<script>
import CloudContent from './CloudContent'
export default {
  name: 'DashboardLayout',
  components: {
    CloudContent,
  },
  props: {
    hideSidebar: {
      type: Boolean,
      default: false,
    },
    appName: {
      type: String,
      default: 'Protal',
    },
  },
}
</script>
<style lang="scss" scoped>

.main-container {
  height: 100%;
  flex-direction: column;
}
</style>
