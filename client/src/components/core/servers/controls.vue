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
  <b-button-group>
    <b-button variant="success" v-on:click="start">Start</b-button>
    <b-button variant="warning" v-on:click="stop">Stop</b-button>
    <b-button variant="danger" v-on:click="kill">Kill</b-button>
  </b-button-group>
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
    }
  }
}
</script>
