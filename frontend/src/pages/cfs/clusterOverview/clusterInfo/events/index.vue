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
  <div v-auth="'BLOBSTORE_CONFIG_LIST'">
    <div class="mg-bt-s">
      <div>
        <el-radio-group v-model="activeName" style="margin-bottom: 10px;">
          <el-radio-button label="blobEvent" :disabled="!ebsClusterList || !ebsClusterList.length">纠删码后台任务</el-radio-button>
        </el-radio-group>
      </div>
    </div>
    <el-card>
      <blobEvent v-if="activeName === 'blobEvent'" />
    </el-card>
  </div>
</template>
<script>
import blobEvent from './components/blobStoreEvent.vue'
import mixin from '@/pages/cfs/clusterOverview/mixin'
export default {
  name: '',
  components: { blobEvent },
  mixins: [mixin],
  data () {
    return {
      activeName: '',
    }
  },
  computed: {},
  watch: {},
  created() {
    if (this.ebsClusterList && this.ebsClusterList.length) {
      this.activeName = 'blobEvent'
    }
  },
  beforeMount() {},
  mounted() {},
  methods: {},
}
</script>
<style lang='scss' scoped>
</style>
