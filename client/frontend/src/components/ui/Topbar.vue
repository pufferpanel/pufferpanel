<script setup>
import { ref, inject, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import md5 from 'js-md5'
import Icon from './Icon.vue'

const props = defineProps({
  user: { type: Object, default: () => undefined }
})

const emit = defineEmits(['toggleSidebar'])

const api = inject('api')
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
    <router-link v-if="props.user" v-hotkey="'g a'" :to="{ name: 'Self' }"><img class="avatar" :src="getAvatarLink()" /></router-link>
  </header>
</template>
