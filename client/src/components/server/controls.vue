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
  <v-slide-group class="pl-12 max-width-100">
    <v-slide-item
      v-if="server.permissions.startServer || isAdmin()"
      v-slot="{}"
    >
      <v-btn
        v-if="!running"
        class="mr-4"
        color="success"
        :loading="restarting"
        @click="action('start')"
      >
        <v-icon left>
          mdi-play
        </v-icon>
        {{ $t('servers.Start') }}
      </v-btn>
      <v-btn
        v-else
        class="mr-4"
        color="success"
        :loading="restarting"
        @click="restart()"
      >
        <v-icon left>
          mdi-reload
        </v-icon>
        {{ $t('servers.Restart') }}
      </v-btn>
    </v-slide-item>
    <v-slide-item
      v-if="server.permissions.stopServer || isAdmin()"
      v-slot="{}"
    >
      <v-btn
        class="mr-4"
        color="warning"
        @click="action('stop')"
      >
        <v-icon left>
          mdi-stop
        </v-icon>
        {{ $t('servers.Stop') }}
      </v-btn>
    </v-slide-item>
    <v-slide-item
      v-if="server.permissions.stopServer || isAdmin()"
      v-slot="{}"
    >
      <v-btn
        class="mr-4"
        color="error"
        @click="action('kill')"
      >
        <v-icon left>
          mdi-skull
        </v-icon>
        {{ $t('servers.Kill') }}
      </v-btn>
    </v-slide-item>
    <v-slide-item
      v-if="server.permissions.installServer || isAdmin()"
      v-slot="{}"
    >
      <v-btn
        color="error"
        @click="action('install')"
      >
        <v-icon left>
          mdi-package-down
        </v-icon>
        {{ $t('servers.Install') }}
      </v-btn>
    </v-slide-item>
  </v-slide-group>
</template>

<script>
export default {
  props: {
    server: { type: Object, default: function () { return {} } }
  },
  data () {
    return {
      running: false,
      restarting: false
    }
  },
  mounted () {
    this.$api.addServerListener(this.server.id, 'status', event => {
      this.running = event.running

      if (this.restarting && !event.running) {
        setTimeout(() => { this.action('start') }, 1000)
        return
      }

      if (this.restarting && event.running) this.restarting = false
    })

    this.$api.requestServerStatus(this.server.id)
  },
  methods: {
    restart () {
      this.restarting = true
      this.action('stop')
    },
    action (action) {
      this.$api.sendServerAction(this.server.id, action)
    }
  }
}
</script>
