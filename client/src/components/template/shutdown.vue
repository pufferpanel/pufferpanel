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
  <v-row>
    <v-col cols="12" class="pb-0">
      <v-btn-toggle v-model="stopType" borderless dense mandatory>
        <v-btn value="command" v-text="$t('templates.stop.Command')" />
        <v-btn value="signal" v-text="$t('templates.stop.Signal')" />
      </v-btn-toggle>
    </v-col>
    <v-col cols="12">
      <v-text-field v-if="stopType === 'command'" v-model="value.stop" :label="$t('templates.stop.Command')" outlined hide-details />
      <v-text-field v-else v-model="value.stopCode" :label="$t('templates.stop.Signal')" type="number" outlined hide-details />
    </v-col>
  </v-row>
</template>

<script>
export default {
  props: {
    value: { type: Object, default: () => {} }
  },
  data () {
    return {
      stopType: this.value.stop ? 'command' : 'signal'
    }
  }
}
</script>
