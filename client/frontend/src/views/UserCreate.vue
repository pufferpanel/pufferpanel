<script setup>
import { ref, inject } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import TextField from '@/components/ui/TextField.vue'

const { t } = useI18n()
const router = useRouter()
const api = inject('api')
const validate = inject('validate')

const username = ref('')
const email = ref('')
const password = ref('')

const usernameError = ref('')
const emailError = ref('')
const passwordError = ref('')

function canSubmit() {
  return validate.username(username.value) &&
    validate.email(email.value) &&
    validate.password(password.value)
}

async function submit() {
  if (!canSubmit()) return false
  const id = await api.user.create(username.value, email.value, password.value)
  router.push({ name: 'UserView', params: { id } })
}
</script>

<template>
  <div class="usercreate">
    <h1 v-text="t('users.Create')" />
    <form @submit.prevent="submit()">
      <text-field
        v-model="username"
        autofocus
        :label="t('users.Username')"
        icon="account"
        :error="usernameError"
        @blur="usernameError = validate.username(username) ? '' : t('errors.ErrUsernameRequirements')"
      />
      <text-field
        v-model="email"
        :label="t('users.Email')"
        type="email"
        icon="email"
        :error="emailError"
        @blur="emailError = validate.email(email) ? '' : t('errors.ErrEmailInvalid')"
      />
      <text-field
        v-model="password"
        :label="t('users.Password')"
        type="password"
        icon="lock"
        :error="passwordError"
        @blur="passwordError = validate.password(password) ? '' : t('errors.ErrPasswordRequirements')"
      />
      <btn color="primary" :disabled="!canSubmit()" @click="submit()"><icon name="save" />{{ t('users.Create') }}</btn>
    </form>
  </div>
</template>
