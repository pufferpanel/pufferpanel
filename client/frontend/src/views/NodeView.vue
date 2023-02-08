<script setup>
import { ref, inject, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import markdown from '@/utils/markdown'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import Overlay from '@/components/ui/Overlay.vue'
import TextField from '@/components/ui/TextField.vue'
import Toggle from '@/components/ui/Toggle.vue'

const api = inject('api')
const toast = inject('toast')
const events = inject('events')
const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const deploymentOpen = ref(false)
let deploymentData = {}
const withPrivateHost = ref(false)
const name = ref('')
const publicHost = ref('')
const publicPort = ref('8080')
const privateHost = ref('')
const privatePort = ref('8080')
const sftpPort = ref('5657')
const currentStep = ref(1)

onMounted(async () => {
  const node = await api.node.get(route.params.id)
  name.value = node.name
  publicHost.value = node.publicHost
  publicPort.value = node.publicPort
  privateHost.value = node.privateHost
  privatePort.value = node.privatePort
  sftpPort.value = node.sftpPort
  withPrivateHost.value = !(node.publicHost === node.privateHost && node.publicPort === node.privatePort)
  deploymentData = await api.node.deployment(route.params.id)
  if (route.query.created) {
    deploymentOpen.value = true
  }
})

function canSubmit() {
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

async function submit() {
  if (!canSubmit()) return
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
  const id = await api.node.update(node)
  toast.success(t('nodes.Updated'))
}

async function deleteNode() {
  events.emit(
    'confirm',
    t('nodes.ConfirmDelete', { name: name.value }),
    {
      text: t('nodes.Delete'),
      icon: 'remove',
      color: 'error',
      action: async () => {
        await api.node.delete(route.params.id)
        toast.success(t('nodes.Deleted'))
        router.push({ name: 'NodeList' })
      }
    },
    {
      color: 'primary'
    }
  )
}

function getDeployConfig() {
  const config = {
    logs: '/var/log/pufferpanel',
    web: {
      host: `0.0.0.0:${privatePort.value}`
    },
    token: {
      public: location.origin + '/auth/publickey'
    },
    panel: {
      enable: false
    },
    daemon: {
      auth: {
        url: location.origin + '/oauth2/token',
        ...deploymentData
      },
      data: {
        cache: '/var/lib/pufferpanel/cache',
        servers: '/var/lib/pufferpanel/servers'
      },
      sftp: {
        host: `0.0.0.0:${sftpPort.value}`
      }
    }
  }
  return JSON.stringify(config, undefined, 2)
}
</script>

<template>
  <div class="nodeview">
    <h1 v-text="t('nodes.Edit')" />
    <text-field v-model="name" class="name" :label="t('common.Name')" />
    <text-field v-model="publicHost" class="public-host" :label="t('nodes.PublicHost')" />
    <text-field v-model="publicPort" class="public-port" :label="t('nodes.PublicPort')" type="number" />
    <toggle v-model="withPrivateHost" class="private-toggle" :label="t('nodes.WithPrivateAddress')" :hint="t('nodes.WithPrivateAddressHint')" />
    <text-field v-if="withPrivateHost" v-model="privateHost" class="private-host" :label="t('nodes.PrivateHost')" />
    <text-field v-if="withPrivateHost" v-model="privatePort" class="private-port" :label="t('nodes.PrivatePort')" type="number" />
    <text-field v-model="sftpPort" class="sftp-port" :label="t('nodes.SftpPort')" type="number" />
    <btn :disabled="!canSubmit()" color="primary" @click="submit()"><icon name="save" />{{ t('nodes.Update') }}</btn>
    <btn color="error" @click="deleteNode()"><icon name="remove" />{{ t('nodes.Delete') }}</btn>
    <btn @click="deploymentOpen = true" v-text="t('nodes.Deploy')" />
    <overlay v-model="deploymentOpen" closable :title="t('nodes.Deploy')" @close="currentStep = 1">
      <!-- eslint-disable-next-line vue/no-v-html -->
      <div v-html="markdown(t(`nodes.deploy.Step${currentStep}`, { config: getDeployConfig() }))" />
      <btn v-if="currentStep < 5" @click="currentStep += 1" v-text="t('common.Next')" />
      <btn v-else @click="deploymentOpen = false; currentStep = 1" v-text="t('common.Close')" />
    </overlay>
  </div>
</template>
