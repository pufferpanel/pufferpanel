<script setup>
const props = defineProps({
  disabled: { type: Boolean, default: () => false },
  variant: { type: String, default: () => 'raised' },
  color: { type: String, default: () => 'neutral' },
  tooltip: { type: String, default: () => undefined }
})
</script>

<template>
  <button
    :disabled="props.disabled"
    :class="['button', props.variant, props.color, props.disabled ? 'disabled' : '']"
    :data-tooltip="props.tooltip"
    type="button"
  >
    <slot />
  </button>
</template>
