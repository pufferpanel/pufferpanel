<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const status = ref(null)

const props = defineProps({
  server: { type: Object, required: true }
})

onMounted(async () => {
  props.server.on('status', e => status.value = e.running)
  props.server.status()
})
</script>

<template>
  <span
    :class="['status', status === true ? 'online' : status === false ? 'offline' : 'unknown']"
    :data-hint="t(status === true ? 'common.Online' : status === false ? 'common.Offline' : 'common.Unknown')"
  />
</template>
