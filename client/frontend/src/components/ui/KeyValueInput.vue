<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from './Btn.vue'
import Icon from './Icon.vue'
import TextField from './TextField.vue'

const props = defineProps({
  label: { type: String, default: () => undefined },
  addLabel: { type: String, default: () => undefined },
  keyLabel: { type: String, default: () => '' },
  valueLabel: { type: String, default: () => '' },
  hint: { type: String, default: () => undefined },
  error: { type: String, default: () => undefined },
  modelValue: { type: Object, default: () => {} }
})

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()

const entries = ref([])
const errors = ref({})

onMounted(() => {
  for (let field in props.modelValue) {
    entries.value.push({ key: field, value: props.modelValue[field] })
  }
})

function emitUpdate() {
  const result = {}
  entries.value.map(entry => {
    if (!entry.key) return

    if (!result[entry.key] && result[entry.key] !== '') {
      result[entry.key] = entry.value
      errors.value[entry.key] = undefined
    } else {
      errors.value[entry.key] = t('errors.DuplicateKey')
    }
  })

  emit('update:modelValue', result)
}

function addEntry() {
  entries.value.push({ key: '', value: '' })
}

function onKeyInput(entry, event) {
  entry.key = event
  emitUpdate()
}

function onValueInput(entry, event) {
  entry.value = event
  emitUpdate()
}

function removeEntry(item) {
  entries.value = entries.value.filter(entry => entry !== item)
  errors.value = {}
  emitUpdate()
}
</script>

<template>
  <div class="key-value-input">
    <div v-if="label" class="label" v-text="label" />
    <div v-if="error" class="error" v-text="error" />
    <div v-else-if="hint" class="hint" v-text="hint" />
    <div v-for="(entry, index) in entries" :key="index" class="entry">
      <div class="fields">
        <text-field :model-value="entry.key" :error="errors[entry.key]" :label="keyLabel" @update:modelValue="onKeyInput(entry, $event)" />
        <text-field :model-value="entry.value" :label="valueLabel" @update:modelValue="onValueInput(entry, $event)" />
      </div>
      <btn variant="icon" @click="removeEntry(entry)"><icon name="remove" /></btn>
    </div>
    <btn variant="text" @click="addEntry()"><icon name="plus" />{{ addLabel || t('common.Add') }}</btn>
  </div>
</template>
