<script setup>
import { useI18n } from 'vue-i18n'
import Dropdown from './Dropdown.vue'

const props = defineProps({
  modelValue: { type: Object, required: true }
})

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()

function onInput(event) {
  emit('update:modelValue', { ...props.modelValue, current: event })
}

function onNativeInput(event) {
  emit('update:modelValue', { ...props.modelValue, current: event.target.value })
}

function getSettingLabel(setting) {
  const fallback = setting.label || undefined
  if (setting.tkey) {
    return t(setting.tkey, fallback)
  } else if (setting.tlabels) {
    return setting.tlabels[locale.value] || setting.tlabels[fallbackLocale.value] || fallback
  } else {
    return fallback
  }
}

function withNormalizedLabels(options) {
  return options.map(option => {
    return { ...option, label: getSettingLabel(option) }
  })
}
</script>

<template>
  <div class="theme-setting-wrapper">
    <dropdown v-if="modelValue.type === 'class'" :model-value="modelValue.current" :options="withNormalizedLabels(modelValue.options)" :label="getSettingLabel(modelValue)" @update:modelValue="onInput($event)" />
    <label v-if="modelValue.type === 'color'" class="color-input">
      <span class="label"><span v-text="getSettingLabel(modelValue)" /></span>
      <input type="color" :value="modelValue.current" @input="onNativeInput" />
    </label>
  </div>
</template>
