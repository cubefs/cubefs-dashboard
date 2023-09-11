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
    <el-card>
      <el-tabs v-model="activeName">
        <el-tab-pane :label="$t('common.cluster') + $t('common.list')" name="first">
          <el-button
            v-auth="'CLUSTER_CREATE'"
            class="fl-rt mg-bt-s"
            type="primary"
            @click="showDialog('add')"
          >{{ $t('common.add') }}{{ $t('common.cluster') }}</el-button>
          <o-page-table :columns="tableColumns" :form-data="formData" :data="tableData"></o-page-table>
        </el-tab-pane>
      </el-tabs>
      <el-dialog
        :title="`${dataId ? $t('common.edit') : $t('common.add')}` + $t('common.cluster')"
        :visible.sync="dialogFormVisible"
        width="800px"

        @closed="clearData"
      >
        <o-form
          ref="form"
          class="mid-block"
          :form-value.sync="formValue"
          :form-list="formList"
          :label-width="200"
          content-width="70%"
          :options="{
            'validate-on-rule-change': false
          }"
          not-ctrl
        ></o-form>
        <div slot="footer" class="dialog-footer">
          <el-button ref="pol" type="primary" @click="doCheck">{{ $t('button.submit') }}</el-button>
          <el-button ref="pol" type="primary" @click="close">{{ $t('button.cancel') }}</el-button>
        </div>
      </el-dialog>
    </el-card>
  </div>
