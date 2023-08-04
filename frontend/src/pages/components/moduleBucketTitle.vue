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
    <el-row>
      <el-col
        :span="12"
      ><h1 class="title-1 color3">{{ bucket }}</h1></el-col>
      <el-col
        :span="12"
        class="txt-r color6"
      ><el-button
        class="border-e5 font14"
        plain
        size="medium"
        @click="goBack"
      >返回列表</el-button></el-col>
    </el-row>
    <div class="bucket_msg color9">
      <slot></slot>
    </div>
  </div>
</template>
<script>
export default {
  props: {
    bucket: {
      type: String,
      default: '',
      required: true,
    },
    goBack: {
      type: Function,
      default: () => ({}),
      required: true,
    },
  },
}
</script>
<style lang="scss" scoped>
.title-1 {
  font-weight: bold;
  margin-top: 8px;
  margin-bottom: 16px;
  font-size: 24px;
}
.border-e5 {
  border: 1px solid #e5e5e5;
}
.bucket_msg {
  margin-bottom: 40px;
}
.txt-r {
  text-align: right;
}
</style>
