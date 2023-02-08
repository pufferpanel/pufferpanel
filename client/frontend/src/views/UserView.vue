<script setup>
import { ref, inject, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import TextField from '@/components/ui/TextField.vue'
import Toggle from '@/components/ui/Toggle.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const api = inject('api')
const toast = inject('toast')
const validate = inject('validate')
const events = inject('events')

const username = ref('')
const email = ref('')
const password = ref('')

const usernameError = ref('')
const emailError = ref('')
const passwordError = ref('')

function canSubmitDetails() {
  return validate.username(username.value) &&
    validate.email(email.value) &&
    (validate.password(password.value) || password.value.length === 0)
}

async function submitDetails() {
  await api.user.update(route.params.id, {
    username: username.value,
    email: email.value,
    password: password.value || undefined
  })
  toast.success(t('users.UpdateSuccess'))
}

async function submitPermissions() {
  if (!canSubmitDetails()) return false
  await api.user.updatePermissions(route.params.id, permissions.value)
  toast.success(t('users.UpdateSuccess'))
}

async function deleteUser() {
  events.emit(
    'confirm',
    t('users.ConfirmDelete', { name: username.value }),
    {
      text: t('users.Delete'),
      icon: 'remove',
      color: 'error',
      action: async () => {
        await api.user.delete(route.params.id)
        toast.success(t('users.DeleteSuccess'))
        router.push({ name: 'UserList' })
      }
    },
    {
      color: 'primary'
    }
  )
}

const scopes = [
  'viewServers',
  'createServers',
  'deleteServers',
  'editServerAdmin',
  'viewNodes',
  'editNodes',
  'deployNodes',
  'panelSettings',
  'viewTemplates',
  'editTemplates',
  'viewUsers',
  'editUsers'
]

const permissions = ref({ admin: false })

scopes.map(scope => {
  permissions.value[scope] = false
})

onMounted(async () => {
  const user = await api.user.get(route.params.id)
  username.value = user.username
  email.value = user.email
  const perms = await api.user.getPermissions(route.params.id)
  for (const scope in perms) {
    permissions.value[scope] = perms[scope]
  }
})

function scopeLabel(scope) {
  return t('scopes.' + scope[0].toUpperCase() + scope.substring(1))
}
</script>

<template>
  <div class="userview">
    <div class="details">
      <h1 v-text="t('users.Details')" />
      <form @submit.prevent="submitDetails()">
        <text-field
          v-model="username"
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
          @blur="passwordError = (validate.password(password) || password.length === 0) ? '' : t('error.PasswordInvalid')"
        />
        <btn color="primary" :disabled="!canSubmitDetails()" @click="submitDetails()"><icon name="save" />{{ t('users.UpdateDetails') }}</btn>
        <btn color="error" @click="deleteUser()"><icon name="remove" />{{ t('users.Delete') }}</btn>
      </form>
    </div>
    <div class="permissions">
      <h1 v-text="t('users.Permissions')" />
      <toggle v-model="permissions.admin" :label="t('scopes.Admin')" />
      <toggle v-for="scope in scopes" :key="scope" v-model="permissions[scope]" :disabled="permissions.admin" :label="scopeLabel(scope)" />
      <btn color="primary" @click="submitPermissions()"><icon name="save" />{{ t('users.UpdatePermissions') }}</btn>
    </div>
  </div>
</template>
