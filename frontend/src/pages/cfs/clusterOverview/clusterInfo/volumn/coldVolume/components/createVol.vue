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
  <el-dialog
    :title="$t('common.create') + $t('common.volume')"
    :visible.sync="dialogFormVisible"
    width="800px"
    @closed="clearData"
  >
    <el-form
      ref="form"
      :model="forms"
      :rules="rules"
      label-width="25%"
      class="mid-block"
    >
      <el-form-item :label="$t('common.cluster') + ':'" prop="clusterName">
        <el-input v-model="forms.clusterName" class="input" disabled></el-input>
      </el-form-item>
      <el-form-item :label="$t('common.volumename') + ':'" prop="volName">
        <el-input
          v-model="forms.volName"
          class="input"
          :placeholder="$t('volume.inputvolume')"
        ></el-input>
        <el-tooltip
          class="item"
          effect="dark"
          :content="$t('volume.volumenamerule')"
          placement="top"
        >
          <i class="el-icon-question fontS16"></i>
        </el-tooltip>
        <!-- <el-checkbox v-model="forms.isCrossZone">跨Zone</el-checkbox>
        <el-checkbox
          v-model="forms.defaultPriority"
          :disabled="errDisable"
        >故障域</el-checkbox> -->
      </el-form-item>
      <el-form-item :label="$t('common.size') + ':'" prop="size">
        <el-input
          v-model.number="forms.size"
          class="input"
          :placeholder="$t('volume.inputsize')"
        ></el-input>&nbsp; GB
      </el-form-item>
      <el-form-item :label="$t('volume.cachesize') + ':'" prop="cache_cap">
        <el-input
          v-model.number="forms.cache_cap"
          class="input"
          :placeholder="$t('volume.inputcachesize')"
        ></el-input>&nbsp; GB
      </el-form-item>
      <el-form-item label="owner:" prop="owner">
        <el-select v-model="forms.owner">
          <el-option v-for="item in userList" :label="item.user_id" :value="item.user_id" :key="item.user_id"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item :label="$t('common.business')+ ':'" prop="work">
        <el-input
          v-model="forms.work"
          type="textarea"
          :rows="2"
          :placeholder="this.$t('volume.inputbusiness')"
          class="input"
        ></el-input>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button ref="pol" type="primary" @click="doCheck">{{ $t('button.submit') }}</el-button>
      <el-button ref="pol" type="primary" @click="close">{{ $t('button.cancel') }}</el-button>
    </div>
  </el-dialog>
</template>
<script>
import { createVol, getClusterIsErrArea, getUserList } from '@/api/cfs/cluster'
import Mixin from '@/pages/cfs/clusterOverview/mixin'
export default {
  mixins: [Mixin],
  data() {
    return {
      clusterList: [],
      forms: {
        clusterName: '',
        isCrossZone: '',
        defaultPriority: '',
        volName: '',
        size: '',
        cache_cap: '',
        work: '',
        owner: '',
      },
      dialogFormVisible: false,
      zoneList: [],
      errDisable: true,
      userList: [],
    }
  },
  computed: {
    rules() {
      return {
        clusterName: [
          {
            required: true,
            message: this.$t('cfsclusteroverview.inputclustername'),
            trigger: 'blur',
          },
        ],
        volName: [
          {
            required: true,
            trigger: 'blur',
            validator: (rule, value, cb) => {
              const reg = /^[a-zA-Z][a-zA-Z_0-9-]*$/
              if (!value) {
                cb(new Error(this.$t('volume.inputvolume')))
              } else if (!reg.test(value)) {
                cb(new Error(this.$t('volume.volumenamerule')))
              }
              cb()
            },
          },
        ],
        size: [
          {
            required: true,
            message: this.$t('volume.inputsize'),
            trigger: 'blur',
          },
        ],
        cache_cap: [
          {
            required: true,
            message: this.$t('volume.inputcachesize'),
            trigger: 'blur',
          },
        ],
        work: [
          {
            required: true,
            message: this.$t('volume.inputbusiness'),
            trigger: 'blur',
          },
        ],
        owner: [
          {
            required: true,
            message: this.$t('volume.selectowner'),
            trigger: 'blur',
          },
        ]
      }
    },
  },
  watch: {
    async 'forms.isCrossZone'(val) {
      if (val) {
        const res = await getClusterIsErrArea({
          cluster_name: this.clusterName,
        })
        if (res?.data?.DomainOn) {
          this.errDisable = false
          return
        }
        this.errDisable = true
      } else {
        this.errDisable = true
      }
    },
  },
  created() {
    this.getClusterList()
    this.queryUserList()
  },
  methods: {
    async queryUserList() {
      const res = await getUserList({
        cluster_name: this.clusterName,
      })
      this.userList = res.data.users || []
    },
    async getClusterList() {
      this.forms.clusterName = this.clusterName
    },
    initForm(val) {
      this.forms = { ...val }
    },
    open() {
      this.dialogFormVisible = true
    },
    clearData() {
      Object.keys(this.forms).forEach((item) => {
        if (item !== 'clusterName') this.forms[item] = ''
      })
    },
    async doCheck() {
      await this.$refs.form.validate()
      const {
        clusterName,
        volName,
        size,
        work,
        owner,
        cache_cap,
      } = this.forms
      await createVol({
        cluster_name: clusterName,
        name: volName,
        cache_cap: +cache_cap || 0,
        capacity: +size || 0,
        business: work,
        owner,
        vol_type: 1,
      })
      this.$message.success(this.$t('common.create') + this.$t('common.xxsuc'))
      this.$emit('refresh')
      this.close()
    },
    close() {
      this.clearData()
      this.dialogFormVisible = false
    },
  },
}
</script>
<style lang="scss" scoped>
.input {
  width: 60%;
}

.dialog-footer {
  text-align: center;
}
</style>
