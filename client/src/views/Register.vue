<template>
  <v-layout
    align-center
    justify-center>
    <v-flex
      xs12
      sm8
      md4>
      <material-card
        color="blue"
        title="Register">
        <v-container>
          <v-card-text>
            <v-card-text v-if="username && !validUsername">
              <material-notification color="error">Username must be at least 8 characters and only contain
                alphanumerics, _, or -.
              </material-notification>
            </v-card-text>
            <v-card-text v-if="email && !validEmail">
              <material-notification color="error">Email is not valid</material-notification>
            </v-card-text>
            <v-form>
              <v-text-field
                v-model.trim="username"
                prepend-icon="mdi-account"
                name="username"
                label="Username"
                type="text"/>
              <v-text-field
                v-model.trim="email"
                prepend-icon="mdi-account"
                name="email"
                label="Email"
                type="text"/>
            </v-form>
          </v-card-text>
          <v-card-text>
            <v-card-text v-if="password && confirmPassword && !samePassword">
              <material-notification color="error">Passwords must be the same</material-notification>
            </v-card-text>
            <v-card-text v-if="password && !validPassword">
              <material-notification color="error">Passwords must be at least 8 characters</material-notification>
            </v-card-text>
            <v-form>
              <v-text-field
                v-model="password"
                prepend-icon="mdi-lock"
                name="password"
                label="Password"
                type="password"/>
              <v-text-field
                v-model="confirmPassword"
                prepend-icon="mdi-lock"
                name="confirmPassword"
                label="Confirm Password"
                type="password"
                @keyup.enter="submit"/>
            </v-form>
          </v-card-text>
        </v-container>
        <v-container>
          <v-card-text v-if="error">
            <material-notification
              color="error"
              v-text="error"/>
          </v-card-text>
          <v-card-actions>
            <router-link to="/auth/login">Login</router-link>
            <v-spacer/>
            <v-btn
              :disabled="!canComplete"
              color="blue"
              @click="submit">Register
            </v-btn>
          </v-card-actions>
        </v-container>
      </material-card>
    </v-flex>
  </v-layout>
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
      error: ''
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
        this.error = 'Username is required'
        return
      }
      if (!this.email) {
        this.error = 'Email is required'
        return
      }
      if (!this.password) {
        this.error = 'Password is required'
        return
      }
      if (!this.validPassword) {
        this.error = 'Password is not valid'
      }
      this.registerDisabled = true
      this.axios.post('/auth/register', {
        data: {
          email: this.email,
          password: this.password,
          username: this.username
        }
      }).then(function (response) {
        let responseData = response.data
        if (responseData.success) {
          this.$router.push('/auth/login')
        } else {
          this.error = responseData.msg
          this.registerDisabled = false
        }
      }).catch(function (error) {
        let msg = error.response.data.msg
        if (!msg) {
          msg = error
        }
        this.error = msg
        this.registerDisabled = false
      })
    }
  }
}
</script>
