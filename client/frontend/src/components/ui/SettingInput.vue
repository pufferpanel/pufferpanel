<script setup>
import Dropdown from './Dropdown.vue'
import Toggle from './Toggle.vue'
import Suggestion from './Suggestion.vue'
import TextField from './TextField.vue'

const props = defineProps({
  modelValue: { type: Object, required: true }
})

const emit = defineEmits(['update:modelValue'])

function onInput(event) {
  emit('update:modelValue', { ...props.modelValue, value: event })
}
</script>

<template>
  <div class="setting-input-wrapper">
    <toggle v-if="modelValue.type === 'boolean'" :model-value="modelValue.value" class="setting-input" :label="modelValue.display" :hint="modelValue.desc" @update:modelValue="onInput($event)" />
    <dropdown v-else-if="modelValue.type === 'option'" :model-value="modelValue.value" label-prop="display" class="setting-input" :options="modelValue.options" :label="modelValue.display" :hint="modelValue.desc" @update:modelValue="onInput($event)" />
    <suggestion v-else-if="modelValue.options" :model-value="modelValue.value" label-prop="display" class="setting-input" :options="modelValue.options" :label="modelValue.display" :hint="modelValue.desc" @update:modelValue="onInput($event)" />
    <text-field v-else :model-value="modelValue.value" class="setting-input" :label="modelValue.display" :required="modelValue.required" :type="modelValue.type === 'integer' ? 'number' : 'text'" :hint="modelValue.desc" @update:modelValue="onInput($event)" />
  </div>
</template>
