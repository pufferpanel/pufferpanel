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
      ctx.users = []
      ctx.$http.get('/api/templates').then(function (response) {
        ctx.templates = response.data
        ctx.loading = false
      }).catch(function (error) {
        let msg = 'errors.ErrUnknownError'
        if (error && error.response && error.response.data.error) {
          if (error.response.data.error.code) {
            msg = 'errors.' + error.response.data.error.code
          } else {
            msg = error.response.data.error.msg
          }
        }

        ctx.$toast.error(ctx.$t(msg))
      })
    }
  }
}
</script>
