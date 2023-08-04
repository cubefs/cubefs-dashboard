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

import Ajax from '../ajax'
import { AES, enc, pad, mode } from 'crypto-js'

const key = enc.Utf8.parse('6&b#Yfc&94jGLpc4(BDqLx2]8.O4PLe.')
const iv = enc.Utf8.parse('6&b#Yfc&94jGLpc4')

export const encryption = (text) => AES.encrypt(text, key, { iv, mode: mode.CBC, padding: pad.Pkcs7 }).toString()

export const decrypted = (text) => AES.decrypt(text, key, { iv, mode: mode.CBC, padding: pad.Pkcs7 }).toString(enc.Utf8)

const base = '/api/cubefs/console/auth'

export const userLogin = (param) => {
  const newParam = JSON.parse(JSON.stringify(param))
  newParam.password = encryption(newParam.password)
  return Ajax.post(base + '/login', newParam)
}

export const userLogout = (param) => {
  return Ajax.post(base + '/logout', param)
}

export const selfPasswordUpdate = (param) => {
  const newParam = JSON.parse(JSON.stringify(param))
  newParam.old_password = encryption(newParam.old_password)
  newParam.new_password = encryption(newParam.new_password)
  return Ajax.put(base + '/user/password/self/update', newParam)
}

export const userCreate = (param) => {
  const newParam = JSON.parse(JSON.stringify(param))
  if (newParam.password) {
    newParam.password = encryption(newParam.password)
  }
  return Ajax.post(base + '/user/create', newParam)
}

export const userList = (param) => {
  return Ajax.get(base + '/user/list', param)
}

export const userAuth = (param) => {
  return Ajax.get(base + '/user/permission', param)
}

export const userUpdate = (param) => {
  return Ajax.put(base + '/user/update', param)
}

export const userDelete = (param) => {
  return Ajax.delete(base + '/user/delete', {}, { data: param })
}

export const getAuthList = (param) => {
  return Ajax.get(base + '/permission/list', param)
}

export const getRoleList = (param) => {
  return Ajax.get(base + '/role/list', param)
}

export const roleCreate = (param) => {
  return Ajax.post(base + '/role/create', param)
}

export const roleUpdate = (param) => {
  return Ajax.put(base + '/role/update', param)
}

export const roleDelete = (param) => {
  return Ajax.delete(base + '/role/delete', {}, { data: param })
}

export const passwordUpdate = (param) => {
  const newParam = JSON.parse(JSON.stringify(param))
  newParam.password = encryption(newParam.password)
  return Ajax.put(base + '/user/password/update', newParam)
}
