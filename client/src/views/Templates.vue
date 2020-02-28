<template>
  <v-container>
    <h1 v-text="$t('templates.Templates')" />
    <v-row>
      <v-col>
        <v-list elevation="1">
          <div v-for="(template, index) in templates" :key="template.name">
            <v-list-item :to="(hasScope('templates.edit') || isAdmin()) ? {name: 'Template', params: {id: template.name}} : undefined">
              <v-list-item-content>
                <v-list-item-title v-text="template.display" />
              </v-list-item-content>
            </v-list-item>
            <v-divider v-if="index !== templates.length - 1" />
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
  </v-container>
</template>

<script>
import { handleError } from '@/utils/api'

export default {
  data () {
    return {
      loading: true,
      templates: []
    }
  },
  mounted () {
    this.loadData()
  },
  methods: {
    loadData () {
      const ctx = this
      ctx.loading = true
      ctx.templates = []
      ctx.$http.get('/api/templates').then(response => {
        ctx.templates = response.data
        for (let i = 0; i < ctx.templates.length; i++) {
          const t = ctx.templates[i]
          if (!t.display) {
            t.display = t.name
          }
        }
        ctx.loading = false
      }).catch(handleError(ctx))
    }
  }
}
</script>
