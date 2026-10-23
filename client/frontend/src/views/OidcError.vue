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
</script>

<template>
  <div class="login oidc-error">
    <h1 v-text="t('users.OidcFailedTitle')" />
    <p v-text="errorMessage" />
    <a href="/auth/oidc/login"><btn color="primary" v-text="t('users.OidcTryAgain')" /></a>
    <btn variant="text" @click="$router.push({ name: 'Login' })" v-text="t('users.OidcBackToLogin')" />
  </div>
</template>
