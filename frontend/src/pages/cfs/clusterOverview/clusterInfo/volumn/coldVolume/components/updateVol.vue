<template>
  <el-dialog
    :title="$t('common.updata') + $t('common.volume')"
    :visible.sync="dialogFormVisible"
    width="800px"
    :destroy-on-close="true"
    @closed="close"
  >
    <el-form
      ref="form"
      :model="formData"
      :rules="rules"
      label-width="25%"
      class="mid-block"
    >
      <el-form-item prop="CacheCapacity" :label="$t('volume.volumecachesize')">
        <el-input v-model.number="formData.CacheCapacity"></el-input>
      </el-form-item>
      <el-form-item prop="CacheThreshold" label="cacheThreeshold">
        <el-input v-model.number="formData.CacheThreshold"></el-input>
      </el-form-item>
      <el-form-item prop="CacheTtl" :label="$t('volume.volumecachettl')">
        <el-input v-model.number="formData.CacheTtl"></el-input>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button ref="pol" type="primary" :loading="loading" @click="save">{{ $t('button.submit') }}</el-button>
      <el-button ref="pol" type="primary" @click="close">{{ $t('button.cancel') }}</el-button>
    </div>
  </el-dialog>
</template>
<script>
import { updateVol, getVolDetail } from '@/api/cfs/cluster'
import mixin from '../../../mixin'
export default {
  mixins: [mixin],
  data() {
    return {
      dialogFormVisible: false,
      rowDetail: {},
      volDetail: {},
      formData: {
        CacheCapacity: undefined,
        CacheThreshold: undefined,
        CacheTtl: undefined,
      },
      rules: {
        CacheCapacity: [
          { required: true, message: this.$t('common.required'), trigger: 'blur' },
        ],
        CacheThreshold: [
          { required: true, message: this.$t('common.required'), trigger: 'blur' },
        ],
        CacheTtl: [
          { required: true, message: this.$t('common.required'), trigger: 'blur' },
        ],
      },
      loading: false,
    }
  },
  computed: {

  },
  watch: {

  },
  created() {

  },
  methods: {
    close() {
      this.dialogFormVisible = false
    },
    open() {
      this.dialogFormVisible = true
    },
    async init(rowDetail, clusterName) {
      this.clusterName = clusterName
      this.rowDetail = rowDetail
      this.open()
      await this.getVolumnDetail()
      this.formData = {
        CacheCapacity: this.volDetail.CacheCapacity,
        CacheThreshold: this.volDetail.CacheThreshold,
        CacheTtl: this.volDetail.CacheTtl,
      }
    },
    async getVolumnDetail() {
      const res = await getVolDetail({
        name: this.rowDetail.name,
        cluster_name: this.clusterName,
      })
      this.volDetail = res.data || {}
    },
    async save() {
      try {
        await this.$refs.form.validate()
        const params = {
          ...this.formData,
          name: this.rowDetail.name,
          cluster_name: this.clusterName,
        }
        this.loading = true
        await updateVol(params)
        this.loading = false
        this.$message.success(this.$t('common.update') + this.$t('common.xxsuc'))
        this.$emit('refresh')
        this.close()
      } catch (e) {
        this.loading = false
      }
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
