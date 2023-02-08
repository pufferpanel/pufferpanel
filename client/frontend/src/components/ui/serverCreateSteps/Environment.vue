<script setup>
import { ref, inject, onMounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import Dropdown from '@/components/ui/Dropdown.vue'
import EnvironmentConfig from '@/components/ui/EnvironmentConfig.vue'
import TextField from '@/components/ui/TextField.vue'

const props = defineProps({
  env: { type: Object, default: undefined },
  supported: { type: Array, default: undefined }
})

const { t } = useI18n()
const api = inject('api')
const emit = defineEmits(['back', 'confirm'])

const name = ref('')
const nameError = ref(null)

const nodes = ref([])
const node = ref(undefined)
const nodeFeatures = ref([])

const selectedEnv = ref({ type: 'unsupported' })
const availableEnvs = ref([])
const envError = ref(null)

const msEnv = ref(null)

onMounted(async () => {
  nodes.value = (await api.node.list()).map(n => { return { value: n.id, label: n.name } })
  node.value = nodes.value[0] ? nodes.value[0].value : null

  // TODO: get node features
  nodeFeatures.value = ['standard', 'tty', 'docker']

  // if no default is given, default to first supported
  // if there are no supported envs either, null to display error
  if (!props.env || !props.env.type) {
    if (props.supported && props.supported.length > 0) {
      selectedEnv.value = props.supported[0]
    } else {
      selectedEnv.value = null
    }
  } else {
    selectedEnv.value = props.env
  }

  // if no supported envs are given, set default as only supported
  if (props.supported && props.supported.length > 0) {
    availableEnvs.value = props.supported
  } else if (selectedEnv.value !== null) {
    availableEnvs.value = [selectedEnv.value]
  }

  // if the default env is not part of the supported envs, add it
  if (!availableEnvs.value.find(e => e.type === selectedEnv.value.type)) {
    availableEnvs.value.push(selectedEnv.value)
  }

  // make available envs useful for dropdown
  availableEnvs.value = availableEnvs.value.map(env => {
    return { value: env, label: t(`env.${env.type}.name`) }
  })

  nextTick(() => checkEnvironmentAvailability())
})

function checkEnvironmentAvailability() {
  envError.value = null
  availableEnvs.value = availableEnvs.value.map(env => {
    return { ...env, disabled: nodeFeatures.value.indexOf(env.value.type) === -1 }
  })

  if (nodeFeatures.value.indexOf(selectedEnv.value.type) === -1) {
    // selected env is not a valid choice on this node
    const candidate = availableEnvs.value.find(env => !env.disabled)
    if (candidate) {
      msEnv.value.select(candidate)
    } else {
      msEnv.value.select({ value: { type: 'unsupported' }, label: '' })
      envError.value = t('servers.NoSupportedEnvironmentOnSelectedNode')
    }
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

function canSubmit() {
  return validateName() && validateEnvironment()
}

function confirm() {
  if (!canSubmit()) return
  emit('confirm', name.value.trim(), node.value, selectedEnv.value)
}
</script>

<template>
  <div>
    <div v-if="selectedEnv === null" class="no-environment">
      <div v-text="t('errors.TemplateMissingEnvironment')" />
      <btn color="error" @click="emit('back')" v-text="t('common.Back')" />
    </div>
    <div v-else class="environment">
      <text-field v-model="name" :label="t('servers.Name')" :error="nameError" autofocus @blur="nameBlur()" @change="nameError = undefined" />
      <dropdown v-model="node" :options="nodes" :label="t('nodes.Node')" />
      <dropdown ref="msEnv" v-model="selectedEnv" :options="availableEnvs" :label="t('servers.Environment')" :error="envError" />
      <environment-config :model-value="selectedEnv" @update:modelValue="updateEnvironment" />
      <btn color="error" @click="emit('back')"><icon name="back" />{{ t('common.Back') }}</btn>
      <btn color="primary" :disabled="!canSubmit()" @click="confirm()"><icon name="check" />{{ t('common.Next') }}</btn>
    </div>
  </div>
</template>
