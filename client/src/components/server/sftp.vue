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
    <v-card-title v-text="$t('servers.SFTPInfo')" />
    <v-card-text class="body-1 text--primary">
      <v-row>
        <v-col
          cols="12"
          sm="6"
          md="2"
          v-text="$t('common.Host') + ':' + $t('common.Port')"
        />
        <v-col
          cols="12"
          sm="6"
          md="10"
        >
          <input
            ref="host"
            :value="host"
            readonly
          >
          <v-btn
            icon
            @click="copyHost"
          >
            <v-icon>mdi-content-copy</v-icon>
          </v-btn>
          <v-chip
            v-if="copiedHost"
            color="success"
            class="mx-2"
            v-text="$t('common.Copied')"
          />
	  {{ host }}
        </v-col>
      </v-row>
      <v-divider />
      <v-row>
        <v-col
          cols="12"
          sm="6"
          md="2"
          v-text="$t('users.Username')"
        />
        <v-col
          cols="12"
          sm="6"
          md="10"
        >
          <input
            ref="username"
            :value="username"
            readonly
	    style=""
          >
          <v-btn
            icon
            @click="copyUsername"
          >
            <v-icon>mdi-content-copy</v-icon>
          </v-btn>
          <v-chip
            v-if="copiedUsername"
            color="success"
            class="mx-2"
            v-text="$t('common.Copied')"
          />
	  {{ username }}
        </v-col>
      </v-row>
      <v-divider />
      <v-row>
        <v-col
          cols="12"
          sm="6"
          md="2"
          v-text="$t('users.Password')"
        />
        <!-- 00A0 is the unicode code point for a non breaking space and required here because js makes &nbsp; print as literal text and not using a non breaking space makes it behave extra dumb on small devices... -->
        <v-col
          cols="12"
          sm="6"
          md="10"
          v-text="$t('users.AccountPassword').replace(' ', '\u00A0')"
        />
      </v-row>
    </v-card-text>
  </v-card>
</template>

<style>
input {
  width: 0;
}
</style>

<script>
export default {
  prop: {
    server: { type: Object, default: function () { return {} } }
  },
  data () {
    return {
      host: '',
      username: '',
      copiedHost: false,
      copiedUsername: false
    }
  },
  mounted () {
    this.host = this.$attrs.server.node.publicHost + ':' + this.$attrs.server.node.sftpPort
    const vue = this
    this.$http.get('/api/self').then(function (data) {
      const user = data.data
      vue.username = user.email + '|' + vue.$attrs.server.id
    })
  },
  methods: {
    copyHost () {
      const vue = this
      vue.$refs.host.select()
      document.execCommand('copy')
      vue.copiedUsername = false
      vue.copiedHost = true
      setTimeout(function () {
        vue.copiedHost = false
      }, 6000)
    },
    copyUsername () {
      const vue = this
      vue.$refs.username.select()
      document.execCommand('copy')
      vue.copiedHost = false
      vue.copiedUsername = true
      setTimeout(function () {
        vue.copiedUsername = false
      }, 6000)
    }
  }
}
</script>
