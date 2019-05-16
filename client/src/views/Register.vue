<template>
  <b-col
    md="6"
    sm="8"
    offset-md="3"
    offset-sm="2">
    <b-card
      header-tag="header"
      footer-tag="footer">
      <h6
        slot="header"
        class="mb-0"
        align="center">Register - PufferPanel</h6>
      <b-form>
        <b-row>
          <b-col sm="1"/>
          <b-col sm="10">
            <b-alert
              v-if="username && !validUsername"
              fade
              show
              variant="danger"
            >Username must be at least 8 characters and only contain alphanumerics, _, or -.
            </b-alert>
          </b-col>
        </b-row>
        <b-row>
          <b-col sm="1"/>
          <b-col sm="10">
            <b-alert
              v-if="email && !validEmail"
              fade
              show
              variant="danger"
            >Email is not valid
            </b-alert>
          </b-col>
        </b-row>

        <b-row class="my-1">
          <b-col sm="1"/>
          <b-col sm="10">
            <b-form-input
              id="username"
              v-model.trim="username"
              prepend-icon="mdi-account"
              name="username"
              placeholder="Username"
              type="text"/>
          </b-col>
        </b-row>

        <b-row class="my-1">
          <b-col sm="1"/>
          <b-col sm="10">
            <b-form-input
              id="email"
              v-model.trim="email"
              prepend-icon="mdi-account"
              name="email"
              placeholder="Email"
              type="text"/>
          </b-col>
        </b-row>

        <b-row>
          <b-col sm="1"/>
          <b-col sm="10">
            <b-alert
              v-if="password && !validPassword"
              fade
              show
              variant="danger"
            >Passwords must be at least 8 characters
            </b-alert>
          </b-col>
        </b-row>
        <b-row>
          <b-col sm="1"/>
          <b-col sm="10">
            <b-alert
              v-if="password && confirmPassword && !samePassword"
              fade
              show
              variant="danger"
            >Passwords must be the same
            </b-alert>
          </b-col>
        </b-row>

        <b-row class="my-1">
          <b-col sm="1"/>
          <b-col sm="10">
            <b-form-input
              id="password"
              v-model="password"
              prepend-icon="mdi-lock"
              name="password"
              placeholder="Password"
              type="password"/>
          </b-col>
        </b-row>

        <b-row class="my-1">
          <b-col sm="1"/>
          <b-col sm="10">
            <b-form-input
              id="confirmPassword"
              v-model="confirmPassword"
              prepend-icon="mdi-lock"
              name="confirmPassword"
              placeholder="Confirm Password"
              type="password"
              @keyup.enter="submit"/>
          </b-col>
        </b-row>
      </b-form>
      <div slot="footer">
        <b-row>
          <b-col sm="1"/>
          <b-col sm="2">
            <b-link :to="{name: 'Login'}">Login</b-link>
          </b-col>
          <b-col sm="6"/>
          <b-col sm="2">
            <b-btn
              :disabled="!canComplete"
              variant="primary"
              size="sm"
              @click="submit">Register
            </b-btn>
          </b-col>
          <b-col sm="1"/>
        </b-row>
        <b-row v-if="errorMsg">
          <b-col sm="1"/>
          <b-col sm="10">
            <b-alert
              :show="dismissCountDown"
              fade
              dismissible
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
export default {
  data () {
    return {
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
      registerDisabled: false,
      errorMsg: '',
      dismissCountDown: 5
    }
  },
  computed: {
    validPassword: function () {
      return this.password.length >= 8
    },
    samePassword: function () {
      return this.password && this.confirmPassword && this.password === this.confirmPassword
    },
    validEmail: function () {
      return this.email && /^\S+@\S+\.\S+$/.test(this.email)
    },
    validUsername: function () {
      return this.username && /^([0-9A-Za-z_-]){8,}$/.test(this.username)
    },
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
    }
  },
  methods: {
    submit: function () {
      if (!this.username) {
        this.errorMsg = 'Username is required'
        return
      }
      if (!this.email) {
        this.errorMsg = 'Email is required'
        return
      }
      if (!this.password) {
        this.errorMsg = 'Password is required'
        return
      }
      if (!this.validPassword) {
        this.errorMsg = 'Password is not valid'
      }
      this.registerDisabled = true

      let data = this

      this.axios.post(this.$route.path, {
        data: {
          email: this.email,
          password: this.password,
          username: this.username
        }
      }).then(function (response) {
        let responseData = response.data
        if (responseData.success) {
          data.$router.push({ name: 'Login' })
        } else {
          data.error = responseData.msg
          data.registerDisabled = false
        }
      }).catch(function (error) {
        let msg = error.response.data.msg + ''
        if (!msg) {
          msg = error + ''
        }
        data.errorMsg = msg
        data.registerDisabled = false
      })
    }
  }
}
</script>
