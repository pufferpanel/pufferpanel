<script setup>
import { ref, onUpdated } from 'vue'
import { useI18n } from 'vue-i18n'
import { operators } from '@/utils/operators.js'
import Ace from '@/components/ui/Ace.vue'
import Dropdown from '@/components/ui/Dropdown.vue'
import ListInput from '@/components/ui/ListInput.vue'
import TextField from '@/components/ui/TextField.vue'
import Toggle from '@/components/ui/Toggle.vue'

const props = defineProps({
  modelValue: { type: Object, required: true }
})

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()

const typeSelect = ref({})
const model = ref({ ...props.modelValue })
const types = Object.keys(operators).map(key => {
  return { value: key, label: t(`operators.${key}.generic`) }
})

function update() {
  emit('update:modelValue', model.value)
}

function typeChanged(newType) {
  const m = { type: newType }
  operators[newType].map(field => {
    m[field.name] = field.default
  })
  model.value = m
  update()
}

function getLabel(field) {
  return field.label ? t(field.label) : t(`operators.${model.value.type}.${field.name}`)
}
</script>

<template>
  <div :class="['operator', typeSelect.isOpen ? 'typeselect-open' : '']">
    <dropdown ref="typeSelect" v-model="model.type" :options="types" @update:modelValue="typeChanged" />
    <div v-for="field in operators[model.type]" :key="field.name">
      <text-field v-if="field.type === 'text'" v-model="model[field.name]" :label="getLabel(field)" @update:modelValue="update" />
      <toggle v-if="field.type === 'boolean'" v-model="model[field.name]" :label="getLabel(field)" @update:modelValue="update" />
      <ace v-if="field.type === 'textarea'" :id="`var-${field.name}-editor`" v-model="model[field.name]" :file="field.modeFile ? model[field.modeFile] : undefined" @update:modelValue="update" />
      <list-input v-if="field.type === 'list'" v-model="model[field.name]" :label="getLabel(field)" allow-swap @update:modelValue="update" />
    </div>
  </div>
</template>
