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

import { Message } from 'element-ui'


function checkValue(arr, key, value) {
  for (let i = 0; i < arr.length; i++) {
    if (arr[i][key] === value) {
      return true
    }
  }
  return false
}

export function authMixin (Vue, { router, store }) {
  router.beforeEach((to, from, next) => {
    if (to.name === 'loginPage') {
      next()
    } else if (!localStorage.getItem('userInfo')) {
      next({ name: 'loginPage' })
    } else {
      const { allAuth } = store.state.moduleUser
      if (!allAuth) {
        store.dispatch('moduleUser/setAuth').then(() => {
          next()
        }).catch(() => {
          next({ name: 'loginPage' })
        })
      } else {
        next()
      }
    }
  })
  Vue.directive('auth', {
    inserted: (el, binding) => {
      const { allAuth } = store.state.moduleUser
      const frontEndAuth = allAuth?.filter(item => item.is_check)
      if (!checkValue(frontEndAuth, 'auth_code', binding.value)) {
        el.parentNode.removeChild(el)
      }
    },
  })
}

export function getCodeList(that) {
  return [
    {
      title: that.$t('auth.cluster'),
      children: ['CLUSTER_CREATE', 'CLUSTER_UPDATE', 'CLUSTER_DELETE', 'CFS_USERS_LIST', 'CFS_USERS_CREATE', 'CFS_USERS_DELETE', 'CFS_USERS_VOLS_TRANSFER'],
    },
    {
      title: that.$t('auth.volume'),
      children: ['CFS_VOLS_CREATE', 'CFS_VOLS_UPDATE', 'CFS_VOLS_EXPAND', 'CFS_VOLS_SHRINK', 'CFS_USERS_POLICIES', 'CFS_USERS_POLICIES_DELETE'],
    },
    {
      title: that.$t('auth.multi_replicas'),
      children: ['CFS_DATANODE_DECOMMISSION', 'CFS_DATANODE_MIGRATE', 'CFS_DISKS_DECOMMISSION', 'CFS_DATAPARTITION_DECOMMISSION'],
    },
    {
      title: that.$t('auth.ec'),
      children: ['BLOBSTORE_NODES_ACCESS', 'BLOBSTORE_DISKS_ACCESS', 'BLOBSTORE_DISKS_SET', 'BLOBSTORE_DISKS_PROBE'],
    },
    {
      title: that.$t('auth.meta'),
      children: ['CFS_METANODE_DECOMMISSION', 'CFS_METANODE_MIGRATE', 'CFS_METAPARTITION_DECOMMISSION'],
    },
    {
      title: that.$t('auth.file'),
      children: ['CFS_S3_DIRS_CREATE', 'CFS_S3_FILES_DOWNLOAD_SIGNEDURL', 'CFS_S3_FILES_UPLOAD_SIGNEDURL'],
    },
    {
      title: that.$t('auth.clusterevents'),
      children: ['BLOBSTORE_CONFIG_SET'],
    },
    {
      title: that.$t('auth.user'),
      children: ['AUTH_USER_CREATE', 'AUTH_USER_UPDATE', 'AUTH_USER_DELETE', 'AUTH_USER_PASSWORD_UPDATE'],
    },
    {
      title: that.$t('auth.role'),
      children: ['AUTH_ROLE_CREATE', 'AUTH_ROLE_UPDATE', 'AUTH_ROLE_DELETE'],
    },
  ]
}

export const backendAuthids = [3, 4, 7, 11, 13, 14, 15, 16, 17, 18, 19, 24, 25, 26, 31, 32, 36, 41, 42, 43, 44, 48, 49, 53, 54, 58, 59, 63, 64, 66, 67, 70, 71]
