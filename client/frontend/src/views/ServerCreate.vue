<script setup>
import { ref, inject } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import SelectTemplate from '@/components/ui/serverCreateSteps/SelectTemplate.vue'
import Environment from '@/components/ui/serverCreateSteps/Environment.vue'
import Settings from '@/components/ui/serverCreateSteps/Settings.vue'

const router = useRouter()
const { t } = useI18n()
const api = inject('api')
const step = ref('environment')
const environment = ref({})
const users = ref([])
const template = ref({})

function envConfirmed(name, nodeId, nodeOs, nodeArch, envConfig, u) {
  users.value = u
  environment.value = { name, nodeId, nodeOs, nodeArch, envConfig }
  step.value = 'template'
}

function templateBack() {
  environment.value = {}
  users.value = []
  step.value = 'environment'
}

function templateSelected(selected) {
  template.value = selected
  step.value = 'settings'
}

function settingsBack() {
  template.value = {}
  step.value = 'template'
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
    <div v-if="$api.auth.hasScope('nodes.view') && $api.auth.hasScope('templates.view')">
      <div :class="['progress', 'on-step-' + step]">
        <div :class="['step', 'step-environment', step === 'environment' ? 'step-current' : '']" />
        <div :class="['step', 'step-template', step === 'template' ? 'step-current' : '']" />
        <div :class="['step', 'step-settings', step === 'settings' ? 'step-current' : '']" />
      </div>
      <Environment v-if="step === 'environment'" :nouser="!$api.auth.hasScope('users.info.search')" @confirm="envConfirmed" />
      <select-template
        v-if="step === 'template'"
        :env="environment.envConfig.type"
        :os="environment.nodeOs"
        :arch="environment.nodeArch"
        @back="templateBack()"
        @selected="templateSelected"
      />
      <settings
        v-if="step === 'settings'"
        :data="template.data"
        :groups="template.groups"
        @back="settingsBack()"
        @confirm="settingsConfirmed"
      />
    </div>
    <div v-else v-text="t('servers.CreateMissingPermissions')" />
  </div>
</template>
