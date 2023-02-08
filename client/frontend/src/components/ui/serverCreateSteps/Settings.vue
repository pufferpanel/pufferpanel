<script setup>
import { ref, inject, onMounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import SettingInput from '@/components/ui/SettingInput.vue'

const { t } = useI18n()
const api = inject('api')
const emit = defineEmits(['back', 'confirm'])
const props = defineProps({
  data: { type: Object, default: () => { return {} } }
})

const settings = ref({})

onMounted(async () => {
  // prevent prop mutation by cloning to local state
  // also ensure that whatever happens, this is at least some value
  settings.value = props.data ? { ...props.data } : {}

  Object.keys(settings.value).map(setting => {
    // ensure booleans are actually booleans
    if (settings.value[setting].type === 'boolean') {
      settings.value[setting].value =
        settings.value[setting].value !== 'false' &&
        settings.value[setting].value !== false
    }
  })
})

function canSubmit() {
  for (let key in settings.value) {
    if (settings.value[key].required && settings.value[key].type !== 'boolean') {
      if (!settings.value[key].value) return false
    }
  }
  return true
}

function confirm() {
  emit('confirm', settings.value)
}
</script>

<template>
  <div class="settings">
    <div v-for="(setting, name) in settings" :key="name">
      <setting-input v-model="settings[name]" />
    </div>
    <div v-if="Object.keys(settings).length === 0" v-text="t('servers.NoSettings')" />
    <btn color="error" @click="emit('back')"><icon name="back" />{{ t('common.Back') }}</btn>
    <btn color="primary" :disabled="!canSubmit()" @click="confirm()"><icon name="save" />{{ t('servers.Create') }}</btn>
  </div>
</template>
