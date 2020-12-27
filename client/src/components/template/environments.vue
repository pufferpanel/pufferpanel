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
        v-for="(env, i) in value"
        :key="i"
        @click="startEdit(i)"
      >
        <v-list-item-content v-text="env.type" />
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
      v-if="availableEnvs.length > 0"
      text
      block
      @click="startAdd()"
      v-text="$t('common.Add')"
    />
    <ui-overlay
      v-model="edit"
      :title="add ? $t('common.Add') : currentEnv.type"
      card
      closable
      @close="reset()"
    >
      <v-row v-if="add">
        <v-col cols="12">
          <ui-select
            v-model="currentEnv.type"
            :items="availableEnvs"
            :label="$t('templates.Environment')"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12">
          <ui-env-config
            v-model="currentEnv"
            :no-fields-text="$t('env.NoEnvFields')"
          />
        </v-col>
      </v-row>
      <v-row>
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
const envs = ['standard', 'tty', 'docker']

export default {
  props: {
    value: { type: Array, default: () => [] }
  },
  data () {
    return {
      add: false,
      edit: false,
      editIndex: 0,
      currentEnv: { type: 'standard' }
    }
  },
  computed: {
    availableEnvs () {
      return envs.filter(elem => {
        return this.value.map(v => { return v.type === elem }).indexOf(true) === -1
      })
    }
  },
  methods: {
    addEnv () {
      this.value.push({ type: this.newEnv })
      this.newEnv = 'standard'
      this.addingEnv = false
    },

    remove (i) {
      const changed = [...this.value]
      changed.splice(i, 1)
      this.$emit('input', changed)
    },
    reset () {
      this.add = false
      this.edit = false
      this.editIndex = 0
      this.currentEnv = { type: this.availableEnvs[0] }
    },
    startAdd () {
      this.currentEnv = { type: this.availableEnvs[0] }
      this.add = true
      this.edit = true
    },
    startEdit (i) {
      this.currentEnv = { ...this.value[i] }
      this.editIndex = i
      this.add = false
      this.edit = true
    },
    save () {
      const changed = [...this.value]
      if (this.add) {
        changed.push(this.currentEnv)
      } else {
        changed[this.editIndex] = this.currentEnv
      }
      this.$emit('input', changed)
      this.reset()
    }
  }
}
</script>
