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
const toast = inject('toast')
const router = useRouter()

const email = ref('')
const password = ref('')
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
</script>

<template>
  <div class="login">
    <h1 v-text="t('users.Login')" />
    <form @keydown.enter="login()">
      <text-field v-model="email" type="email" name="email" :label="t('users.Email')" icon="email" autofocus />
      <text-field v-model="password" type="password" name="password" :label="t('users.Password')" icon="lock" />
      <btn color="primary" @click="login()" v-text="t('users.Login')" />
      <btn v-if="$config.registrationEnabled" variant="text" @click="$router.push({ name: 'Register' })" v-text="t('users.RegisterLink')" />
    </form>
    <overlay v-model="otpNeeded" :title="t('users.OtpNeeded')" closable @close="resetOtp()">
      <text-field v-model="token" autofocus />
      <btn color="primary" @click="submitOtp()" v-text="t('users.Login')" />
    </overlay>
  </div>
</template>
