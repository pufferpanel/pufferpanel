<script setup>
import { ref, inject } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import TextField from '@/components/ui/TextField.vue'
import Btn from '@/components/ui/Btn.vue'
import defaultRoute from '@/router/defaultRoute'

const { t } = useI18n()
const api = inject('api')
const events = inject('events')
const validate = inject('validate')
const router = useRouter()
const route = useRoute()

const username = ref('')
const email = ref(route.query.email)
const password = ref('')
const confirmPassword = ref('')

const errorUsername = ref('')
const errorEmail = ref('')
const errorPassword = ref('')
const errorConfirmPassword = ref('')

function isValidConfirmPassword() {
  return confirmPassword.value === password.value
}

function checkUsername() {
  if (!validate.username(username.value)) errorUsername.value = t('errors.ErrUsernameRequirements')
}

function checkEmail() {
  if (!validate.email(email.value)) errorEmail.value = t('errors.ErrEmailInvalid')
}

function checkPassword() {
  if (!validate.password(password.value)) errorPassword.value = t('errors.ErrPasswordRequirements')
}

function checkConfirmPassword() {
  if (!isValidConfirmPassword()) errorConfirmPassword.value = t('errors.ErrPasswordsNotIdentical')
}

function canSubmit() {
  return validate.username(username.value) &&
    validate.email(email.value) &&
    validate.password(password.value) &&
    isValidConfirmPassword()
}

async function register() {
  if (!canSubmit()) return

  await api.auth.login(email.value, route.query.token)
  await api.self.updateDetails(username.value, email.value, route.query.token)
  await api.self.changePassword(route.query.token, password.value)
  router.push(defaultRoute(api))
  events.emit('login')
}
</script>

<template>
  <div class="register">
    <h1 v-text="t('users.Register')" />
    <text-field
      v-model="username"
      autofocus
      :label="t('users.Username')"
      icon="account"
      :error="errorUsername"
      @change="errorUsername = ''"
      @blur="checkUsername()"
    />
    <text-field
      v-model="email"
      disabled
      type="email"
      :label="t('users.Email')"
      icon="email"
      :error="errorEmail"
      @change="errorEmail = ''"
      @blur="checkEmail()"
    />
    <text-field
      v-model="password"
      type="password"
      :label="t('users.Password')"
      icon="lock"
      :error="errorPassword"
      @change="errorPassword = ''"
      @blur="checkPassword()"
    />
    <text-field
      v-model="confirmPassword"
      type="password"
      :label="t('users.ConfirmPassword')"
      icon="lock"
      :error="errorConfirmPassword"
      @change="errorConfirmPassword = ''"
      @blur="checkConfirmPassword()"
    />
    <btn color="primary" :disabled="!canSubmit()" @click="register()" v-text="t('users.Register')" />
    <btn variant="text" @click="$router.push({ name: 'Login' })" v-text="t('users.LoginLink')" />
  </div>
</template>
