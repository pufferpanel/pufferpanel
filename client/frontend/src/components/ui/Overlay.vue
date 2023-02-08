<script setup>
import Btn from './Btn.vue'
import Icon from './Icon.vue'

const props = defineProps({
  title: { type: String, required: false, default: () => undefined },
  closable: { type: Boolean, required: false },
  modelValue: { type: Boolean, required: true }
})

const emit = defineEmits(['update:modelValue', 'close'])
</script>

<template>
  <div v-if="modelValue" class="overlay-backdrop">
    <div class="overlay">
      <div v-if="title" class="overlay-header">
        <h1 class="title" v-text="title" />
        <btn v-if="closable" v-hotkey="'Escape'" variant="icon" @click="$emit('update:modelValue', false); $emit('close')"><icon name="close" /></btn>
      </div>
      <slot />
    </div>
  </div>
</template>
