<script setup>
import { ref, onUpdated } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import Toggle from '@/components/ui/Toggle.vue'
import EnvironmentConfig from '@/components/ui/EnvironmentConfig.vue'

const props = defineProps({
  modelValue: { type: String, required: true }
})

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()
const template = ref(JSON.parse(props.modelValue))

function updateEnv(env, updated) {
  const e = template.value.supportedEnvironments.find(e => e.type === env)
  Object.keys(e).map(f => delete e[f])
  Object.keys(updated).map(f => e[f] = updated[f])
  if (isDefault(env)) {
    template.value.environment = e
  }
  emit('update:modelValue', JSON.stringify(template.value, undefined, 4))
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
  <div class="environment">
    <environment-config :model-value="template.environment" :no-fields-message="t('env.NoEnvFields')" @update:modelValue="updateEnv(env, $event)" />
  </div>
</template>
