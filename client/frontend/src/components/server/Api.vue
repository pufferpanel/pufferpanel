<script setup>
import { useI18n } from 'vue-i18n'
import OAuth from '@/components/ui/OAuth.vue'

const { t } = useI18n()

const props = defineProps({
  server: { type: Object, required: true }
})
</script>

<template>
  <div class="oauth">
    <h2 v-text="t('servers.API')" />
    <o-auth :server="server.id" />
  </div>
</template>
