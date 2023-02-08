<script setup>
import { ref, onMounted, nextTick } from 'vue'
import Icon from './Icon.vue'
import markdown from '@/utils/markdown.js'

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
  disabled: { type: Boolean, default: () => false },
  modelValue: { type: [String, Number], default: () => '' }
})

const emit = defineEmits(['change', 'blur', 'focus', 'update:modelValue'])

const input = ref(null)
const showPassword = ref(false)

onMounted(() => {
  if (props.autofocus) {
    nextTick(() => {
      input.value.focus()
    })
  }
})

function onInput(e) {
  emit('update:modelValue', e.target.value)
  emit('change', e)
}

function onBlur(e) {
  emit('blur', e)
}

function onFocus(e) {
  emit('focus', e)
}
</script>

<template>
  <div class="input-field-wrapper">
    <div :class="['input-field', 'input-' + type, error ? 'error' : '', disabled ? 'disabled' : '']">
      <icon v-if="icon" class="pre" :name="icon" />
      <input :id="id" ref="input" :value="modelValue" :type="showPassword ? 'text' : type" :placeholder="label" :name="name" :disabled="disabled" @input="onInput($event)" @blur="onBlur($event)" @focus="onFocus($event)" />
      <icon v-if="type === 'password'" class="post" :name="showPassword ? 'eye-off' : 'eye'" @click="showPassword = !showPassword" />
      <label v-if="label" :for="id"><span v-text="label" /></label>
    </div>
    <span v-if="error" class="error" v-text="error" />
    <!-- eslint-disable-next-line vue/no-v-html -->
    <span v-if="hint && !error" class="hint" v-html="markdown(hint)" />
  </div>
</template>
