<script setup>
import { ref } from 'vue'
import TextField from './TextField.vue'

const props = defineProps({
  // default to random id as labels need target ids to exist exactly once
  id: { type: String, default: () => (Math.random() + 1).toString(36).substring(2) },
  label: { type: String, default: () => undefined },
  hint: { type: String, default: () => undefined },
  error: { type: String, default: () => undefined },
  name: { type: String, default: () => undefined },
  type: { type: String, default: () => 'text' },
  icon: { type: String, default: () => undefined },
  autofocus: { type: Boolean, default: () => false },
  labelProp: { type: String, default: () => 'label' },
  options: { type: Array, default: () => [] },
  modelValue: { type: [String, Number], default: () => '' }
})

const emit = defineEmits(['blur', 'focus', 'update:modelValue'])

const open = ref(false)

function onInput(e) {
  emit('update:modelValue', e)
}

function onBlur(e) {
  // closing too fast makes clicking on an option not fire correctly
  // as the field is losing focus first, so wait for click events for a moment
  setTimeout(() => open.value = false, 250)
  emit('blur', e)
}

function onFocus(e) {
  open.value = true
  emit('focus', e)
}
</script>

<template>
  <div :class="['suggestions', open ? 'open' : 'closed']">
    <text-field :id="id" :model-value="modelValue" :label="label" :hint="hint" :error="error" :name="name" :type="type" :icon="icon" :autofocus="autofocus" @blur="onBlur" @focus="onFocus" @update:modelValue="onInput" />
    <div class="suggestions-list">
      <div v-for="option in options" :key="typeof option === 'object' ? option.value : option" class="suggestion" @click="onInput(typeof option === 'object' ? option.value : option)" v-text="typeof option === 'object' ? option[labelProp] || option.value : option" />
    </div>
  </div>
</template>
