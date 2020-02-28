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
  <div>
    <v-row>
      <v-col cols="12">
        <v-text-field v-model="value.program" :label="$t('templates.Program')" outlined hide-details />
      </v-col>
      <v-col cols="12" class="pb-2 pt-0">
        <span v-text="$t('templates.Arguments')" class="body-1" />
      </v-col>
      <v-col cols="12" v-for="(argument, i) in value.arguments" :key="i" class="pb-0 pt-0">
        <v-text-field v-model="value.arguments[i]" label="" dense outlined hide-details append-outer-icon="mdi-close-circle" @click:append-outer="$delete(value.arguments, i)" />
      </v-col>
      <v-col cols="12">
        <v-btn v-text="$t('templates.AddArgument')" text block @click="addArgument()" />
      </v-col>
    </v-row>
  </div>
</template>

<script>
export default {
  props: {
    value: { type: Object, default: () => {} }
  },
  methods: {
    addArgument () {
      this.value.arguments.push('')
    }
  }
}
</script>
