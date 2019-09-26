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
    <h1
      style="float: left;"
      v-text="server.name"
    />
    <div style="float: right;">
      <core-servers-controls :server="server" />
    </div>
    <div style="clear: both;" />
    <v-row v-if="hasScope('servers.console', server.id)">
      <v-col>
        <core-servers-console :server="server" />
      </v-col>
    </v-row>
    <v-row v-if="hasScope('servers.stats', server.id)">
      <v-col md="6">
        <core-servers-cpu />
      </v-col>
      <v-col md="6">
        <core-servers-memory />
      </v-col>
    </v-row>
    <v-row v-if="hasScope('servers.files.get', server.id)">
      <v-col>
        <core-servers-files :server="server" />
      </v-col>
    </v-row>
    <v-row v-if="hasScope('servers.files', server.id)">
      <v-col>
        <core-servers-sftp :server="server" />
      </v-col>
    </v-row>
    <v-row v-if="hasScope('servers.edit.users', server.id)">
      <v-col>
        <core-servers-users :server="server" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  props: {
    server: { type: Object, default: function () { return {} } }
  }
}
</script>
