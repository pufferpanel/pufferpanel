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
    <h6 slot="header" class="mb-0" v-text="$t('common.Controls')"></h6>
    <b-button-group>
      <b-button variant="success" v-on:click="start" v-text="$t('common.Start')"></b-button>
      <b-button variant="warning" v-on:click="stop" v-text="$t('common.Stop')"></b-button>
      <b-button variant="danger" v-on:click="kill" v-text="$t('common.Kill')"></b-button>
      <b-button variant="danger" v-on:click="install" v-text="$t('common.Install')"></b-button>
    </b-button-group>
  </b-card>
</template>

<script>
export default {
  methods: {
    start () {
      this.$socket.sendObj({ type: 'start' })
    },
    stop () {
      this.$socket.sendObj({ type: 'stop' })
    },
    kill () {
      this.$socket.sendObj({ type: 'kill' })
    },
    install () {
      this.$socket.sendObj({ type: 'install' })
    }
  }
}
</script>
