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
  <div class="mb-2">
    <v-expansion-panels multiple class="mb-2">
      <v-expansion-panel v-for="(env, i) in value" :key="i">
        <v-expansion-panel-header v-text="env.type" />
        <v-expansion-panel-content>
          <v-text-field v-if="env.type === 'docker'" v-model="env.image" :label="$t('templates.DockerImage')" outlined hide-details class="mb-2" />
          <v-btn color="error" v-text="$t('common.Delete')" block @click="$delete(value, i)" />
        </v-expansion-panel-content>
      </v-expansion-panel>
    </v-expansion-panels>
    <v-btn v-if="!addingEnv" color="primary" v-text="$t('templates.AddEnvironment')" block @click="addingEnv = true" />
    <v-row v-else>
      <v-col cols="12" md="6">
        <v-select v-model="newEnv" :label="$t('templates.Environment')" :items="possibleEnvironments" dense outlined hide-details />
      </v-col>
      <v-col cols="6" md="3">
        <v-btn color="primary" v-text="$t('templates.AddEnvironment')" block @click="addEnv()" />
      </v-col>
      <v-col cols="6" md="3">
        <v-btn color="error" v-text="$t('common.Cancel')" block @click="newEnv = 'standard'; addingEnv = false" />
      </v-col>
    </v-row>
  </div>
</template>

<script>
export default {
  props: {
    value: { type: Array, default: () => [] }
  },
  data () {
    return {
      addingEnv: false,
      newEnv: 'standard',
      possibleEnvironments: [
        'standard',
        'tty',
        'docker'
      ]
    }
  },
  methods: {
    addEnv () {
      this.value.push({ type: this.newEnv })
      this.newEnv = 'standard'
      this.addingEnv = false
    }
  }
}
</script>
