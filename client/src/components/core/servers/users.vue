<!--
  - Copyright 2019 Padduck, LLC
  -  Licensed under the Apache License, Version 2.0 (the "License");
  -  you may not use this file except in compliance with the License.
  -  You may obtain a copy of the License at
  -  	http://www.apache.org/licenses/LICENSE-2.0
  -  Unless required by applicable law or agreed to in writing, software
  -  distributed under the License is distributed on an "AS IS" BASIS,
  -  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  -  See the License for the specific language governing permissions and
  -  limitations under the License.
  -->

<template>
  <b-card
    header-tag="header">
    <h6 slot="header" class="mb-0">
      <span v-text="$t('common.Users')"></span>
    </h6>
    <b-card v-for="user in users" header-tag="header">
      <h6 slot="header" class="mb-0">
        <span v-text="user.username"></span>
      </h6>
      <b-card-text>
        <b-form-group>
          <b-form-checkbox-group
            v-model="user.scopes"
            :options="options"
            switches
            stacked
          ></b-form-checkbox-group>
          <b-button size="sm" v-on:click="updateUser(user)" v-text="$t('common.Update')" variant="primary"></b-button>
        </b-form-group>
      </b-card-text>
    </b-card>
  </b-card>

</template>

<script>
export default {
  data () {
    return {
      users: [],
      options: [
        { text: this.$t('scopes.ServersView'), value: 'servers.view' },
        { text: this.$t('scopes.ServersEdit'), value: 'servers.edit' },
        { text: this.$t('scopes.ServersInstall'), value: 'servers.install' },
        { text: this.$t('scopes.ServersConsole'), value: 'servers.console' },
        { text: this.$t('scopes.ServersConsoleSend'), value: 'servers.console.send' },
        { text: this.$t('scopes.ServersStop'), value: 'servers.stop' },
        { text: this.$t('scopes.ServersStart'), value: 'servers.start' },
        { text: this.$t('scopes.ServersStat'), value: 'servers.stat' },
        { text: this.$t('scopes.ServersFiles'), value: 'servers.files' },
        { text: this.$t('scopes.ServersFilesGet'), value: 'servers.files.get' },
        { text: this.$t('scopes.ServersFilesPut'), value: 'servers.files.put' },
        { text: this.$t('scopes.ServersEditUsers'), value: 'servers.edit.users' },
      ]
    }
  },
  mounted () {
    this.loadUsers()
  },
  methods: {
    loadUsers () {
      let vue = this
      this.$http.get('/api/servers/' + this.$route.params.id + '/user').then(function (response) {
        vue.users = response.data.data
      })
    },
    updateUser (user) {
      let vue = this
      if (user.scopes.length === 0) {
        this.$http.delete('/api/servers/' + this.$route.params.id + '/user/' + user.username).then(function (response) {
          vue.loadUsers()
        })
      } else {
        this.$http.put('/api/servers/' + this.$route.params.id + '/user/' + user.username, user).then(function (response) {
          vue.loadUsers()
        })
      }
    }
  }
}
</script>
