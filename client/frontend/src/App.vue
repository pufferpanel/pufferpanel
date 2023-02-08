<script setup>
import { ref, inject, watch, onMounted, onUnmounted, onUpdated } from 'vue'
import { RouterView, useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Hotkeys from '@/components/ui/Hotkeys.vue'
import Icon from '@/components/ui/Icon.vue'
import Overlay from '@/components/ui/Overlay.vue'
import Topbar from '@/components/ui/Topbar.vue'
import Sidebar from '@/components/ui/Sidebar.vue'

const api = inject('api')
const toast = inject('toast')
const events = inject('events')
const { t } = useI18n()
const themeClasses = inject('themeClasses')
const sidebarClosedBelow = inject('sidebarClosedBelow')
const route = useRoute()
const router = useRouter()
const allowSidebar = ref(false)
const user = ref(undefined)
const sidebarClosed = ref(window.innerWidth < sidebarClosedBelow.value)
const ltr = inject('ltr')
const error = ref('')
const showError = ref(false)
const showHotkeys = ref(false)
const confirmOpen = ref(false)
const confirm = ref({})
let lastWidth = window.innerWidth

function showErrorDetails(e) {
  const getCircularReplacer = () => {
    const seen = new WeakSet()
    return (key, value) => {
      if (key === 'password') return '[password]'
      if (typeof value === 'string') {
        try {
          const json = JSON.parse(value)
          if (typeof json === 'object' && json !== null) {
            if (Object.keys(json).indexOf('password') !== -1) {
              json.password = '[password]'
            }
            return JSON.stringify(json)
          } else { return value }
        } catch { return value }
      }
      if (typeof value === 'object' && value !== null) {
        if (seen.has(value)) {
          return
        }
        seen.add(value)
      }
      return value
    }
  }

  let statusMessage = `${e.status} ${e.statusText}`
  switch (e.status) {
    case 401:
      statusMessage = 'Not logged in (401)'
      break
    case 403:
      statusMessage = 'Permission denied (403)'
      break
    case 404:
      statusMessage = 'Endpoint not found, is something blocking access to the API? (404)'
      break
    case 500:
      statusMessage = 'Server error (500)'
      break
    case 502:
      statusMessage = 'Bad Gateway, is PufferPanel running? (502)'
      break
  }

  let auth = e.request.headers.Authorization
  if (auth) {
    auth = auth.substring(0, 16)
    if (auth.length === 16) auth = auth + '...'
  } else auth = 'Not present'

  let body = e.request.data
  if (body) {
    body = JSON.stringify(JSON.parse(body), getCircularReplacer(), 2)
  }

  error.value = `${statusMessage}

Endpoint: ${e.request.method} ${e.request.url}

Authorization Header: ${auth}

${body ? 'Request Body: ' + body : ''}`
    .replace(/>/g, '&gt;')
    .replace(/</g, '&lt;')
    .replace(/ /g, '&nbsp;')
    .replace(/\n/g, '<br />')

  showError.value = true
}

api._errorHandler = e => {
  if (e.status === 401) {
    // session expired, save route and return to login
    sessionStorage.setItem('returnTo', JSON.stringify({
      name: route.name,
      params: route.params,
      hash: route.hash,
      query: route.query
    }))
    // make sure expired session is cleaned up properly
    api.auth.logout()
    toast.error(t('errors.ErrSessionTimedOut'))
    router.push({ name: 'Login' })
  } else if (e.code === 'ErrGeneric' && e.msg) {
    toast.error(t(e.msg))
  } else if (e.code === 'ErrUnknownError') {
    toast.error(t('errors.ErrUnknownError'), () => showErrorDetails(e), t('common.Details'))
  } else {
    toast.error(t('errors.' + e.code))
  }
}

onMounted(async () => {
  window.addEventListener('resize', onResize)
  events.on('confirm', handleConfirm)

  document.documentElement.style.setProperty('--inner-height', `${window.innerHeight}px`)
  allowSidebar.value = !route.meta.noAuth
  setInterval(() => {
    if (api.auth.isLoggedIn()) api.auth.reauth()
  }, 1000 * 60 * 15)
  if (api.auth.isLoggedIn()) {
    user.value = await api.self.get()
    api.auth.reauth()
  }
})

onUpdated(async () => {
  if (api.auth.isLoggedIn() && user.value == undefined) {
    user.value = await api.self.get()
  }
  if (!api.auth.isLoggedIn() && user.value) {
    user.value = undefined
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', onResize)
})

watch(
  () => route.meta,
  async newMeta => {
    maybeCloseSidebar()
    allowSidebar.value = !newMeta.noAuth
  }
)

function onResize() {
  if (lastWidth < sidebarClosedBelow.value && window.innerWidth >= sidebarClosedBelow.value) {
    sidebarClosed.value = false
  } else if (lastWidth >= sidebarClosedBelow.value && window.innerWidth < sidebarClosedBelow.value) {
    sidebarClosed.value = true
  }

  lastWidth = window.innerWidth

  document.documentElement.style.setProperty('--inner-height', `${window.innerHeight}px`)
}

function maybeCloseSidebar() {
  if (window.innerWidth < sidebarClosedBelow.value) {
    sidebarClosed.value = true
  }
}

function handleConfirm(title, ok, cancel) {
  if (!title || !ok.action) {
    console.warn('ignoring confirm request, no title or confirm action given')
    return
  }

  if (!ok.text) ok.text = t('common.Confirm')
  if (!ok.icon) ok.icon = 'apply'
  if (!ok.color) ok.color = 'primary'

  if (!cancel) cancel = {}
  if (!cancel.action) cancel.action = () => {}
  if (!cancel.text) cancel.text = t('common.Cancel')
  if (!cancel.icon) cancel.icon = 'close'
  if (!cancel.color) cancel.color = 'error'

  ok.handle = () => {
    ok.action()
    confirmOpen.value = false
  }

  cancel.handle = () => {
    cancel.action()
    confirmOpen.value = false
  }

  confirm.value = {
    title,
    confirm: ok,
    cancel
  }

  confirmOpen.value = true
}
</script>

<template>
  <div id="root" v-hotkey="['?', 'Shift+?']" :class="themeClasses" @hotkey="showHotkeys = !showHotkeys">
    <topbar :class="allowSidebar ? 'sidebar-exists' : ''" :user="user" @toggleSidebar="sidebarClosed = !sidebarClosed" />
    <sidebar v-if="allowSidebar" :closed="sidebarClosed" :right="!ltr" />
    <main class="main" @click="maybeCloseSidebar()">
      <router-view />
    </main>
    <overlay v-model="showError" :title="t('common.ErrorDetails')" closable>
      <!-- eslint-disable-next-line vue/no-v-html -->
      <div dir="ltr" v-html="error" />
    </overlay>
    <overlay v-model="showHotkeys" class="hotkeys" :title="t('hotkeys.Title')" closable>
      <hotkeys />
    </overlay>
    <overlay v-model="confirmOpen" class="confirm-overlay" :title="confirm.title">
      <div class="actions">
        <btn :color="confirm.cancel.color" @click="confirm.cancel.handle">
          <icon :name="confirm.cancel.icon" />
          {{ confirm.cancel.text }}
        </btn>
        <btn :color="confirm.confirm.color" @click="confirm.confirm.handle">
          <icon :name="confirm.confirm.icon" />
          {{ confirm.confirm.text }}
        </btn>
      </div>
    </overlay>
    <div id="toasts" />
  </div>
</template>
