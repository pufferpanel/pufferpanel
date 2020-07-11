<template>
  <v-container>
    <div class="d-flex">
      <h1 class="flex-grow-1" v-text="$t('templates.Templates')" />
      <v-tooltip left>
        <template v-slot:activator="{ on }">
          <v-btn icon v-on="on" @click="loadTemplateImporter"><v-icon>mdi-import</v-icon></v-btn>
        </template>
        <span v-text="$t('templates.ImportDesc')" />
      </v-tooltip>
    </div>
    <v-row>
      <v-col>
        <v-list subheader elevation="1">
          <div v-for="(elements, type, i) in templates" :key="type">
            <v-subheader v-text="type" :class="isDark() ? headerClasses.dark : headerClasses.light" />
            <div v-for="(template, index) in templates[type]" :key="template.name">
              <v-list-item :to="(hasScope('templates.edit') || isAdmin()) ? {name: 'Template', params: {id: template.name}} : undefined">
                <v-list-item-content>
                  <v-list-item-title v-text="template.display" />
                </v-list-item-content>
              </v-list-item>
            </div>
            <v-divider v-if="i !== Object.keys(templates).length - 1" />
          </div>
          <div class="pt-2 text-center text--disabled" v-if="Object.keys(templates).length === 0" v-text="$t('templates.NoTemplates')" />
        </v-list>
        <v-btn
          v-show="hasScope('templates.edit') || isAdmin()"
          color="primary"
          bottom
          right
          fixed
          fab
          dark
          large
          :to="{name: 'AddTemplate'}"
        >
          <v-icon>mdi-plus</v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <common-overlay v-model="showTemplateImporter" card closable :title="$t('templates.Import')">
      <v-alert border="bottom" text type="warning" prominent>
        {{ $t('templates.OverrideWarning') }}
      </v-alert>
      <v-list flat>
        <v-list-item-group multiple>
          <v-list-item v-for="key in importableTemplates" :key="key">
            <v-list-item-action>
              <v-checkbox
                v-model="templatesToImport[key]"
              ></v-checkbox>
            </v-list-item-action>
            <v-list-item-content>
              <v-list-item-title v-text="key">
            </v-list-item-content>
          </v-list-item>
        </v-list-item-group>
      </v-list>
      <v-row>
        <v-col>
          <v-btn block color="error" v-text="$t('common.Cancel')" @click="showTemplateImporter = false" />
        </v-col>
        <v-col>
          <v-btn block color="success" v-text="$t('templates.Import')" @click="doImports()" />
        </v-col>
      </v-row>
    </common-overlay>
  </v-container>
</template>

<script>
import { isDark } from '@/utils/dark'
import { handleError } from '@/utils/api'

export default {
  data () {
    return {
      loading: true,
      templates: {},
      headerClasses: {
        light: 'body-1 grey font-weight-bold lighten-4 black--text',
        dark: 'body-1 grey font-weight-bold darken-4 white--text'
      },
      showTemplateImporter: false,
      importableTemplates: [],
      templatesToImport: {},
      unproxiedTheme: this.$vuetify.theme
    }
  },
  mounted () {
    this.loadData()
    const proxy = new Proxy(this.$vuetify.theme, {
      set: (target, key, value) => {
        this.$forceUpdate()
        target[key] = value
        return true
      }
    })
    this.$vuetify.theme = proxy
  },
  beforeDestroy () {
    this.$vuetify.theme = this.unproxiedTheme
  },
  methods: {
    loadData () {
      const ctx = this
      ctx.loading = true
      ctx.templates = {}
      ctx.$http.get('/api/templates').then(response => {
        response.data.map(template => {
          if (!template.display) template.display = template.name
          if (!template.type) template.type = 'none'
          if (!ctx.templates[template.type]) ctx.templates[template.type] = []
          ctx.templates[template.type].push(template)
        })

        const keys = Object.keys(ctx.templates)
        const index = keys.indexOf('other')
        if (index !== -1) this.$delete(keys, index)
        keys.map(key => {
          if (ctx.templates[key].length === 1) {
            if (!ctx.templates.other) ctx.templates.other = []
            ctx.templates.other.push(ctx.templates[key][0])
            delete ctx.templates[key]
          }
        })

        ctx.templates = { ...ctx.templates }
        ctx.loading = false
      }).catch(handleError(ctx))
    },
    loadTemplateImporter () {
      const ctx = this
      ctx.importableTemplates = []
      ctx.templatesToImport = {}
      ctx.$http.post('/api/templates/import').then(response => {
        ctx.importableTemplates = response.data
        ctx.showTemplateImporter = true
      }).catch(handleError(ctx))
    },
    doImports () {
      this.showTemplateImporter = false
      this.$toast.info(this.$t('templates.ImportStarted'))
      Object.keys(this.templatesToImport)
        .filter(elem => this.templatesToImport[elem])
        .map(elem => this.importTemplate(elem))
    },
    importTemplate (template) {
      const ctx = this
      ctx.$http.post(`/api/templates/import/${template}`).then(response => {
        ctx.$toast.success(ctx.$t('templates.ImportSuccessful', { template }))
        ctx.loadData()
      }).catch(handleError(ctx))
    },
    isDark
  }
}
</script>
