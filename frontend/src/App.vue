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
  <div id="app">
    <Noop v-if="layout === 'Noop'"></Noop>
    <Common v-else app-name="CubeFS控制台" :hide-psa-side-bar="false"></Common>
  </div>
</template>
<script>
import { Common } from '@/layouts/Common'
import Noop from '@/layouts/Noop.vue'
export default {
  name: 'App',
  components: {
    Noop,
    Common,
  },
  computed: {
    layout() {
      return this.$route.meta.layout || 'Common'
    },
  },
}
</script>

<style lang="scss" src="@/assets/styles/index.scss"></style>
