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
    >
      <server-status :server="server" />
      {{server.name}}
    </h1>
    <div style="float: right;">
      <server-controls :server="server" />
    </div>
    <div style="clear: both;" />
    <v-row v-if="server.permissions.viewServerConsole || isAdmin()">
      <v-col>
        <server-console :server="server" />
      </v-col>
    </v-row>
    <v-row v-if="server.permissions.viewServerStats || isAdmin()">
      <v-col cols="12" md="6">
        <server-cpu />
      </v-col>
      <v-col cols="12" md="6">
        <server-memory />
      </v-col>
    </v-row>
    <v-row v-if="server.permissions.viewServerFiles || isAdmin()">
      <v-col>
        <server-files :server="server" />
      </v-col>
    </v-row>
    <v-row v-if="server.permissions.sftpServer || isAdmin()">
      <v-col>
        <server-sftp :server="server" />
      </v-col>
    </v-row>
    <v-row v-if="server.permissions.editServerUsers || isAdmin()">
      <v-col>
        <server-users :server="server" />
      </v-col>
    </v-row>
    <v-row v-if="server.permissions.editServer || isAdmin()">
      <v-col>
        <server-settings :server="server" />
      </v-col>
    </v-row>
    <v-row v-if="server.permissions.deleteServer || isAdmin()">
      <v-col>
        <server-admin :server="server" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  props: {
    server: { type: Object, default: () => {} }
  }
}
</script>
