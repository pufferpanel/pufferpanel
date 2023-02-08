<script setup>
import generic from './generic.vue'

const props = defineProps({
  server: { type: Object, required: true }
})
</script>

<template>
  <generic :server="server" />
</template>
