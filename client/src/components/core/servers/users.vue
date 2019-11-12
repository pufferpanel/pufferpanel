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
              <v-text-field hide-details outlined v-model="user.email" type="email" :label="$t('common.Email')" prepend-inner-icon="mdi-email" @keyup.enter="updateUser(user)" v-if="user.new" />
              <span v-text="user.username || user.email" v-if="!user.new" />
              <v-btn v-if="!user.new"
                icon
                @click="toggleEdit(user.username)"
              >
                <v-icon v-text="editUsers.indexOf(user.username) > -1 ? 'mdi-close' : 'mdi-pencil'" />
              </v-btn>
            </v-card-title>
            <v-card-text v-if="user.new || (editUsers.indexOf(user.username) > -1)">
              <v-row>
                <v-col
                  v-for="option in options"
                  class="pt-1 pb-0"
                  cols="12"
                  md="6"
                >
                  <v-switch
                    v-model="user[option.value]"
                    hide-details
                    :label="option.text"
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
                    v-text="user.new ? $t('common.Save') : $t('common.Update')"
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
                    @click="deleteUser(user)"
                    v-text="user.new ? $t('common.Cancel') : $t('common.Delete')"
                  />
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-btn large block color="success" @click="addUser()" v-text="$t('common.AddUser')" />
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
        { text: this.$t('scopes.ServersEdit'), value: 'editServerData' },
        { text: this.$t('scopes.ServersInstall'), value: 'installServer' },
        { text: this.$t('scopes.ServersConsole'), value: 'viewServerConsole' },
        { text: this.$t('scopes.ServersConsoleSend'), value: 'sendServerConsole' },
        { text: this.$t('scopes.ServersStop'), value: 'stopServer' },
        { text: this.$t('scopes.ServersStart'), value: 'startServer' },
        { text: this.$t('scopes.ServersStat'), value: 'viewServerStats' },
        { text: this.$t('scopes.ServersFiles'), value: 'sftpServer' },
        { text: this.$t('scopes.ServersFilesGet'), value: 'viewServerFiles' },
        { text: this.$t('scopes.ServersFilesPut'), value: 'putServerFiles' },
        { text: this.$t('scopes.ServersEditUsers'), value: 'editServerUsers' }
      ]
    }
  },
  mounted () {
    this.loadUsers()
  },
  methods: {
    loadUsers () {
      const ctx = this
      this.$http.get('/api/servers/' + this.$route.params.id + '/user').then(function (response) {
        ctx.users = response.data
      })
    },
    addUser () {
      const newUser = {}
      newUser.new = true
      newUser.email = ''
      for (const option in this.options) {
        newUser[this.options[option].value] = false
      }
      this.users.push(newUser)
    },
    updateUser (user) {
      if (user.new && (!user.email || user.email === '')) {
        this.$toast.error(this.$t('common.NoEmailGiven'))
        return
      }
      const ctx = this
      for (const key of Object.keys(user)) {
        user[key] = (user[key] === true) ? 'true' : (user[key] === false) ? 'false' : user[key]
      }
      this.$http.put('/api/servers/' + this.$route.params.id + '/user/' + user.email, user).then(function (response) {
        ctx.$toast.success(ctx.$t('common.SavedUsers'))
        ctx.loadUsers()
      })
    },
    toggleEdit (username) {
      if (this.editUsers.indexOf(username) > -1) {
        this.editUsers.splice(this.editUsers.indexOf(username), 1)
      } else {
        this.editUsers.push(username)
      }
    },
    deleteUser (user) {
      if (user.new) {
        this.$delete(this.users, this.users.indexOf(user))
        return
      }
      const ctx = this
      this.$http.delete('/api/servers/' + this.$route.params.id + '/user/' + user.email).then(function (response) {
        ctx.loadUsers()
      })
    }
  }
}
</script>
