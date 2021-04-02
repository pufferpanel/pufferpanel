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
              disabled
              :label="$t('users.Email')"
              icon="mdi-email"
              type="email"
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
import validate from '@/utils/validate'

export default {
  data () {
    return {
      username: '',
      email: this.$route.query.email,
      password: '',
      confirmPassword: '',
      registerDisabled: false,
      errors: {
        username: '',
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
    }
  },
  mounted () {
    if (this.hasAuth()) this.$router.push({ name: 'Servers' })
  },
  methods: {
    async submit () {
      this.errors.username = ''
      this.errors.password = ''

      if (!this.username) {
        this.errors.username = this.$t('errors.ErrFieldRequired', { field: this.$t('users.Username') })
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

      let step = 'verifyToken'
      try {
        await this.$api.login(this.email, this.$route.query.token, { silent: true })
        step = 'update'
        await this.$api.updateSelf(this.username, this.email, this.$route.query.token, { noToast: true })
        await this.$api.updatePassword(this.$route.query.token, this.password, { noToast: true })
        step = 'login'
        await this.$api.login(this.email, this.password, { noToast: true })
        this.$emit('logged-in')
        this.$router.push({ name: 'Servers' })
      } catch {
        if (step === 'verifyToken') {
          this.$toast.error(this.$t('errors.ErrInviteLinkInvalid'))
        } else if (step === 'update') {
          this.$toast.error(this.$t('errors.ErrSavingInvitedUser'))
        } else {
          this.$toast.success(this.$t('users.RegisterSuccess'))
          this.$router.push({ name: 'Login' })
        }
      } finally {
        this.registerDisabled = false
      }
    }
  }
}
</script>
