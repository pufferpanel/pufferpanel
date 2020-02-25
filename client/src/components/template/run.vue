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
        <v-btn v-text="$t('templates.AddArgument')" block color="primary" @click="addArgument()" />
      </v-col>
      <v-col cols="12" class="py-2">
        <span v-text="$t('templates.Shutdown')" class="title" />
      </v-col>
      <v-col cols="12" class="pb-0">
        <v-btn-toggle v-model="stopType" borderless dense mandatory>
          <v-btn value="command" v-text="$t('templates.stop.Command')" />
          <v-btn value="signal" v-text="$t('templates.stop.Signal')" />
        </v-btn-toggle>
      </v-col>
      <v-col cols="12">
        <v-text-field v-if="stopType === 'command'" v-model="value.stop" :label="$t('templates.stop.Command')" dense outlined hide-details />
        <v-text-field v-else v-model="value.stopCode" :label="$t('templates.stop.Signal')" type="number" dense outlined hide-details />
      </v-col>
      <v-col cols="12" class="py-1">
        <span v-text="$t('templates.PreHook')" class="title" />
      </v-col>
      <v-col cols="12" class="py-1">
        <template-processors v-model="value.pre" name="pre" />
      </v-col>
      <v-col cols="12" class="py-1">
        <span v-text="$t('templates.PostHook')" class="title" />
      </v-col>
      <v-col cols="12" class="py-1">
        <template-processors v-model="value.post" name="post" />
      </v-col>
      <v-col cols="12">
        <span v-text="$t('templates.EnvVars')" class="title" />
      </v-col>
    </v-row>
    <v-row v-for="(env, i) in value.environmentVars" :key="i">
      <v-col cols="4" class="my-0 py-0">
        <v-subheader v-text="i" />
      </v-col>
      <v-col cols="8" class="my-0 py-0">
        <v-text-field v-model="value.environmentVars[i]" dense outlined hide-details append-outer-icon="mdi-close-circle" @click:append-outer="delete value.environmentVars[i]; $forceUpdate()" />
      </v-col>
    </v-row>
    <v-btn v-if="!addingEnvVar" color="primary" block v-text="$t('templates.AddEnvVar')" @click="addingEnvVar = true" />
    <v-row v-else>
      <v-col cols="12" md="6">
        <v-text-field v-model="newEnvVar" :label="$t('common.Name')" dense outlined hide-details />
      </v-col>
      <v-col cols="6" md="3">
        <v-btn color="primary" v-text="$t('templates.AddEnvVar')" block @click="addEnvVar()" />
      </v-col>
      <v-col cols="6" md="3">
        <v-btn color="error" v-text="$t('common.Cancel')" block @click="newEnvVar = ''; addingEnvVar = false" />
      </v-col>
    </v-row>
  </div>
</template>

<script>
export default {
  props: {
    value: { type: Object, default: () => {} }
  },
  data () {
    return {
      addingEnvVar: false,
      newEnvVar: '',
      stopType: this.value.stop ? 'command' : 'signal'
    }
  },
  methods: {
    addEnvVar () {
      if (this.newEnvVar.trim().length > 0 && this.newEnvVar.trim().indexOf(' ') === -1) {
        this.value.environmentVars[this.newEnvVar.trim()] = ''
        this.newEnvVar = ''
        this.addingEnvVar = false
      } else {
        this.$toast.error(this.$t('templates.VarNameNoSpaces'))
      }
    },
    addArgument () {
      this.value.arguments.push('')
    }
  }
}
</script>
