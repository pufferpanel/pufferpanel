<script setup>
import { computed, inject } from 'vue'
import SettingInput from '@/components/ui/SettingInput.vue'

const props = defineProps({
  modelValue: { type: Object, required: true },
  disabled: { type: Boolean, default: () => false }
})

const emit = defineEmits(['update:modelValue'])

const conditions = inject('conditions')

// ensure groups are sorted correctly
if (Array.isArray(props.modelValue.groups)) {
  const isSorted = !!props.modelValue.groups.reduce((acc, curr) => {
    if (acc === false || acc.order > curr.order) return false
    return curr
  })
  if (!isSorted) {
    const v = { ...props.modelValue, groups: [ ...props.modelValue.groups ] }
    v.groups.sort((a, b) => a.order > b.order ? 1 : -1)
    emit('update:modelValue', v)
  }
}

const grouplessVars = computed(() => {
  if (Array.isArray(props.modelValue.groups)) {
    return Object.keys(props.modelValue.data).filter(varname => {
      return props.modelValue.groups.map(g => g.variables).flat().indexOf(varname) === -1
    })
  } else {
    return Object.keys(props.modelValue.data)
  }
})

function updateValue(name, event) {
  const v = { ...props.modelValue }
  v.data[name].value = event.value
  emit('update:modelValue', v)
}

function visibleGroups() {
  const data = {}
  Object.keys(props.modelValue.data).map(name => {
    data[name] = props.modelValue.data[name].value
  })
  return props.modelValue.groups.filter(group => {
    if (group.if) {
      return conditions(group.if, data)
    }
    return true
  })
}

function filtered(group) {
  return group.variables.filter(v => {
    return props.modelValue.data[v]
  })
}
</script>

<template>
  <div v-if="modelValue.groups && modelValue.groups.length > 0">
    <div v-for="group in visibleGroups()" :key="group.order">
      <div class="group-header">
        <div class="title">
          <h3 v-text="group.display" />
          <div class="hint" v-text="group.description" />
        </div>
      </div>
      <div v-for="name in filtered(group)" :key="name">
        <setting-input :model-value="modelValue.data[name]" :disabled="disabled" @update:modelValue="updateValue(name, $event)" />
      </div>
    </div>
    <div v-if="grouplessVars.length > 0">
      <div class="group-header">
        <h3 class="title" v-text="t('templates.NoGroup')" />
      </div>
      <div v-for="name in grouplessVars" :key="name">
        <setting-input :model-value="modelValue.data[name]" :disabled="disabled" @update:modelValue="updateValue(name, $event)" />
      </div>
    </div>
  </div>
  <div v-else>
    <div v-for="(_, name) in modelValue.data" :key="name">
      <setting-input :model-value="modelValue.data[name]" :disabled="disabled" @update:modelValue="updateValue(name, $event)" />
    </div>
  </div>
</template>