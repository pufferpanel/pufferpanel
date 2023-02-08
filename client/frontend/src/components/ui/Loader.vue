<script setup>
import { useI18n } from 'vue-i18n'
import Icon from './Icon.vue'

const props = defineProps({
  small: { type: Boolean, default: () => false }
})

const { t } = useI18n()
</script>

<template>
  <div :class="['loader', small ? 'small' : '' ]">
    <icon name="loading" spin /><span class="text" v-text="t('common.Loading')" />
  </div>
</template>
