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
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <span v-text="$t('users.Otp')" />
            <span class="flex-grow-1" />
            <v-btn
              :loading="otpLoading"
              :disabled="otpLoading"
              :color="getOtpBtnColor()"
              @click="toggleOtp"
              v-text="otpActive ? $t('users.OtpDisable') : $t('users.OtpEnable')"
            />
          </v-card-title>
          <v-card-text v-text="$t('users.OtpHint')" />
        </v-card>
      </v-col>
    </v-row>
    <ui-overlay v-model="otpEnroll.started" card closable :title="$t('users.OtpSetup')">
      <v-row>
        <v-col>
          <v-row>
            <v-col>
              <h3 v-text="$t('users.OtpSetupHint')" />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <h4 v-text="$t('users.OtpSecret')" />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <span v-text="otpEnroll.secret" />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <h3 v-text="$t('users.OtpConfirm')" />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <ui-input v-model="otpConfirmToken" autofocus @keyup.enter="confirmOtpEnroll" />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-btn color="success" block v-text="$t('users.OtpEnable')" @click="confirmOtpEnroll" />
            </v-col>
          </v-row>
        </v-col>
        <v-col cols="12" sm="4" md="3">
          <img style="width:100%;max-width:300px;" :src="otpEnroll.qrCode">
        </v-col>
      </v-row>
    </ui-overlay>
    <ui-overlay v-model="otpDisabling" card closable :title="$t('users.OtpDisable')">
      <v-row>
        <v-col>
          <ui-input v-model="otpConfirmToken" autofocus @keyup.enter="confirmOtpDisable" />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-btn color="error" block v-text="$t('users.OtpDisable')" @click="confirmOtpDisable" />
        </v-col>
      </v-row>
    </ui-overlay>
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
      confirmNewPassword: '',
      otpLoading: true,
      otpActive: false,
      otpConfirmToken: '',
      otpDisabling: false,
      otpEnroll: {
        started: false,
        secret: '',
        qrCode: ''
      }
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
    this.otpActive = await this.$api.getOtp()
    this.otpLoading = false
  },
  methods: {
    async submitInfoChange () {
      await this.$api.updateSelf(this.username, this.email, this.confirmPassword)
      this.$toast.success(this.$t('users.InfoChanged'))
    },
    async submitPassChange () {
      await this.$api.updatePassword(this.oldPassword, this.newPassword)
      this.$toast.success(this.$t('users.PasswordChanged'))
    },
    getOtpBtnColor () {
      if (this.otpLoading) return undefined
      if (this.otpActive) return 'error'
      return 'success'
    },
    async toggleOtp () {
      if (this.otpActive) {
        this.otpDisabling = true
      } else {
        const otpData = await this.$api.startOtpEnroll()
        this.otpEnroll = { started: true, secret: otpData.secret, qrCode: otpData.img }
      }
    },
    async confirmOtpEnroll () {
      const ok = await this.$api.validateOtpEnroll(this.otpConfirmToken)
      if (ok) {
        this.otpEnroll.started = false
        this.otpActive = true
      }
      this.otpConfirmToken = ''
    },
    async confirmOtpDisable () {
      const ok = await this.$api.disableOtp(this.otpConfirmToken)
      if (ok) {
        this.otpDisabling = false
        this.otpActive = false
      }
      this.otpConfirmToken = ''
    }
  }
}
</script>
