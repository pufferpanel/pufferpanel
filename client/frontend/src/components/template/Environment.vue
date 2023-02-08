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
const envs = ['standard', 'docker']
if (hasTTY()) envs.push('tty')
const envDefaults = {
  standard: { type: 'standard' },
  tty: { type: 'tty' },
  docker: { type: 'docker', image: 'pufferpanel/generic' }
}

function update() {
  emit('update:modelValue', JSON.stringify(template.value, undefined, 4))
}

function isDefault(env) {
  return env === template.value.environment.type
}

function setDefault(env) {
  template.value.environment = template.value.supportedEnvironments.find(e => e.type === env)
  update()
}

function isEnabled(env) {
  return !!template.value.supportedEnvironments.find(e => e.type === env)
}

function toggleEnv(env) {
  if (template.value.supportedEnvironments.find(e => e.type === env)) {
    template.value.supportedEnvironments = template.value.supportedEnvironments.filter(e => e.type !== env)
  } else {
    template.value.supportedEnvironments.push(envDefaults[env])
  }
  update()
}

function updateEnv(env, updated) {
  const e = template.value.supportedEnvironments.find(e => e.type === env)
  Object.keys(e).map(f => delete e[f])
  Object.keys(updated).map(f => e[f] = updated[f])
  if (isDefault(env)) {
    template.value.environment = e
  }
  update()
}

function hasTTY() {
  return !!template.value.supportedEnvironments.find(e => e.type === 'tty')
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
    <div v-for="env in envs" :key="env" :class="['environment', env, isDefault(env) ? 'default' : '']" :data-warn="env === 'tty' ? t('templates.TTYDeprecated') : undefined">
      <h2 v-text="t(`env.${env}.name`)" />
      <toggle :model-value="isEnabled(env)" :label="t('templates.EnvEnabled')" :disabled="isDefault(env)" @update:modelValue="toggleEnv(env)" />
      <btn v-if="isEnabled(env)" variant="icon" :disabled="isDefault(env)" @click="setDefault(env)">
        <icon name="default" />
      </btn>
      <environment-config v-if="isEnabled(env)" :model-value="template.supportedEnvironments.find(e => e.type === env)" @update:modelValue="updateEnv(env, $event)" />
    </div>
  </div>
</template>