</template>
<script>
import { getClusterList, upDateCluster, createCluster } from '@/api/cfs/cluster'
import { getClusterList as getEbsClusterList } from '@/api/ebs/ebs'
import { initCfsClusterRoute } from '@/router/index'
import { mapMutations } from 'vuex'
export default {
  name: '',
  components: {},
  data() {
    return {
      activeName: 'first',
      dataId: '',
      formData: {},
      formValue: {
        consul_addr: '',
      },
      tableData: [],
      dialogFormVisible: false,
    }
  },
  computed: {
    tableColumns() {
      return [
        {
          title: this.$t('cfsclusteroverview.clustername'),
          key: 'name',
          render: (h, { row }) => {
            return (<a
              class="primary-color"
              onClick={() => {
                this.toDetail(row)
              }}
            >{row.name}</a>)
          },
        },
        {
          title: this.$t('cfsclusteroverview.metasize'),
          key: 'meta_total',
        },
        {
          title: this.$t('cfsclusteroverview.metaused'),
          key: 'meta_used',
        },
        {
          title: this.$t('cfsclusteroverview.metausage'),
          key: 'meta_used_ratio',
        },
        {
          title: this.$t('cfsclusteroverview.datasize'),
          key: 'data_total',
        },
        {
          title: this.$t('cfsclusteroverview.dataused'),
          key: 'data_used',
        },
        {
          title: this.$t('cfsclusteroverview.datausage'),
          key: 'data_used_ratio',
        },
        {
          title: this.$t('cfsclusteroverview.masteraddress'),
          key: 'domain',
        },
        {
          title: this.$t('cfsclusteroverview.s3address'),
          key: 's3_endpoint',
        },
        {
          title: this.$t('cfsclusteroverview.tag'),
          key: 'tag',
        },
        {
          title: this.$t('common.action'),
          render: (h, { row }) => {
            return (
              <div>
                <el-button
                  v-auth="CLUSTER_UPDATE"
                  type="text"
                  size="medium"
                  class="ft-16"
                  icon="el-icon-edit-outline"
                  title={ this.$t('common.edit') }
                  onClick={() => {
                    this.showDialog(row)
                  }}>
                </el-button>
              </div>
            )
          },
        },
      ]
    },
    formList() {
      const consulAddr = {
        title: this.$t('cfsclusteroverview.blobstoreaddr'),
        key: 'consul_addr',
        type: 'input',
        rule: {
          required: true,
          trigger: 'blur',
          validator: (rule, value, cb) => {
            const regIp =
                  /^https?:\/\/(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\:([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-5]{2}[0-3][0-5])$/
            if (!value) {
              cb(new Error(this.$t('cfsclusteroverview.inputbsaddr')))
            } else {
              const ipArr = value.split('\n').filter((item) => {
                return !!item
              })
              for (const ip of ipArr) {
                if (!regIp.test(ip)) {
                  cb(new Error(this.$t('cfsclusteroverview.illegalip')))
                  break
                }
              }
              cb()
            }
            cb()
          },
        },
        props: {
          type: 'textarea',
          placeholder: this.$t('cfsclusteroverview.bsaddrrule'),
          autosize: {
            minRows: 2,
            maxRows: 3,
          },
        },
      }
      return {
        children: [
          {
            title: this.$t('cfsclusteroverview.clustername'),
            key: 'name',
            type: 'input',
            rule: {
              required: true,
              message: this.$t('cfsclusteroverview.inputclustername'),
              trigger: 'blur',
            },
            props: {
              placeholder: !this.dataId ? this.$t('cfsclusteroverview.inputclustername') : '',
            },
          },
          {
            title: this.$t('cfsclusteroverview.masteraddress'),
            key: 'master_addr',
            type: 'input',
            rule: {
              required: true,
              // message: '请输入Master地址',
              trigger: 'blur',
              validator: (rule, value, cb) => {
                const regIp =
                  /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\:([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-5]{2}[0-3][0-5])$/
                if (!value) {
                  cb(new Error(this.$t('cfsclusteroverview.inputmasteraddr')))
                } else {
                  const ipArr = value.split('\n').filter((item) => {
                    return !!item
                  })
                  for (const ip of ipArr) {
                    if (!regIp.test(ip)) {
                      cb(new Error(this.$t('cfsclusteroverview.illegalip')))
                      break
                    }
                  }
                  cb()
                }
                cb()
              },
            },
            props: {
              type: 'textarea',
              placeholder: this.$t('cfsclusteroverview.inputmasteraddrrule'),
              autosize: {
                minRows: 2,
                maxRows: 3,
              },
            },
          },
          {
            title: this.$t('cfsclusteroverview.attacheccluster'),
            key: 'vol_type',
            renderContent: (h, item, formData) => {
              return (
                <el-radio-group v-model={formData.vol_type} disabled={Boolean(this.dataId)}>
                  <el-radio label={1}>{ this.$t('common.yes') }</el-radio>
                  <el-radio label={0}>{ this.$t('common.no') }</el-radio>
                </el-radio-group>
              )
            },
            rule: {
              required: true,
              message: this.$t('cfsclusteroverview.volumetype'),
              trigger: 'change',
            },
            defaultValue: 0,
          },
          this.formValue.vol_type === 1 ? consulAddr : undefined,
          {
            title: this.$t('common.idc'),
            key: 'idc',
            type: 'input',
            rule: {
              required: true,
              message: this.$t('cfsclusteroverview.inputidc'),
              trigger: 'blur',
            },
            props: {
              placeholder: !this.dataId ? this.$t('cfsclusteroverview.inputidc') : '',
            },
          },
          {
            title: 'cli:',
            key: 'cli',
            type: 'input',
            rule: {
              required: false,
              message: this.$t('cfsclusteroverview.inputcli'),
              trigger: 'blur',
            },
            props: {
              placeholder: !this.dataId ? this.$t('cfsclusteroverview.inputcli') : '',
            },
          },
          {
            title: this.$t('cfsclusteroverview.masteraddress'),
            key: 'domain',
            type: 'input',
            rule: {
              required: false,
              message: this.$t('cfsclusteroverview.inputmasteraddr'),
              trigger: 'blur',
            },
            props: {
              placeholder: !this.dataId ? this.$t('cfsclusteroverview.inputmasteraddr') : '',
            },
          },
          {
            title: 's3 endpoint:',
            key: 's3_endpoint',
            type: 'input',
            props: {
              type: 'textarea',
              placeholder: this.$t('cfsclusteroverview.s3addrrule'),
              autosize: {
                minRows: 2,
                maxRows: 3,
              },
            },
            rule: {
              trigger: 'blur',
              validator: (rule, value, cb) => {
                const regIp =
                  /^https?:\/\/(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\:([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-5]{2}[0-3][0-5])$/
                  const regDomain = /^(?=^.{3,255}$)http(s)?:\/\/?(www\.)?[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+(:\d+)*(\/\w+\.\w+)*$/
                if (!value) {
                  cb()
                } else {
                  const ipArr = value.split('\n').filter((item) => {
                    return !!item
                  })
                  for (const ip of ipArr) {
                    if (!regIp.test(ip) && !regDomain.test(ip) ) {
                      cb(new Error(this.$t('cfsclusteroverview.illegaladdr')))
                      break
                    }
                  }
                  cb()
                }
                cb()
              },
            },
          },
        ],
      }
    },
  },
  watch: {},
  created() {
    this.getClusterList()
  },
  beforeMount() { },
  mounted() { },
  methods: {
    ...mapMutations('clusterInfoModule', ['setClusterInfo']),
    async getClusterList() {
      const res = await getClusterList()
      this.tableData = res.data.clusters || []
    },
    showDialog(type) {
      this.dialogFormVisible = true
      if (type !== 'add') {
        const {
          name,
          idc,
          master_addr: masterAddr,
          id,
          cli,
          domain,
          vol_type: volType,
          consul_addr: consulAddr,
          s3_endpoint: s3Endpoint,
        } = type
        this.initForm(
          {
            name,
            idc: idc,
            consul_addr: consulAddr,
            master_addr: masterAddr.join('\n'),
            cli,
            domain,
            vol_type: volType,
            s3_endpoint: s3Endpoint,
          },
          id,
        )
      }
    },
    async doCheck() {
      await this.$refs.form.validate()
      const {
        name,
        master_addr: masterAddr,
        idc, cli, domain,
        vol_type: volType,
        s3_endpoint: s3Endpoint,
        consul_addr: consulAddr
      } = this.formValue
      const params = {
        s3_endpoint: s3Endpoint,
        consul_addr: consulAddr,
        name,
        master_addr: masterAddr.split('\n').filter((item) => {
          return !!item
        }),
        idc,
        cli,
        domain,
        vol_type: volType,
      }
      let publishCluster = createCluster
      if (this.dataId) {
        params.id = this.dataId
        publishCluster = upDateCluster
      }
      await publishCluster(params)
      this.$message.success(`${this.dataId ? this.$t('common.edit') : this.$t('common.add')}this.$t('common.xxsuc')`)
      this.getClusterList()
      this.close()
    },
    initForm(val, id) {
      this.formValue = { ...val }
      this.dataId = id
    },
    handleClick() {

    },
    clearData() {
      this.dataId = null
      this.$refs.form.reset()
    },
    close() {
      this.clearData()
      this.dialogFormVisible = false
    },
    async getEbsClusterList(region) {
      const res = await getEbsClusterList({
        region,
      })
      return res.data
    },
    async toDetail(item) {
      let ebsClusterInfo = []
      try {
        const data = await this.getEbsClusterList(item.name)
        ebsClusterInfo = data
      } catch (e) {

      } finally {
        this.setClusterInfo({
          clusterName: item.name,
          masterAddr: item.master_addr,
          leaderAddr: item.leader_addr,
          cli: item.cli,
          domainName: item.domain,
          clusterInfo: item,
          ebsClusterInfo,
        })
        initCfsClusterRoute()
        this.$router.push({
          name: 'clusterInfo',
        })
      }
    },
  },
}
</script>
<style lang='scss' scoped>

</style>
