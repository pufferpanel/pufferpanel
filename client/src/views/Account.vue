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
                <ui-input
                  v-model="username"
                  icon="mdi-account"
                  :label="$t('users.Username')"
                />
              </v-col>
              <v-col cols="12">
                <ui-input
                  v-model="email"
                  icon="mdi-email"
                  :label="$t('users.Email')"
                  type="email"
                />
              </v-col>
              <v-col cols="12">
                <ui-password-input
                  v-model="confirmPassword"
                  :label="$t('users.ConfirmPassword')"
                />
              </v-col>
              <v-col cols="12">
                <v-btn
                  large
                  block
                  color="success"
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
                <ui-password-input
                  v-model="oldPassword"
                  :label="$t('users.OldPassword')"
                />
              </v-col>
              <v-col cols="12">
                <ui-password-input
                  v-model="newPassword"
                  :label="$t('users.NewPassword')"
                />
              </v-col>
              <v-col cols="12">
                <ui-password-input
                  v-model="confirmNewPassword"
                  :label="$t('users.ConfirmPassword')"
                />
              </v-col>
              <v-col cols="12">
                <v-btn
                  large
                  block
                  color="success"
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
  async mounted () {
    const user = await this.$api.getSelf()
    this.username = user.username
    this.email = user.email
  },
  methods: {
    async submitInfoChange () {
      await this.$api.updateSelf(this.username, this.email, this.confirmPassword)
      this.$toast.success(this.$t('users.InfoChanged'))
    },
    async submitPassChange () {
      await this.$api.updatePassword(this.oldPassword, this.newPassword)
      this.$toast.success(this.$t('users.PasswordChanged'))
    }
  }
}
</script>
