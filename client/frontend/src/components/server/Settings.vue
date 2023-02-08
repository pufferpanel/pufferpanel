<script setup>
import { ref, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import SettingInput from '@/components/ui/SettingInput.vue'

const { t } = useI18n()
const toast = inject('toast')

const props = defineProps({
  server: { type: Object, required: true }
})

const items = ref({})

onMounted(async () => {
  items.value = (await props.server.getData()) || {}
})

async function save() {
  await props.server.updateData(items.value)
  toast.success(t('servers.SettingsSaved'))
}
</script>

<template>
  <div>
    <h2 v-text="t('servers.Settings')" />
    <div v-for="(item, name) in items" :key="name">
      <setting-input v-model="items[name]" />
    </div>
    <span v-if="Object.keys(items).length === 0" v-text="t('servers.NoSettings')" />
    <btn v-else color="primary" @click="save()"><icon name="save" />{{ t('servers.SaveSettings') }}</btn>
  </div>
</template>
