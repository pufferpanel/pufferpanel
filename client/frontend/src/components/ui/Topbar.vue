<script setup>
import { inject } from 'vue'
import { RouterLink } from 'vue-router'
import md5 from 'js-md5'
import Icon from './Icon.vue'
import PanelSearch from './PanelSearch.vue'

const props = defineProps({
  user: { type: Object, default: () => undefined }
})

defineEmits(['toggleSidebar'])

const config = inject('config')
const name = config.branding.name

function getAvatarLink() {
  return 'https://www.gravatar.com/avatar/' + md5(props.user.email.trim().toLowerCase()) + '?d=mp'
}
</script>

<template>
  <header class="topbar">
    <icon class="sidebar-toggle" name="nav-menu" @click="$emit('toggleSidebar')" />
    <div :data-name="name" class="name">
      {{ name }}
    </div>
    <panel-search />
    <router-link v-if="props.user" v-hotkey="'g a'" :to="{ name: 'Self' }"><img class="avatar" :src="getAvatarLink()" /></router-link>
  </header>
</template>
