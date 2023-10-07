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
    :title="title"
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
      <el-form-item :label="$t('common.cluster')+ ':'" prop="volName">
        <el-input v-model="forms.volName" disabled class="input"></el-input>
      </el-form-item>
      <el-form-item :label="$t('volume.cursize') + ':'" prop="size">
        <el-input v-model="forms.size" disabled class="input"></el-input>&nbsp;
        GB
      </el-form-item>
      <el-form-item :label="$t('volume.tarsize') + ':'" prop="targetSize">
        <el-input v-model.number="forms.targetSize" class="input"></el-input>&nbsp; GB
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button ref="pol" type="primary" @click="doCheck">{{ $t('button.submit') }}</el-button>
      <el-button ref="pol" type="primary" @click="close">{{ $t('button.cancel') }}</el-button>
    </div>
  </el-dialog>
</template>
<script>
import { expandVol, shrinkVol } from '@/api/cfs/cluster'
import { getGBSize_B } from '@/utils'
import Mixin from '@/pages/cfs/clusterOverview/mixin'
export default {
  mixins: [Mixin],
  data() {
    return {
      type: 'expansion',
      forms: {
        volName: '',
        size: '',
        targetSize: '',
      },
      dialogFormVisible: false,
      zoneList: [
        {
          text: this.$t('volume.dgzx'),
          value: 'tests',
        },
      ],
    }
  },
  computed: {
    title() {
      return this.type === 'expansion' ? this.$t('common.scaleup') : this.$t('common.down')
    },
    rules() {
      return {
        volName: [
          {
            required: true,
            message: this.$t('volume.inputvolume'),
            trigger: 'blur',
          },
        ],
        targetSize: [
          {
            required: true,
            trigger: 'blur',
            validator: (rule, value, callback) => {
              const isError =
                this.type === 'expansion'
                  ? +value < +this.forms.size
                  : +value > +this.forms.size
              if (!value) {
                callback(new Error(this.$t('volume.inputsize')))
              } else if (isError) {
                callback(
                  new Error(
                    `${
                      this.type === 'expansion' ? this.$t('volume.inputsize2less') : this.$t('volume.inputsize2greater')
                    }`,
                  ),
                )
              } else {
                callback()
              }
            },
          },
        ],
      }
    },
  },
  methods: {
    initForm(val) {
      // 初始化默认值和单位
      const { size } = val
      this.forms = { ...this.forms, ...val }
      this.forms.targetSize = getGBSize_B(size)
      this.forms.size = getGBSize_B(size)
    },
    open(type = 'expansion') {
      this.dialogFormVisible = true
      this.type = type
    },
    clearData() {
      this.forms = {
        volName: '',
        size: '',
        targetSize: '',
      }
    },
    async doCheck() {
      await this.$refs.form.validate()
      const { volName, targetSize } = this.forms
      const params = {
        name: volName,
        capacity: +targetSize || 0,
        cluster_name: this.clusterName,
      }
      if (this.title === this.$t('common.scaleup')) {
        await expandVol(params)
      } else {
        await shrinkVol(params)
      }
      this.$message.success(this.title + this.$t('common.xxsuc'))
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
  width: 70%;
}
.dialog-footer {
  text-align: center;
}
</style>
