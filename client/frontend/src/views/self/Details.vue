<script setup>
import { ref, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import TextField from '@/components/ui/TextField.vue'
import Loader from '@/components/ui/Loader.vue'

const { t } = useI18n()
const api = inject('api')
const toast = inject('toast')

const loading = ref(false)
const username = ref('')
const email = ref('')
const password = ref('')

onMounted(async () => {
  loading.value = true
  const data = await api.self.get()
  username.value = data.username
  email.value = data.email
  loading.value = false
})

function isValidUsername(u) {
  return u.length >= 5
}

function isValidEmail(e) {
  return e.match(/.+@.+\..{2,}/)
}

function isValidPassword(p) {
  return p.length >= 8
}

function canSubmit() {
  return isValidUsername(username.value) && isValidEmail(email.value) && isValidPassword(password.value)
}

async function submit() {
  if (!canSubmit()) return
  await api.self.updateDetails(username.value, email.value, password.value)
  toast.success(t('users.InfoChanged'))
}
</script>

<template>
  <div v-if="loading" class="accountdetails loading"><loader /></div>
  <div v-else class="accountdetails">
    <h1 v-text="t('users.ChangeInfo')" />
    <form>
      <text-field v-model="username" icon="account" :label="t('users.Username')" />
      <text-field v-model="email" icon="email" type="email" :label="t('users.Email')" />
      <text-field v-model="password" icon="lock" type="password" :label="t('users.ConfirmPassword')" />
      <btn :disabled="!canSubmit()" color="primary" @click="submit()"><icon name="save" />{{ t('users.ChangeInfo') }} </btn>
    </form>
  </div>
</template>
