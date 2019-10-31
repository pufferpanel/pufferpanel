<template>
  <v-col
    lg="4"
    md="6"
    sm="8"
    offset-lg="4"
    offset-md="3"
    offset-sm="2"
  >
    <v-card :loading="loginDisabled">
      <v-card-title class="d-flex justify-center">
        <p v-text="$t('common.Login')" />
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12">
            <v-form>
              <v-text-field
                id="email"
                v-model.trim="email"
                outlined
                :label="$t('common.Email')"
                :error-messages="errors.email"
                prepend-inner-icon="mdi-account"
                name="email"
                type="email"
                @keyup.enter="submit"
              />
              <v-text-field
                id="password"
                v-model="password"
                outlined
                :label="$t('common.Password')"
                :error-messages="errors.password"
                prepend-inner-icon="mdi-lock"
                :append-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                name="password"
                :type="!showPassword ? 'password' : 'text'"
                @click:append="showPassword = !showPassword"
                @keyup.enter="submit"
              />
            </v-form>
            <v-btn
              color="primary"
              class="body-1 mb-5"
              large
              block
              @click="submit"
              v-text="$t('common.Login')"
            />
            <v-btn
              text
              block
              :to="{name: 'Register'}"
              v-text="$t('common.RegisterLink')"
            />
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </v-col>
</template>

<script>
import Cookies from 'js-cookie'
import { hasAuth } from '@/utils/auth'

function getReauthReason () {
  const reason = localStorage.getItem('reauth_reason') || ''
  localStorage.removeItem('reauth_reason')
  return reason
}

function getRegistered () {
  const registered = !!((localStorage.getItem('registered') || ''))
  localStorage.removeItem('registered')
  return registered
}

export default {
  data () {
    return {
      email: '',
      password: '',
      errors: {
        email: '',
        password: ''
      },
      loginDisabled: false,
      reauthReason: '',
      registered: false,
      showPassword: false
    }
  },
  computed: {
    canSubmit: function () {
      return !(this.loginDisabled || this.email === '' || this.password === '')
    }
  },
  mounted () {
    if (hasAuth()) this.$router.push({ name: 'Servers' })
    if (getRegistered()) this.$notify(this.$t('common.RegisterSuccess'), 'success')
    const reauthReason = getReauthReason()
    if (reauthReason === 'session_timed_out') this.$notify(this.$t('errors.ErrSessionTimedOut'), 'error')
  },
  methods: {
    submit () {
      const data = this
      data.errors.form = ''
      data.errors.email = ''
      data.errors.password = ''

      if (!data.email) {
        data.errors.email = this.$t('errors.ErrFieldRequired', { field: this.$t('common.Email') })
        return
      }

      if (!data.password) {
        data.errors.password = this.$t('errors.ErrFieldRequired', { field: this.$t('common.Password') })
        return
      }

      this.loginDisabled = true

      this.axios.post(this.$route.path, {
        email: this.email,
        password: this.password
      }).then(function (response) {
        if (response.status >= 200 && response.status < 300) {
          Cookies.set('puffer_auth', response.data.session)
          localStorage.setItem('admin', response.data.admin)
          data.$emit('logged-in')
          data.$router.push({ name: 'Servers' })
        } else {
          data.$notify(response.data.msg, 'error')
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

        data.$notify(data.$t(msg), 'error')
      }).finally(function () {
        data.loginDisabled = false
      })
    }
  }
}
</script>
