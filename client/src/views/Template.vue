<template>
  <div>
    <v-row>
      <v-col>
        <v-card>
          <v-card-title>
            <span v-text="$t(create ? 'templates.New' : 'templates.Edit')" />
            <div class="flex-grow-1" />
            <v-btn-toggle
              v-model="mode"
              borderless
              dense
              mandatory
            >
              <v-btn
                :disabled="loading"
                value="editor"
                v-text="$t('templates.Editor')"
              />
              <v-btn
                :disabled="loading"
                value="json"
                v-text="$t('templates.Json')"
              />
            </v-btn-toggle>
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col>
                <ui-input
                  v-model="template.name"
                  :label="$t('common.Name')"
                  :disabled="!create"
                />
              </v-col>
            </v-row>
            <v-row v-if="mode === 'editor'">
              <v-col
                cols="12"
                md="6"
              >
                <ui-input
                  v-model="template.display"
                  :label="$t('templates.Display')"
                />
              </v-col>
              <v-col
                cols="12"
                md="6"
              >
                <ui-input
                  v-model="template.type"
                  :label="$t('templates.Type')"
                />
              </v-col>
            </v-row>
            <v-row v-if="loading">
              <v-col
                offset="5"
                cols="2"
              >
                <v-progress-circular
                  indeterminate
                  class="mr-2"
                />
                <strong v-text="$t('common.Loading')" />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <div v-if="!loading">
      <template-editor
        v-if="mode === 'editor'"
        v-model="template"
      />
      <v-row v-else>
        <v-col>
          <ace
            ref="editor"
            v-model="templateJson"
            :editor-id="'edit' + template.name"
            height="50vh"
            lang="json"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-btn
            color="success"
            large
            block
            :disabled="!template.name || template.name.trim() === ''"
            @click="save"
            v-text="$t('common.Save')"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-btn
            v-if="!create"
            color="error"
            block
            @click="remove()"
            v-text="$t('common.Delete')"
          />
        </v-col>
      </v-row>
    </div>

    <ui-overlay
      v-model="offerV1Convert"
      card
    >
      <v-row class="mb-8">
        <v-col cols="12">
          <h2 v-text="$t('templates.OfferV1Convert')" />
        </v-col>
      </v-row>
      <v-row>
        <v-col
          cols="12"
          md="4"
        >
          <v-btn
            block
            color="error"
            @click="offerV1Convert = false"
            v-text="$t('common.Cancel')"
          />
        </v-col>
        <v-col
          cols="12"
          md="4"
        >
          <v-btn
            block
            @click="offerV1Convert = false; skipNextV1Check = true"
            v-text="$t('common.Ignore')"
          />
        </v-col>
        <v-col
          cols="12"
          md="4"
        >
          <v-btn
            block
            color="success"
            @click="convertV1()"
            v-text="$t('templates.Convert')"
          />
        </v-col>
      </v-row>
    </ui-overlay>
  </div>
</template>

<script>
export default {
  data () {
    return {
      loading: false,
      create: this.$route.params.id === undefined,
      mode: 'editor',
      templateJson: '',
      template: {
        name: '',
        type: '',
        display: '',
        command: '',
        stop: {},
        vars: [],
        install: [],
        pre: [],
        post: [],
        envVars: {},
        defaultEnv: {},
        supportedEnvs: []
      },
      offerV1Convert: false,
      skipNextV1Check: false
    }
  },
  watch: {
    async mode (newVal) {
      if (newVal === 'editor') {
        if (this.catchV1()) return
        this.template = this.$api.templateFromApiJson(this.templateJson)
      } else {
        this.templateJson = this.$api.templateToApiJson(this.template)
        if (this.$refs.editor && this.$refs.editor.ready) this.$refs.editor.setValue(this.templateJson)
      }
    }
  },
  mounted () {
    if (!this.create) this.loadData()
  },
  methods: {
    async loadData () {
      this.loading = true
      this.template = await this.$api.getTemplate(this.$route.params.id)
      this.loading = false
    },
    async remove () {
      if (this.create) return
      await this.$api.deleteTemplate(this.template.name)
      this.$toast.success(this.$t('templates.Deleted'))
      this.$router.push({ name: 'Templates' })
    },
    async save () {
      // console.log('1', this.$api.templateToApiJson(this.template))
      if (!this.template.name || this.template.name.trim() === '') return
      if (this.mode === 'json') {
        if (this.catchV1()) return
        this.template = this.$api.templateFromApiJson(this.templateJson)
      }
      // console.log('2', this.$api.templateToApiJson(this.template))
      await this.$api.saveTemplate(this.template.name, this.template)
      // console.log('3', this.$api.templateToApiJson(this.template))
      this.$toast.success(this.$t('templates.SaveSuccess'))
      if (this.create) this.$router.push({ name: 'Template', params: { id: this.name } })
      // console.log('4', this.$api.templateToApiJson(this.template))
    },
    catchV1 () {
      const skipCheck = this.skipNextV1Check
      this.skipNextV1Check = false
      if (Object.prototype.toString.call((JSON.parse(this.templateJson) || {}).pufferd) === '[object Object]' && !skipCheck) {
        this.offerV1Convert = true
        return true
      }
      return false
    },
    convertV1 () {
      const template = JSON.parse(this.templateJson).pufferd
      template.install = template.install.commands
      template.run.command = template.run.program + ' ' + template.run.arguments.join(' ')
      delete template.run.program
      delete template.run.arguments
      template.supportedEnvironments = [template.environment]
      template.name = template.display.toLowerCase().replace(/ /g, '-').replace(/[^a-z0-9]/g, '-').replace(/-+/g, '-')
      while (template.name.substring(0, 1) === '-') {
        template.name = template.name.substring(1)
      }
      while (template.name.substring(template.name.length - 1) === '-') {
        template.name = template.name.substring(0, template.name.length - 1)
      }
      // eslint-disable-next-line no-template-curly-in-string
      this.templateJson = JSON.stringify(template, undefined, 2).replace(/\$\{rootdir}/g, '${rootDir}')
      if (this.$refs.editor && this.$refs.editor.ready) this.$refs.editor.setValue(this.templateJson)
      this.template = this.$api.templateFromApiJson(this.templateJson)
      this.offerV1Convert = false
    }
  }
}
</script>
