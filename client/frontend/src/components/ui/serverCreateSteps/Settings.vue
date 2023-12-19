<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import EnvironmentConfig from '@/components/ui/EnvironmentConfig.vue'
import Icon from '@/components/ui/Icon.vue'
import Variables from '@/components/ui/Variables.vue'

const { t } = useI18n()
const emit = defineEmits(['back', 'confirm'])
const props = defineProps({
  data: { type: Object, default: () => { return {} } },
  groups: { type: Array, default: undefined },
  env: { type: Object, required: true }
})

const settings = ref({})
const environment = ref(null)

onMounted(async () => {
  // prevent prop mutation by cloning to local state
  // also ensure that whatever happens, this is at least some value
  settings.value = props.data ? { ...props.data } : {}
  environment.value = { ...props.env }

  Object.keys(settings.value).map(setting => {
    // ensure booleans are actually booleans
    if (settings.value[setting].type === 'boolean') {
      settings.value[setting].value =
        settings.value[setting].value !== 'false' &&
        settings.value[setting].value !== false
    }
  })
})

function updateSettings(event) {
  settings.value = event.data
}

function canSubmit() {
  for (let key in settings.value) {
    if (settings.value[key].required && settings.value[key].type !== 'boolean') {
      if (!settings.value[key].value) return false
    }
  }
  return true
}

function confirm() {
  emit('confirm', settings.value, environment.value)
}
</script>

<template>
  <div class="settings">
    <environment-config v-if="environment" v-model="environment" />
    <variables :model-value="{ data: settings, groups }" @update:modelValue="updateSettings" />
    <div v-if="Object.keys(settings).length === 0" v-text="t('servers.NoSettings')" />
    <btn color="error" @click="emit('back')"><icon name="back" />{{ t('common.Back') }}</btn>
    <btn color="primary" :disabled="!canSubmit()" @click="confirm()"><icon name="save" />{{ t('servers.Create') }}</btn>
  </div>
</template>
