<script setup>
import Icon from '@/components/ui/Icon.vue'

const props = defineProps({
  icon: { type: String, default: () => undefined }
})
</script>

<template>
  <div class="list-item-icon">
    <icon v-if="icon" :name="icon" />
    <slot v-else />
  </div>
</template>
