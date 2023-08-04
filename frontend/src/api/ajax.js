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
const base = () => process.env.NODE_ENV === 'production' ? window.location.origin : ''
export default {
  async get(url, params, opt, flag = false) {
    return Vue.prototype.$ajax.get(base() + url, params, opt)
  },
  async post(url, params, opt, flag = false) {
    return Vue.prototype.$ajax.post(base() + url, params, opt)
  },
  async delete(url, params, opt) {
    return Vue.prototype.$ajax.delete(base() + url, params, opt)
  },
  async put(url, params, opt) {
    return Vue.prototype.$ajax.put(base() + url, params, opt)
  },
}
