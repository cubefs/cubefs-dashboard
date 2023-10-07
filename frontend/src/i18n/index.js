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

/* ------------------------------------------------------ */
/* ------- auto generated code by vue-i18n-transform------- */
import Vue from 'vue'
import locale from 'element-ui/lib/locale'
import language from './lang'
import VueI18n from 'vue-i18n'

Vue.prototype.isRtl = localStorage.getItem('language') === 'zh'

Vue.use(VueI18n)

export const i18n = new VueI18n({
  locale: localStorage.getItem('language') || 'zh',
  fallbackLocale: localStorage.getItem('language') || 'zh',
  messages: language,
  formatFallbackMessages: true,
  silentTranslationWarn: true, // 浏览器控制台不显示warning
})

locale.i18n((key, value) => i18n.t(key, value))

export default i18n
