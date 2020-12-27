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
        v-for="(entry, i) in value"
        :key="i"
        link
        @click="startEdit(i)"
      >
        <v-list-item-content v-text="entry.type" />
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
                @click.stop="remove(i)"
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
      @click="startAdd()"
      v-text="$t('common.Add')"
    />
    <ui-overlay
      v-model="edit"
      card
      closable
      @close="reset()"
    >
      <template-processor v-model="value[editIndex]" />
    </ui-overlay>
  </div>
</template>

<script>
export default {
  props: {
    value: { type: Array, default: () => [] }
  },
  data () {
    return {
      template: { type: 'command' },
      edit: false,
      editIndex: 0
    }
  },
  methods: {
    startAdd () {
      const changed = [...this.value]
      const newIndex = changed.length
      changed.push(this.template)
      this.$emit('input', changed)
      this.startEdit(newIndex)
    },
    startEdit (index) {
      this.editIndex = index
      this.editModel = this.value[index]
      this.add = false
      this.edit = true
    },
    remove (index) {
      const changed = [...this.value]
      changed.splice(index, 1)
      this.$emit('input', changed)
    },
    reset () {
      this.edit = false
      this.add = false
      this.editIndex = 0
      this.editModel = {}
    },
    save () {
      const changed = [...this.value]
      if (this.add) {
        changed.push(this.editModel)
      } else {
        changed[this.editIndex] = this.editModel
      }
      this.$emit('input', changed)
      this.reset()
    }
  }
}
</script>
