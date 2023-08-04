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
    <el-radio-group v-model="activeName" style="margin-bottom: 10px;">
      <el-radio-button label="multiple">多副本卷</el-radio-button>
      <el-radio-button label="cold" :disabled="!ebsClusterList || !ebsClusterList.length">纠删码卷</el-radio-button>
    </el-radio-group>
    <Multiple v-if="activeName === 'multiple'" />
    <Cold v-if="activeName === 'cold'" />
  </div>
</template>

<script>
import Cold from './coldVolume/index.vue'
import Multiple from './multipleVolume/index.vue'
import mixin from '@/pages/cfs/clusterOverview/mixin'
export default {
  name: '',
  components: { Cold, Multiple },
  mixins: [mixin],
  data () {
    return {
      activeName: 'multiple',
    }
  },
  computed: {},
  watch: {},
  created() {},
  beforeMount() {},
  mounted() {},
  methods: {},
}
</script>
<style scoped lang='scss'>
</style>
