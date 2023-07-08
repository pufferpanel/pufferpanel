<script setup>
import { ref, inject, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import markdown from '@/utils/markdown'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import Loader from '@/components/ui/Loader.vue'
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
const featuresFetched = ref(null)
const features = ref({})

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

  fetchFeatures()
})

async function fetchFeatures() {
  featuresFetched.value = null
  features.value = {}
  try {
    const f = await api.node.features(route.params.id)
    console.log(f)
    features.value.envs = [ ...new Set(f.environments) ].map(e => t(`env.${e}.name`))
    features.value.docker = f.features.indexOf('docker') !== -1
    features.value.os = f.os
    features.value.arch = f.arch
    featuresFetched.value = true
  } catch(e) {
    featuresFetched.value = false
  }
}

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

function closeDeploy() {
  deploymentOpen.value = false
  currentStep.value = 1
  fetchFeatures()
}
</script>

<template>
  <div class="nodeview">
    <h1 v-text="name" />
    <loader v-if="featuresFetched === null" />
    <div v-else-if="featuresFetched === false" class="features">
      <div class="unreachable" v-text="t('nodes.Unreachable')" />
    </div>
    <div v-else class="features">
      <div class="reachable" v-text="t('nodes.Reachable')" />
      <div class="os">
        <span v-text="t('nodes.features.os.label')" />
        <span v-text="t('nodes.features.os.' + features.os)" />
      </div>
      <div class="arch">
        <span v-text="t('nodes.features.arch.label')" />
        <span v-text="t('nodes.features.arch.' + features.arch)" />
      </div>
      <div class="env">
        <span v-text="t('nodes.features.envs')" />
        <span v-text="features.envs.join(', ')" />
      </div>
      <div class="docker">
        <span v-text="t('env.docker.name')" />
        <span v-text="t('nodes.features.docker.' + features.docker)" />
      </div>
    </div>
    <h2 v-text="t('nodes.Edit')" />
    <div v-if="route.params.id > 0" class="edit">
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
    </div>
    <div v-else class="edit" v-html="markdown(t('nodes.LocalNodeEdit'))" />
    <overlay v-model="deploymentOpen" closable :title="t('nodes.Deploy')" @close="closeDeploy()">
      <!-- eslint-disable-next-line vue/no-v-html -->
      <div v-html="markdown(t(`nodes.deploy.Step${currentStep}`, { config: getDeployConfig() }))" />
      <btn v-if="currentStep < 5" @click="currentStep += 1" v-text="t('common.Next')" />
      <btn v-else @click="closeDeploy()" v-text="t('common.Close')" />
    </overlay>
  </div>
</template>
