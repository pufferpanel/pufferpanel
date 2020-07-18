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
      confirmDeleteOpen: false
    }
  },
  methods: {
    deleteConfirmed () {
      this.loading = true
      const ctx = this
      this.$http.delete(`/api/servers/${this.server.id}`).then(response => {
        ctx.$toast.success(ctx.$t('servers.Deleted'))
        ctx.$router.push({ name: 'Servers' })
      }).catch(handleError(ctx))
    }
  }
}
</script>
