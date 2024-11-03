<script setup>
import { ref, inject, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import TextField from '@/components/ui/TextField.vue'
import Toggle from '@/components/ui/Toggle.vue'

const { t, te, locale } = useI18n()
const route = useRoute()
const router = useRouter()
const api = inject('api')
const toast = inject('toast')
const validate = inject('validate')
const events = inject('events')

const username = ref('')
const email = ref('')
const password = ref('')
const otpactive = ref('')

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
  await api.user.updatePermissions(route.params.id, { scopes: permissions.value })
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

const scopes = {
  general: [
    'admin',
    'login',
    'self.edit',
    'self.clients',
    'settings.edit'
  ],
  servers: [
    'server.create'
  ],
  nodes: [
    'nodes.view',
    'nodes.create',
    'nodes.edit',
    'nodes.deploy',
    'nodes.delete'
  ],
  users: [
    'users.info.search',
    'users.info.view',
    'users.info.edit',
    'users.perms.view',
    'users.perms.edit'
  ],
  templates: [
    'templates.view',
    'templates.local.edit',
    'templates.repo.view',
    'templates.repo.add',
    'templates.repo.remove'
  ]
}

const permissions = ref([])

onMounted(async () => {
  const user = await api.user.get(route.params.id)
  username.value = user.username
  email.value = user.email
  permissions.value = await api.user.getPermissions(route.params.id)
  otpactive.value = user.otpactive !== undefined ? user.otpactive : false
})

function scopeLabel(scope) {
  return t('scopes.name.' + scope.replace(/\./g, '-'))
}

function scopeHint(scope) {
  if (te('scopes.hint.' + scope.replace(/\./g, '-'), locale)) {
    return t('scopes.hint.' + scope.replace(/\./g, '-'))
  }
}

function permissionCategoryHeading(category) {
  const map = {
    servers: 'servers.Servers',
    nodes: 'nodes.Nodes',
    users: 'users.Users',
    templates: 'templates.Templates'
  }
  return t(map[category])
}

function permissionDisabled(scope) {
  if (!api.auth.hasScope('user.perms.edit')) return true
  if (scope === 'admin' && api.auth.hasScope('admin')) return false
  if (scope === 'admin') return true

  if (permissions.value.indexOf('admin') >= 0) return true
  return !api.auth.hasScope(scope)
}

function togglePermission(scope) {
  if (permissions.value.indexOf(scope) === -1) {
    permissions.value = [...permissions.value, scope]
  } else {
    permissions.value = permissions.value.filter(e => e !== scope)
  }
}
</script>

<template>
  <div class="userview">
    <div v-if="$api.auth.hasScope('users.info.view')" class="details">
      <h1 v-text="t('users.Details')" />
      <form @submit.prevent="submitDetails()">
        <text-field
          v-model="username"
          :label="t('users.Username')"
          icon="account"
          :error="usernameError"
          :disabled="!$api.auth.hasScope('users.info.edit')"
          @blur="usernameError = validate.username(username) ? '' : t('errors.ErrUsernameRequirements')"
        />
        <text-field
          v-model="email"
          :label="t('users.Email')"
          type="email"
          icon="email"
          :error="emailError"
          :disabled="!$api.auth.hasScope('users.info.edit')"
          @blur="emailError = validate.email(email) ? '' : t('errors.ErrEmailInvalid')"
        />
        <text-field
          v-if="$api.auth.hasScope('users.info.edit')"
          v-model="password"
          :label="t('users.Password')"
          type="password"
          icon="lock"
          :error="passwordError"
          @blur="passwordError = (validate.password(password) || password.length === 0) ? '' : t('error.PasswordInvalid')"
        />
        <h2>{{ t('users.Authentication') }}: {{ otpactive ? t('users.TwoFactor') : t('users.Password') }}</h2>
        <btn v-if="$api.auth.hasScope('users.info.edit')" color="primary" :disabled="!canSubmitDetails()" @click="submitDetails()"><icon name="save" />{{ t('users.UpdateDetails') }}</btn>
        <btn v-if="$api.auth.hasScope('users.info.edit')" color="error" @click="deleteUser()"><icon name="remove" />{{ t('users.Delete') }}</btn>
      </form>
    </div>

    <div v-if="$api.auth.hasScope('users.perms.view')" class="permissions">
      <h1 v-text="t('users.Permissions')" />
      <div v-for="(scopeCat, catName) in scopes" :key="scopeCat">
        <h3 v-if="catName !== 'general'" v-text="permissionCategoryHeading(catName)" />
        <toggle
          v-for="scope in scopeCat"
          :key="scope"
          :model-value="permissions.indexOf(scope) >= 0"
          :disabled="permissionDisabled(scope)"
          :label="scopeLabel(scope)"
          :hint="scopeHint(scope)"
          @update:modelValue="togglePermission(scope)"
        />
      </div>
      <btn v-if="api.auth.hasScope('user.perms.edit')" color="primary" @click="submitPermissions()"><icon name="save" />{{ t('users.UpdatePermissions') }}</btn>
    </div>
  </div>
</template>
