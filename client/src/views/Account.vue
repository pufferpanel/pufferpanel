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
        cols="12"
        md="6"
      >
        <v-card>
          <v-card-title v-text="$t('users.ChangeInfo')" />
          <v-card-text>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  v-model="username"
                  outlined
                  prepend-inner-icon="mdi-account"
                  :label="$t('users.Username')"
                />
                <v-text-field
                  v-model="email"
                  outlined
                  prepend-inner-icon="mdi-email"
                  :label="$t('users.Email')"
                />
                <v-text-field
                  v-model="confirmPassword"
                  outlined
                  prepend-inner-icon="mdi-lock"
                  type="password"
                  :label="$t('users.ConfirmPassword')"
                />
                <v-btn
                  large
                  block
                  color="primary"
                  :disabled="!canSubmitInfoChange"
                  @click="submitInfoChange"
                  v-text="$t('users.ChangeInfo')"
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col
        cols="12"
        md="6"
      >
        <v-card>
          <v-card-title v-text="$t('users.ChangePassword')" />
          <v-card-text>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  v-model="oldPassword"
                  outlined
                  prepend-inner-icon="mdi-lock"
                  type="password"
                  :label="$t('users.OldPassword')"
                />
                <v-text-field
                  v-model="newPassword"
                  outlined
                  prepend-inner-icon="mdi-lock"
                  type="password"
                  :label="$t('users.NewPassword')"
                />
                <v-text-field
                  v-model="confirmNewPassword"
                  outlined
                  prepend-inner-icon="mdi-lock"
                  type="password"
                  :label="$t('users.ConfirmPassword')"
                />
                <v-btn
                  large
                  block
                  color="primary"
                  :disabled="!canSubmitPassChange"
                  @click="submitPassChange"
                  v-text="$t('users.ChangePassword')"
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
import { handleError } from '@/utils/api'

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
    validPassword () {
      return validate.validPassword(this.newPassword)
    },
    validUsername () {
      return validate.validUsername(this.username)
    },
    validEmail () {
      return validate.validEmail(this.email)
    },
    canSubmitInfoChange () {
      return this.validUsername && this.validEmail && this.confirmPassword
    },
    canSubmitPassChange () {
      return this.oldPassword && this.validPassword && this.newPassword === this.confirmNewPassword
    }
  },
  mounted () {
    const ctx = this
    this.$http.get('/api/self').then(data => {
      const user = data.data
      ctx.username = user.username
      ctx.email = user.email
    }).catch(handleError(ctx))
  },
  methods: {
    submitInfoChange () {
      const ctx = this
      this.$http.put('/api/self', {
        username: this.username,
        email: this.email,
        password: this.confirmPassword
      }).then(result => {
        ctx.$toast.success(ctx.$t('users.InfoChanged'))
      }).catch(handleError(ctx))
    },
    submitPassChange () {
      const ctx = this
      this.$http.put('/api/self', {
        password: this.oldPassword,
        newPassword: this.newPassword
      }).then(result => {
        ctx.$toast.success(ctx.$t('users.PasswordChanged'))
      }).catch(handleError(ctx))
    }
  }
}
</script>
