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
  <v-card>
    <v-card-title>
      <span v-text="$t('servers.Tasks')" />
    </v-card-title>
    <v-card-text>
      <v-row>
        <v-col v-if="tasks.length === 0">
          <span v-text="$t('servers.NoTasks')" />
        </v-col>
        <v-col>
          <v-list>
            <v-list-item
              v-for="(task, id) in tasks"
              :key="id"
              @click="edit(id)"
            >
              <v-list-item-content>
                <v-list-item-title>{{ task.name }}</v-list-item-title>
                <v-list-item-subtitle>{{ describe(task.cronSchedule) }}</v-list-item-subtitle>
              </v-list-item-content>
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
                      @click.stop="trigger(id)"
                    >
                      <v-icon>mdi-play</v-icon>
                    </v-btn>
                  </template>
                  <span v-text="$t('servers.RunTask')" />
                </v-tooltip>
                <v-tooltip bottom>
                  <template v-slot:activator="{ on }">
                    <v-btn
                      icon
                      v-on="on"
                      @click.stop="remove(id)"
                    >
                      <v-icon>mdi-close</v-icon>
                    </v-btn>
                  </template>
                  <span v-text="$t('common.Delete')" />
                </v-tooltip>
              </v-list-item-action>
            </v-list-item>
          </v-list>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-btn
            block
            text
            @click="add = true"
            v-text="$t('common.Add')"
          />
        </v-col>
      </v-row>

      <ui-overlay
        v-model="add"
        closable
        card
        :title="$t('servers.NewTask')"
        @close="reset"
      >
        <v-row>
          <v-col>
            <ui-input
              v-model="newTask.name"
              :label="$t('common.Name')"
              autofocus
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <ui-input
              v-model="newTask.cronSchedule"
              :label="$t('servers.Schedule')"
              :hint="describe(newTask.cronSchedule)"
              :error-messages="!isValidSchedule(newTask.cronSchedule) ? $t('servers.TaskScheduleInvalid') : undefined"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <template-processors v-model="newTask.operations" />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-btn
              block
              color="success"
              :disabled="!canSave('new')"
              @click="saveNew"
              v-text="$t('common.Save')"
            />
          </v-col>
        </v-row>
      </ui-overlay>
      <ui-overlay
        :value="editId !== false"
        closable
        card
        :title="$t('servers.EditTask')"
        @close="reset"
      >
        <v-row>
          <v-col>
            <ui-input
              v-model="editTask.name"
              :label="$t('common.Name')"
              autofocus
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <ui-input
              v-model="editTask.cronSchedule"
              :label="$t('servers.Schedule')"
              :hint="describe(editTask.cronSchedule)"
              :error-messages="!isValidSchedule(editTask.cronSchedule) ? $t('servers.TaskScheduleInvalid') : undefined"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <template-processors v-model="editTask.operations" />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-btn
              block
              color="success"
              :disablef="!canSave('edit')"
              @click="saveEdit"
              v-text="$t('common.Save')"
            />
          </v-col>
        </v-row>
      </ui-overlay>
    </v-card-text>
  </v-card>
</template>

<script>
import cronstrue from 'cronstrue/i18n'
import { isValidCron } from 'cron-validator'

export default {
  props: {
    server: { type: Object, default: () => {} }
  },
  data () {
    return {
      add: false,
      editId: false,
      tasks: {},
      newTask: {
        name: '',
        cronSchedule: '',
        operations: []
      },
      editTask: {
        name: '',
        cronSchedule: '',
        operations: []
      }
    }
  },
  mounted () {
    this.loadData()
  },
  methods: {
    async loadData () {
      this.tasks = await this.$api.getServerTasks(this.server.id) || {}
    },
    async save () {
      await this.$api.updateServerTasks(this.server.id, this.tasks)
      this.$toast.success(this.$t('common.Saved'))
      this.reset()
      this.loadData()
    },
    async saveNew () {
      const id = await this.$api.createServerTask(this.server.id, this.newTask)
      this.$toast.success(this.$t('common.Saved'))
      this.tasks[id] = this.newTask
      this.reset()
      this.loadData()
    },
    edit (id) {
      this.editTask = { ...this.tasks[id] }
      this.editId = id
    },
    async saveEdit () {
      await this.$api.editServerTask(this.server.id, this.editId, this.editTask)
      this.$toast.success(this.$t('common.Saved'))
      this.tasks[this.editId] = this.editTask
      this.reset()
      this.loadData()
    },
    async remove (id) {
      await this.$api.deleteServerTask(this.server.id, id)
      this.$toast.success(this.$t('servers.TaskDeleted'))
      delete this.tasks[id]
      this.loadData()
    },
    async trigger (id) {
      await this.$api.runServerTask(this.server.id, id)
      this.$toast.success(this.$t('servers.TaskTriggered'))
    },
    reset () {
      this.newTask = { name: '', cronSchedule: '', operations: [] }
      this.add = false
      this.editTask = { name: '', cronSchedule: '', operations: [] }
      this.editId = false
    },
    describe (schedule) {
      if (!schedule || schedule === '') return this.$t('servers.TaskScheduleManual')
      if (!this.isValidSchedule(schedule)) return

      let locale = 'en'
      if (this.$i18n.locale === 'zh_TW' || this.$i18n.locale === 'zh_HK') {
        locale = 'zh_TW'
      } else if (this.$i18n.locale === 'zh_CN') {
        locale = 'zh_CN'
      } else {
        locale = this.$i18n.locale.split('_')[0] || 'en'
      }
      return cronstrue.toString(schedule, { verbose: true, locale })
    },
    canSave (mode) {
      const subject = mode === 'edit' ? this.editTask : this.newTask
      if (!subject.name || subject.name.trim() === '') return false
      if (!this.isValidSchedule(subject.cronSchedule)) return false
      return true
    },
    isValidSchedule (schedule) {
      if (!schedule || schedule === '') return true
      return isValidCron(schedule, { alias: true, allowBlankDay: true })
    }
  }
}
</script>
