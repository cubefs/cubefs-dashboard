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

const moduleCluster = {
  namespaced: true,
  state: () => ({
    clusterName: '',
    leaderAddr: '',
    masterAddr: [],
    cli: '',
    ebsClusterInfo: [],
    clusterInfo: {},
    domainName: ''
  }),
  mutations: {
    setClusterInfo(state, payload) {
      const { clusterName, leaderAddr, masterAddr, cli, domainName, clusterInfo, ebsClusterInfo } = payload
      state.clusterName = clusterName
      state.leaderAddr = leaderAddr
      state.masterAddr = masterAddr
      state.cli = cli
      state.domainName = domainName
      state.clusterInfo = clusterInfo
      state.ebsClusterInfo = ebsClusterInfo
      sessionStorage.setItem('clusterInfo', JSON.stringify({ ...payload }))
    },
    setEbsClusterInfo(state, payload) {
      state.clusterInfo = payload.clusterInfo
      sessionStorage.setItem('clusterInfo', JSON.stringify(payload.clusterInfo))
    },
  },
  actions: {},
  getters: {
    clusterInfog(state) {
      const { clusterName, leaderAddr, masterAddr, cli, domainName, clusterInfo, ebsClusterInfo } = state
      if (clusterName && leaderAddr) {
        return {
          clusterName, leaderAddr, masterAddr, cli, domainName, clusterInfo, ebsClusterInfo
        }
      } else {
        return JSON.parse(sessionStorage.getItem('clusterInfo'))
      }
    },
    ebsClusterInfog(state) {
      const { ebsClusterInfo, clusterName } = state
      if (clusterName) {
        return ebsClusterInfo;
      } else {
        return JSON.parse(sessionStorage.getItem('clusterInfo')).ebsClusterInfo
      }
    },
    ebsIdcList(state) {
      const clusterInfo = state.clusterInfo || JSON.parse(sessionStorage.getItem('clusterInfo'))
      return clusterInfo.idc.split(',').map(item => ({ label: item, value: item }))
    },
  },
}
export default moduleCluster
