<script setup>
import { ref, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import TextField from '@/components/ui/TextField.vue'
import Toggle from '@/components/ui/Toggle.vue'

const { t } = useI18n()
const toast = inject('toast')

const users = ref([])
const newEmail = ref('')

const perms = [
  { label: t('scopes.ServersEdit'), name: 'editServerData' },
  { label: t('scopes.ServersInstall'), name: 'installServer' },
  { label: t('scopes.ServersConsole'), name: 'viewServerConsole' },
  { label: t('scopes.ServersConsoleSend'), name: 'sendServerConsole' },
  { label: t('scopes.ServersStop'), name: 'stopServer' },
  { label: t('scopes.ServersStart'), name: 'startServer' },
  { label: t('scopes.ServersStat'), name: 'viewServerStats' },
  { label: t('scopes.ServersFiles'), name: 'sftpServer' },
  { label: t('scopes.ServersFilesGet'), name: 'viewServerFiles' },
  { label: t('scopes.ServersFilesPut'), name: 'putServerFiles' },
  { label: t('scopes.ServersEditUsers'), name: 'editServerUsers' }
]

const props = defineProps({
  server: { type: Object, required: true }
})

async function sendInvite() {
  const newUser = { email: newEmail.value }
  await props.server.updateUser(newUser)
  toast.success(t('users.UserInvited'))
  loadUsers()
}

async function togglePerm(user, perm) {
  user[perm] = !user[perm]
  await props.server.updateUser(user)
  toast.success(t('users.UpdateSuccess'))
}

async function deleteUser(user) {
  await props.server.deleteUser(user.email)
  loadUsers()
}

async function loadUsers() {
  const u = await props.server.getUsers()
  u.map(user => {
    perms.map(p => {
      user[p.name] = false
    })
  })
  users.value = u
}

onMounted(async () => {
  loadUsers()
})
</script>

<template>
  <div>
    <h2 v-text="t('users.Users')" />
    <div v-for="user in users" :key="user.email" :class="['user', user.open ? 'open' : 'closed']">
      <h3 @click="user.open = !user.open" v-text="user.username" />
      <toggle v-for="perm in perms" :key="perm.name" v-model="user[perm.name]" :label="perm.label" @click="togglePerm(user, perm.name)" />
      <btn color="error" @click="deleteUser(user)" v-text="t('users.Delete')" />
    </div>
    <div v-if="users.length === 0" class="no-users" v-text="t('servers.NoUsers')" />
    <div class="invite">
      <text-field v-model="newEmail" type="email" icon="email" :label="t('users.Email')" />
      <btn color="primary" @click="sendInvite()"><icon name="plus" />{{ t('servers.InviteUser') }}</btn>
    </div>
  </div>
</template>
