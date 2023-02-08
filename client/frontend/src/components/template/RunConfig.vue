<script setup>
import { ref, onUpdated } from 'vue'
import { useI18n } from 'vue-i18n'
import Dropdown from '@/components/ui/Dropdown.vue'
import KeyValueInput from '@/components/ui/KeyValueInput.vue'
import Suggestion from '@/components/ui/Suggestion.vue'
import TextField from '@/components/ui/TextField.vue'

const props = defineProps({
  modelValue: { type: String, required: true }
})

const emit = defineEmits(['update:modelValue', 'valid'])

const { t } = useI18n()

const template = ref(JSON.parse(props.modelValue))
const commandInvalid = ref(false)
const stopType = ref(stopTypeFromValue(template.value.run))
const stopTypes = [
  { value: 'command', label: t('templates.StopCommand') },
  { value: 'signal', label: t('templates.StopSignal') }
]
const signalSuggestions = [
  { value: '1', label: t('templates.signals.1') },
  { value: '2', label: t('templates.signals.2') },
  { value: '9', label: t('templates.signals.9') },
  { value: '15', label: t('templates.signals.15') }
]

function stopTypeFromValue(run) {
  if (run.stop || run.stop === '') return 'command'
  return 'signal'
}

function update() {
  const t = template.value
  if (stopType.value === 'signal') {
    t.run.stopCode = Number(t.run.stopCode)
  }
  emit('update:modelValue', JSON.stringify(t, undefined, 4))
}

function updateStopType() {
  const t = template.value
  if (stopType.value === 'command') {
    delete t.run.stopCode
    t.run.stop = ''
  } else {
    delete t.run.stop
    t.run.stopCode = '2'
  }
  template.value = t
  update()
}

function validate() {
  commandInvalid.value = false

  if (!template.value.run.command || template.value.run.command.trim() === '') {
    commandInvalid.value = true
  }

  emit('valid', !commandInvalid.value)
}

onUpdated(() => {
  try {
    const u = JSON.parse(props.modelValue)
    // reserializing to avoid issues due to formatting
    if (JSON.stringify(template.value) !== JSON.stringify(u)) {
      template.value = u
      stopType.value = stopTypeFromValue(u.run)
      validate()
    }
  } catch {
    // expected failure caused by json editor producing invalid json during modification
  }
})
</script>

<template>
  <div>
    <text-field v-model="template.run.command" icon="start" :label="t('templates.Command')" :hint="t('templates.description.Command')" :error="commandInvalid ? t('templates.errors.CommandInvalid') : undefined" @update:modelValue="update" @blur="validate()" />
    <text-field v-model="template.run.workingDirectory" icon="folder" :label="t('templates.WorkingDirectory')" @update:modelValue="update" />
    <dropdown v-model="stopType" :options="stopTypes" @update:modelValue="updateStopType" />
    <text-field v-if="stopType === 'command'" v-model="template.run.stop" icon="stop" :label="t('templates.StopCommand')" :hint="t('templates.description.StopCommand')" @update:modelValue="update" />
    <suggestion v-else v-model="template.run.stopCode" icon="stop" :label="t('templates.StopSignal')" :options="signalSuggestions" :hint="t('templates.description.StopSignal')" @update:modelValue="update" />
    <key-value-input v-model="template.run.environmentVars" :label="t('templates.EnvVars')" :add-label="t('templates.AddEnvVar')" @update:modelValue="update" />
  </div>
</template>
