<template>
  <v-row>
    <v-col>
      <v-card>
        <v-card-title
          v-text="$t(create ? 'templates.New' : 'templates.Edit')"
        />
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
          <v-row>
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
          <v-row>
            <v-col>
              <v-btn
                color="primary"
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
      loading: false,
      create: this.$route.params.id === undefined,
      name: this.$route.params.id === undefined ? '' : this.$route.params.id,
      template: ''
    }
  },
  mounted () {
    if (!this.create) this.loadData()
  },
  methods: {
    loadData () {
      this.loading = true
      const ctx = this
      ctx.$http.get(`/api/templates/${ctx.$route.params.id}`).then(response => {
        const data = response.data
        data.readme = undefined
        ctx.template = JSON.stringify(data, undefined, 4)
        if (ctx.$refs.editor.ready) ctx.$refs.editor.setValue(ctx.template)
        ctx.loading = false
      }).catch(handleError(ctx))
    },
    save () {
      const ctx = this
      ctx.$http.put(`/api/templates/${ctx.name}`, ctx.template).then(response => {
        ctx.$toast.success(ctx.$t('templates.SaveSuccess'))
        if (ctx.create) ctx.$router.push({ name: 'Template', params: { id: ctx.name } })
      }).catch(handleError(ctx, { 400: 'errors.ErrInvalidJson' }))
    },
    isDark
  }
}
</script>
