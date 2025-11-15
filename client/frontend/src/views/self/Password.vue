<script setup>
import { ref, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import TextField from '@/components/ui/TextField.vue'

const { t } = useI18n()
const api = inject('api')
const toast = inject('toast')

const oldPass = ref('')
const newPass = ref('')
const confirmPass = ref('')

function isValidPassword(p) {
  return p.length >= 8
}

function canSubmit() {
  return isValidPassword(oldPass.value) && isValidPassword(newPass.value) && newPass.value === confirmPass.value
}

async function submit() {
  if (!canSubmit()) return
  await api.self.changePassword(oldPass.value, newPass.value)
  toast.success(t('users.PasswordChanged'))
}
</script>

<template>
  <div class="changepassword">
    <h1 v-text="t('users.ChangePassword')" />
    <form>
      <text-field v-model="oldPass" icon="lock" type="password" :label="t('users.OldPassword')" />
      <text-field v-model="newPass" icon="lock" type="password" :label="t('users.NewPassword')" />
      <text-field v-model="confirmPass" icon="lock" type="password" :label="t('users.ConfirmPassword')" />
      <btn :disabled="!canSubmit()" color="primary" @click="submit()"><icon name="save" />{{ t('users.ChangePassword') }}</btn>
    </form>
  </div>
</template>
