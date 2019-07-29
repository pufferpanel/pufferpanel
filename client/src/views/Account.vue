<!--
  - Copyright 2019 Padduck, LLC
  -  Licensed under the Apache License, Version 2.0 (the "License");
  -  you may not use this file except in compliance with the License.
  -  You may obtain a copy of the License at
  -  	http://www.apache.org/licenses/LICENSE-2.0
  -  Unless required by applicable law or agreed to in writing, software
  -  distributed under the License is distributed on an "AS IS" BASIS,
  -  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  -  See the License for the specific language governing permissions and
  -  limitations under the License.
  -->

<template>
  <b-container>
    <b-row v-if="errorMsg">
      <b-col sm="1"/>
      <b-col sm="10">
        <b-alert
          :show="dismissCountDown"
          fade
          dismissible
          variant="danger"
        >{{ errorMsg }}
        </b-alert>
      </b-col>
    </b-row>
    <b-row v-if="successMsg">
      <b-col sm="1"/>
      <b-col sm="10">
        <b-alert
          :show="dismissCountDown"
          fade
          dismissible
          variant="primary"
        >{{ successMsg }}
        </b-alert>
      </b-col>
    </b-row>

    <b-card header-tag="header"
            footer-tag="footer">
      <b-form-group label-cols="4" label-cols-lg="2" label-size="sm" :label="$t('common.Username')"
                    label-for="input-1">
        <b-form-input v-model="username" id="input-1" size="sm"></b-form-input>
      </b-form-group>

      <b-form-group label-cols="4" label-cols-lg="2" label-size="sm" :label="$t('common.Email')" label-for="input-2">
        <b-form-input v-model="email" id="input-2" size="sm"></b-form-input>
      </b-form-group>

      <b-form-group label-cols="4" label-cols-lg="2" label-size="sm" :label="$t('common.ConfirmPassword')"
                    label-for="input-3">
        <b-form-input v-model="confirmPassword" id="input-3" size="sm" type="password"></b-form-input>
      </b-form-group>

      <b-button slot="footer" variant="primary" size="sm" v-text="$t('common.Update')"
                v-bind:disabled="!canSubmitInfoChange" @click="submitInfoChange"></b-button>
    </b-card>


    <!-- -->

    <b-card header-tag="header"
            footer-tag="footer">
      <b-form-group label-cols="4" label-cols-lg="2" label-size="sm" :label="$t('common.OldPassword')"
                    label-for="input-4">
        <b-form-input v-model="oldPassword" id="input-4" size="sm" type="password"></b-form-input>
      </b-form-group>

      <b-form-group label-cols="4" label-cols-lg="2" label-size="sm" :label="$t('common.NewPassword')"
                    label-for="input-5">
        <b-form-input v-model="newPassword" id="input-5" size="sm" type="password"></b-form-input>
      </b-form-group>

      <b-form-group label-cols="4" label-cols-lg="2" label-size="sm" :label="$t('common.ConfirmPassword')"
                    label-for="input-6">
        <b-form-input v-model="newPassword2" id="input-6" size="sm" type="password"></b-form-input>
      </b-form-group>

      <b-button slot="footer" variant="primary" size="sm" v-text="$t('common.ChangePassword')"
                v-bind:disabled="!canSubmitPassChange" @click="submitPassChange"></b-button>
    </b-card>
  </b-container>
</template>

<script>
import validate from '@/utils/validate'

export default {
  data () {
    return {
      username: '',
      email: '',
      confirmPassword: '',
      oldPassword: '',
      newPassword: '',
      newPassword2: '',
      errorMsg: '',
      successMsg: '',
      dismissCountDown: 5
    }
  },
  computed: {
    validPassword: function () {
      return validate.validPassword(this.newPassword)
    },
    samePassword: function () {
      return validate.samePassword(this.newPassword, this.newPassword2)
    },
    validUsername: function () {
      return validate.validUsername(this.username)
    },
    validEmail: function () {
      return validate.validEmail(this.email)
    },
    canSubmitInfoChange: function () {
      return this.validUsername && this.validEmail && this.confirmPassword
    },
    canSubmitPassChange: function () {
      return this.oldPassword && this.validPassword && this.samePassword
    }
  },
  methods: {
    submitInfoChange () {
      let vue = this
      this.$http.post('/api/users', {
        username: this.username,
        email: this.email
      }).then(function (result) {
        if (result.data.success) {
          vue.successMsg = data.$t('common.InfoChanged')
        } else {
          let msg = 'errors.ErrUnknownError'
          if (result.data.error.code) {
            msg = 'errors.' + result.data.error.code
          } else {
            msg = result.data.error.msg
          }
          vue.errorMsg = data.$t(msg)
        }
      }).catch(function (error) {
        let msg = 'errors.ErrUnknownError'
        if (error && error.response && error.response.data.error) {
          if (error.response.data.error.code) {
            msg = 'errors.' + error.response.data.error.code
          } else {
            msg = error.response.data.error.msg
          }
        }

        vue.errorMsg = data.$t(msg)
      })
    },
    submitPassChange () {
      let vue = this
      this.$http.put('/api/users', {
        password: this.newPassword
      }).then(function (result) {
        if (result.data.success) {
          vue.successMsg = data.$t('common.PasswordChanged')
        } else {
          let msg = 'errors.ErrUnknownError'
          if (result.data.error.code) {
            msg = 'errors.' + result.data.error.code
          } else {
            msg = result.data.error.msg
          }
          vue.errorMsg = data.$t(msg)
        }
      }).catch(function (error) {
        let msg = 'errors.ErrUnknownError'
        if (error && error.response && error.response.data.error) {
          if (error.response.data.error.code) {
            msg = 'errors.' + error.response.data.error.code
          } else {
            msg = error.response.data.error.msg
          }
        }

        vue.errorMsg = data.$t(msg)
      })
    }
  },
  mounted () {
    let vue = this
    this.$http.get('/api/users').then(function (data) {
      let user = data.data.data
      vue.username = user.username
      vue.email = user.email
    })
  }
}
</script>
