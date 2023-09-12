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
  <div class="container">
    <div class="language">
          <el-select v-model="language" @change="changeLang">
            <el-option value="zh" label="简体中文"></el-option>
            <el-option value="en" label="English"></el-option>
          </el-select>
    </div>
    <Logo class="logo" />
    <div class="description">
      Infinite Storage, Limitless Value
    </div>
    <div v-if="!isSignup" class="login">
      <el-form
        ref="loginForm"
        key="login"
        :model="loginForm"
        :rules="loginRules"
        label-width="0px"
        style="width: 240px; margin-top: 45px;"
      >
        <div class="label">
          {{ $t('common.account') }}
          <el-tooltip class="item" effect="dark" :content="$t('usermgt.accountrule')" placement="right">
            <i class="el-icon-question"></i>
          </el-tooltip>
        </div>
        <el-form-item prop="user_name">
          <el-input v-model="loginForm.user_name" :placeholder="$t('usermgt.inputaccount')" />
        </el-form-item>
        <div class="label">
          {{ $t('common.passwd') }}
          <el-tooltip class="item" effect="dark" placement="right">
            <div slot="content">{{ $t('usermgt.passwdlen') }}<br />{{ $t('usermgt.passwdrule') }}</div>
            <i class="el-icon-question"></i>
          </el-tooltip>
        </div>
        <el-form-item prop="password">
          <el-input v-model="loginForm.password" type="password" :placeholder="$t('usermgt.inputpasswd')" />
        </el-form-item>
      </el-form>
      <div class="btn-group">
        <div class="signup-btn" @click="goToSignup">{{ $t('common.signup') }}</div>
        <div class="signup-btn password-btn" @click="goToPassword">{{ $t('common.changepasswd') }}</div>
        <el-button v-loading="loading" class="login-btn" :class="{ active: canLogin }" @click="login">{{ $t('button.signin') }}</el-button>
      </div>
    </div>
    <div v-if="isSignup" class="signup">
      <i class="el-icon-back back-icon" @click="goToLogin"></i>
      <el-form
        ref="signupForm"
        key="signup"
        :model="signupForm"
        label-width="0px"
        :rules="signupRules"
        style="width: 240px; margin-top: 45px;"
      >
        <div class="label">
          {{ $t('common.account') }}
          <el-tooltip class="item" effect="dark" :content="$t('usermgt.accountrule')" placement="right">
            <i class="el-icon-question"></i>
          </el-tooltip>
        </div>
        <el-form-item prop="user_name">
          <el-input v-model="signupForm.user_name" autocomplete="off" placeholder="$t('usermgt.inputaccount')" />
        </el-form-item>
        <div v-if="!isPassword">
          <div class="label">
            {{ $t('common.passwd') }}
            <el-tooltip class="item" effect="dark" placement="right">
              <div slot="content">{{ $t('usermgt.passwdlen') }}<br />{{ $t('usermgt.passwdrule') }}</div>
              <i class="el-icon-question"></i>
            </el-tooltip>
          </div>
          <el-form-item prop="password">
            <el-input v-model="signupForm.password" autocomplete="new-password" type="password" :placeholder="$t('usermgt.inputpasswd')" />
          </el-form-item>
          <div class="label">{{ $t('common.phone') }}</div>
          <el-form-item prop="phone">
            <el-input v-model="signupForm.phone" :placeholder="$t('usermgt.inputphone')" />
          </el-form-item>
          <div class="label">{{ $t('common.email') }}</div>
          <el-form-item prop="email">
            <el-input v-model="signupForm.email" :placeholder="$t('usermgt.inputemail')" />
          </el-form-item>
        </div>
        <div v-else>
          <div class="label">
            {{ $t('common.old') }}{{ $('common.passwd') }}
            <el-tooltip class="item" effect="dark" placement="right">
              <div slot="content">{{ $t('usermgt.passwdlen') }}<br />{{ $t('usermgt.passwdrule') }}</div>
              <i class="el-icon-question"></i>
            </el-tooltip>
          </div>
          <el-form-item prop="old_password">
            <el-input v-model="signupForm.old_password" type="password" :placeholder="$t('usermgt.inputoldpasswd')" />
          </el-form-item>
          <div class="label">
            {{ $t('common.new') }}{{ $('common.passwd') }}
            <el-tooltip class="item" effect="dark" placement="right">
              <div slot="content">{{ $t('usermgt.passwdlen') }}<br />{{ $t('usermgt.passwdrule') }}</div>
              <i class="el-icon-question"></i>
            </el-tooltip>
          </div>
          <el-form-item prop="new_password">
            <el-input v-model="signupForm.new_password" autocomplete="new-password" type="password" :placeholder="$t('usermgt.inputnewpasswd')" />
          </el-form-item>
        </div>
      </el-form>
      <el-button v-loading="signupLoading" class="signup-btn-big" :class="{ active: canSignup }" @click="signup">{{ isPassword ? $t('common.changepasswd') : $t('common.signup') }}</el-button>
    </div>
    <div class="copyright">© 2023 The CubeFS Authors</div>
  </div>
