<!--
  - Copyright 2019 Padduck, LLC
  -  Licensed under the Apache License, Version 2.0 (the "License");
  -  you may not use this file except in compliance with the License.
  -  You may obtain a copy of the License at
  -          http://www.apache.org/licenses/LICENSE-2.0
  -  Unless required by applicable law or agreed to in writing, software
  -  distributed under the License is distributed on an "AS IS" BASIS,
  -  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  -  See the License for the specific language governing permissions and
  -  limitations under the License.
  -->

<template>
  <v-container>
    <v-row>
      <v-col
        offset-md="3"
        md="6"
      >
        <v-card>
          <v-card-title v-text="$t('common.ChangeInfo')" />
          <v-card-text>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  v-model="username"
                  outlined
                  prepend-inner-icon="mdi-account"
                  :label="$t('common.Username')"
                />
                <v-text-field
                  v-model="email"
                  outlined
                  prepend-inner-icon="mdi-email"
                  :label="$t('common.Email')"
                />
                <v-text-field
                  v-model="confirmPassword"
                  outlined
                  prepend-inner-icon="mdi-lock"
                  type="password"
                  :label="$t('common.ConfirmPassword')"
                />
                <v-btn
                  large
                  block
                  color="primary"
                  :disabled="!canSubmitInfoChange"
                  @click="submitInfoChange"
                  v-text="$t('common.ChangeInfo')"
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col
        offset-md="3"
        md="6"
      >
        <v-card>
          <v-card-title v-text="$t('common.ChangePassword')" />
          <v-card-text>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  v-model="oldPassword"
                  outlined
                  prepend-inner-icon="mdi-lock"
                  type="password"
                  :label="$t('common.OldPassword')"
                />
                <v-text-field
                  v-model="newPassword"
                  outlined
                  prepend-inner-icon="mdi-lock"
                  type="password"
                  :label="$t('common.NewPassword')"
                />
                <v-text-field
                  v-model="confirmNewPassword"
                  outlined
                  prepend-inner-icon="mdi-lock"
                  type="password"
                  :label="$t('common.ConfirmPassword')"
                />
                <v-btn
                  large
                  block
                  color="primary"
                  :disabled="!canSubmitPassChange"
                  @click="submitPassChange"
                  v-text="$t('common.ChangePassword')"
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
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
      confirmNewPassword: ''
    }
  },
  computed: {
    validPassword: function () {
      return validate.validPassword(this.newPassword)
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
      return this.oldPassword && this.validPassword && this.newPassword === this.confirmNewPassword
    }
  },
  mounted () {
    const ctx = this
    this.$http.get('/api/self').then(function (data) {
      const user = data.data
      ctx.username = user.username
      ctx.email = user.email
    })
  },
  methods: {
    submitInfoChange () {
      const ctx = this
      this.$http.put('/api/self', {
        username: this.username,
        email: this.email,
        password: this.confirmPassword
      }).then(function (result) {
        if (result.status >= 200 && result.status < 300) {
          ctx.$toast.success(ctx.$t('common.InfoChanged'))
        } else {
          let msg = 'errors.ErrUnknownError'
          if (result.data.error.code) {
            msg = 'errors.' + result.data.error.code
          } else {
            msg = result.data.error.msg
          }
          ctx.$toast.error(ctx.$t(msg))
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

        ctx.$toast.error(ctx.$t(msg))
      })
    },
    submitPassChange () {
      const ctx = this
      this.$http.put('/api/self', {
        password: this.oldPassword,
        newPassword: this.newPassword
      }).then(function (result) {
        if (result.status >= 200 && result.status < 300) {
          ctx.$toast.success(ctx.$t('common.PasswordChanged'))
        } else {
          let msg = 'errors.ErrUnknownError'
          if (result.data.error.code) {
            msg = 'errors.' + result.data.error.code
          } else {
            msg = result.data.error.msg
          }
          ctx.$toast.error(ctx.$t(msg))
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

        ctx.$toast.error(ctx.$t(msg))
      })
    }
  }
}
</script>
