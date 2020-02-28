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
      <v-switch v-model="value.required" :label="$t('templates.Required')" hide-details class="mt-0" />
      <v-switch v-model="value.userEdit" :label="$t('templates.UserEditable')" hide-details class="mt-0" />
    </v-col>
    <v-col cols="12" md="6">
      <v-text-field v-model="value.display" :label="$t('templates.Display')" outlined hide-details />
    </v-col>
    <v-col cols="12" md="6">
      <v-text-field v-model="value.desc" :label="$t('templates.Description')" outlined hide-details />
    </v-col>
    <v-col cols="12" :md="value.type === 'option' ? '12' : '6'">
      <v-select v-model="value.type" :label="$t('templates.DataType')" :items="possibleTypes" outlined hide-details />
    </v-col>
    <v-col cols="12" md="6" v-if="value.type !== 'option'">
      <v-text-field v-model="value.value" :label="$t('templates.DefaultValue')" outlined hide-details />
    </v-col>
    <v-col cols="12" v-else>
      <v-row v-for="(option, i) in value.options" :key="i">
        <v-col cols="12" md="6">
          <v-text-field v-model="option.value" :label="$t('templates.Value')" outlined hide-details />
        </v-col>
        <v-col cols="12" md="6">
          <v-text-field v-model="option.display" :label="$t('templates.Display')" outlined hide-details append-outer-icon="mdi-close-circle" @click:append-outer="$delete(value.options, i)" />
        </v-col>
      </v-row>
      <v-btn block text v-text="$t('templates.AddOption')" @click="value.options.push({ value: '', display: '' })"/>
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
  }
}
</script>
