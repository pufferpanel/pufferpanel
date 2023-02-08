<script setup>
import { ref, inject } from 'vue'

const props = defineProps({
  title: { type: String, required: true },
  icon: { type: String, default: () => undefined },
  id: { type: String, default: () => undefined }
})

const activeKey = inject('activeKey')
const key = props.id || props.title.toLowercase().replace(/ /g, '-')
</script>

<template>
  <div :id="id" :class="['tab', key === activeKey ? 'active' : 'inactive']">
    <slot />
  </div>
</template>
