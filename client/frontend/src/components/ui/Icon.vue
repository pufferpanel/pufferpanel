<script setup>
const props = defineProps({
  name: { type: String, required: true },
  spin: { type: Boolean, required: false }
})
</script>

<template>
  <span :class="['icon', 'icon-' + props.name, spin ? 'icon-spin' : '']" />
</template>
