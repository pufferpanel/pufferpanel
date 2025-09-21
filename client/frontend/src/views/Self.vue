<script setup>
import { inject } from 'vue'
import { useI18n } from 'vue-i18n'
import OAuth from '@/components/ui/OAuth.vue'
import Tab from '@/components/ui/Tab.vue'
import Tabs from '@/components/ui/Tabs.vue'

import Password from './self/Password.vue'
import Preferences from './self/Preferences.vue'
import Security from './self/Security.vue'
import UserDetails from './self/Details.vue'

const { t } = useI18n()
const api = inject('api')
</script>

<template>
  <div class="self">
    <tabs anchors>
      <tab id="preferences" :title="t('users.Preferences')" icon="settings" hotkey="t s">
        <preferences />
      </tab>
      <tab v-if="api.auth.hasScope('self.edit')" id="account" :title="t('users.ChangeInfo')" icon="account" hotkey="t a">
        <user-details />
      </tab>
      <tab v-if="api.auth.hasScope('self.edit')" id="changepassword" :title="t('users.ChangePassword')" icon="lock" hotkey="t p">
        <password />
      </tab>
      <tab v-if="api.auth.hasScope('self.edit')" id="security" :title="t('users.2fa')" icon="2fa" hotkey="t 2">
        <security />
      </tab>
      <tab v-if="api.auth.hasScope('self.clients')" id="oauth" :title="t('oauth.Clients')" icon="api" hotkey="t o">
        <div class="oauth">
          <h1 v-text="t('oauth.Clients')" />
          <o-auth />
        </div>
      </tab>
    </tabs>
  </div>
</template>
