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
  unbindEvent = props.server.on('status', e => {
    if (e.installing) {
      status.value = 'installing'
    } else if (e.running) {
      status.value = 'online'
    } else {
      status.value = 'offline'
    }
  })

  task = props.server.startTask(async () => {
    if (props.server.needsPolling() && props.server.hasScope('server.status')) {
      status.value = await props.server.getStatus()
    }
  }, 5000)

  if (props.server.hasScope('server.status'))
    status.value = await props.server.getStatus()
})

onUnmounted(() => {
  if (unbindEvent) unbindEvent()
  if (task) props.server.stopTask(task)
})
</script>

<template>
  <span
    v-if="server.hasScope('server.status')"
    :class="['status', status]"
  >
    <span
      class="tooltip"
      :data-tooltip="t(status === 'online' ? 'common.Online' : status === 'offline' ? 'common.Offline' : status === 'installing' ? 'common.Installing' : 'common.Unknown')"
    />
  </span>
</template>
