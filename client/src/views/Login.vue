<template>
  <v-layout align-center justify-center>
    <v-flex xs12 sm8 md4>
      <material-card color="blue" title="Login">
        <v-container>
          <v-card-text>
            <v-form>
              <v-text-field v-model.trim="email" prepend-icon="mdi-account" name="email" label="Email"
                            type="text"></v-text-field>
              <v-text-field v-model="password" id="password" prepend-icon="mdi-lock" name="password"
                            label="Password" type="password"></v-text-field>
            </v-form>
          </v-card-text>
          <material-notification v-if="error" v-text="error" color="error"></material-notification>
        </v-container>
        <v-container>
          <v-card-actions>
            <a href="/auth/register">Register</a>
            <v-spacer></v-spacer>
            <v-btn :disabled="loginDisabled" color="blue" v-on:click="submit">Login</v-btn>
          </v-card-actions>
        </v-container>
      </material-card>
    </v-flex>
  </v-layout>
</template>

<script>
import axios from 'axios'
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

      axios.post('/auth/login', {
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
