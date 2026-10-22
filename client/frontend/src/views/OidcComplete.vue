<script setup>
import { inject, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Loader from '@/components/ui/Loader.vue'
import defaultRoute from '@/router/defaultRoute'

const api = inject('api')
const events = inject('events')
const router = useRouter()

// the OIDC callback sets the session cookies; load the scopes before anything renders, as a password login does
onMounted(async () => {
  await api.auth.reauth()
  events.emit('login')
  let target = null
  try {
    target = JSON.parse(sessionStorage.getItem('returnTo'))
    sessionStorage.removeItem('returnTo')
  } catch {
    target = null
  }
  router.replace(target || defaultRoute(api))
})
</script>

<template>
  <loader />
</template>
