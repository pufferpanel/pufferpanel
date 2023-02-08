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
if (!template.value.run) template.value.run = {}
if (!template.value.run.pre) template.value.run.pre = []
if (!template.value.run.post) template.value.run.post = []

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
  <div class="hooks">
    <div class="pre-run">
      <h2 v-text="t('templates.PreRunHook')" />
      <div class="hint" v-text="t('templates.description.PreRunHook')" />
      <operator-list v-model="template.run.pre" :add-label="t('templates.AddPreStep')" @update:modelValue="update()" />
    </div>
    <div class="post-run">
      <h2 v-text="t('templates.PostRunHook')" />
      <div class="hint" v-text="t('templates.description.PostRunHook')" />
      <operator-list v-model="template.run.post" :add-label="t('templates.AddPostStep')" @update:modelValue="update()" />
    </div>
  </div>
</template>
