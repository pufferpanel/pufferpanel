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
        align="center"
        v-text="$t('common.Register')"></h6>
      <b-form>
        <b-row>
          <b-col sm="1"/>
          <b-col sm="10">
            <b-alert
              v-if="username && !validUsername"
              v-text="$t('errors.ErrUsernameRequirements')"
              fade
              show
              variant="danger"
            >
            </b-alert>
          </b-col>
        </b-row>
        <b-row>
          <b-col sm="1"/>
          <b-col sm="10">
            <b-alert
              v-if="email && !validEmail"
              v-text="$t('errors.ErrFieldNotEmail', {'field': $t('common.Email')})"
              fade
              show
              variant="danger"
            >
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
              :placeholder="$t('common.Username')"
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
              :placeholder="$t('common.Email')"
              type="text"/>
          </b-col>
        </b-row>

        <b-row>
          <b-col sm="1"/>
          <b-col sm="10">
            <b-alert
              v-if="password && !validPassword"
              v-text="$t('errors.ErrPasswordRequirements')"
              fade
              show
              variant="danger"
            >
            </b-alert>
          </b-col>
        </b-row>
        <b-row>
          <b-col sm="1"/>
          <b-col sm="10">
            <b-alert
              v-if="password && confirmPassword && !samePassword"
              v-text="$t('errors.ErrPasswordsNotIdentical')"
              fade
              show
              variant="danger"
            >
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
              :placeholder="$t('common.Password')"
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
              :placeholder="$t('common.ConfirmPassword')"
              type="password"
              @keyup.enter="submit"/>
          </b-col>
        </b-row>
      </b-form>
      <div slot="footer">
        <b-row>
          <b-col sm="1"/>
          <b-col sm="2">
            <b-link :to="{name: 'Login'}" v-text="$t('common.Login')"></b-link>
          </b-col>
          <b-col sm="6"/>
          <b-col sm="2">
            <b-btn
              :disabled="!canComplete"
              variant="primary"
              size="sm"
              v-text="$t('common.Register')"
              @click="submit">
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
import { validator } from '@/plugins'

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
      return validator.validPassword(this.password)
    },
    samePassword: function () {
      return validator.samePassword(this.password, this.confirmPassword)
    },
    validUsername: function () {
      return validator.validUsername(this.username)
    },
    validEmail: function() {
      return validator.validEmail(this.email)
    }
  },
  methods: {
    //real methods
    submit: function () {
      if (!this.username) {
        this.errorMsg = this.$t('errors.ErrFieldRequired', { 'field': this.$t('common.Username') })
        return
      }
      if (!this.email) {
        this.errorMsg = this.$t('errors.ErrFieldRequired', { 'field': this.$t('common.Email') })
        return
      }
      if (!this.password) {
        this.errorMsg = this.$t('errors.ErrFieldRequired', { 'field': this.$t('common.Password') })
        return
      }
      if (!validator.validPassword(this.password)) {
        return
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
        let msg = 'errors.ErrUnknownError'
        if (error && error.response && error.response.data.error) {
          if (error.response.data.error.code) {
            msg = 'errors.' + error.response.data.error.code
          } else {
            msg = error.response.data.error.msg
          }
        }

        data.errorMsg = data.$(msg)
        data.registerDisabled = false
      })
    }
  }
}
</script>
