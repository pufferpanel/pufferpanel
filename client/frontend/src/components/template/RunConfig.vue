<script setup>
import { ref, onUpdated } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Dropdown from '@/components/ui/Dropdown.vue'
import Icon from '@/components/ui/Icon.vue'
import KeyValueInput from '@/components/ui/KeyValueInput.vue'
import Overlay from '@/components/ui/Overlay.vue'
import Suggestion from '@/components/ui/Suggestion.vue'
import TextField from '@/components/ui/TextField.vue'

const props = defineProps({
  modelValue: { type: String, required: true }
})

const emit = defineEmits(['update:modelValue', 'valid'])

const { t } = useI18n()

const template = ref(JSON.parse(props.modelValue))
const commandInvalid = ref(false)
const editOpen = ref(false)
const editIndex = ref(0)
const edit = ref({})
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

function addCommand() {
  if (!template.value.run.command) {
    template.value.run.command = []
  } else if (!Array.isArray(template.value.run.command)) {
    const tmp = template.value.run.command
    template.value.run.command = [{ command: tmp, if: '' }]
  }
  editIndex.value = template.value.run.command.length
  edit.value = { command: '', if: '' }
  editOpen.value = true
}

function swap(i1, i2) {
  const x = template.value.run.command[i1]
  template.value.run.command[i1] = template.value.run.command[i2]
  template.value.run.command[i2] = x
  update()
}

function remove(i) {
  template.value.run.command.splice(i, 1)
  update()
}

function startEdit(index) {
  editIndex.value = index
  edit.value = template.value.run.command[index]
  editOpen.value = true
}

function cancelEdit() {
  editOpen.value = false
  editIndex.value = 0
  edit.value = {}
}

function confirmEdit() {
  editOpen.value = false
  template.value.run.command[editIndex.value] = edit.value
  editIndex.value = 0
  edit.value = {}
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
    <ul v-if="Array.isArray(template.run.command)" class="list">
      <li v-for="(cmd, index) in template.run.command" :key="index" class="list-item clickable">
        <div class="list-item-content" @click="startEdit(index)">
          <div class="text">
            <span class="title" v-text="cmd.command" />
            <span v-if="cmd.if" class="subline" v-text="cmd.if" />
          </div>
          <btn :disabled="index === 0" variant="icon" @click.stop="swap(index, index-1)"><icon name="up" /></btn>
          <btn :disabled="index === template.run.command.length - 1" variant="icon" @click.stop="swap(index, index+1)"><icon name="down" /></btn>
          <btn variant="icon" @click.stop="remove(index)"><icon name="remove" /></btn>
        </div>
      </li>
    </ul>
    <text-field v-else v-model="template.run.command" icon="start" :label="t('templates.Command')" :hint="t('templates.description.Command')" :error="commandInvalid ? t('templates.errors.CommandInvalid') : undefined" @update:modelValue="update" @blur="validate()" />
    <btn variant="text" @click="addCommand()"><icon name="plus" />{{ t('templates.AddCommand') }}</btn>
    <text-field v-model="template.run.workingDirectory" icon="folder" :label="t('templates.WorkingDirectory')" @update:modelValue="update" />
    <dropdown v-model="stopType" :options="stopTypes" @update:modelValue="updateStopType" />
    <text-field v-if="stopType === 'command'" v-model="template.run.stop" icon="stop" :label="t('templates.StopCommand')" :hint="t('templates.description.StopCommand')" @update:modelValue="update" />
    <suggestion v-else v-model="template.run.stopCode" icon="stop" :label="t('templates.StopSignal')" :options="signalSuggestions" :hint="t('templates.description.StopSignal')" @update:modelValue="update" />
    <key-value-input v-model="template.run.environmentVars" :label="t('templates.EnvVars')" :add-label="t('templates.AddEnvVar')" @update:modelValue="update" />

    <overlay v-model="editOpen">
      <text-field v-model="edit.if" :label="t('common.Condition')" :hint="t('templates.CommandConditionHint')" />
      <text-field v-model="edit.command" :label="t('templates.Command')" />
      <div class="actions">
        <btn color="error" @click="cancelEdit()"><icon name="close" />{{ t('common.Cancel') }}</btn>
        <btn color="primary" @click="confirmEdit()"><icon name="save" />{{ t('common.Save') }}</btn>
      </div>
    </overlay>
  </div>
</template>
