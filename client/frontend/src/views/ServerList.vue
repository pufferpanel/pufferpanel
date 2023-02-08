<script setup>
import { ref, inject, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/ui/Icon.vue'
import Btn from '@/components/ui/Btn.vue'
import Loader from '@/components/ui/Loader.vue'

const api = inject('api')
const { t } = useI18n()

const servers = ref([])
let lastPage = 0
let loadingPage = false
const allServersLoaded = ref(false)
const loaderRef = ref(null)
const firstEntry = ref(null)
let interval = null

function addServers(newServers) {
  newServers.map(server => servers.value.push(server))
  refreshServerStatus()
}

async function refreshServerStatus() {
  servers.value.map(async s => s.online = await api.server.getStatus(s.id))
}

function isLoaderVisible() {
  if (!loaderRef.value) return false
  const vw = window.innerWidth || document.documentElement.clientWidth
  const vh = window.innerHeight || document.documentElement.clientHeight
  const rect = loaderRef.value.$el.getBoundingClientRect()
  return rect.top >= 0 && rect.left >= 0 && rect.bottom <= vh && rect.right <= vw
}

async function loadPage(page = 1) {
  loadingPage = true
  const data = await api.server.list(page)
  addServers(data.servers)
  lastPage = data.paging.page
  allServersLoaded.value = data.paging.page * data.paging.pageSize >= (data.paging.total || 0)
  nextTick(() => {
    loadingPage = false
    if (!allServersLoaded.value && isLoaderVisible()) loadPage(lastPage + 1)
  })
}

function onScroll() {
  if (!loadingPage && isLoaderVisible()) loadPage(lastPage + 1)
}

onMounted(() => {
  interval = setInterval(refreshServerStatus, 30 * 1000)
  nextTick(() => {
    loadPage()
    window.addEventListener('scroll', onScroll)
  })
})

onUnmounted(() => {
  clearInterval(interval)
  window.removeEventListener('scroll', onScroll)
})

function getServerAddress(server) {
  let ip = server.node.publicHost
  if (server.ip && server.ip !== '0.0.0.0') {
    ip = server.ip
  }
  return ip + (server.port ? ':' + server.port : '')
}

function setFirstEntry(ref) {
  if (!firstEntry.value) firstEntry.value = ref
}

function focusList() {
  firstEntry.value.$el.focus()
}
</script>

<template>
  <div class="serverlist">
    <h1 v-text="t('servers.Servers')" />
    <div v-hotkey="'l'" class="list" @hotkey="focusList()">
      <div v-for="server in servers" :key="server.id" class="list-item">
        <router-link :ref="setFirstEntry" :to="{ name: 'ServerView', params: { id: server.id } }">
          <div
            :class="['server', 'server-' + server.type]"
            :data-online="server.online === true ? 'online' : server.online === false ? 'offline' : 'loading'"
          >
            <span class="title" :title="server.name">{{server.name}}</span>
            <span class="type">{{server.type}}</span>
            <span class="subline">{{getServerAddress(server)}} @ {{server.node.name}}</span>
          </div>
        </router-link>
      </div>
      <div v-if="!allServersLoaded" class="list-item">
        <loader ref="loaderRef" small />
      </div>
      <div v-if="$api.auth.hasScope('servers.create')" class="list-item">
        <router-link v-hotkey="'c'" :to="{ name: 'ServerCreate' }">
          <div class="createLink"><icon name="plus" />{{ t('servers.Add') }}</div>
        </router-link>
      </div>
    </div>
  </div>
</template>
