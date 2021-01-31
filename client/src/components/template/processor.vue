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
    <v-col>
      <ui-select
        v-model="value.type"
        :items="processorTypes"
        :label="$t('templates.Type')"
        @change="changeType()"
      />
    </v-col>
    <v-col
      v-for="field in fields"
      :key="value.type + field.name"
      cols="12"
    >
      <h4
        v-if="field.type === 'map' || field.type === 'list' || field.headline"
        v-text="getLabel(field)"
      />

      <component
        :is="field.component"
        v-if="field.type === 'custom'"
        :value="value[field.name] || field.default"
        @input="onInput(field.name, $event)"
      />
      <ui-map-input
        v-else-if="field.type === 'map'"
        :value="value[field.name] || field.default"
        :key-label="field.keyLabel ? $t(field.keyLabel) : undefined"
        :value-label="field.valueLabel ? $t(field.valueLabel) : undefined"
        @input="onInput(field.name, $event)"
      />
      <ui-list-input
        v-else-if="field.type === 'list'"
        :value="value[field.name] || field.default"
        @input="onInput(field.name, $event)"
      />
      <ui-switch
        v-else-if="field.type === 'boolean'"
        :label="getLabel(field)"
        :value="value[field.name] || field.default"
        @input="onInput(field.name, $event)"
      />
      <ui-input-suggestions
        v-else-if="field.options !== undefined"
        :type="field.type"
        :label="getLabel(field)"
        :items="field.options"
        :value="value[field.name] || field.default"
        @input="onInput(field.name, $event)"
      />
      <v-textarea
        v-else-if="field.type === 'textarea'"
        :value="value[field.name] || field.default"
        :label="getLabel(field)"
        outlined
        hide-details
        @input="onInput(field.name, $event)"
      />
      <ui-input
        v-else
        :type="field.type"
        :label="getLabel(field)"
        :value="value[field.name] || field.default"
        @input="onInput(field.name, $event)"
      />
    </v-col>
  </v-row>
</template>

<script>
const processors = {
  download: [
    {
      name: 'files',
      type: 'list',
      default: []
    }
  ],
  command: [
    {
      name: 'commands',
      type: 'list',
      default: []
    }
  ],
  alterfile: [
    {
      name: 'file',
      type: 'text',
      label: 'templates.Filename',
      default: ''
    },
    {
      name: 'regex',
      type: 'boolean',
      default: true
    },
    {
      name: 'search',
      type: 'text',
      default: ''
    },
    {
      name: 'replace',
      type: 'text',
      default: ''
    }
  ],
  writefile: [
    {
      name: 'target',
      type: 'text',
      label: 'templates.Filename',
      default: ''
    },
    {
      name: 'text',
      type: 'textarea',
      default: ''
    }
  ],
  move: [
    {
      name: 'source',
      type: 'text',
      default: ''
    },
    {
      name: 'target',
      type: 'text',
      default: ''
    }
  ],
  mkdir: [
    {
      name: 'target',
      type: 'text',
      label: 'common.Name',
      default: ''
    }
  ],
  mojangdl: [
    {
      name: 'version',
      type: 'text',
      label: 'templates.Version',
      default: ''
    },
    {
      name: 'target',
      type: 'text',
      label: 'templates.Filename',
      default: ''
    }
  ],
  forgedl: [
    {
      name: 'version',
      type: 'text',
      label: 'templates.Version',
      default: ''
    },
    {
      name: 'filename',
      type: 'text',
      label: 'templates.Filename',
      default: ''
    }
  ],
  spongeforgedl: [
    {
      name: 'releaseType',
      type: 'text',
      default: ''
    }
  ],
  fabricdl: [
    {
      name: 'targetFile',
      type: 'text',
      label: 'templates.Filename',
      default: ''
    }
  ]
}

export default {
  props: {
    value: { type: Object, default: () => {} }
  },
  computed: {
    fields () {
      return processors[this.value.type]
    },
    processorTypes () {
      return Object.keys(processors).map(elem => {
        return { value: elem, text: this.$t(`templates.processors.${elem}.ProcessorName`) }
      })
    }
  },
  mounted () {
    const defaulted = {}
    this.fields.map(elem => {
      if (!this.value[elem.name]) {
        defaulted[elem.name] = elem.default
      }
    })
    this.$emit('input', { ...this.value, ...defaulted })
  },
  methods: {
    changeType () {
      const changed = { ...this.value }
      Object.keys(changed).map(elem => {
        if (elem !== 'type') changed[elem] = undefined
      })
      this.$emit('input', changed)
    },
    onInput (field, event) {
      this.$emit('input', { ...this.value, [field]: event })
    },
    getLabel (field) {
      return field.label ? this.$t(field.label) : this.$t(`templates.processors.${this.value.type}.${field.name}`)
    }
  }
}
</script>
