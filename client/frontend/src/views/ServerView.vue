<script setup>
import { shallowRef, inject, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import Loader from '@/components/ui/Loader.vue'

const api = inject('api')
const route = useRoute()

const server = shallowRef(null)
const serverComponent = shallowRef(Loader)

onMounted(async () => {
  server.value = await api.server.get(route.params.id)
  try {
    serverComponent.value = (await import(`../components/serverTypes/${server.value.type}.vue`)).default
  } catch {
    serverComponent.value = (await import('../components/serverTypes/generic.vue')).default
  }
})

onUnmounted(() => {
  server.value.close()
})
</script>

<template>
  <div v-if="server != null" :class="['serverview', server ? server.type : '']">
    <component :is="serverComponent" :server="server" />
  </div>
</template>
