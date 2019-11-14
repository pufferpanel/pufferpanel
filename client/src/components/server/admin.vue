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
      <span v-text="$t('common.AdminControls')" />
    </v-card-title>
    <v-card-text>
      <v-dialog v-model="confirmDeleteOpen" max-width="600">
        <v-card>
          <v-card-title v-text="$t('common.ConfirmDeleteServer')" />
          <v-card-actions>
            <v-spacer />
            <v-btn v-text="$t('common.Cancel')" @click="confirmDeleteOpen = false" color="error" />
            <v-btn v-text="$t('common.Confirm')" @click="deleteConfirmed()" color="success" />
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-btn
        block
        large
        color="error"
        @click="confirmDeleteOpen = true"
        v-text="$t('common.DeleteServer')"
      />
    </v-card-text>
  </v-card>
</template>

<script>
export default {
  props: {
    server: { type: Object, default: function () { return {} } }
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
      this.$http.delete(`/api/servers/${this.server.id}`).then(function (response) {
        ctx.$toast.success(ctx.$t('common.ServerDeleted'))
        ctx.$router.push({ name: 'Servers' })
      }).catch(function (error) {
        let msg = 'errors.ErrUnknownError'
        if (error && error.response && error.response.data.error) {
          if (error.response.data.error.code) {
            msg = 'errors.' + error.response.data.error.code
          } else {
            msg = error.response.data.error.msg
          }
        }

        ctx.$toast.error(ctx.$t(msg))
      })
    }
  }
}
</script>

<style scoped>
  a {
    color: inherit;
  }

  .input-small {
    width: 200px;
    display: inline-block;
  }
</style>
