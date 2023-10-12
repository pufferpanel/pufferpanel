<script setup>
import { ref, inject, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import SettingInput from '@/components/ui/SettingInput.vue'
import Toggle from '@/components/ui/Toggle.vue'

const { t, te, locale } = useI18n()
const toast = inject('toast')

const props = defineProps({
  server: { type: Object, required: true }
})

const variables = ref({})
const flags = ref({})
const anyItems = computed(() => {
  if (Object.keys(variables.value).length > 0) return true
  if (Object.keys(flags.value).length > 0) return true
  return false
})

onMounted(async () => {
  if (props.server.hasScope('server.data.view'))
    variables.value = (await props.server.getData()) || {}
  if (props.server.hasScope('server.flags.view'))
    flags.value = (await props.server.getFlags()) || {}
})

async function save() {
  if (props.server.hasScope('server.data.edit'))
    await props.server.updateData(variables.value)
  if (props.server.hasScope('server.flags.edit'))
    await props.server.setFlags(flags.value)
  toast.success(t('servers.SettingsSaved'))
}

function getFlagHint(name) {
  if (te(`servers.flags.hint.${name}`, locale))
    return t(`servers.flags.hint.${name}`)
}
</script>

<template>
  <div>
    <h2 v-text="t('servers.Settings')" />
    <div v-for="(_, name) in variables" :key="name">
      <setting-input v-model="variables[name]" :disabled="!server.hasScope('server.data.edit')" />
    </div>
    <div v-for="(_, name) in flags" :key="name">
      <toggle v-model="flags[name]" :disabled="!server.hasScope('server.flags.edit')" :label="t(`servers.flags.${name}`)" :hint="getFlagHint()" />
    </div>
    <span v-if="!anyItems" v-text="t('servers.NoSettings')" />
    <btn v-else color="primary" @click="save()"><icon name="save" />{{ t('servers.SaveSettings') }}</btn>
  </div>
</template>
