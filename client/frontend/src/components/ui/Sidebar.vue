<script setup>
import { ref, inject, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { List, ListItem, ListItemContent, ListItemIcon } from '@/components/ui/list'

const props = defineProps({
  right: { type: Boolean, default: () => false },
  closed: { type: Boolean, default: () => false }
})

const api = inject('api')
const events = inject('events')
const { t } = useI18n()
const router = useRouter()
const routes = router.getRoutes().filter(e => {
  if (e.meta.permission === true) return true
  return e.meta.permission && api.auth.isLoggedIn() && api.auth.hasScope(e.meta.permission)
})

const mini = ref(localStorage.getItem('sidebar.mini') === 'true')

async function logout() {
  await api.auth.logout()
  router.push({ name: 'Login' })
  events.emit('logout')
}

function toggleMini() {
  mini.value = !mini.value
  localStorage.setItem('sidebar.mini', mini.value ? 'true' : 'false')
}
</script>

<template>
  <nav :class="['sidebar', props.right ? 'right' : 'left', mini ? 'mini' : '', closed ? 'closed' : 'open']">
    <div tabindex="-1" class="sidebar-content-top">
      <list>
        <list-item
          v-for="route in routes"
          :key="route.name"
          v-hotkey="route.meta.hotkey"
          :to="route"
        >
          <list-item-icon v-if="route.meta.icon" :icon="route.meta.icon" />
          <list-item-content v-text="t(route.meta.tkey ? route.meta.tkey : 'common.navigation.' + route.name)" />
        </list-item>
      </list>
    </div>
    <div tabindex="-1" class="sidebar-content-bottom">
      <list>
        <list-item tabindex="0" class="collapse-toggle" @click="toggleMini()">
          <list-item-icon :icon="mini ? 'chevron-right' : 'chevron-left'" />
          <list-item-content v-text="t('common.' + (mini ? 'Expand' : 'Collapse'))" />
        </list-item>
        <list-item tabindex="0" @click="logout()">
          <list-item-icon icon="logout" />
          <list-item-content v-text="t('users.Logout')" />
        </list-item>
      </list>
    </div>
  </nav>
</template>
