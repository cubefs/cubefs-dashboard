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

const RouterViewHoc = {
  name: 'RouterViewHoc',
  render(h) {
    return h('router-view')
  },
}

export const clusterDefaultRoutes = {
  path: 'cluster',
  name: 'cluster',
  meta: {
    title: '集群管理',
  },
  component: RouterViewHoc,
  children: [
    {
      path: 'list',
      name: 'clusterList',
      component: () => import('@/pages/cfs/clusterOverview'),
      meta: {
        title: '集群列表',
      },
    },
  ],
}

export const clusterDetailRoutes = {
  path: 'clusterDetail',
  name: 'clusterDetail',
  meta: {
    title: '集群详情',
  },
  component: RouterViewHoc,
  children: [],
}

export const clusterDetailChildren = [
  {
    path: 'clusterInfo',
    name: 'clusterInfo',
    component: () => import('@/pages/cfs/clusterOverview/clusterInfo'),
    meta: {
      title: '集群概览',
    },
  },
  {
    path: 'volManage',
    name: 'volManage',
    meta: {
      title: '卷管理',
    },
    component: () => import('@/pages/cfs/clusterOverview/clusterInfo/volumn/index.vue'),
  },
  {
    path: 'dataManage',
    name: 'dataManage',
    meta: {
      title: '数据管理',
    },
    component: () => import('@/pages/cfs/clusterOverview/clusterInfo/dataManage/index.vue'),
  },
  {
    path: 'metaDataManage',
    name: 'metaDataManage',
    meta: {
      title: '元数据管理',
    },
    component: () => import('@/pages/cfs/clusterOverview/clusterInfo/metaDataManage/metaPartition.vue'),
  },
  {
    path: 'resourceManage',
    name: 'resourceManage',
    meta: {
      title: '节点管理',
    },
    component: () => import('@/pages/cfs/clusterOverview/clusterInfo/resourceManage/index.vue'),
  },
  {
    path: 'fileManage',
    name: 'fileManage',
    meta: {
      title: '文件管理',
      menuLinkGroup: true,
    },
    component: RouterViewHoc,
    children: [
      {
        path: '',
        name: 'fileManageList',
        meta: {
          title: '卷列表',
        },
        component: () => import('@/pages/cfs/fileManage/index'),
      },
      {
        path: 'fileList',
        name: 'fileList',
        meta: {
          title: '文件上传列表',
        },
        component: () => import('@/pages/cfs/fileManage/fileList/index'),
      },
    ],
  },
  {
    path: 'clusterEvent',
    name: 'clusterEvent',
    meta: {
      title: '集群事件',
    },
    component: () => import('@/pages/cfs/clusterOverview/clusterInfo/events/index.vue'),
  },
]

export default [
  clusterDefaultRoutes,
  clusterDetailRoutes,
]
