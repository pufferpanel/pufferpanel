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
        <v-row v-show="currentTab === 'console'">
          <v-col>
            <server-console :server="server" />
          </v-col>
        </v-row>
        <v-row v-show="currentTab === 'stats'">
          <v-col>
            <server-cpu :server="server" />
          </v-col>
        </v-row>
        <v-row v-show="currentTab === 'stats'">
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
        <!--<v-row v-show="currentTab === 'tasks'">
          <v-col>
            <server-tasks :server="server" />
          </v-col>
        </v-row>
        -->
        <v-row v-show="currentTab === 'settings'">
          <v-col>
            <server-settings :server="server" />
          </v-col>
        </v-row>
        <v-row v-show="currentTab === 'users'">
          <v-col>
            <server-users :server="server" />
          </v-col>
        </v-row>
        <v-row v-show="currentTab === 'admin'">
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
              <v-btn value="console">
                <span>{{ $t('servers.Console') }}</span>
                <v-icon>mdi-console-line</v-icon>
              </v-btn>
            </v-slide-item>
            <v-slide-item
              v-if="server.permissions.viewServerStats || isAdmin()"
              v-slot="{}"
            >
              <v-btn value="stats">
                <span>{{ $t('servers.Statistics') }}</span>
                <v-icon>mdi-chart-line</v-icon>
              </v-btn>
            </v-slide-item>
            <v-slide-item
              v-if="server.permissions.viewServerFiles || server.permissions.sftpServer || isAdmin()"
              v-slot="{}"
            >
              <v-btn value="files">
                <span>{{ $t('servers.Files') }}</span>
                <v-icon>mdi-file</v-icon>
              </v-btn>
            </v-slide-item>
            <!--<v-slide-item
              v-if="server.permissions.editServerData || isAdmin()"
              v-slot="{}"
            >
              <v-btn value="tasks">
                <span>{{ $t('servers.Tasks') }}</span>
                <v-icon>mdi-timer</v-icon>
              </v-btn>
            </v-slide-item>
            -->
            <v-slide-item
              v-if="server.permissions.editServerData || isAdmin()"
              v-slot="{}"
            >
              <v-btn value="settings">
                <span>{{ $t('servers.Settings') }}</span>
                <v-icon>mdi-cog</v-icon>
              </v-btn>
            </v-slide-item>
            <v-slide-item
              v-if="server.permissions.editServerUsers || isAdmin()"
              v-slot="{}"
            >
              <v-btn value="users">
                <span>{{ $t('users.Users') }}</span>
                <v-icon>mdi-account-multiple</v-icon>
              </v-btn>
            </v-slide-item>
            <v-slide-item
              v-if="server.permissions.deleteServer || isAdmin()"
              v-slot="{}"
            >
              <v-btn value="admin">
                <span>{{ $t('servers.Admin') }}</span>
                <v-icon>mdi-account-star</v-icon>
              </v-btn>
            </v-slide-item>
          </v-slide-group>
        </v-bottom-navigation>
      </v-col>
    </v-row>
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
      currentTab: 'console'
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
  }
}
</script>
