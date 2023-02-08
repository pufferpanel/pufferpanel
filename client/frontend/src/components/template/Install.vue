<script setup>
import { ref, onUpdated } from 'vue'
import { useI18n } from 'vue-i18n'
import OperatorList from './OperatorList.vue'

const props = defineProps({
  modelValue: { type: String, required: true }
})

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()

const template = ref(JSON.parse(props.modelValue))
if (!template.value.install) template.value.install = []

function update() {
  emit('update:modelValue', JSON.stringify(template.value, undefined, 4))
}

onUpdated(() => {
  try {
    const u = JSON.parse(props.modelValue)
    // reserializing to avoid issues due to formatting
    if (JSON.stringify(template.value) !== JSON.stringify(u))
      template.value = u
  } catch {
    // expected failure caused by json editor producing invalid json during modification
  }
})
</script>

<template>
  <div class="install">
    <div class="hint" v-text="t('templates.description.Install')" />
    <operator-list v-model="template.install" :add-label="t('templates.AddInstallStep')" @update:modelValue="update()" />
  </div>
</template>
