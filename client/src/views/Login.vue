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
        v-text="$t('common.Login')"
        slot="header"
        class="mb-0"
        align="center"></h6>
      <b-form>
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
        <b-row class="my-1">
          <b-col sm="1"/>
          <b-col sm="10">
            <b-form-input
              id="password"
              v-model="password"
              prepend-icon="mdi-lock"
              name="password"
              :placeholder="$t('common.Password')"
              type="password"
              @keyup.enter="submit"/>
          </b-col>
        </b-row>
      </b-form>
      <div slot="footer">
        <b-row>
          <b-col sm="1"/>
          <b-col sm="2">
            <b-link :to="{name: 'Register'}" v-text="$t('common.Register')"></b-link>
          </b-col>
          <b-col sm="6"/>
          <b-col sm="2">
            <b-btn
              :disabled="!canSubmit"
              v-text="$t('common.Login')"
              variant="primary"
              size="sm"
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
import Cookies from 'js-cookie'

export default {
  data () {
    return {
      email: '',
      password: '',
      loginDisabled: false,
      errorMsg: '',
      dismissCountDown: 5
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
        data.errorMsg = this.$t('errors.ErrFieldRequired', { 'field': this.$t('common.Email') })
        return
      }

      if (!data.password) {
        data.errorMsg = this.$t('errors.ErrFieldRequired', { 'field': this.$t('common.Password') })
        return
      }

      this.loginDisabled = true

      this.axios.post(this.$route.path, {
        data: {
          email: this.email,
          password: this.password
        }
      }).then(function (response) {
        let responseData = response.data
        if (responseData.success) {
          Cookies.set('puffer_auth', responseData.data)
          data.$router.push({ name: 'Servers' })
        } else {
          data.errorMsg = responseData.msg + ''
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

        data.errorMsg = data.$t(msg)
      }).finally(function () {
        data.loginDisabled = false
      })
    }
  }
}
</script>
