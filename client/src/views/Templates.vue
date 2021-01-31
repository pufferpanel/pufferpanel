<template>
  <v-container>
    <h1 v-text="$t('templates.Templates')" />
    <v-row>
      <v-col>
        <v-list
          subheader
          elevation="1"
          class="pt-2"
        >
          <div
            v-for="(elements, type, i) in templates"
            :key="type"
          >
            <v-subheader
              v-if="Object.keys(templates).length !== 1"
              :class="isDark() ? headerClasses.dark : headerClasses.light"
              v-text="type"
            />
            <div
              v-for="template in templates[type]"
              :key="template.name"
            >
              <v-list-item :to="(hasScope('templates.edit') || isAdmin()) ? {name: 'Template', params: {id: template.name}} : undefined">
                <v-list-item-content>
                  <v-list-item-title v-text="template.display" />
                </v-list-item-content>
              </v-list-item>
            </div>
            <v-divider v-if="i !== Object.keys(templates).length - 1" />
          </div>
          <div
            v-if="Object.keys(templates).length === 0"
            class="pt-2 text-center text--disabled"
          >
            <span
              v-if="loading"
              v-text="$t('common.Loading')"
            />
            <span
              v-else
              v-text="$t('templates.NoTemplates')"
            />
          </div>
        </v-list>
        <div style="position: fixed; bottom: 16px; right: 16px; display: flex; flex-direction: column; align-items: center;">
          <v-tooltip left>
            <template v-slot:activator="{ on }">
              <v-btn
                v-show="hasScope('templates.edit') || isAdmin()"
                color="primary"
                class="mb-4"
                fab
                dark
                small
                v-on="on"
                @click="loadTemplateImporter"
              >
                <v-icon>mdi-import</v-icon>
              </v-btn>
            </template>
            <span v-text="$t('templates.import.Tooltip')" />
          </v-tooltip>
          <v-btn
            v-show="hasScope('templates.edit') || isAdmin()"
            color="primary"
            fab
            dark
            large
            :to="{name: 'AddTemplate'}"
          >
            <v-icon>mdi-plus</v-icon>
          </v-btn>
	</div>
      </v-col>
    </v-row>
    <ui-overlay
      v-model="showTemplateImporter"
      card
      closable
      :title="$t('templates.import.Import')"
    >
      <v-alert
        border="bottom"
        text
        type="warning"
        dense
      >
        {{ $t('templates.import.OverrideWarning') }}
      </v-alert>
      <v-autocomplete
        v-model="templatesToImport"
        :items="importableTemplates"
        chips
        :label="$t('templates.import.Select')"
        multiple
        clearable
        deletable-chips
        hide-selected
        outlined
        open-on-clear
      />
      <v-row>
        <v-col>
          <v-btn
            block
            color="error"
            @click="showTemplateImporter = false"
            v-text="$t('common.Cancel')"
          />
        </v-col>
        <v-col>
          <v-btn
            block
            color="success"
            @click="doImports()"
            v-text="$t('templates.import.Import')"
          />
        </v-col>
      </v-row>
    </ui-overlay>
    <ui-overlay
      v-model="offerImport"
      card
      closable
      :title="$t('templates.import.NoTemplatesTitle')"
    >
      <v-row>
        <v-col>
          <!-- eslint-disable-next-line vue/no-v-html -->
          <span v-html="$t('templates.import.NoTemplatesContent')" />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-btn
            block
            color="error"
            @click="importDeclined()"
            v-text="$t('common.Cancel')"
          />
        </v-col>
        <v-col>
          <v-btn
            block
            color="success"
            @click="offerImport = false; loadTemplateImporter()"
            v-text="$t('templates.import.Import')"
          />
        </v-col>
      </v-row>
    </ui-overlay>
  </v-container>
</template>

<script>
import { isDark } from '@/utils/dark'

export default {
  data () {
    return {
      loading: true,
      templates: {},
      headerClasses: {
        light: 'body-1 grey font-weight-bold lighten-4 black--text',
        dark: 'body-1 grey font-weight-bold darken-4 white--text'
      },
      offerImport: false,
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
    async loadData () {
      this.loading = true
      this.templates = {}
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
      this.loading = false
      if (templates.length === 0 && localStorage.getItem('offerTemplateImport') !== 'false') {
        this.offerImport = true
      }
    },
    async loadTemplateImporter () {
      this.importableTemplates = []
      this.templatesToImport = []
      this.importableTemplates = await this.$api.getImportableTemplates()
      this.showTemplateImporter = true
    },
    doImports () {
      const ctx = this
      ctx.loading = true
      ctx.showTemplateImporter = false
      ctx.$toast.info(this.$t('templates.import.Started'))
      ctx.templatesToImport.reduce((prev, next) => {
        return prev.finally(() => {
          return ctx.importTemplate(ctx, next)
        })
      }, Promise.resolve()).finally(() => { ctx.loadData() })
    },
    async importTemplate (ctx, template) {
      await ctx.$api.importTemplate(template)
      ctx.$toast.success(ctx.$t('templates.import.Successful', { template }))
    },
    importDeclined () {
      this.offerImport = false
      localStorage.setItem('offerTemplateImport', 'false')
    },
    isDark
  }
}
</script>
