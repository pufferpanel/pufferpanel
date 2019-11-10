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
    <h6 slot="header" class="mb-0" v-text="$t('common.SFTPInfo')"></h6>
    <b-input-group :prepend="$t('common.HostPort')" class="mt-3">
      <b-form-input readonly v-bind:value="host"></b-form-input>
    </b-input-group>
    <b-input-group :prepend="$t('common.Username')" class="mt-3" >
      <b-form-input readonly v-bind:value="username"></b-form-input>
    </b-input-group>
    <!--<b-input-group :prepend="$t('common.Password')" class="mt-3">
      <b-form-input readonly v-bind:value="$t('common.AccountPassword')"></b-form-input>
    </b-input-group>-->
  </b-card>
</template>

<script>
export default {
  prop: {
    server: Object
  },
  data () {
    return {
      host: this.$attrs.server.node.publicHost + ":" + this.$attrs.server.node.sftpPort,
      username: ""
    }
  },
  mounted () {
    let vue = this
    this.$http.get('/api/users').then(function (data) {
      let user = data.data.data
      vue.username = user.email + "|" + vue.$attrs.server.id
    })
  }
}
</script>
