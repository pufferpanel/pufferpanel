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
            <ui-input
              v-model.trim="username"
              autofocus
              :label="$t('users.Username')"
              :error-messages="(username && !validUsername) ? $t('errors.ErrUsernameRequirements', { min: 5 }) : errors.username"
              icon="mdi-account"
              @keyup.enter="submit"
            />
          </v-col>
          <v-col cols="12">
            <ui-input
              v-model.trim="email"
              :label="$t('users.Email')"
              :error-messages="errors.email"
              icon="mdi-email"
              type="email"
              @keyup.enter="submit"
            />
          </v-col>
          <v-col cols="12">
            <ui-password-input
              v-model="password"
              :label="$t('users.Password')"
              :error-messages="(password && !validPassword) ? $t('errors.ErrPasswordRequirements', { min: 8 }) : errors.password"
              @keyup.enter="submit"
            />
          </v-col>
          <v-col cols="12">
            <ui-password-input
              v-model="confirmPassword"
              :label="$t('users.ConfirmPassword')"
              :error-messages="(confirmPassword !== '' && !samePassword) ? $t('errors.ErrPasswordsNotIdentical') : ''"
              @keyup.enter="submit"
            />
          </v-col>
          <v-col cols="12">
            <v-btn
              color="primary"
              large
              block
              :disabled="!canComplete"
              @click="submit"
              v-text="$t('users.Register')"
            />
          </v-col>
          <v-col cols="12">
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
import { handleError } from '@/utils/api'
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
    canComplete () {
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
    validPassword () {
      return validate.validPassword(this.password)
    },
    samePassword () {
      return validate.samePassword(this.password, this.confirmPassword)
    },
    validUsername () {
      return validate.validUsername(this.username)
    },
    validEmail () {
      return validate.validEmail(this.email)
    }
  },
  mounted () {
    if (hasAuth()) this.$router.push({ name: 'Servers' })
  },
  methods: {
    // real methods
    submit () {
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

      const ctx = this

      this.axios.post(this.$route.path, {
        email: this.email,
        password: this.password,
        username: this.username
      }).then(response => {
        if (response.data.token && response.data.token !== '') {
          Cookies.set('puffer_auth', response.data.token, { sameSite: 'strict' })
          localStorage.setItem('scopes', JSON.stringify(response.data.scopes || []))
          ctx.$emit('logged-in')
          ctx.$router.push({ name: 'Servers' })
        } else {
          this.$toast.success(this.$t('users.RegisterSuccess'))
          ctx.$router.push({ name: 'Login' })
        }
      })
        .catch(handleError(ctx))
        .finally(() => { ctx.registerDisabled = false })
    }
  }
}
</script>
