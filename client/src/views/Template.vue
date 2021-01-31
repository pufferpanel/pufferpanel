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
                value="editor"
                v-text="$t('templates.Editor')"
              />
              <v-btn
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
      <template-editor v-if="mode === 'editor'" v-model="template" />
      <v-row v-else>
        <v-col>
          <ace
            v-model="templateJson"
            :editor-id="template.name"
            height="50vh"
            lang="json"
            @editorready="$refs.editor.setValue(template)"
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
      }
    }
  },
  watch: {
    mode (newVal) {
      if (newVal === 'editor') {
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
      if (!this.template.name || this.template.name.trim() === '') return
      if (this.mode === 'json') this.template = this.$api.templateFromApiJson(this.templateJson)
      await this.$api.saveTemplate(this.template.name, this.template)
      this.$toast.success(this.$t('templates.SaveSuccess'))
      if (this.create) this.$router.push({ name: 'Template', params: { id: this.name } })
    }
  }
}
</script>
