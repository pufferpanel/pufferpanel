<template>
  <v-row>
    <v-col>
      <v-card>
        <v-card-title>
          <span v-text="$t(create ? 'templates.New' : 'templates.Edit')" class="display-1off" />
          <div class="flex-grow-1" />
          <v-btn-toggle v-model="currentMode" borderless dense mandatory>
            <v-btn value="editor" v-text="$t('templates.Editor')" @click="updateEditor()" />
            <v-btn value="json" v-text="$t('templates.Json')" @click="updateJson()" />
          </v-btn-toggle>
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-col>
              <v-text-field
                v-model="name"
                :label="$t('common.Name')"
                :disabled="!create"
                outlined
                hide-details
              />
            </v-col>
          </v-row>
          <v-row v-if="loading">
            <v-col cols="5" />
            <v-col cols="2">
              <v-progress-circular
                indeterminate
                class="mr-2"
              />
              <strong v-text="$t('common.Loading')" />
            </v-col>
          </v-row>
          <v-row v-if="currentMode === 'json' && !loading">
            <v-col>
              <ace
                v-model="template"
                :editor-id="name + 'json'"
                :theme="isDark() ? 'monokai' : 'github'"
                height="50vh"
                ref="editor"
                file="template.json"
                @editorready="$refs.editor.setValue(template)"
              />
            </v-col>
          </v-row>
          <div v-if="currentMode === 'editor' && !loading">
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="templateObj.display"
                  :label="$t('templates.Display')"
                  outlined
                  hide-details
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="templateObj.type"
                  :label="$t('templates.Type')"
                  outlined
                  hide-details
                />
              </v-col>
            </v-row>
            <p v-text="$t('templates.Variables')" class="title" />
            <template-variables v-model="templateObj.data" />
            <p v-text="$t('templates.Install')" class="title" />
            <template-processors v-model="templateObj.install" name="install" />
            <p v-text="$t('templates.RunConfig')" class="title" />
            <template-run v-model="templateObj.run" />
            <p v-text="$t('templates.SupportedEnvironments')" class="title" />
            <template-environments v-model="templateObj.supportedEnvironments" />
          </div>
          <v-row>
            <v-col>
              <v-btn
                color="success"
                large
                block
                @click="save"
                v-text="$t('common.Save')"
              />
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script>
import { isDark } from '@/utils/dark'
import { handleError } from '@/utils/api'

export default {
  data () {
    return {
      currentMode: 'editor',
      loading: false,
      create: this.$route.params.id === undefined,
      name: this.$route.params.id === undefined ? '' : this.$route.params.id,
      template: '',
      templateObj: {
        data: {},
        run: {
          arguments: [],
          environmentVars: {}
        },
        display: '',
        type: '',
        install: [],
        supportedEnvironments: []
      }
    }
  },
  mounted () {
    if (!this.create) this.loadData()
  },
  methods: {
    log (x, y, z) { console.log(x, y, z) },
    loadData () {
      this.loading = true
      const ctx = this
      ctx.$http.get(`/api/templates/${ctx.$route.params.id}`).then(response => {
        const data = response.data
        data.readme = undefined
        ctx.template = JSON.stringify(data, undefined, 4)
        ctx.templateObj = data
        if (data.run && data.run.stop) {
          ctx.stopType = 'command'
        }
        if (data.run && data.run.stopCode) {
          ctx.stopType = 'signal'
        }
        if (!this.templateObj.data) this.templateObj.data = {}
        Object.keys(this.templateObj.data).forEach(key => {
          if (!this.templateObj.data[key].type) this.templateObj.data[key].type = 'string'
        })
        if (!this.templateObj.run) this.templateObj.run = {}
        if (!this.templateObj.run.environmentVars) this.templateObj.run.environmentVars = {}
        if (!this.templateObj.run.arguments) this.templateObj.run.arguments = []
        if (!this.templateObj.run.pre) this.templateObj.run.pre = []
        if (!this.templateObj.run.post) this.templateObj.run.post = []
        if (!this.templateObj.supportedEnvironments) this.templateObj.supportedEnvironments = []
        if (!this.templateObj.install) this.templateObj.install = []
        this.templateObj.install.map(element => {
          if (element.type === 'download' && typeof element.files === 'string') element.files = [element.files]
          if (element.type === 'command' && typeof element.commands === 'string') element.commands = [element.commands]
          return element
        })
        if (ctx.$refs.editor && ctx.$refs.editor.ready) ctx.$refs.editor.setValue(ctx.template)
        ctx.loading = false
      }).catch(handleError(ctx))
    },
    save () {
      this.updateJson()
      const ctx = this
      ctx.$http.put(`/api/templates/${ctx.name}`, ctx.template).then(response => {
        ctx.$toast.success(ctx.$t('templates.SaveSuccess'))
        if (ctx.create) ctx.$router.push({ name: 'Template', params: { id: ctx.name } })
      }).catch(handleError(ctx, { 400: 'errors.ErrInvalidJson' }))
    },
    updateEditor () {
      const data = JSON.parse(this.template)
      if (data.run && data.run.stop) {
        this.stopType = 'command'
      }
      if (data.run && data.run.stopCode) {
        this.stopType = 'signal'
      }
      this.templateObj = JSON.parse(this.template)
    },
    updateJson () {
      this.templateObj.name = this.name
      this.templateObj.run.stopCode = this.templateObj.run.stopCode * 1
      this.templateObj.run[this.stopType === 'command' ? 'stopCode' : 'stop'] = undefined
      this.template = JSON.stringify(this.templateObj, undefined, 4)
      if (this.$refs.editor && this.$refs.editor.ready) this.$refs.editor.setValue(this.template)
    },
    isDark
  }
}
</script>
