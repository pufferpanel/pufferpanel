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
    <v-card-title v-text="$t('common.Users')" />
    <v-card-text>
      <v-row>
        <v-col
          v-for="user in users"
          cols="12"
          md="6"
        >
          <v-card class="mb-4">
            <v-card-title>
              <span v-text="user.username" />
              <v-btn
                icon
                @click="toggleEdit(user.username)"
              >
                <v-icon v-text="editUsers.indexOf(user.username) > -1 ? 'mdi-close' : 'mdi-pencil'" />
              </v-btn>
            </v-card-title>
            <v-card-text v-if="editUsers.indexOf(user.username) > -1">
              <v-row>
                <v-col
                  v-for="option in options"
                  class="pt-1 pb-0"
                  cols="12"
                  md="6"
                >
                  <v-switch
                    v-model="user.scopes"
                    hide-details
                    :label="option.text"
                    :value="option.value"
                    color="primary"
                  />
                </v-col>
              </v-row>
              <v-row class="mt-2">
                <v-col
                  cols="12"
                  md="6"
                >
                  <v-btn
                    large
                    block
                    color="success"
                    @click="updateUser(user)"
                    v-text="$t('common.Update')"
                  />
                </v-col>
                <v-col
                  cols="12"
                  md="6"
                >
                  <v-btn
                    large
                    block
                    color="error"
                    @click="deleteUser(user.username)"
                    v-text="$t('common.Delete')"
                  />
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script>
export default {
  data () {
    return {
      users: [],
      editUsers: [],
      options: [
        { text: this.$t('scopes.ServersView'), value: 'servers.view' },
        { text: this.$t('scopes.ServersEdit'), value: 'servers.edit' },
        { text: this.$t('scopes.ServersInstall'), value: 'servers.install' },
        { text: this.$t('scopes.ServersConsole'), value: 'servers.console' },
        { text: this.$t('scopes.ServersConsoleSend'), value: 'servers.console.send' },
        { text: this.$t('scopes.ServersStop'), value: 'servers.stop' },
        { text: this.$t('scopes.ServersStart'), value: 'servers.start' },
        { text: this.$t('scopes.ServersStat'), value: 'servers.stats' },
        { text: this.$t('scopes.ServersFiles'), value: 'servers.files' },
        { text: this.$t('scopes.ServersFilesGet'), value: 'servers.files.get' },
        { text: this.$t('scopes.ServersFilesPut'), value: 'servers.files.put' },
        { text: this.$t('scopes.ServersEditUsers'), value: 'servers.edit.users' }
      ]
    }
  },
  mounted () {
    this.loadUsers()
  },
  methods: {
    loadUsers () {
      const vue = this
      this.$http.get('/api/servers/' + this.$route.params.id + '/user').then(function (response) {
        vue.users = response.data.data
      })
    },
    updateUser (user) {
      const vue = this
      if (user.scopes.length === 0) {
        this.deleteUser(user.username)
      } else {
        this.$http.put('/api/servers/' + this.$route.params.id + '/user/' + user.username, user).then(function (response) {
          vue.loadUsers()
        })
      }
    },
    toggleEdit (username) {
      if (this.editUsers.indexOf(username) > -1) {
        this.editUsers.splice(this.editUsers.indexOf(username), 1)
      } else {
        this.editUsers.push(username)
      }
    },
    deleteUser (username) {
      this.$http.delete('/api/servers/' + this.$route.params.id + '/user/' + username).then(function (response) {
        vue.loadUsers()
      })
    }
  }
}
</script>
