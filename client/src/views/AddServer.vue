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
    <b-row>
      <b-col cols="2">
        <label :for="`nodeSelect`"><strong v-text="$t('common.Node')"></strong></label>
      </b-col>
      <b-col cols="10">
        <b-form-select :disabled="loadingNodes" id="nodeSelect" v-model="selectedNode" :options="nodes"></b-form-select>
      </b-col>
    </b-row>

    <b-row>
      <b-col cols="2">
        <label :for="`templateSelect`"><strong v-text="$t('common.Template')"></strong></label>
      </b-col>
      <b-col cols="10">
        <b-form-select :disabled="loadingTemplates" id="templateSelect" v-model="selectedTemplate"
                       :options="templates"></b-form-select>
      </b-col>
    </b-row>
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
    selectedNode: function (newVal, oldVal) {
      this.getTemplates()
    }
  },
  methods: {
    getTemplates () {
      console.log('Getting templates for node ' + this.selectedNode)
      let vue = this
      this.loadingTemplates = true
      this.templates = [{
        value: null,
        disabled: true,
        text: this.$t('common.Loading')
      }]
      this.templateData = {}
      this.selectedTemplate = null
      this.$http.get('/daemon/node/' + this.selectedNode + '/templates').then(function (res) {
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
    }
  },
  mounted () {
    let vue = this
    this.nodes = [{
      value: null,
      disabled: true,
      text: this.$t('common.Loading')
    }]
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
  }
}
</script>

<style scoped>

</style>