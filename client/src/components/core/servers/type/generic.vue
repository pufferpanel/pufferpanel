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
  <b-container>
    <h2 v-text="server.name"/>
    <br>
    <b-row>
      <b-col>
        <core-servers-console v-bind:server="server"/>
      </b-col>
    </b-row>
    <b-row>
      <b-col>
        <core-servers-controls/>
      </b-col>
    </b-row>
    <b-row>
      <b-col>
        <b-card-group>
          <core-servers-cpu/>
          <core-servers-memory/>
        </b-card-group>
      </b-col>
    </b-row>
    <b-row>
      <b-col>
        <core-servers-files v-bind:server="server"/>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
export default {
  props: {
    server: Object
  }
}
</script>
