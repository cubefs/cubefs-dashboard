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
    <el-card>
      <o-form
        v-if="formValue"
        :form-value.sync="formValue"
        :form-list="formList"
        not-ctrl
      ></o-form>
      <div class="btn-group fl-rt">
        <el-button type="default" @click="cancle">返回</el-button>
      </div>
    </el-card>
  </div>
</template>
<script>
export default {
  data() {
    return {
      formValue: null,
      formList: {
        children: [
          [
            {
              render(h) {
                return <h3 class="sub-title">基本信息</h3>
              },
            },
          ],
          [
            {
              title: '应用名称',
              type: 'text',
              key: 'appName',
            },
            {
              title: '所属机房',
              type: 'text',
              key: 'codeName',
            },
            {
              title: '镜像版本',
              type: 'text',
              key: 'imageName',
            },
          ],
          [
            {
              title: '集群名称',
              type: 'text',
              key: 'clusterName',
            },
            {
              title: '集群编码',
              type: 'text',
              key: 'clusterCode',
            },
            {
              title: 'config',
              type: 'text',
              key: 'config',
            },
          ],
          [
            {
              render(h) {
                return <h3 class="sub-title">基本信息</h3>
              },
            },
          ],
          [
            {
              title: '灰度环境',
              type: 'text',
              key: 'grayEnv',
            },
            {
              title: '灰度实例数',
              type: 'text',
              key: 'grayReplicas',
            },
            {
              title: '预估TPS峰值',
              type: 'text',
              key: 'maxTps',
            },
          ],
          [
            {
              title: '配置环境',
              type: 'text',
              key: 'configEnv',
            },
            {
              title: '类型',
              type: 'text',
              key: 'usedMode',
            },
            {
              title: 'configMap名称',
              type: 'text',
              key: 'configName',
            },
          ],
        ],
      },
    }
  },
  async created() {
    this.getDetail()
  },
  methods: {
    async getDetail() {
      this.formValue = await this.$ajax.get(
        '/base/getDetail',
        {},
        { _mock: true },
      )
    },
    cancle() {
      this.$router.go(-1)
    },
  },
}
</script>
