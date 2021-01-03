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
        <p v-text="$t('users.Login')" />
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12">
            <ui-input
              v-model.trim="email"
              autofocus
              :label="$t('users.Email')"
              :error-messages="errors.email"
              icon="mdi-account"
              type="email"
              @keyup.enter="submit"
            />
          </v-col>
          <v-col cols="12">
            <ui-password-input
              v-model="password"
              :label="$t('users.Password')"
              :error-messages="errors.password"
              @keyup.enter="submit"
            />
          </v-col>
          <v-col cols="12">
            <v-btn
              color="primary"
              large
              block
              @click="submit"
              v-text="$t('users.Login')"
            />
          </v-col>
          <v-col cols="12">
            <v-btn
              text
              block
              :to="{name: 'Register'}"
              v-text="$t('users.RegisterLink')"
            />
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </v-col>
</template>

<script>
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
    canSubmit () {
      return !(this.loginDisabled || this.email === '' || this.password === '')
    }
  },
  mounted () {
    if (this.hasAuth()) this.$router.push({ name: 'Servers' })
  },
  methods: {
    async submit () {
      this.errors.form = ''
      this.errors.email = ''
      this.errors.password = ''

      if (!this.email) {
        this.errors.email = this.$t('errors.ErrFieldRequired', { field: this.$t('users.Email') })
        return
      }

      if (!this.password) {
        this.errors.password = this.$t('errors.ErrFieldRequired', { field: this.$t('users.Password') })
        return
      }

      this.loginDisabled = true

      try {
        await this.$api.login(this.email, this.password)
        this.$emit('logged-in')
        this.$router.push({ name: 'Servers' })
      } finally {
        this.loginDisabled = false
      }
    }
  }
}
</script>
