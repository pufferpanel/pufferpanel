<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from './Btn.vue'
import Icon from './Icon.vue'
import TextField from './TextField.vue'

const props = defineProps({
  label: { type: String, default: () => undefined },
  addLabel: { type: String, default: () => undefined },
  hint: { type: String, default: () => undefined },
  error: { type: String, default: () => undefined },
  allowSwap: { type: Boolean, default: () => false },
  modelValue: { type: Array, default: () => [] }
})

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()

const values = ref([ ...props.modelValue ])

function emitUpdate() {
  emit('update:modelValue', values)
}

function addEntry() {
  values.value.push('')
}

function onInput(i, event) {
  values.value[i] = event
  emitUpdate()
}

function swap(i1, i2) {
  const x = values.value[i1]
  values.value[i1] = values.value[i2]
  values.value[i2] = x
  emitUpdate()
}

function removeEntry(i) {
  values.value.splice(i, 1)
  emitUpdate()
}
</script>

<template>
  <div class="list-input">
    <div v-if="label" class="label" v-text="label" />
    <div v-if="error" class="error" v-text="error" />
    <div v-else-if="hint" class="hint" v-text="hint" />
    <div v-for="(entry, index) in values" :key="index" :class="['entry', allowSwap ? 'swap' : '']">
      <text-field :model-value="entry" @update:modelValue="onInput(index, $event)" />
      <btn v-if="allowSwap" :disabled="index === 0" variant="icon" @click="swap(index, index-1)"><icon name="up" /></btn>
      <btn v-if="allowSwap" :disabled="index === values.length-1" variant="icon" @click="swap(index, index+1)"><icon name="down" /></btn>
      <btn variant="icon" @click="removeEntry(index)"><icon name="remove" /></btn>
    </div>
    <btn variant="text" @click="addEntry()"><icon name="plus" />{{ addLabel || t('common.Add') }}</btn>
  </div>
</template>
