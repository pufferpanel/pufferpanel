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
    <v-col cols="12">
      <v-select
        v-model="value.type"
        :items="processorTypes"
        :label="$t('templates.Type')"
        outlined
        hide-details
        @change="changeType"
      />
    </v-col>
    <v-col
      v-if="value.type === 'download'"
      cols="12"
    >
      <v-text-field
        v-for="(e, i) in value.files"
        :key="i"
        v-model="value.files[i]"
        dense
        outlined
        hide-details
        append-outer-icon="mdi-close-circle"
        @click:append-outer="$delete(value.files, i)"
      />
    </v-col>
    <v-col
      v-if="value.type === 'download'"
      cols="12"
    >
      <v-btn
        text
        block
        @click="value.files.push(''); $forceUpdate()"
        v-text="$t('templates.AddFile')"
      />
    </v-col>
    <v-col
      v-if="value.type === 'command'"
      cols="12"
    >
      <v-text-field
        v-for="(e, i) in value.commands"
        :key="i"
        v-model="value.commands[i]"
        dense
        outlined
        hide-details
        append-outer-icon="mdi-close-circle"
        @click:append-outer="$delete(value.commands, i)"
      />
    </v-col>
    <v-col
      v-if="value.type === 'command'"
      cols="12"
    >
      <v-btn
        text
        block
        @click="value.commands.push(''); $forceUpdate()"
        v-text="$t('templates.AddCommand')"
      />
    </v-col>
    <v-col
      v-if="value.type === 'writefile'"
      cols="12"
    >
      <v-text-field
        v-model="value.target"
        outlined
        hide-details
        :label="$t('templates.Filename')"
      />
    </v-col>
    <v-col
      v-if="value.type === 'writefile'"
      cols="12"
    >
      <v-textarea
        v-model="value.text"
        :label="$t('templates.Content')"
        outlined
        hide-details
      />
    </v-col>
    <v-col
      v-if="value.type === 'move'"
      cols="12"
      md="6"
    >
      <v-text-field
        v-model="value.source"
        outlined
        hide-details
        :label="$t('templates.Source')"
      />
    </v-col>
    <v-col
      v-if="value.type === 'move'"
      cols="12"
      md="6"
    >
      <v-text-field
        v-model="value.target"
        outlined
        hide-details
        :label="$t('templates.Target')"
      />
    </v-col>
    <v-col
      v-if="value.type === 'mkdir'"
      cols="12"
    >
      <v-text-field
        v-model="value.target"
        outlined
        hide-details
        :label="$t('common.Name')"
      />
    </v-col>
    <v-col
      v-if="value.type === 'mojangdl'"
      cols="12"
    >
      <v-text-field
        v-model="value.version"
        outlined
        hide-details
        :label="$t('templates.Version')"
      />
    </v-col>
    <v-col
      v-if="value.type === 'mojangdl'"
      cols="12"
    >
      <v-text-field
        v-model="value.target"
        outlined
        hide-details
        :label="$t('templates.Filename')"
      />
    </v-col>
    <v-col
      v-if="value.type === 'forgedl'"
      cols="12"
    >
      <v-text-field
        v-model="value.version"
        outlined
        hide-details
        :label="$t('templates.Version')"
      />
    </v-col>
    <v-col
      v-if="value.type === 'forgedl'"
      cols="12"
    >
      <v-text-field
        v-model="value.filename"
        outlined
        hide-details
        :label="$t('templates.Filename')"
      />
    </v-col>
    <v-col
      v-if="value.type === 'spongeforgedl'"
      cols="12"
    >
      <v-text-field
        v-model="value.releaseType"
        outlined
        hide-details
        :label="$t('templates.ReleaseType')"
      />
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
      processorTypes: [
        {
          value: 'download',
          text: 'Download'
        },
        {
          value: 'command',
          text: 'Run Command'
        },
        {
          value: 'writefile',
          text: 'Write to file'
        },
        {
          value: 'move',
          text: 'Move File'
        },
        {
          value: 'mkdir',
          text: 'Create Directory'
        },
        {
          value: 'mojangdl',
          text: 'Download Minecraft'
        },
        {
          value: 'forgedl',
          text: 'Download Minecraft Forge'
        },
        {
          value: 'spongeforgedl',
          text: 'Download Minecraft SpongeForge'
        }
      ]
    }
  },
  methods: {
    changeType (newType) {
      this.value.files = newType === 'download' ? [] : undefined
      this.value.commands = newType === 'command' ? [] : undefined
      this.value.target = undefined
      this.value.text = undefined
      this.value.source = undefined
      this.value.version = undefined
      this.value.filename = undefined
      this.value.releaseType = undefined
    }
  }
}
</script>
