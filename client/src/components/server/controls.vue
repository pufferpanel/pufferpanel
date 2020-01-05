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
  <v-container>
    <v-btn
      v-if="server.permissions.startServer || isAdmin()"
      class="mr-4"
      color="success"
      @click="action('start')"
      v-text="$t('servers.Start')"
    />
    <v-btn
      v-if="server.permissions.stopServer || isAdmin()"
      class="mr-4"
      color="warning"
      @click="action('stop')"
      v-text="$t('servers.Stop')"
    />
    <v-btn
      v-if="server.permissions.stopServer || isAdmin()"
      class="mr-4"
      color="error"
      @click="action('kill')"
      v-text="$t('servers.Kill')"
    />
    <v-btn
      v-if="server.permissions.installServer || isAdmin()"
      color="error"
      @click="action('install')"
      v-text="$t('servers.Install')"
    />
  </v-container>
</template>

<script>
export default {
  props: {
    server: { type: Object, default: function () { return {} } }
  },
  methods: {
    action (act) {
      this.$socket.sendObj({ type: act })
    }
  }
}
</script>
