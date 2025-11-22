<script setup>
import { ref, inject } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import TextField from '@/components/ui/TextField.vue'
import OtpInput from '@/components/ui/OtpInput.vue'
import Loader from '@/components/ui/Loader.vue'
import Btn from '@/components/ui/Btn.vue'
import defaultRoute from '@/router/defaultRoute'

const { t } = useI18n()
const api = inject('api')
const events = inject('events')
const validate = inject('validate')
const router = useRouter()

// isSecureContext is true on localhost, but localhost isn't allowed for webauthn, using location instead
const passkeysSupported = location.protocol === 'https:' && (location.port === "" || location.port === "443") && window.navigator && navigator.credentials

const loading = ref(false)
const step = ref('email')
const email = ref('')
const emailError = ref(false)
const password = ref('')
const passwordError = ref(false)
const otpRecovery = ref(false)
const token = ref('')
const secondFactorData = ref(null)

function loggedIn() {
  try {
    router.push(JSON.parse(sessionStorage.getItem('returnTo')))
    sessionStorage.removeItem('returnTo')
  } catch {
    router.push(defaultRoute(api))
  }
  events.emit('login')
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

async function attemptPasskeyLogin() {
  loading.value = true
  try {
    const res = await api.auth.passkeyLogin(email.value)
    return res !== false
  } finally {
    loading.value = false
  }
}

async function attemptPasskey2fa(doRedirect = false) {
  loading.value = true
  try {
    const attestation = await navigator.credentials.get(secondFactorData.value.webauthnChallenge)
    const done = await api.auth.validatePasskeyLogin(attestation)
    if (doRedirect && done) loggedIn()
    return done
  } finally {
    loading.value = false
  }
}

async function attemptPasswordLogin() {
  loading.value = true
  try {
    const res = await api.auth.login(email.value, password.value)
    secondFactorData.value = {
      otpEnabled: res.otpEnabled,
      webauthnChallenge: res.webauthnChallenge
    }
    return res.needsSecondFactor
  } finally {
    loading.value = false
  }
}

async function attemptOtpLogin() {
  loading.value = true
  try {
    await api.auth.loginOtp(token.value)
    return true
  } catch(_) {
    return false
  } finally {
    loading.value = false
  }
}

function canNextStep() {
  if (step.value === 'email') {
    return validate.email(email.value)
  } else if (step.value === 'password') {
    return validate.password(password.value)
  } else if (step.value === 'secondFactor') {
    return token.value.length >= 6
  }
}

async function nextStep() {
  if (!canNextStep()) return

  if (step.value === 'email') {
    if (passkeysSupported) {
      const done = await attemptPasskeyLogin()
      if (done) {
        loggedIn()
      } else {
        step.value = 'password'
      }
    } else {
      step.value = 'password'
    }
  } else if (step.value === 'password') {
    const needsSecondFactor = await attemptPasswordLogin()
    if (needsSecondFactor) {
      if (secondFactorData.value.webauthnChallenge) {
        try {
          const done = await attemptPasskey2fa()
          if (done) loggedIn()
        } catch (_) {
          // just catching to let the rest run
        }
      }
      step.value = 'secondFactor'
    } else {
      loggedIn()
    }
  } else if (step.value === 'secondFactor') {
    const done = await attemptOtpLogin()
    if (done) {
      loggedIn()
    }
  }
}
</script>

<template>
  <div :class="['login', step]">
    <h1 v-text="t('users.Login')" />
    <form @keydown.enter="nextStep()">
      <text-field
        v-model="email"
        type="email"
        name="email"
        :disabled="step !== 'email' || loading"
        :label="t('users.Email')"
        :error="emailErrorMsg()"
        icon="email"
        autofocus
        @blur="validateEmail"
        @change="validateEmail(true)"
      />
      <text-field
        v-if="step !== 'email'"
        v-model="password"
        type="password"
        name="password"
        :disabled="step !== 'password' || loading"
        :label="t('users.Password')"
        :error="passwordErrorMsg()"
        icon="lock"
        autofocus
        @blur="validatePassword"
        @change="validatePassword(true)"
      />
      <div v-if="step === 'secondFactor' && secondFactorData.otpEnabled">
        <otp-input
          v-if="!otpRecovery"
          :disabled="loading"
          @update:modelValue="token = $event"
          @complete="token = $event; nextStep()"
        />
        <text-field
          v-else
          v-model="token"
          autofocus
          :disabled="loading"
        />
        <btn
          :disabled="loading"
          variant="text"
          @click="otpRecovery = !otpRecovery; token = ''"
          v-text="otpRecovery ? t('users.OtpUseAuthenticator') : t('users.OtpUseRecovery')"
        />
      </div>
      <btn
        v-if="step === 'secondFactor' && secondFactorData.webauthnChallenge && passkeysSupported"
        :disabled="loading"
        variant="text"
        @click="attemptPasskey2fa(true)"
        v-text="t('users.RetryPasskey')"
      />
      <btn color="primary" :disabled="loading || !canNextStep()" @click="nextStep()">
        <loader v-if="loading" />
        <span v-else v-text="t('common.Next')" />
      </btn>
      <btn
        v-if="$config.registrationEnabled"
        variant="text"
        @click="$router.push({ name: 'Register' })"
        v-text="t('users.RegisterLink')"
      />
    </form>
  </div>
</template>
