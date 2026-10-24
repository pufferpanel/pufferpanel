<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import { oidcErrorMessage } from '@/utils/oidc'

const { t } = useI18n()
const route = useRoute()

const errorMessage = computed(() => oidcErrorMessage(t, route.params.reason))
</script>

<template>
  <div class="login oidc-error">
    <h1 v-text="t('users.OidcFailedTitle')" />
    <p v-text="errorMessage" />
    <a href="/auth/oidc/login"><btn color="primary" v-text="t('users.OidcTryAgain')" /></a>
    <btn variant="text" @click="$router.push({ name: 'Login' })" v-text="t('users.OidcBackToLogin')" />
  </div>
</template>
