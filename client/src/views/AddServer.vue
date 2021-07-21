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
  <v-container>
    <h1 v-text="$t('servers.Add')" />
    <v-stepper v-model="currentStep">
      <v-stepper-header>
        <v-stepper-step
          step="1"
          :complete="currentStep > 1"
        />
        <v-divider />
        <v-stepper-step
          step="2"
          :complete="currentStep > 2"
        />
        <v-divider />
        <v-stepper-step
          step="3"
          :complete="currentStep > 3"
        />
        <v-divider />
        <v-stepper-step
          step="4"
          :complete="currentStep > 4"
        />
      </v-stepper-header>
      <v-stepper-items>
        <v-stepper-content step="1">
          <div v-if="Object.keys(templates).length > 0">
            <h3 v-text="$t('servers.SelectTemplate')" />
            <ui-input
              v-model="templateFilter"
              look="material"
              autofocus
              :placeholder="$t('common.Search')"
            />
          </div>
          <div
            v-else
            class="text-center text--disabled"
            v-text="$t('templates.NoTemplates')"
          />
          <v-expansion-panels>
            <fragment
              v-for="(elements, type) in templates"
              :key="type"
            >
              <v-expansion-panel
                v-if="filterTemplates(elements, templateFilter).length > 0"
                disabled
              >
                <v-expansion-panel-header v-text="type" />
              </v-expansion-panel>
              <v-expansion-panel
                v-for="template in filterTemplates(elements, templateFilter)"
                :key="template.name"
                @click="loadTemplateData(template.name)"
              >
                <v-expansion-panel-header v-text="template.display" />
                <v-expansion-panel-content>
                  <v-row v-if="templateData[template.name] === undefined">
                    <v-col cols="5" />
                    <v-col cols="2">
                      <v-progress-circular
                        indeterminate
                        class="mr-2"
                      />
                      <strong v-text="$t('common.Loading')" />
                    </v-col>
                  </v-row>
                  <!-- eslint-disable vue/no-v-html -->
                  <span
                    v-else
                    v-html="markdown(templateData[template.name].readme || $t('servers.NoReadme'))"
                  />
                  <!-- eslint-enable -->
                  <br>
                  <v-btn
                    color="success"
                    large
                    block
                    @click="selectTemplate(template.name)"
                    v-text="$t('servers.SelectThisTemplate')"
                  />
                </v-expansion-panel-content>
              </v-expansion-panel>
            </fragment>
          </v-expansion-panels>
        </v-stepper-content>

        <v-stepper-content step="2">
          <v-row>
            <v-col cols="12">
              <h3
                class="mb-4"
                v-text="$t('servers.Name')"
              />
              <ui-input
                v-model="serverName"
                autofocus
                @keyup.enter="step2Continue()"
              />
            </v-col>
          </v-row>

          <v-row>
            <v-col cols="12">
              <h3
                class="mb-4"
                v-text="$t('nodes.Node')"
              />
              <ui-select
                v-model="selectedNode"
                :disabled="loadingNodes"
                :items="nodes"
                :no-data-text="$t('errors.ErrNoNodes')"
                :placeholder="$t('servers.SelectNode')"
              />
            </v-col>
          </v-row>

          <v-row>
            <v-col cols="12">
              <h3 v-text="$t('servers.Environment')" />
            </v-col>
            <v-col cols="12">
              <ui-select
                v-model="selectedEnvironment"
                :disabled="loadingTemplates"
                :items="environments"
                :placeholder="$t('servers.SelectEnvironment')"
              />
            </v-col>
            <v-col cols="12">
              <ui-env-config
                v-if="selectedEnvironment && environments[selectedEnvironment]"
                v-model="environments[selectedEnvironment]"
              />
            </v-col>
          </v-row>

          <v-row>
            <v-col>
              <v-btn
                large
                block
                color="error"
                @click="currentStep = 1"
                v-text="$t('common.Back')"
              />
            </v-col>
            <v-col>
              <v-btn
                large
                block
                color="success"
                :disabled="!step2CanContinue()"
                @click="step2Continue()"
                v-text="$t('common.Next')"
              />
            </v-col>
          </v-row>
        </v-stepper-content>

        <v-stepper-content step="3">
          <v-row>
            <v-col cols="12">
              <h3
                class="mb-4"
                v-text="$t('users.Users')"
              />
              <ui-input
                v-model="userInput"
                autofocus
                :placeholder="$t('servers.TypeUsername')"
                @keyup.enter="step3Continue()"
              />
              <v-list v-if="users.length > 0 || selectedUsers.length > 0">
                <v-subheader
                  v-if="users.length > 0"
                  v-text="$t('users.Add')"
                />
                <v-list-item-group v-if="users.length > 0">
                  <v-list-item
                    v-for="user in users"
                    :key="user.value"
                    @click="selectUser(user.value)"
                  >
                    <v-list-item-icon>
                      <v-icon>mdi-plus</v-icon>
                    </v-list-item-icon>
                    <v-list-item-content>
                      <v-list-item-title v-text="user.text" />
                    </v-list-item-content>
                  </v-list-item>
                </v-list-item-group>
                <v-subheader
                  v-if="selectedUsers.length > 0"
                  v-text="$t('users.Users')"
                />
                <v-list-item-group v-if="selectedUsers.length > 0">
                  <v-list-item
                    v-for="user in selectedUsers"
                    :key="user"
                    @click="removeUser(user)"
                  >
                    <v-list-item-icon>
                      <v-icon>mdi-minus</v-icon>
                    </v-list-item-icon>
                    <v-list-item-content>
                      <v-list-item-title v-text="user" />
                    </v-list-item-content>
                  </v-list-item>
                </v-list-item-group>
              </v-list>
            </v-col>
          </v-row>

          <v-row>
            <v-col>
              <v-btn
                large
                block
                color="error"
                @click="currentStep = 2"
                v-text="$t('common.Back')"
              />
            </v-col>
            <v-col>
              <v-btn
                large
                block
                color="success"
                :disabled="!step3CanContinue()"
                @click="step3Continue()"
                v-text="$t('common.Next')"
              />
            </v-col>
          </v-row>
        </v-stepper-content>

        <v-stepper-content step="4">
          <v-row v-if="Object.keys(formData).length > 0">
            <v-col cols="12">
              <v-card-title v-text="$t('common.Options')" />
              <v-row>
                <v-col
                  v-for="(item, name) in filteredFormData"
                  :key="name"
                  cols="12"
                >
                  <ui-variable-input v-model="formData[name]" />
                </v-col>
              </v-row>
            </v-col>
          </v-row>

          <v-row>
            <v-col>
              <v-btn
                large
                block
                color="error"
                @click="currentStep = 3"
                v-text="$t('common.Back')"
              />
            </v-col>
            <v-col>
              <v-btn
                large
                block
                color="success"
                :disabled="!canCreate"
                @click="submitCreate"
                v-text="$t('common.Create')"
              />
            </v-col>
          </v-row>
        </v-stepper-content>
      </v-stepper-items>
    </v-stepper>
  </v-container>
