<template>
  <v-col
    lg="4"
    md="6"
    sm="8"
    offset-lg="4"
    offset-md="3"
    offset-sm="2"
  >
    <v-card :loading="registerDisabled">
      <v-card-title class="d-flex justify-center">
        <p v-text="$t('users.Register')" />
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12">
            <v-form>
              <v-text-field
                id="username"
                v-model.trim="username"
                outlined
                :label="$t('users.Username')"
                :error-messages="errors.username"
                prepend-inner-icon="mdi-account"
                name="username"
                type="text"
                @keyup.enter="submit"
              />
              <v-text-field
                id="email"
                v-model.trim="email"
                outlined
                :label="$t('users.Email')"
                :error-messages="errors.email"
                prepend-inner-icon="mdi-email"
                name="email"
                type="text"
                @keyup.enter="submit"
              />
              <v-text-field
                id="password"
                v-model="password"
                outlined
                :label="$t('users.Password')"
                :error-messages="(password && !validPassword) ? $t('errors.ErrPasswordRequirements') : errors.password"
                prepend-inner-icon="mdi-lock"
                :append-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                name="password"
                :type="!showPassword ? 'password' : 'text'"
                @click:append="showPassword = !showPassword"
                @keyup.enter="submit"
              />
              <v-text-field
                id="confirmPassword"
                v-model="confirmPassword"
                outlined
                :label="$t('users.ConfirmPassword')"
                :error-messages="(confirmPassword !== '' && !samePassword) ? $t('errors.ErrPasswordsNotIdentical') : ''"
                prepend-inner-icon="mdi-lock"
                :append-icon="showConfirmPassword ? 'mdi-eye-off' : 'mdi-eye'"
                name="password"
                :type="!showConfirmPassword ? 'password' : 'text'"
                @click:append="showConfirmPassword = !showConfirmPassword"
                @keyup.enter="submit"
              />
            </v-form>
            <v-btn
              color="primary"
              class="body-1 mb-5"
              large
              block
              :disabled="!canComplete"
              @click="submit"
              v-text="$t('users.Register')"
            />
            <v-btn
              text
              block
              :to="{name: 'Login'}"
              v-text="$t('users.LoginLink')"
            />
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </v-col>
</template>

<script>
import Cookies from 'js-cookie'
import validate from '@/utils/validate'
import { hasAuth } from '@/utils/auth'

export default {
  data () {
    return {
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
      registerDisabled: false,
      showPassword: false,
      showConfirmPassword: false,
      errors: {
        username: '',
        email: '',
        password: ''
      }
    }
  },
  computed: {
    canComplete: function () {
      if (this.registerDisabled) {
        return false
      }
      if (!this.username || !this.validUsername) {
        return false
      }
      if (!this.email || !this.validEmail) {
        return false
      }

      return !(!this.validPassword || !this.samePassword)
    },
    validPassword: function () {
      return validate.validPassword(this.password)
    },
    samePassword: function () {
      return validate.samePassword(this.password, this.confirmPassword)
    },
    validUsername: function () {
      return validate.validUsername(this.username)
    },
    validEmail: function () {
      return validate.validEmail(this.email)
    }
  },
  mounted () {
    if (hasAuth()) this.$router.push({ name: 'Servers' })
  },
  methods: {
    // real methods
    submit: function () {
      this.$toast.clearQueue()
      if (this.$toast.getCmp()) this.$toast.getCmp().close()
      this.errors.username = ''
      this.errors.email = ''
      this.errors.password = ''
      if (!this.username) {
        this.errors.username = this.$t('errors.ErrFieldRequired', { field: this.$t('users.Username') })
        return
      }
      if (!this.email) {
        this.errors.email = this.$t('errors.ErrFieldRequired', { field: this.$t('users.Email') })
        return
      }
      if (!this.password) {
        this.errors.password = this.$t('errors.ErrFieldRequired', { field: this.$t('users.Password') })
        return
      }
      if (!validate.validPassword(this.password)) {
        return
      }
      this.registerDisabled = true

      const vue = this

      this.axios.post(this.$route.path, {
        email: this.email,
        password: this.password,
        username: this.username
      }).then(function (response) {
        if (response.status >= 200 && response.status < 300) {
          if (response.data.token && response.data.token !== '') {
            Cookies.set('puffer_auth', response.data.token)
            localStorage.setItem('scopes', JSON.stringify(response.data.scopes || []))
            vue.$emit('logged-in')
            vue.$router.push({ name: 'Servers' })
          } else {
            this.$toast.success(this.$t('users.RegisterSuccess'))
            vue.$router.push({ name: 'Login' })
          }
        } else {
          vue.$toast.error(response.data.msg)
          vue.registerDisabled = false
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

        vue.$toast.error(vue.$t(msg))
        vue.registerDisabled = false
      })
    }
  }
}
</script>
