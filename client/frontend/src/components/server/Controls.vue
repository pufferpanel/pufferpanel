<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'

const { t } = useI18n()

const menuOpen = ref(false)

function hideMenu() {
  if (!menuOpen.value) return
  menuOpen.value = false
}

function toggleMenu() {
  const state = menuOpen.value
  setTimeout(() => menuOpen.value = !state, 100)
}

const props = defineProps({
  server: { type: Object, required: true }
})

const hotkeys = {
  'r r': () => props.server.start(),
  'r s': () => props.server.stop(),
  'r k': () => props.server.kill(),
  'r i': () => props.server.install()
}

function onHotkey(keys) {
  if (hotkeys[keys]) hotkeys[keys]()
}

const showMenu =
  props.server.hasScope('server.start') || 
  props.server.hasScope('server.stop') || 
  props.server.hasScope('server.kill') || 
  props.server.hasScope('server.install')
</script>

<template>
  <span v-hotkey="Object.keys(hotkeys)" class="server-controls" @hotkey="onHotkey">
    <btn v-if="server.hasScope('server.start')" class="start" @click="server.start()">
      <icon name="play" />
      <span class="text">{{ t('servers.Start') }}</span>
    </btn>
    <btn v-if="server.hasScope('server.stop')" class="stop" @click="server.stop()">
      <icon name="stop" />
      <span class="text">{{ t('servers.Stop') }}</span>
    </btn>
    <btn v-if="server.hasScope('server.kill')" class="kill" @click="server.kill()">
      <icon name="kill" />
      <span class="text">{{ t('servers.Kill') }}</span>
    </btn>
    <btn v-if="server.hasScope('server.install')" class="install" @click="server.install()">
      <icon name="install" />
      <span class="text">{{ t('servers.Install') }}</span>
    </btn>
    <btn class="menu" @click="toggleMenu()">
      <icon name="menu" />
    </btn>
    <div v-if="showMenu" v-click-outside="hideMenu" :class="['menu', menuOpen ? 'open' : 'closed']">
      <btn v-if="server.hasScope('server.start')" class="start" @click="menuOpen = false; server.start()">
        <icon name="play" />
        <span class="text">{{ t('servers.Start') }}</span>
      </btn>
      <btn v-if="server.hasScope('server.stop')" class="stop" @click="menuOpen = false; server.stop()">
        <icon name="stop" />
        <span class="text">{{ t('servers.Stop') }}</span>
      </btn>
      <btn v-if="server.hasScope('server.kill')" class="kill" @click="menuOpen = false; server.kill()">
        <icon name="kill" />
        <span class="text">{{ t('servers.Kill') }}</span>
      </btn>
      <btn v-if="server.hasScope('server.install')" class="install" @click="menuOpen = false; server.install()">
        <icon name="install" />
        <span class="text">{{ t('servers.Install') }}</span>
      </btn>
    </div>
  </span>
</template>
