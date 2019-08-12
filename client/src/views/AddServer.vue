<!--
  - Copyright 2019 Padduck, LLC
  -  Licensed under the Apache License, Version 2.0 (the "License");
  -  you may not use this file except in compliance with the License.
  -  You may obtain a copy of the License at
  -  	http://www.apache.org/licenses/LICENSE-2.0
  -  Unless required by applicable law or agreed to in writing, software
  -  distributed under the License is distributed on an "AS IS" BASIS,
  -  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  -  See the License for the specific language governing permissions and
  -  limitations under the License.
  -->

<template>
  <b-container>
    <b-card header-tag="header" footer-tag="footer">
      <h6 slot="header" class="mb-0"><span v-text="$t('common.AddServer')"></span></h6>
      <b-btn slot="footer" variant="primary" :disabled="!canCreate" v-on:click="submitCreate"
             v-text="$t('common.Create')"></b-btn>
      <b-card-text>
        <b-row>
          <b-col sm="12" md="10">
            <b-card header-tag="header">
              <h6 slot="header" class="mb-0"><span v-text="$t('common.Node')"></span></h6>
              <b-card-text>
                <b-form-select :disabled="loadingNodes" id="nodeSelect" v-model="selectedNode"
                               :options="nodes"></b-form-select>
              </b-card-text>
            </b-card>
          </b-col>
        </b-row>

        <b-row>
          <b-col sm="12" md="10">
            <b-card header-tag="header">
              <h6 slot="header" class="mb-0"><span v-text="$t('common.Template')"></span></h6>
              <b-card-text>
                <b-form-select :disabled="loadingTemplates" id="templateSelect" v-model="selectedTemplate"
                               :options="templates"></b-form-select>
              </b-card-text>
            </b-card>
          </b-col>
        </b-row>

        <b-row>
          <b-col sm="12" md="10">
            <b-card header-tag="header">
              <h6 slot="header" class="mb-0"><span v-text="$t('common.Users')"></span></h6>
              <b-card-text>
              </b-card-text>
            </b-card>
          </b-col>
        </b-row>

        <b-row>
          <b-col sm="12" md="10">
            <b-card header-tag="header">
              <h6 slot="header" class="mb-0"><span v-text="$t('common.Options')"></span></h6>
              <b-card-text>
                <div v-if="selectedTemplate" v-for="item in formData">
                  <b-card header-tag="header" v-if="!item.internal">
                    <h6 slot="header" v-text="item.display"></h6>
                    <b-card-text>
                      <span v-html="item.desc"></span>
                      <b-form-input v-if="item.type === 'integer'" type="number" v-model="item.value"
                                    :required="item.required"></b-form-input>
                      <b-form-checkbox v-else-if="item.type === 'boolean'" v-model="item.value"
                                       :required="item.required"></b-form-checkbox>
                      <b-form-select v-else-if="item.type === 'option'" :options="item.options" v-model="item.value">
                      </b-form-select>
                      <b-form-input v-else v-model="item.value" :required="item.required"></b-form-input>
                    </b-card-text>
                  </b-card>
                </div>
              </b-card-text>
            </b-card>
          </b-col>
        </b-row>
      </b-card-text>
    </b-card>
  </b-container>
</template>

<script>
export default {
  data () {
    return {
      nodes: [],
      selectedNode: null,
      templates: [],
      templateData: {},
      selectedTemplate: null,
      formData: {},
      readme: '',

      loadingNodes: true,
      loadingTemplates: true
    }
  },
  watch: {
    selectedTemplate: function (newVal, oldVal) {
      this.formData = this.templateData[newVal].data
    }
  },
  computed: {
    canCreate: function () {
      if (this.loadingTemplates || this.loadingNodes) {
        return false
      }

      if (!this.selectedTemplate) {
        return false
      }

      for (let k in this.templateData[this.selectedTemplate].data) {
        let data = this.templateData[this.selectedTemplate].data[k]
        if (data.type === 'boolean') {
          continue
        }
        if (data.required && !data.value) {
          return false
        }
      }

      return true
    }
  },
  methods: {
    getTemplates () {
      let vue = this
      this.loadingTemplates = true
      this.templates = [{
        value: null,
        disabled: true,
        text: this.$t('common.Loading')
      }]
      this.templateData = {}
      this.selectedTemplate = null
      this.$http.get('/api/templates').then(function (res) {
        let response = res.data
        if (response.success) {
          vue.templates = [{
            value: null,
            disabled: true,
            text: vue.$t('common.SelectTemplate')
          }]
          vue.templateData = response.data
          for (let k in vue.templateData) {
            vue.templates.push({
              text: vue.templateData[k].display,
              value: k
            })
          }

          if (vue.templates.length === 2) {
            vue.selectedTemplate = vue.templates[1].value
          } else {
          }

          vue.loadingTemplates = false
        }
      })
    },
    getNodes () {
      let vue = this
      this.$http.get('/api/nodes').then(function (res) {
        let callResult = res.data
        if (callResult.success) {
          for (let i = 0; i < callResult.data.length; i++) {
            let node = callResult.data[i]
            vue.nodes = [{
              value: null,
              disabled: true,
              text: vue.$t('common.SelectNode')
            }]
            vue.nodes.push({
              value: node.id,
              text: node.name
            })
          }

          if (vue.nodes.length === 2) {
            vue.selectedNode = vue.nodes[1].value
          }

          vue.loadingNodes = false
        }
      })
    },
    submitCreate () {
      let vue = this
      let data = this.templateData[this.selectedTemplate]
      data.node = this.selectedNode
      this.$http.post('/api/servers', data).then(function (response) {
        console.log(response)
      }).catch(function (response) {
        console.log(response)
      })
    }
  },
  mounted () {
    this.nodes = [{
      value: null,
      disabled: true,
      text: this.$t('common.Loading')
    }]
    this.getTemplates()
    this.getNodes()
  }
}
</script>
