<script setup>
import generic from './generic.vue'

defineProps({
  server: { type: Object, required: true }
})
</script>

<template>
  <generic :server="server" />
</template>
