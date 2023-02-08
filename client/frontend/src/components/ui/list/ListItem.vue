<script setup>
import { RouterLink } from 'vue-router'

const props = defineProps({
  // eslint-disable-next-line vue/require-prop-types
  to: { default: () => undefined }
})
</script>

<template>
  <router-link v-if="to" :to="to">
    <li class="list-item">
      <slot />
    </li>
  </router-link>
  <li v-else :class="['list-item', $attrs.onClick ? 'clickable' : '']">
    <slot />
  </li>
</template>
