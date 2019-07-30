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
        <label :for="`nodeSelect`" v-text="$t('common.Node')"></label>
      </b-col>
      <b-col cols="10">
        <b-form-select id="nodeSelect" v-model="selectedNode" :options="nodes"></b-form-select>
      </b-col>
    </b-row>

    <b-row>
      <b-col cols="2">
        <label :for="`templateSelect`" v-text="$t('common.Template')"></label>
      </b-col>
      <b-col cols="10">
        <b-form-select id="templateSelect" v-model="selectedTemplate" :options="templates"></b-form-select>
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
      selectedTemplate: null,
      formData: {},
      readme: ''
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
    }
  },
  mounted () {
    let vue = this
    this.$http.get('/api/nodes').then(function (res) {
      let callResult = res.data
      if (callResult.success) {
        for (let i = 0; i < callResult.data.length; i++) {
          let node = callResult.data[i]
          vue.nodes.push({
            value: node.id,
            text: node.name
          })
        }

        if (vue.nodes.length === 1) {
          vue.selectedNode = vue.nodes[0].value
        }
      }
    })
  }
}
</script>

<style scoped>

</style>