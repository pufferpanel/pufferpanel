<script setup>
import { ref, inject, onMounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import Multiselect from '@vueform/multiselect'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import Dropdown from '@/components/ui/Dropdown.vue'
import EnvironmentConfig from '@/components/ui/EnvironmentConfig.vue'
import TextField from '@/components/ui/TextField.vue'

const { t } = useI18n()
const api = inject('api')
const emit = defineEmits(['confirm'])

const name = ref('')
const nameError = ref(null)

const nodes = ref([])
const node = ref(undefined)
const nodeFeatures = ref({})

const selectedEnv = ref({ type: 'unsupported' })
const availableEnvs = ref([])
const envError = ref(null)
const msEnv = ref(null)

const users = ref([])
const msUsers = ref(null)

defineProps({
  nouser: { type: Boolean, default: () => true }
})

onMounted(async () => {
  const self = await api.self.get()
  nextTick(() => {
    msUsers.value.select({ value: self.username, label: self.username })
  })

  nodes.value = (await api.node.list()).map(n => { return { value: n.id, label: n.name } })
  node.value = nodes.value[0] ? nodes.value[0].value : null

  nodeChanged()
})

async function nodeChanged() {
  try {
    envError.value = null
    nodeFeatures.value = await api.node.features(node.value)
    nodeFeatures.value.environments = nodeFeatures.value.environments.filter(env => {
      if (env === 'docker') {
        // only allow if the node is actually configured to talk to docker
        return nodeFeatures.value.features.indexOf('docker') >= 0
      } else if (env === 'tty' || env === 'standard') {
        return false
      } else {
        return true
      }
    })

    availableEnvs.value = nodeFeatures.value.environments.map(env => {
      return { value: { type: env }, label: t(`env.${env}.name`) }
    })

    if (availableEnvs.value.length > 0) {
      selectedEnv.value = availableEnvs.value[0].value
    } else {
      msEnv.value.select({ value: { type: 'unsupported' }, label: '' })
      envError.value = t('servers.NoSupportedEnvironmentOnSelectedNode')
    }
  } catch (e) {
    // error asking for node features
    // general error handling takes care of informing user
    // we just need to set a proper error state
    availableEnvs.value = []
    msEnv.value.select({ value: { type: 'unsupported' }, label: '' })
    envError.value = t('servers.CannotFetchNodeEnvironments')
  }
}

function updateEnvironment(event) {
  const environment = availableEnvs.value.find(env => env.value.type === event.type)
  // copy fields while keeping reference so multiselect does not get confused
  for (let field in event) {
     environment.value[field] = event[field]
  }
  msEnv.value.select(environment)
}

function validateEnvironment() {
  return selectedEnv.value.type !== 'unsupported'
}

function validateName() {
  if (name.value.trim().match(/^[\x20-\x7e]+$/)) {
    return true
  }
  return false
}

function nameBlur() {
  if (validateName()) {
    nameError.value = undefined
  } else {
    nameError.value = t('servers.NameInvalid')
  }
}

async function searchUsers(query) {
  const res = await api.user.search(query)
  return res.map(u => {
    return {
      value: u.username,
      label: u.username
    }
  })
}

function validateUsers() {
  return users.value.length > 0
}

function canSubmit() {
  return validateName() && validateEnvironment() && validateUsers()
}

function confirm() {
  if (!canSubmit()) return
  emit('confirm',
    name.value.trim(),
    node.value,
    nodeFeatures.value.os,
    nodeFeatures.value.arch,
    selectedEnv.value,
    users.value
  )
}
</script>

<template>
  <div class="environment">
    <text-field v-model="name" :label="t('servers.Name')" :error="nameError" autofocus @blur="nameBlur()" @change="nameError = undefined" />
    <div class="dropdown-wrapper">
      <div :class="['dropdown', !validateUsers() ? 'error' : '']">
        <multiselect
          id="userselect"
          ref="msUsers"
          v-model="users"
          mode="tags"
          placeholder="t('server.SearchUsers')"
          :close-on-select="false"
          :can-clear="false"
          :filter-results="false"
          :min-chars="1"
          :resolve-on-load="false"
          :delay="500"
          :searchable="true"
          :options="searchUsers"
          :disabled="nouser"
        />
        <label for="userselect" @click="msUsers.open()">{{ t('users.Users') }}</label>
      </div>
      <span v-if="!validateUsers()" class="error" v-text="t('servers.AtLeastOneUserRequired')" />
    </div>
    <dropdown v-model="node" :options="nodes" :label="t('nodes.Node')" @change="nodeChanged()" />
    <dropdown ref="msEnv" v-model="selectedEnv" :options="availableEnvs" :label="t('servers.Environment')" :error="envError" />
    <environment-config :model-value="selectedEnv" @update:modelValue="updateEnvironment" />
    <btn color="primary" :disabled="!canSubmit()" @click="confirm()"><icon name="check" />{{ t('common.Next') }}</btn>
  </div>
</template>
