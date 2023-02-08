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
</script>

<template>
  <span v-hotkey="Object.keys(hotkeys)" class="server-controls" @hotkey="onHotkey">
    <btn class="start" @click="server.start()">
      <icon name="play" />
      <span class="text">{{ t('servers.Start') }}</span>
    </btn>
    <btn class="stop" @click="server.stop()">
      <icon name="stop" />
      <span class="text">{{ t('servers.Stop') }}</span>
    </btn>
    <btn class="kill" @click="server.kill()">
      <icon name="kill" />
      <span class="text">{{ t('servers.Kill') }}</span>
    </btn>
    <btn class="install" @click="server.install()">
      <icon name="install" />
      <span class="text">{{ t('servers.Install') }}</span>
    </btn>
    <btn class="menu" @click="toggleMenu()">
      <icon name="menu" />
    </btn>
    <div v-click-outside="hideMenu" :class="['menu', menuOpen ? 'open' : 'closed']">
      <btn class="start" @click="menuOpen = false; server.start()">
        <icon name="play" />
        <span class="text">{{ t('servers.Start') }}</span>
      </btn>
      <btn class="stop" @click="menuOpen = false; server.stop()">
        <icon name="stop" />
        <span class="text">{{ t('servers.Stop') }}</span>
      </btn>
      <btn class="kill" @click="menuOpen = false; server.kill()">
        <icon name="kill" />
        <span class="text">{{ t('servers.Kill') }}</span>
      </btn>
      <btn class="install" @click="menuOpen = false; server.install()">
        <icon name="install" />
        <span class="text">{{ t('servers.Install') }}</span>
      </btn>
    </div>
  </span>
</template>
