<!--
  - Copyright 2019 Padduck, LLC
  -  Licensed under the Apache License, Version 2.0 (the "License");
  -  you may not use this file except in compliance with the License.
  -  You may obtain a copy of the License at
  -          http://www.apache.org/licenses/LICENSE-2.0
  -  Unless required by applicable law or agreed to in writing, software
  -  distributed under the License is distributed on an "AS IS" BASIS,
  -  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  -  See the License for the specific language governing permissions and
  -  limitations under the License.
  -->

<template>
  <v-card>
    <v-card-title>
      <span v-text="$t('servers.AdminControls')" />
    </v-card-title>
    <v-card-text>
      <v-switch
        v-model="autostart"
        :loading="loading"
        :disabled="loading"
        hide-details
        :label="$t('servers.Autostart')"
        @click="toggleSwitch('autostart')"
      />
      <v-switch
        v-model="autorestart"
        :loading="loading"
        :disabled="loading"
        hide-details
        :label="$t('servers.Autorestart')"
        @click="toggleSwitch('autorestart')"
      />
      <v-switch
        v-model="autorecover"
        :loading="loading"
        :disabled="loading"
        hide-details
        :label="$t('servers.Autorecover')"
        class="mb-4"
        @click="toggleSwitch('autorecover')"
      />
      <v-btn
        block
        color="primary"
        class="mb-4"
        @click="reloadServer()"
        v-text="$t('servers.Reload')"
      />
      <v-dialog
        v-model="confirmDeleteOpen"
        max-width="600"
      >
        <v-card>
          <v-card-title v-text="$t('servers.ConfirmDelete')" />
          <v-card-actions>
            <v-spacer />
            <v-btn
              color="error"
              @click="confirmDeleteOpen = false"
              v-text="$t('common.Cancel')"
            />
            <v-btn
              color="success"
              @click="deleteConfirmed()"
              v-text="$t('common.Confirm')"
            />
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-btn
        block
        color="error"
        @click="confirmDeleteOpen = true"
        v-text="$t('servers.Delete')"
      />
    </v-card-text>
  </v-card>
</template>

<script>
import { handleError } from '@/utils/api'

export default {
  props: {
    server: { type: Object, default: () => {} }
  },
  data () {
    return {
      confirmDeleteOpen: false,
      loading: true,
      autostart: false,
      autorestart: false,
      autorecover: false
    }
  },
  mounted () {
    this.loadData()
  },
  methods: {
    loadData () {
      const ctx = this
      ctx.loading = true
      ctx.$http.get(`/proxy/daemon/server/${ctx.server.id}`).then(response => {
        ctx.autostart = !!response.data.run.autostart
        ctx.autorestart = !!response.data.run.autorestart
        ctx.autorecover = !!response.data.run.autorecover
        ctx.loading = false
      }).catch(handleError(ctx))
    },
    toggleSwitch (field) {
      const ctx = this
      ctx.loading = true
      const body = { run: {} }
      body.run[field] = this[field]
      ctx.$http.post(`/proxy/daemon/server/${ctx.server.id}`, body).then(response => {
        ctx.loadData()
      })
    },
    reloadServer () {
      const ctx = this
      ctx.$http.post(`/proxy/daemon/server/${ctx.server.id}/reload`).then(response => {
        ctx.$toast.success(ctx.$t('servers.Reloaded'))
      }).catch(handleError(ctx))
    },
    deleteConfirmed () {
      const ctx = this
      this.$http.delete(`/api/servers/${this.server.id}`).then(response => {
        ctx.$toast.success(ctx.$t('servers.Deleted'))
        ctx.$router.push({ name: 'Servers' })
      }).catch(handleError(ctx))
    }
  }
}
</script>
