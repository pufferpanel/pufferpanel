<script setup>
import { ref, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import TextField from '@/components/ui/TextField.vue'
import Toggle from '@/components/ui/Toggle.vue'

const { t, te, locale } = useI18n()
const toast = inject('toast')

const users = ref([])
const newEmail = ref('')

const perms = [
  'server.view',
	'server.admin',
	'server.delete',
	'server.definition.view',
	'server.definition.edit',
	'server.data.view',
	'server.data.edit',
	'server.flags.view',
	'server.flags.edit',
	'server.name.edit',
	'server.users.view',
	'server.users.create',
	'server.users.edit',
	'server.users.delete',
	//'server.tasks.view',
	//'server.tasks.run',
	//'server.tasks.create',
	//'server.tasks.delete',
	//'server.tasks.edit',
	'server.start',
	'server.stop',
	'server.kill',
	'server.install',
	'server.files.view',
	'server.files.edit',
	'server.sftp',
	'server.console',
	'server.console.send',
	'server.stats',
	'server.status'
].map(scope => {
  const res = {
    label: t('scopes.name.' + scope.replace(/\./g, '-')),
    name: scope
  }
  if (te('scopes.hint.' + scope.replace(/\./g, '-'), locale))
    res.hint = t('scopes.hint.' + scope.replace(/\./g, '-'))
  return res
})

const props = defineProps({
  server: { type: Object, required: true }
})

async function sendInvite() {
  const newUser = { email: newEmail.value }
  await props.server.updateUser(newUser)
  toast.success(t('users.UserInvited'))
  loadUsers()
}

async function updatePerms(user) {
  const scopes = Object.keys(user.scopes).filter(p => user.scopes[p])
  const update = { ...user, scopes }
  await props.server.updateUser(update)
  toast.success(t('users.UpdateSuccess'))
}

async function deleteUser(user) {
  await props.server.deleteUser(user.email)
  loadUsers()
}

async function loadUsers() {
  const u = await props.server.getUsers()
  users.value = u.map(user => {
    const scopes = {}
    perms.map(p => {
      scopes[p.name] = user.scopes.indexOf(p.name) > -1
    })
    user.scopes = scopes
    return user
  })
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
      <toggle
        v-for="perm in perms"
        :key="perm.name"
        v-model="user.scopes[perm.name]"
        :disabled="!server.hasScope('server.users.edit')"
        :label="perm.label"
        :hint="perm.hint"
        @update:modelValue="updatePerms(user)"
      />
      <btn v-if="server.hasScope('server.users.delete')" color="error" @click="deleteUser(user)" v-text="t('users.Delete')" />
    </div>
    <div v-if="users.length === 0" class="no-users" v-text="t('servers.NoUsers')" />
    <div v-if="server.hasScope('server.users.create')" class="invite">
      <text-field v-model="newEmail" type="email" icon="email" :label="t('users.Email')" />
      <btn color="primary" @click="sendInvite()"><icon name="plus" />{{ t('servers.InviteUser') }}</btn>
    </div>
  </div>
</template>
