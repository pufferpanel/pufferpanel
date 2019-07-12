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
                v-bind:disabled="!canSubmitInfoChange"></b-button>
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
                v-bind:disabled="!canSubmitPassChange"></b-button>
    </b-card>
  </b-container>
</template>

<script>
import validate from '@/plugins/validate.js'

export default {
  data () {
    return {
      username: '',
      email: '',
      confirmPassword: '',
      oldPassword: '',
      newPassword: '',
      newPassword2: ''
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
