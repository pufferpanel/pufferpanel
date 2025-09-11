<script setup>
import { ref } from 'vue'

defineProps({
  title: { type: String, default: () => '' }
})

const open = ref(false)
</script>

<template>
  <div :class="['collapse', open ? 'open' : 'closed']">
    <div class="title" @click="open = !open" v-text="title" />
    <div class="content">
      <slot />
    </div>
  </div>
</template>
