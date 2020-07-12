<template>
  <v-container>
    <div class="d-flex">
      <h1 class="flex-grow-1" v-text="$t('templates.Templates')" />
      <v-tooltip left>
        <template v-slot:activator="{ on }">
          <v-btn icon v-on="on" @click="loadTemplateImporter"><v-icon>mdi-import</v-icon></v-btn>
        </template>
        <span v-text="$t('templates.import.Tooltip')" />
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
          <div class="pt-2 text-center text--disabled" v-if="Object.keys(templates).length === 0">
            <span v-if="loading" v-text="$t('common.Loading')" />
            <span v-else v-text="$t('templates.NoTemplates')" />
          </div>
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
    <common-overlay v-model="showTemplateImporter" card closable :title="$t('templates.import.Import')">
      <v-alert border="bottom" text type="warning" dense>
        {{ $t('templates.import.OverrideWarning') }}
      </v-alert>
      <v-autocomplete v-model="templatesToImport" :items="importableTemplates" chips :label="$t('templates.import.Select')" multiple clearable deletable-chips solo open-on-clear />
      <v-row>
        <v-col>
          <v-btn block color="error" v-text="$t('common.Cancel')" @click="showTemplateImporter = false" />
        </v-col>
        <v-col>
          <v-btn block color="success" v-text="$t('templates.import.Import')" @click="doImports()" />
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
      templatesToImport: [],
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
      ctx.templatesToImport = []
      ctx.$http.post('/api/templates/import').then(response => {
        ctx.importableTemplates = response.data
        ctx.templatesToImport = response.data
        ctx.showTemplateImporter = true
      }).catch(handleError(ctx))
    },
    doImports () {
      const ctx = this
      ctx.loading = true
      ctx.showTemplateImporter = false
      ctx.$toast.info(this.$t('templates.import.Started'))
      Promise.all(
        ctx.templatesToImport
          .map(elem => ctx.importTemplate(ctx, elem))
      ).then(x => { ctx.loadData() })
    },
    importTemplate (ctx, template) {
      return ctx.$http.post(`/api/templates/import/${template}`).then(response => {
        ctx.$toast.success(ctx.$t('templates.import.Successful', { template }))
      }).catch(handleError(ctx))
    },
    isDark
  }
}
</script>
