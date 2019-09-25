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
        <v-alert
          v-if="registered"
          dense
          outlined
          dismissible
          type="success"
        >
          <span v-text="$t('common.RegisterSuccess')" />
        </v-alert>
        <v-alert
          v-if="reauthReason"
          dense
          outlined
          dismissible
          type="error"
        >
          <span
            v-if="reauthReason == 'session_timed_out'"
            v-text="$t('errors.ErrSessionTimedOut')"
          />
        </v-alert>
        <v-alert
          v-if="errors.form"
          dense
          outlined
          dismissible
          type="error"
        >
          {{ errors.form }}
        </v-alert>
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
        password: '',
        form: ''
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
    this.reauthReason = getReauthReason()
    this.registered = getRegistered()
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
        data: {
          email: this.email,
          password: this.password
        }
      }).then(function (response) {
        const responseData = response.data
        if (responseData.success) {
          Cookies.set('puffer_auth', responseData.data.session)
          localStorage.setItem('scopes', JSON.stringify(responseData.data.scopes))
          data.$emit('logged-in')
          data.$router.push({ name: 'Servers' })
        } else {
          data.errors.form = responseData.msg + ''
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

        data.errors.form = data.$t(msg)
      }).finally(function () {
        data.loginDisabled = false
      })
    }
  }
}
</script>
