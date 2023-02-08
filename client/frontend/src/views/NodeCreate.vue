<script setup>
import { ref, inject } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import TextField from '@/components/ui/TextField.vue'
import Toggle from '@/components/ui/Toggle.vue'

const api = inject('api')
const toast = inject('toast')
const { t } = useI18n()
const router = useRouter()

const withPrivateHost = ref(false)
const name = ref('')
const publicHost = ref('')
const publicPort = ref('8080')
const privateHost = ref('')
const privatePort = ref('8080')
const sftpPort = ref('5657')

function canCreate() {
  if (!name.value) return false
  if (!publicHost.value) return false
  if (!publicPort.value) return false
  if (!sftpPort.value) return false
  if (withPrivateHost.value) {
    if (!privateHost.value) return false
    if (!privatePort.value) return false
  }
  return true
}

async function create() {
  if (!canCreate()) return
  const node = {
    name: name.value,
    publicHost: publicHost.value,
    publicPort: publicPort.value,
    sftpPort: sftpPort.value
  }
  if (withPrivateHost.value) {
    node.privateHost = privateHost.value
    node.privatePort = privatePort.value
  } else {
    node.privateHost = publicHost.value
    node.privatePort = publicPort.value
  }
  const id = await api.node.create(node)
  toast.success(t('nodes.Created'))
  router.push({ name: 'NodeView', params: { id }, query: { created: true } })
}
</script>

<template>
  <div class="nodecreate">
    <h1 v-text="t('nodes.Create')" />
    <text-field v-model="name" autofocus class="name" :label="t('common.Name')" />
    <text-field v-model="publicHost" class="public-host" :label="t('nodes.PublicHost')" />
    <text-field v-model="publicPort" class="public-port" :label="t('nodes.PublicPort')" type="number" />
    <toggle v-model="withPrivateHost" class="private-toggle" :label="t('nodes.WithPrivateAddress')" :hint="t('nodes.WithPrivateAddressHint')" />
    <text-field v-if="withPrivateHost" v-model="privateHost" class="private-host" :label="t('nodes.PrivateHost')" />
    <text-field v-if="withPrivateHost" v-model="privatePort" class="private-port" :label="t('nodes.PrivatePort')" type="number" />
    <text-field v-model="sftpPort" class="sftp-port" :label="t('nodes.SftpPort')" type="number" />
    <btn :disabled="!canCreate()" color="primary" @click="create()"><icon name="save" />{{ t('nodes.Create') }}</btn>
  </div>
</template>
