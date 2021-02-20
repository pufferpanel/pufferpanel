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
      <ui-switch
        v-model="autostart"
        :loading="loading"
        :disabled="loading"
        :label="$t('servers.Autostart')"
        @click="toggleSwitch('autostart')"
      />
      <ui-switch
        v-model="autorestart"
        :loading="loading"
        :disabled="loading"
        :label="$t('servers.Autorestart')"
        @click="toggleSwitch('autorestart')"
      />
      <ui-switch
        v-model="autorecover"
        :loading="loading"
        :disabled="loading"
        :label="$t('servers.Autorecover')"
        class="mb-4"
        @click="toggleSwitch('autorecover')"
      />
      <v-btn
        block
        color="primary"
        class="mb-4"
        @click="editServerDefinition()"
        v-text="$t('servers.EditDefinition')"
      />
      <v-btn
        block
        color="primary"
        class="mb-4"
        @click="reloadServer()"
        v-text="$t('servers.Reload')"
      />
      <v-btn
        block
        color="error"
        @click="confirmDeleteOpen = true"
        v-text="$t('servers.Delete')"
      />
      <ui-overlay
        v-model="editDefinition"
        card
        closable
      >
        <template v-slot:title>
          <span v-text="$t('servers.EditDefinition')" />
          <div style="flex-grow:50;" />
          <v-btn-toggle
            v-model="editMode"
            borderless
            dense
            mandatory
          >
            <v-btn
              value="editor"
              v-text="$t('templates.Editor')"
            />
            <v-btn
              value="json"
              v-text="$t('templates.Json')"
            />
          </v-btn-toggle>
        </template>
        <template-editor
          v-if="editMode === 'editor'"
          v-model="definition"
          server
        />
        <ace
          v-else
          ref="editor"
          v-model="defJson"
          :editor-id="server.id"
          height="75vh"
          lang="json"
        />
        <v-btn
          block
          color="success"
          class="mt-4"
          @click="saveServerDefinition()"
          v-text="$t('common.Save')"
        />
      </ui-overlay>
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
    </v-card-text>
  </v-card>
</template>

<script>
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
      autorecover: false,
      editDefinition: false,
      editMode: 'editor',
      definition: {},
      defJson: ''
    }
  },
  watch: {
    editMode (newVal) {
      if (newVal === 'editor') {
        this.definition = this.$api.templateFromApiJson(this.defJson, true)
      } else {
        this.defJson = this.$api.templateToApiJson(this.definition)
        if (this.$refs.editor && this.$refs.editor.ready) this.$refs.editor.setValue(this.defJson)
      }
    }
  },
  mounted () {
    this.loadData()
  },
  methods: {
    async loadData () {
      this.loading = true
      const def = await this.$api.getServerDefinition(this.server.id)
      this.autostart = !!def.run.autostart
      this.autorestart = !!def.run.autorestart
      this.autorecover = !!def.run.autorecover
      this.loading = false
    },
    async toggleSwitch (field) {
      this.loading = true
      const body = { run: {} }
      body.run[field] = this[field]
      await this.$api.updateServerDefinition(this.server.id, body)
      this.loadData()
    },
    async reloadServer () {
      await this.$api.reloadServer(this.server.id)
      this.$toast.success(this.$t('servers.Reloaded'))
    },
    async deleteConfirmed () {
      await this.$api.deleteServer(this.server.id)
      this.$toast.success(this.$t('servers.Deleted'))
      this.$router.push({ name: 'Servers' })
    },
    async editServerDefinition () {
      this.definition = this.$api.templateFromApiJson(await this.$api.getServerDefinition(this.server.id), true)
      this.editDefinition = true
    },
    async saveServerDefinition () {
      await this.$api.updateServerDefinition(this.server.id, this.$api.templateToApiJson(this.definition))
      this.editDefinition = false
    }
  }
}
</script>
