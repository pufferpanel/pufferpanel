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
    <v-expansion-panels multiple class="mb-2">
      <v-expansion-panel v-for="(item, i) in value" :key="i">
        <v-expansion-panel-header v-text="i" />
        <v-expansion-panel-content>
          <template-variable v-model="value[i]" />
          <v-btn color="error" block v-text="$t('common.Delete')" @click="delete value[i]; $forceUpdate()" />
        </v-expansion-panel-content>
      </v-expansion-panel>
    </v-expansion-panels>
    <v-btn v-if="!addingVariable" text v-text="$t('templates.AddVariable')" block @click="addingVariable = true" />
    <v-row v-else>
      <v-col cols="12" md="6">
        <v-text-field v-model="newVarName" :label="$t('common.Name')" dense outlined hide-details />
      </v-col>
      <v-col cols="6" md="3">
        <v-btn color="primary" v-text="$t('templates.AddVariable')" block @click="addVariable()" />
      </v-col>
      <v-col cols="6" md="3">
        <v-btn color="error" v-text="$t('common.Cancel')" block @click="newVarName = ''; addingVariable = false" />
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
      addingVariable: false,
      newVarName: '',
      variableTemplate: {
        required: true,
        userEdit: true,
        display: '',
        desc: '',
        type: null,
        value: ''
      }
    }
  },
  methods: {
    addVariable () {
      if (this.newVarName.trim().length > 0 && this.newVarName.trim().indexOf(' ') === -1) {
        this.value[this.newVarName.trim()] = { ...this.variableTemplate }
        this.newVarName = ''
        this.addingVariable = false
      } else {
        this.$toast.error(this.$t('templates.VarNameNoSpaces'))
      }
    }
  }
}
</script>
