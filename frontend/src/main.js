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

import Vue from 'vue'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import OIcon from '@/components/common/OIcon'
import OPageTable from '@/components/common/OPageTable.vue'
import OForm from '@/components/common/OForm.vue'
import ODropdown from '@/components/common/ODropdown.vue'
import App from './App.vue'
import router, { initMenu } from './router'
import store from './store'
import i18n from './i18n'
import '@/assets/styles/index.scss'

import { ajaxMixin } from './utils/ajax.js'
import { authMixin } from './utils/auth.js'

const components = {
  OIcon,
  OPageTable,
  OForm,
  ODropdown,
}

// 全局注册组件
Object.values(components).forEach(component => {
  Vue.component(component.name, component)
})

Vue.config.productionTip = false
Vue.use(ElementUI, {
  size: 'small',
})

initMenu()

ajaxMixin(Vue, { router })
authMixin(Vue, { router, store })

new Vue({
  router,
  store,
  i18n,
  render: h => h(App),
}).$mount('#app')
