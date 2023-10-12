<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const status = ref(null)

const props = defineProps({
  server: { type: Object, required: true }
})

let unbindEvent = null
let task = null
onMounted(async () => {
  unbindEvent = props.server.on('status', e => status.value = e.running)

  task = props.server.startTask(async () => {
    if (props.server.needsPolling()) {
      status.value = await props.server.getStatus()
    }
  }, 5000)

  status.value = await props.server.getStatus()
})

onUnmounted(() => {
  if (unbindEvent) unbindEvent()
  if (task) props.server.stopTask(task)
})
</script>

<template>
  <span
    :class="['status', status === true ? 'online' : status === false ? 'offline' : 'unknown']"
    :data-hint="t(status === true ? 'common.Online' : status === false ? 'common.Offline' : 'common.Unknown')"
  />
</template>
