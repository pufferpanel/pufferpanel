<template>
  <b-col
    md="6"
    sm="8"
    offset-md="3"
    offset-sm="2">
    <b-card
      header-tag="header"
      footer-tag="footer">
      <h6 slot="header" class="mb-0" align="center">Login</h6>
      <b-form>
        <b-row class="my-1">
          <b-col sm="1"></b-col>
          <b-col sm="10">
            <b-form-input
              v-model.trim="email"
              prepend-icon="mdi-account"
              name="email"
              id="email"
              placeholder="Email"
              type="text"/>
          </b-col>
        </b-row>
        <b-row class="my-1">
          <b-col sm="1"></b-col>
          <b-col sm="10">
            <b-form-input
              v-model="password"
              prepend-icon="mdi-lock"
              name="password"
              placeholder="Password"
              type="password"
              id="password"
              @keyup.enter="submit"/>
          </b-col>
        </b-row>
      </b-form>
      <div slot="footer">
        <b-row>
          <b-col sm="1"></b-col>
          <b-col sm="2">
            <router-link to="/auth/register">Register</router-link>
          </b-col>
          <b-col sm="6"></b-col>
          <b-col sm="2">
            <b-btn
              :disabled="!canSubmit"
              variant="primary"
              size="sm"
              @click="submit">Login
            </b-btn>
          </b-col>
          <b-col sm="1"></b-col>
        </b-row>
        <b-row v-if="errorMsg">
          <b-col sm="1"></b-col>
          <b-col sm="10">
            <b-alert
              fade
              dismissible
              show="5"
              variant="danger"
            >{{ errorMsg }}
            </b-alert>
          </b-col>
        </b-row>
      </div>
    </b-card>
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
      errorMsg: ''
    }
  },
  computed: {
    canSubmit: function () {
      return !(this.loginDisabled || this.email === '' || this.password === '')
    }
  },
  methods: {
    submit () {
      let data = this
      data.errorMsg = ''

      if (!data.email) {
        data.errorMsg = 'Email required'
        return
      }

      if (!data.password) {
        data.errorMsg = 'Password required'
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
          this.$router.push('/server')
        } else {
          data.errorMsg = responseData.msg + ''
        }
      }).catch(function (error) {
        let msg = error.message + ''
        if (error && error.response && error.response.data && error.response.data.msg) {
          msg = error.response.data.msg + ''
        }

        data.errorMsg = msg
      }).finally(function () {
        data.loginDisabled = false
      })
    }
  }
}
</script>