</template>

<script>
import axios from 'axios'
import { Fragment } from 'vue-fragment'
import markdown from '@/utils/markdown'

const CancelToken = axios.CancelToken

export default {
  components: { Fragment },
  data () {
    return {
      nodes: [],
      selectedNode: null,
      templateFilter: '',
      templates: {},
      templateData: {},
      selectedTemplate: '',
      formData: {},

      loadingNodes: true,
      loadingTemplates: true,

      searchingUsers: true,
      users: [],
      selectedUser: null,
      selectedUsers: [],
      userInput: null,
      userCancelSearch: CancelToken.source(),

      serverName: '',

      selectedEnvironment: null,
      environments: [],

      currentStep: 1
    }
  },
  computed: {
    canCreate () {
      if (this.loadingTemplates || this.loadingNodes) {
        return false
      }

      if (!this.selectedTemplate || this.selectedTemplate === '') {
        return false
      }

      if (this.selectedUsers.length === 0) {
        return false
      }

      if (!this.selectedEnvironment) {
        return false
      }

      for (const k in this.templateData[this.selectedTemplate].data) {
        const data = this.templateData[this.selectedTemplate].data[k]
        if (data.type === 'boolean') {
          continue
        }
        if (data.required && !data.value) {
          return false
        }
      }

      return true
    },
    filteredFormData () {
      const remove = Object.keys(this.formData).filter(elem => this.formData[elem].internal)
      const filtered = { ...this.formData }
      remove.map(elem => {
        delete filtered[elem]
      })
      return filtered
    },
    environmentKeys () {
      return Object.keys(this.environments[this.selectedEnvironment]).filter(elem => ['type', 'value', 'text'].indexOf(elem) === -1)
    }
  },
  watch: {
    selectedTemplate (newVal) {
      if (!newVal || newVal === '') {
        return
      }
      this.formData = this.templateData[newVal].vars
      this.environments = []
      for (const k in this.templateData[newVal].supportedEnvs) {
        const env = this.templateData[newVal].supportedEnvs[k]
        this.environments.push({
          value: k,
          text: this.$t('env.' + env.type + '.name'),
          ...env
        })
      }

      const env = this.templateData[newVal].defaultEnv
      let def = null
      if (env && env.type) {
        def = env.type
      }

      if (def) {
        for (const k in this.environments) {
          if (this.environments[k].type === def) {
            this.selectedEnvironment = k
            break
          }
        }
      } else {
        this.selectedEnvironment = null
      }

      for (const key in this.formData) {
        if (this.formData[key].type === 'boolean') {
          this.formData[key].value = this.formData[key].value === 'true'
        }
      }
    },
    userInput (newVal) {
      if (!newVal || newVal === '') {
        this.users = []
      } else {
        // eslint-disable-next-line no-new
        new Promise((resolve, reject) => {
          this.findUsers(newVal)
          resolve()
        })
      }
    }
  },
  async mounted () {
    this.nodes = [{
      value: null,
      disabled: true,
      text: this.$t('common.Loading')
    }]
    this.getTemplates()
    this.getNodes()
    const self = await this.$api.getSelf()
    this.selectedUsers.push(self.username)
  },
  methods: {
    async getTemplates () {
      this.loadingTemplates = true
      this.templateData = {}
      this.selectedTemplate = null
      const templates = await this.$api.getTemplates()
      templates.map(template => {
        if (!template.display) template.display = template.name
        if (!template.type) template.type = 'none'
        if (!this.templates[template.type]) this.templates[template.type] = []
        this.templates[template.type].push(template)
      })

      const keys = Object.keys(this.templates)
      const index = keys.indexOf('other')
      if (index !== -1) this.$delete(keys, index)
      keys.map(key => {
        if (this.templates[key].length === 1) {
          if (!this.templates.other) this.templates.other = []
          this.templates.other.push(this.templates[key][0])
          delete this.templates[key]
        }
      })

      this.templates = { ...this.templates }
      this.loadingTemplates = false
    },
    async getNodes () {
      const nodes = await this.$api.getNodes()
      this.nodes = nodes.map(node => {
        return { value: node.id, text: node.name }
      })

      if (this.nodes.length > 0) {
        this.selectedNode = this.nodes[0].value
      }

      this.loadingNodes = false
    },
    async findUsers () {
      this.searchingUsers = true
      this.userCancelSearch.cancel()
      this.userCancelSearch = CancelToken.source()
      const users = await this.$api.searchUsers(this.userInput, this.userCancelSearch.token)
      this.users = users.map(user => {
        return { value: user.username, text: `${user.username} <${user.email}>` }
      }).sort()
      this.searchingUsers = false
    },
    async submitCreate () {
      const data = this.$api.templateToApiJson(this.templateData[this.selectedTemplate], false)
      for (const item in data.data) {
        switch (data.data[item].type) {
          case 'integer':
            data.data[item].value = Number(data.data[item].value)
            break
          case 'boolean':
            data.data[item].value = Boolean(data.data[item].value)
            break
        }
      }
      data.node = this.selectedNode
      data.users = this.selectedUsers
      data.name = this.serverName !== '' ? this.serverName : undefined
      data.environment = {
        ...this.environments[this.selectedEnvironment]
      }
      delete data.environment.text
      delete data.environment.value
      const id = await this.$api.createServer(data)
      this.$router.push({ name: 'Server', params: { id: id } })
    },
    selectUser (username) {
      if (!username || username === '') {
        return
      }
      for (let i = 0; i < this.selectedUsers.length; i++) {
        if (this.selectedUsers[i] === username) {
          return
        }
      }
      this.userInput = null
      this.selectedUsers.push(username)
      this.selectedUsers.sort()
      this.selectedUser = null
      this.users = []
    },
    removeUser (username) {
      for (let i = 0; i < this.selectedUsers.length; i++) {
        if (this.selectedUsers[i] === username) {
          this.selectedUsers.splice(i, 1)
          return
        }
      }
    },
    async loadTemplateData (template) {
      if (!template) return
      if (!this.templateData[template]) {
        this.templateData[template] = await this.$api.getTemplate(template)
        this.templateData = { ...this.templateData }
      }
    },
    selectTemplate (template) {
      this.loadTemplateData(template)
      this.selectedTemplate = template
      this.currentStep = 2
    },
    filterTemplates (templates, filter) {
      return templates.filter(t => {
        if (filter.trim() === '') {
          return true
        } else {
          let name = t.display
          if (!name) {
            name = t.name
          }
          return name.toLowerCase().indexOf(filter.trim().toLowerCase()) > -1
        }
      })
    },
    step2CanContinue () {
      return this.selectedNode && this.selectedEnvironment && this.serverName && this.servername !== ''
    },
    step2Continue () {
      if (!this.step2CanContinue()) return
      this.currentStep = 3
    },
    step3CanContinue () {
      return this.selectedUsers.length > 0
    },
    step3Continue () {
      if (!this.step3CanContinue()) return
      this.currentStep = 4
    },
    markdown
  }
}
</script>
