<script setup>
import { ref, inject, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import Toggle from '@/components/ui/Toggle.vue'
import Variables from '@/components/ui/Variables.vue'

const { t, te, locale } = useI18n()
const toast = inject('toast')

const props = defineProps({
  server: { type: Object, required: true }
})

const vars = ref({})
const flags = ref({})
const anyItems = computed(() => {
  if (Object.keys(vars.value).length > 0) return true
  if (Object.keys(flags.value).length > 0) return true
  return false
})

onMounted(async () => {
  if (props.server.hasScope('server.data.view'))
    vars.value = (await props.server.getData()) || {}
  if (props.server.hasScope('server.flags.view'))
    flags.value = (await props.server.getFlags()) || {}
})

async function save() {
  if (props.server.hasScope('server.data.edit')) {
    const data = {}
    Object.keys(vars.value.data).map(name => {
      data[name] = vars.value.data[name].value
    })
    await props.server.updateData(data)
  }
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
    <variables v-model="vars" :disabled="!server.hasScope('server.data.edit')" />
    <div class="group-header">
      <div class="title">
        <h3 v-text="t('servers.FlagsHeader')" />
      </div>
    </div>
    <div v-for="(_, name) in flags" :key="name">
      <toggle v-model="flags[name]" :disabled="!server.hasScope('server.flags.edit')" :label="t(`servers.flags.${name}`)" :hint="getFlagHint()" />
    </div>
    <span v-if="!anyItems" v-text="t('servers.NoSettings')" />
    <btn v-else color="primary" @click="save()"><icon name="save" />{{ t('servers.SaveSettings') }}</btn>
  </div>
</template>
