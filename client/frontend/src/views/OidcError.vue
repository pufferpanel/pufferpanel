<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'

const { t } = useI18n()
const route = useRoute()

// the OIDC callback redirects here when the SSO login did not end in a session
const errors = {
  notlinked: 'users.OidcErrorNotLinked',
  denied: 'users.OidcErrorDenied',
  state: 'users.OidcErrorState',
  unavailable: 'users.OidcErrorUnavailable',
  nologin: 'users.OidcErrorNoLogin'
}

const errorMessage = computed(() => t(errors[route.params.reason] || 'users.OidcErrorGeneric'))

function retry() {
  window.location.href = '/auth/oidc/login'
}
</script>

<template>
  <div class="login oidc-error">
    <h1 v-text="t('users.OidcFailedTitle')" />
    <p v-text="errorMessage" />
    <btn color="primary" @click="retry()" v-text="t('users.OidcTryAgain')" />
    <btn variant="text" @click="$router.push({ name: 'Login' })" v-text="t('users.OidcBackToLogin')" />
  </div>
</template>
