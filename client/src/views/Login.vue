<template>
  <b-col md="6" sm="8" offset-md="3" offset-sm="2">
    <v-card
      title="Login">
      <v-container>
        <v-card-text>
          <v-form>
            <v-text-field
              v-model.trim="email"
              prepend-icon="mdi-account"
              name="email"
              label="Email"
              type="text"/>
            <v-text-field
              v-model="password"
              prepend-icon="mdi-lock"
              name="password"
              label="Password"
              type="password"
              @keyup.enter="submit"/>
          </v-form>
        </v-card-text>
        <v-alert
          v-if="error"
          color="error"
          v-text="error"
        />
      </v-container>
      <v-container>
        <v-card-actions>
          <a href="/auth/register">Register</a>
          <v-spacer/>
          <b-btn
            :disabled="cannotSubmit"
            variant="primary"
            size="sm"
            @click="submit">Login
          </b-btn>
        </v-card-actions>
      </v-container>
    </v-card>
  </b-col>
</template>

<script>
import Cookies from 'js-cookie'

export default {
  data () {
    return {
      email: '',
      password: '',
      loginDisabled: false,
      error: ''
    }
  },
  computed: {
    cannotSubmit: function () {
      return this.loginDisabled || this.email === '' || this.password === ''
    }
  },
  methods: {
    submit () {
      let data = this
      data.error = ''

      if (!data.email) {
        data.error = 'Email required'
        return
      }

      if (!data.password) {
        data.error = 'Password required'
        return
      }

      this.loginDisabled = true

      this.axios.post('/auth/login', {
        data: {
          email: this.email,
          password: this.password
        }
      }).then(function (response) {
        let responseData = response.data
        if (responseData.success) {
          Cookies.set('puffer_auth', responseData.data)
          window.location.href = '/server'
        } else {
          data.error = responseData.msg
          data.loginDisabled = false
        }
      }).catch(function (error) {
        let msg = error.response.data.msg
        if (!msg) {
          msg = error
        }
        data.error = msg
        data.loginDisabled = false
      })
    }
  }
}
</script>
