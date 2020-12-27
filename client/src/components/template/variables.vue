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
    <v-list>
      <v-list-item
        v-for="(item, name) in value"
        :key="name"
        link
        @click="startEdit(name)"
      >
        <v-list-item-content v-text="name" />
        <v-list-item-action class="flex-row">
          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <v-btn
                icon
                v-on="on"
              >
                <v-icon>mdi-pencil</v-icon>
              </v-btn>
            </template>
            <span v-text="$t('common.Edit')" />
          </v-tooltip>
          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <v-btn
                icon
                v-on="on"
                @click.stop="remove(name)"
              >
                <v-icon>mdi-close</v-icon>
              </v-btn>
            </template>
            <span v-text="$t('common.Delete')" />
          </v-tooltip>
        </v-list-item-action>
      </v-list-item>
    </v-list>
    <v-btn
      text
      block
      @click="add = true"
      v-text="$t('common.Add')"
    />
    <ui-overlay
      v-model="add"
      :title="$t('common.Add')"
      card
      closable
      @close="reset()"
    >
      <v-row>
        <v-col cols="12">
          <ui-input
            v-model="newVarName"
            :label="$t('common.Name')"
          />
        </v-col>
        <v-col cols="12">
          <v-btn
            color="success"
            block
            @click="addVariable()"
            v-text="$t('common.Add')"
          />
        </v-col>
      </v-row>
    </ui-overlay>
    <ui-overlay
      v-model="edit"
      :title="currentEdit"
      card
      closable
      @close="reset()"
    >
      <template-variable v-model="value[currentEdit]" />
    </ui-overlay>
  </div>
</template>

<script>
export default {
  props: {
    value: { type: Object, default: () => {} }
  },
  data () {
    return {
      add: false,
      edit: false,
      newVarName: '',
      currentEdit: '',
      variableTemplate: {
        required: true,
        userEdit: true,
        display: '',
        desc: '',
        type: 'string',
        value: ''
      }
    }
  },
  methods: {
    addVariable () {
      const name = this.newVarName.trim()
      if (name.length > 0 && name.indexOf(' ') === -1) {
        const changed = { ...this.value }
        changed[name] = { ...this.variableTemplate }
        this.$emit('input', changed)
        this.reset()
        this.startEdit(name)
      } else {
        this.$toast.error(this.$t('templates.VarNameNoSpaces'))
      }
    },
    startEdit (name) {
      this.currentEdit = name
      this.edit = true
    },
    remove (name) {
      const changed = { ...this.value }
      delete changed[name]
      this.$emit('input', changed)
    },
    reset () {
      this.add = false
      this.edit = false
      this.newVarName = ''
      this.currentEdit = ''
    }
  }
}
</script>
