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
  <div class="mb-4">
    <v-expansion-panels multiple class="mb-2">
      <v-expansion-panel v-for="(entry, i) in value" :key="i">
        <v-expansion-panel-header v-text="entry.type" />
        <v-expansion-panel-content>
          <template-processor v-model="value[i]" />
          <v-btn color="error" block v-text="$t(getDeleteKey())" @click="$delete(value, i)" />
        </v-expansion-panel-content>
      </v-expansion-panel>
    </v-expansion-panels>
    <v-btn color="primary" block v-text="$t(getAddKey())" @click="value.push({ ...template }); $forceUpdate()" />
  </div>
</template>

<script>
export default {
  props: {
    value: { type: Array, default: () => [] },
    name: { type: String, default: () => 'install' }
  },
  data () {
    return {
      template: {
        type: 'command',
        commands: ['']
      }
    }
  },
  methods: {
    getAddKey () {
      const capitalized = this.name.charAt(0).toUpperCase() + this.name.slice(1)
      return `templates.Add${capitalized}Step`
    },
    getDeleteKey () {
      const capitalized = this.name.charAt(0).toUpperCase() + this.name.slice(1)
      return `templates.Delete${capitalized}Step`
    }
  }
}
</script>
