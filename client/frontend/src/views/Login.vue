<script setup>
import { ref, inject } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Overlay from '@/components/ui/Overlay.vue'
import TextField from '@/components/ui/TextField.vue'
import Btn from '@/components/ui/Btn.vue'
import defaultRoute from '@/router/defaultRoute'

const { t } = useI18n()
const api = inject('api')
const events = inject('events')
const validate = inject('validate')
const router = useRouter()

const email = ref('')
const emailError = ref(false)
const password = ref('')
const passwordError = ref(false)
const otpNeeded = ref(false)
const token = ref('')

function loggedIn() {
  try {
    router.push(JSON.parse(sessionStorage.getItem('returnTo')))
    sessionStorage.removeItem('returnTo')
  } catch {
    router.push(defaultRoute(api))
  }
  events.emit('login')
}

async function login() {
  const res = await api.auth.login(email.value, password.value)
  if (res === true) {
    loggedIn()
  } else if (res === 'otp') {
    otpNeeded.value = true
  }
}

function resetOtp() {
  otpNeeded.value = false
  token.value = ''
}

async function submitOtp() {
  await api.auth.loginOtp(token.value)
  loggedIn()
}

function validateEmail(onChange = false) {
  if (!validate.email(email.value)) {
    if (onChange === true) return
    emailError.value = true
  } else {
    emailError.value = false
  }
}

function emailErrorMsg() {
  if (emailError.value) return t('errors.ErrEmailInvalid')
}

function validatePassword(onChange = false) {
  if (!validate.password(password.value)) {
    if (onChange === true) return
    passwordError.value = true
  } else {
    passwordError.value = false
  }
}

function passwordErrorMsg() {
  if (passwordError.value) return t('errors.ErrPasswordRequirements')
}
</script>

<template>
  <div class="login">
    <h1 v-text="t('users.Login')" />
    <form @keydown.enter="login()">
      <text-field v-model="email" type="email" name="email" :label="t('users.Email')" :error="emailErrorMsg()" icon="email" autofocus @blur="validateEmail" @change="validateEmail(true)" />
      <text-field v-model="password" type="password" name="password" :label="t('users.Password')" :error="passwordErrorMsg()" icon="lock" @blur="validatePassword" @change="validatePassword(true)" />
      <btn color="primary" :disabled="emailError || passwordError" @click="login()" v-text="t('users.Login')" />
      <btn v-if="$config.registrationEnabled" variant="text" @click="$router.push({ name: 'Register' })" v-text="t('users.RegisterLink')" />
    </form>
    <overlay v-model="otpNeeded" :title="t('users.OtpNeeded')" closable @close="resetOtp()">
      <text-field v-model="token" autofocus />
      <btn color="primary" @click="submitOtp()" v-text="t('users.Login')" />
    </overlay>
  </div>
</template>
