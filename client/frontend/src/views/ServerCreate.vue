<script setup>
import { ref, inject } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import SelectTemplate from '@/components/ui/serverCreateSteps/SelectTemplate.vue'
import Environment from '@/components/ui/serverCreateSteps/Environment.vue'
import Users from '@/components/ui/serverCreateSteps/Users.vue'
import Settings from '@/components/ui/serverCreateSteps/Settings.vue'

const router = useRouter()
const { t } = useI18n()
const api = inject('api')
const step = ref('template')
const template = ref({})
const environment = ref({})
const users = ref([])

function templateSelected(selected) {
  template.value = selected
  step.value = 'environment'
}

function envBack() {
  template.value = {}
  environment.value = {}
  step.value = 'template'
}

function envConfirmed(name, nodeId, envConfig) {
  environment.value = { name, nodeId, envConfig }
  step.value = 'users'
}

function usersBack() {
  environment.value = {}
  users.value = []
  step.value = 'environment'
}

function usersConfirmed(selected) {
  users.value = selected
  step.value = 'settings'
}

function settingsBack() {
  users.value = []
  step.value = 'users'
}

async function settingsConfirmed(settings) {
  // last step confirmed, create server
  const request = template.value
  request.name = environment.value.name
  request.node = environment.value.nodeId
  request.environment = environment.value.envConfig
  request.users = users.value
  request.data = {}
  for (const setting in settings) {
    request.data[setting] = settings[setting]

    // fix value types
    if (request.data[setting].type === 'boolean') {
      request.data[setting].value =
        request.data[setting].value !== 'false' &&
        request.data[setting].value !== false
    }

    if (request.data[setting].type === 'integer') {
      request.data[setting].value = Number(request.data[setting].value)
    }
  }

  const id = await api.server.create(request)
  router.push({ name: 'ServerView', params: { id }, query: { created: true } })
}
</script>

<template>
  <div class="servercreate">
    <h1 v-text="t('servers.Create')" />
    <div :class="['progress', 'on-step-' + step]">
      <div :class="['step', 'step-template', step === 'template' ? 'step-current' : '']" />
      <div :class="['step', 'step-environment', step === 'environment' ? 'step-current' : '']" />
      <div :class="['step', 'step-users', step === 'users' ? 'step-current' : '']" />
      <div :class="['step', 'step-settings', step === 'settings' ? 'step-current' : '']" />
    </div>
    <select-template v-if="step === 'template'" @selected="templateSelected" />
    <environment v-if="step === 'environment'" :env="template.environment" :supported="template.supportedEnvironments" @back="envBack()" @confirm="envConfirmed" />
    <users v-if="step === 'users'" @back="usersBack()" @confirm="usersConfirmed" />
    <settings v-if="step === 'settings'" :data="template.data" @back="settingsBack()" @confirm="settingsConfirmed" />
  </div>
</template>
