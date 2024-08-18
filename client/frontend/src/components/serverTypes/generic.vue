<script setup>
import { ref, defineAsyncComponent, inject, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

import ServerHeader from '../server/Header.vue'

import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import Loader from '@/components/ui/Loader.vue'
import Tab from '@/components/ui/Tab.vue'
import Tabs from '@/components/ui/Tabs.vue'

const Console = defineAsyncComponent({
  loader: () => import('../server/Console.vue'),
  loadingComponent: Loader
})
const Stats = defineAsyncComponent({
  loader: () => import('../server/Stats.vue'),
  loadingComponent: Loader
})
const Files = defineAsyncComponent({
  loader: () => import('../server/Files.vue'),
  loadingComponent: Loader
})
const Settings = defineAsyncComponent({
  loader: () => import('../server/Settings.vue'),
  loadingComponent: Loader
})
const Users = defineAsyncComponent({
  loader: () => import('../server/Users.vue'),
  loadingComponent: Loader
})
const Tasks = defineAsyncComponent({
  loader: () => import('../server/Tasks.vue'),
  loadingComponent: Loader
})
const Sftp = defineAsyncComponent({
  loader: () => import('../server/Sftp.vue'),
  loadingComponent: Loader
})
const Admin = defineAsyncComponent({
  loader: () => import('../server/Admin.vue'),
  loadingComponent: Loader
})

const { t } = useI18n()
const events = inject('events')
const route = useRoute()
const router = useRouter()
const http = ref(false)
const httpWarnDismissed = ref(false)
let httpCount = 2

const props = defineProps({
  server: { type: Object, required: true }
})

let task = null
onMounted(() => {
  if (route.query.created && props.server.hasScope('server.install')) {
    events.emit(
      'confirm',
      {
        title: t('servers.InstallPrompt'),
        body: t('servers.InstallPromptBody'),
      },
      {
        text: t('servers.Install'),
        icon: 'install',
        action: () => {
          props.server.install()
        }
      },
      {
        color: 'neutral'
      }
    )
    router.push({query: {}, hash: route.hash})
  }

  task = props.server.startTask(() => {
    // prevent the info alert from flickering on and off too much when dealing with a wonky connection
    if (props.server.needsPolling() && httpCount < 3) httpCount += 1
    if (!props.server.needsPolling() && httpCount > 0) httpCount -= 1
    if (httpCount === 3) http.value = true
    if (httpCount === 0) http.value = false
  }, 5000)
})

onUnmounted(() => {
  if (task) props.server.stopTask(task)
})
</script>

<template>
  <div :class="http ? 'http-fallback' : ''">
    <server-header :key="nameUpdateHack" :server="server" />

    <tabs anchors>
      <tab
        v-if="server.hasScope('server.console') || server.hasScope('server.console.send')"
        id="console"
        :title="t('servers.Console')"
        icon="console"
        hotkey="t c"
      >
        <div v-if="http && !httpWarnDismissed" class="alert info">
          <span v-text="t('servers.SocketWarnConsole')" />
          <btn variant="icon" @click="httpWarnDismissed = true"><icon name="close"></icon></btn>
        </div>
        <Console :server="server" />
      </tab>
      <tab
        v-if="server.hasScope('server.stats')"
        id="stats"
        :title="t('servers.Statistics')"
        icon="stats"
        hotkey="t i"
      >
        <stats :server="server" />
      </tab>
      <tab
        v-if="server.hasScope('server.files.view')"
        id="files"
        :title="t('servers.Files')"
        icon="files"
        hotkey="t f"
      >
        <files :server="server" />
      </tab>
      <tab
        v-if="server.hasScope('server.data.view') || server.hasScope('server.flags.view')"
        id="settings"
        :title="t('servers.Settings')"
        icon="settings"
        hotkey="t s"
      >
        <settings :server="server" />
      </tab>
      <tab
        v-if="server.hasScope('server.users.view')"
        id="users"
        :title="t('users.Users')"
        icon="users"
        hotkey="t u"
      >
        <users :server="server" />
      </tab>
      <!-- currently disabled due to tasks being broken -->
      <tab
        v-if="false && server.hasScope('server.tasks.view')"
        id="tasks"
        :title="t('servers.Tasks')"
        icon="tasks"
        hotkey="t t"
      >
        <tasks :server="server" />
      </tab>
      <tab
        v-if="server.hasScope('server.sftp')"
        id="sftp"
        :title="t('servers.SFTPInfo')"
        icon="sftp"
        hotkey="t 6"
      >
        <sftp :server="server" />
      </tab>
      <tab
        v-if="server.hasScope('server.definition.view') || server.hasScope('server.delete')"
        id="admin"
        :title="t('servers.Admin')"
        icon="admin"
        hotkey="t a"
      >
        <admin :server="server" />
      </tab>
    </tabs>
  </div>
</template>
