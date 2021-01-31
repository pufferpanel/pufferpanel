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
      <ui-switch
        v-model="value.required"
        :label="$t('templates.Required')"
      />
      <ui-switch
        v-model="value.userEdit"
        :label="$t('templates.UserEditable')"
      />
    </v-col>
    <v-col
      cols="12"
      md="6"
    >
      <ui-input
        v-model="value.display"
        :label="$t('templates.Display')"
      />
    </v-col>
    <v-col
      cols="12"
      md="6"
    >
      <ui-input
        v-model="value.desc"
        :label="$t('templates.Description')"
      />
    </v-col>
    <v-col
      cols="12"
      :md="value.type === 'option' ? '12' : '6'"
    >
      <ui-select
        v-model="value.type"
        :label="$t('templates.DataType')"
        :items="possibleTypes"
        @change="typeChanged()"
      />
    </v-col>
    <v-col
      v-if="value.type !== 'option'"
      cols="12"
      md="6"
    >
      <ui-input
        v-model="value.value"
        :label="$t('templates.DefaultValue')"
      />
    </v-col>
    <v-col
      v-if="value.type === 'option' || value.type === 'string'"
      cols="12"
    >
      <v-row
        v-for="(option, i) in value.options"
        :key="i"
      >
        <v-col
          cols="12"
          md="6"
        >
          <ui-input
            v-model="option.value"
            :label="$t('templates.Value')"
          />
        </v-col>
        <v-col
          cols="12"
          md="6"
        >
          <ui-input
            v-model="option.display"
            :label="$t('templates.Display')"
            icon-behind="mdi-close"
            @click:append-outer="$delete(value.options, i)"
          />
        </v-col>
      </v-row>
      <v-btn
        block
        text
        @click="value.options.push({ value: '', display: '' })"
        v-text="$t('templates.AddOption')"
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
      possibleTypes: [
        {
          text: 'String',
          value: 'string'
        },
        {
          text: 'Boolean',
          value: 'boolean'
        },
        {
          text: 'Integer',
          value: 'integer'
        },
        {
          text: 'Options',
          value: 'option'
        }
      ]
    }
  },
  methods: {
    typeChanged () {
      if (this.value.type === 'option' && !this.value.options) {
        this.value.options = []
        this.$emit('input', { ...this.value })
      }
    }
  }
}
</script>