</template>

<script>
import Logo from './logo.vue'
import {
  userLogin,
  userCreate,
  selfPasswordUpdate,
} from '@/api/auth'
export default {
  components: {
    Logo,
  },
  data() {
    const validateUser = (rule, value, cb) => {
      if (this.checkUser(value)) {
        cb()
      } else {
        cb(new Error(this.$t('usermgt.illegalaccount')))
      }
    }
    const validatePassword = (rule, value, cb) => {
      if (this.checkPassword(value)) {
        cb()
      } else {
        cb(new Error(this.$t('usermgt.illegalpasswd')))
      }
    }
    return {
      loginForm: {
        user_name: '',
        password: '',
      },
      signupForm: {
        user_name: '',
        password: '',
        email: '',
        phone: '',
        old_password: '',
        new_password: '',
      },
      isSignup: false,
      isPassword: false,
      loginRules: {
        user_name: [{ validator: validateUser, trigger: 'blur' }],
        password: { validator: validatePassword, trigger: 'blur' },
      },
      signupRules: {
        user_name: [{ validator: validateUser, trigger: 'blur' }],
        password: [{ validator: validatePassword, trigger: 'blur' }],
        phone: [{ required: true, message: this.$t('usermgt.inputphone'), trigger: 'blur' }],
        email: [{ required: true, message: this.$t('usermgt.inputemail'), trigger: 'blur' }],
        old_password: [{ validator: validatePassword, trigger: 'blur' }],
        new_password: [
          { validator: validatePassword, trigger: 'blur' },
        ],
      },
      loading: false,
      signupLoading: false,
      language: this.$i18n.locale,
    }
  },
  computed: {
    canLogin() {
      return this.checkUser(this.loginForm.user_name) && this.checkPassword(this.loginForm.password)
    },
    canSignup() {
      if (!this.isPassword) {
        return this.checkUser(this.signupForm.user_name) && this.checkPassword(this.signupForm.password) && !!this.signupForm.phone && !!this.signupForm.email
      } else {
        return this.checkUser(this.signupForm.user_name) && this.checkPassword(this.signupForm.old_password) && this.checkPassword(this.signupForm.new_password)
      }
    },
  },
  watch: {
    isSignup(newVal) {
      if (newVal) {
        this.loginForm = {
          user_name: '',
          password: '',
        }
      } else {
        this.signupForm = {
          user_name: '',
          password: '',
          email: '',
          phone: '',
          old_password: '',
          new_password: '',
        }
      }
    },
  },
  methods: {
    goToSignup() {
      this.isSignup = true
    },
    goToPassword() {
      this.isSignup = true
      this.isPassword = true
    },
    goToLogin() {
      this.isSignup = false
      this.isPassword = false
    },
    login() {
      if (!this.canLogin || this.loading) {
        return
      }
      this.loading = true
      userLogin(this.loginForm).then(async() => {
        this.$message.success(this.$t('common.login') + this.$t('common.xxsuc'))
        localStorage.setItem('userInfo', JSON.stringify({ user_name: this.loginForm.user_name }))
        await this.$store.dispatch('moduleUser/setAuth')
        this.$router.push('/')
      }).finally(() => {
        this.loading = false
      })
    },
    async signup() {
      if (!this.canSignup || this.loading) {
        return
      }
      await this.$refs.signupForm.validate()
      this.signupLoading = true
      if (!this.isPassword) {
        // eslint-disable-next-line camelcase
        const { old_password, new_password, ...params } = this.signupForm
        userCreate(params).then(() => {
          this.$message.success(this.$t('common.signup') + this.$t('common.xxsuc'))
          this.isSignup = false
        }).finally(() => {
          this.signupLoading = false
        })
      } else {
        // eslint-disable-next-line camelcase
        const { old_password, new_password, user_name } = this.signupForm
        selfPasswordUpdate({ old_password, new_password, user_name }).then(() => {
          this.$message.success(this.$t('common.change') + this.$t('common.passwd') + this.$t('common.xxsuc'))
          this.isSignup = false
        }).finally(() => {
          this.signupLoading = false
        })
      }
    },
    checkUser(str) {
      return /^[0-9A-Za-z_]+$/.test(str)
    },
    checkPassword(str) {
      // 密码只能包含数字、大写字母、小写字母、特殊字符（~!@#$%^&*_.?），且至少两种类型以上
      const regList = [/[0-9]/, /[A-Z]/, /[a-z]/, /[~!@#$%^&*_.?]/]
      let num = 0
      regList.forEach(item => {
        if (item.test(str)) {
          num++
        }
      })
      if (str.length > 7 && str.length < 17 && num > 1) {
        return true
      }
      return false
    },
    changeLang() {
      const langType = this.language
      localStorage.setItem("language", langType)
      this.$i18n.locale = langType
    },
  },
}

</script>
<style scoped lang='scss'>
.container {
  width: 100vw;
  height: 100vh;
  background: url('@/assets/images/login.png');
  display: flex;
  flex-direction: column;
  align-items: center;
}
.logo {
  margin-top: 8vh;
  margin-bottom: 4vh;
}
.description {
  width: 800px;
  font-size: 26px;
  line-height: 36px;
  text-align: center;
  color: #FFFFFF;
  margin-bottom: 4vh;
}
.login {
  width: 540px;
  height: 326px;
  background: rgba(0, 0, 0, 0.33);
  display: flex;
  align-items: center;
  border-radius: 8px;
  flex-direction: column;
}
.signup {
  width: 540px;
  background: rgba(0, 0, 0, 0.33);
  display: flex;
  align-items: center;
  border-radius: 8px;
  flex-direction: column;
  position: relative;
}
.label {
  font-family: 'PingFang SC';
  font-style: normal;
  font-weight: 400;
  font-size: 13px;
  line-height: 22px;
  color: #FFFFFF;
  margin-bottom: 5px;
}
.btn-group {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 300px;
}
.signup-btn {
  margin-top: 10px;
  color: #c7c7c7;
  cursor: pointer;
}
.signup-btn:hover {
  color: #FFFFFF
}
.password-btn {
  margin-left: -20px;
}
.login-btn {
  border: none;
  margin-top: 12px;
  width: 120px;
  background: #CCE6DE;
  border-radius: 4px;
  height: 38px;
  font-weight: 550;
  font-size: 14px;
  line-height: 22px;
  color: #FFFFFF;
}
.signup-btn-big {
  border: none;
  margin-top: 12px;
  margin-bottom: 40px;
  width: 240px;
  background: #CCE6DE;
  border-radius: 4px;
  height: 38px;
  font-weight: 550;
  font-size: 14px;
  line-height: 22px;
  color: #FFFFFF;
}
.active {
  background: #008059;
}

.copyright {
  margin-top: auto;
  margin-bottom: 2vh;
  color: #FFFFFF;
  opacity: 39%;
}
.back-icon {
  color: #c7c7c7;
  cursor: pointer;
  position: absolute;
  top: 10px;
  left: 10px;
  font-size: 24px;
}
.back-icon:hover {
  color: #FFFFFF
}

.language {
  margin-left: auto;
  width: 100px;
  margin-right: 20px;
  margin-top: 20px;
}

::v-deep {
  .language .el-select .el-input__inner{
    color: #fff;
    background: none;
    border: 0px;
  }
}
</style>

