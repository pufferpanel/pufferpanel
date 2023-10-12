<script setup>
import { ref, onUpdated } from 'vue'
import { useI18n } from 'vue-i18n'
import Dropdown from '@/components/ui/Dropdown.vue'
import EnvironmentConfig from '@/components/ui/EnvironmentConfig.vue'

const props = defineProps({
  modelValue: { type: String, required: true }
})

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()
const template = ref(JSON.parse(props.modelValue))
const envs = [
  { value: 'host', label: t('env.host.name') },
  { value: 'docker', label: t('env.docker.name') }
]
const envDefaults = {
  host: { type: 'host' },
  docker: { type: 'docker', image: 'pufferpanel/generic' }
}

function update() {
  emit('update:modelValue', JSON.stringify(template.value, undefined, 4))
}

function envChanged(newEnv) {
  template.value.environment = envDefaults[newEnv]
  update()
}

function envValueChanged(newEnv) {
  template.value.environment = newEnv
  update()
}

onUpdated(() => {
  try {
    const u = JSON.parse(props.modelValue)
    // reserializing to avoid issues due to formatting
    if (JSON.stringify(template.value) !== JSON.stringify(u))
      template.value = u
  } catch {
    // expected failure caused by json editor producing invalid json during modification
  }
})
</script>

<template>
  <div class="environments">
    <dropdown :options="envs" :model-value="template.environment.type" :label="t('templates.Environment')" @update:modelValue="envChanged" />
    <environment-config :model-value="template.environment" @update:modelValue="envValueChanged" />
  </div>
</template>
