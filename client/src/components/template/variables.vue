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
        v-for="(item, index) in value"
        :key="item.name"
        link
        @click="startEdit(index)"
      >
        <v-list-item-content v-text="item.display" />
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
                @click.stop="remove(index)"
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
      @click="addVar()"
      v-text="$t('common.Add')"
    />
    <ui-overlay
      v-model="edit"
      card
      closable
      @close="reset()"
    >
      <v-row>
        <v-col cols="12">
          <ui-input
            v-model="currentEdit.name"
            :label="$t('common.Name')"
          />
        </v-col>
        <v-col cols="12">
          <template-variable v-model="currentEdit" />
        </v-col>
        <v-col cols="12">
          <v-btn
            color="success"
            block
            @click="save()"
            v-text="$t('common.Save')"
          />
        </v-col>
      </v-row>
    </ui-overlay>
  </div>
</template>

<script>
const isValidName = name => name !== '' && name.indexOf(' ') === -1

export default {
  props: {
    value: { type: Array, default: () => [] }
  },
  data () {
    return {
      add: false,
      edit: false,
      editIndex: 0,
      currentEdit: {},
      variableTemplate: {
        required: true,
        userEdit: true,
        display: '',
        desc: '',
        type: 'string',
        value: '',
        name: ''
      }
    }
  },
  methods: {
    addVar () {
      const changed = [...this.value]
      changed.push({ ...this.variableTemplate })
      this.$emit('input', changed)
      this.add = true
      this.startEdit(this.value.length)
    },
    startEdit (index) {
      this.editIndex = index
      this.currentEdit = { ...this.value[index] }
      this.edit = true
    },
    remove (index) {
      const changed = [...this.value]
      changed.splice(index, 1)
      this.$emit('input', changed)
    },
    save () {
      if (!isValidName(this.currentEdit.name || '')) {
        this.$toast.error(this.$t('templates.VarNameNoSpaces'))
        return
      }

      const changed = [...this.value]
      changed[this.editIndex] = this.currentEdit
      this.$emit('input', changed)
      this.add = false
      this.reset()
    },
    reset () {
      if (this.add) {
        this.$emit('input', [...this.value].filter(item => item.name !== ''))
      }

      this.add = false
      this.edit = false
      this.editIndex = 0
      this.currentEdit = {}
    }
  }
}
</script>
