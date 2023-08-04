/**
 * Copyright 2023 The CubeFS Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
 * implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

import DashboardLayout from './DashboardLayout'
const layouts = {
  DashboardLayout,
}
export function layoutMixin(Vue) {
  // 全局注册组件
  Object.values(layouts).forEach(layout => {
    Vue.component(layout.name, layout)
  })
}
export const Common = {
  components: {
    DashboardLayout,
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
  render(h) {
    const component = 'DashboardLayout'
    return h(
      component,
      {
        props: this.$props,
      },
    )
  },
}
