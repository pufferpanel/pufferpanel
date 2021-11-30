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
  <div>
    <div style="position:absolute;top:0px;bottom:56px;left:0px;right:0px;overflow-y:auto;">
      <v-container>
        <v-row ref="headerRow">
          <v-col>
            <h1
              style="float: left;"
            >
              <server-status :server="server" />
              {{ server.name }}
              <v-btn
                v-if="server.permissions.editServerData"
                icon
                @click="startRename()"
              >
                <v-icon>mdi-pencil</v-icon>
              </v-btn>
            </h1>
            <div style="float: right;">
              <server-controls :server="server" />
            </div>
            <div
              style="clear: both;"
              class="mb-2"
            />
          </v-col>
        </v-row>
        <v-row
          v-if="socketError"
          ref="alertRow"
        >
          <v-col>
            <v-alert
              border="left"
              text
              type="warning"
            >
              {{ $t('errors.ErrSocketFailed') }}
            </v-alert>
          </v-col>
        </v-row>
        <v-row
          v-show="currentTab === 'console'"
          v-if="server.permissions.viewServerConsole || isAdmin()"
        >
          <v-col>
            <server-console :server="server" />
          </v-col>
        </v-row>
        <v-row
          v-show="currentTab === 'stats'"
          v-if="server.permissions.viewServerStats || isAdmin()"
        >
          <v-col>
            <server-cpu :server="server" />
          </v-col>
        </v-row>
        <v-row
          v-show="currentTab === 'stats'"
          v-if="server.permissions.viewServerStats || isAdmin()"
        >
          <v-col>
            <server-memory :server="server" />
          </v-col>
        </v-row>
        <v-row
          v-show="currentTab === 'files'"
          v-if="server.permissions.viewServerFiles || isAdmin()"
        >
          <v-col>
            <server-files :server="server" />
          </v-col>
        </v-row>
        <v-row
          v-show="currentTab === 'files'"
          v-if="server.permissions.sftpServer || isAdmin()"
        >
          <v-col>
            <server-sftp :server="server" />
          </v-col>
        </v-row>
        <v-row
          v-show="currentTab === 'tasks'"
          v-if="server.permissions.editServerData || isAdmin()"
        >
          <v-col>
            <server-tasks :server="server" />
          </v-col>
        </v-row>
        <v-row
          v-show="currentTab === 'settings'"
          v-if="server.permissions.editServerData || isAdmin()"
        >
          <v-col>
            <server-settings :server="server" />
          </v-col>
        </v-row>
        <v-row
          v-show="currentTab === 'users'"
          v-if="server.permissions.editServerUsers || isAdmin()"
        >
          <v-col>
            <server-users :server="server" />
          </v-col>
        </v-row>
        <v-row v-show="currentTab === 'admin'">
          <v-col>
            <ui-oauth :server="server" />
          </v-col>
        </v-row>
        <v-row
          v-show="currentTab === 'admin'"
          v-if="server.permissions.deleteServer || isAdmin()"
        >
          <v-col>
            <server-admin :server="server" />
          </v-col>
        </v-row>
      </v-container>
    </div>

    <v-row
      ref="bottomRow"
      class="mt-8"
    >
      <v-col>
        <v-bottom-navigation
          v-model="currentTab"
          absolute
          mandatory
          color="primary"
          style="z-index:3;"
        >
          <v-slide-group>
            <v-slide-item
              v-if="server.permissions.viewServerConsole || isAdmin()"
              v-slot="{}"
            >
              <v-btn
                v-hotkey="'t c'"
                value="console"
              >
                <span>{{ $t('servers.Console') }}</span>
                <v-icon>mdi-console-line</v-icon>
              </v-btn>
            </v-slide-item>
            <v-slide-item
              v-if="server.permissions.viewServerStats || isAdmin()"
              v-slot="{}"
            >
              <v-btn
                v-hotkey="'t i'"
                value="stats"
              >
                <span>{{ $t('servers.Statistics') }}</span>
                <v-icon>mdi-chart-line</v-icon>
              </v-btn>
            </v-slide-item>
            <v-slide-item
              v-if="server.permissions.viewServerFiles || server.permissions.sftpServer || isAdmin()"
              v-slot="{}"
            >
              <v-btn
                v-hotkey="'t f'"
                value="files"
              >
                <span>{{ $t('servers.Files') }}</span>
                <v-icon>mdi-file</v-icon>
              </v-btn>
            </v-slide-item>
            <v-slide-item
              v-if="showTasks && (server.permissions.editServerData || isAdmin())"
              v-slot="{}"
            >
              <v-btn
                v-hotkey="'t t'"
                value="tasks"
              >
                <span>{{ $t('servers.Tasks') }}</span>
                <v-icon>mdi-timer</v-icon>
              </v-btn>
            </v-slide-item>
            <v-slide-item
              v-if="server.permissions.editServerData || isAdmin()"
              v-slot="{}"
            >
              <v-btn
                v-hotkey="'t s'"
                value="settings"
              >
                <span>{{ $t('servers.Settings') }}</span>
                <v-icon>mdi-cog</v-icon>
              </v-btn>
            </v-slide-item>
            <v-slide-item
              v-if="server.permissions.editServerUsers || isAdmin()"
              v-slot="{}"
            >
              <v-btn
                v-hotkey="'t u'"
                value="users"
              >
                <span>{{ $t('users.Users') }}</span>
                <v-icon>mdi-account-multiple</v-icon>
              </v-btn>
            </v-slide-item>
            <v-slide-item v-slot="{}">
              <v-btn
                v-hotkey="'t a'"
                value="admin"
              >
                <span>{{ $t('servers.Admin') }}</span>
                <v-icon>mdi-account-star</v-icon>
              </v-btn>
            </v-slide-item>
          </v-slide-group>
        </v-bottom-navigation>
      </v-col>
    </v-row>

    <ui-overlay
      v-model="renameOpen"
      :title="$t('servers.Rename')"
      card
      closable
      @close="resetRename()"
    >
      <v-row>
        <v-col cols="12">
          <ui-input
            v-model="newName"
            :label="$t('common.Name')"
            autofocus
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col
          cols="12"
          md="6"
        >
          <v-btn
            block
            color="error"
            @click="resetRename()"
            v-text="$t('common.Cancel')"
          />
        </v-col>
        <v-col
          cols="12"
          md="6"
        >
          <v-btn
            block
            color="success"
            @click="confirmRename()"
            v-text="$t('common.Save')"
          />
        </v-col>
      </v-row>
    </ui-overlay>
  </div>
</template>

<script>
export default {
  props: {
    server: { type: Object, default: () => {} }
  },
  data () {
    return {
      socketError: false,
      currentTab: 'console',
      showTasks: false,
      renameOpen: false,
      newName: ''
    }
  },
  mounted () {
    this.$api.startServerTask(this.server.id, () => {
      this.$api.requestServerStats(this.server.id)
    }, 3000)

    this.$api.addServerListener(this.server.id, 'error', () => {
      this.$toast.warning(this.$t('errors.ErrSocketFailed'))
      this.socketError = true
    })

    if (process.env.NODE_ENV !== 'production') {
      window.pufferpanel.allowServerTasks = () => { this.showTasks = true }
    }
  },
  methods: {
    startRename () {
      this.newName = this.server.name
      this.renameOpen = true
    },
    resetRename () {
      this.renameOpen = false
      this.newName = ''
    },
    confirmRename () {
      this.$api.updateServerName(this.server.id, this.newName)
      this.server.name = this.newName
      this.resetRename()
    }
  }
}
</script>
