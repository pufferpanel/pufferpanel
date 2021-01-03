// used as a mixin for the ApiClient, use of `this` refers to the ApiClient instance
export const TemplatesApi = {
  getTemplates () {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.get('/api/templates')).data
    })
  },

  getImportableTemplates () {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.post('/api/templates/import')).data
    })
  },

  getTemplate (name) {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.get(`/api/templates/${name}`)).data
    })
  },

  createTemplate (name, template) {
    return this.updateTemplate(name, template)
  },

  updateTemplate (name, template) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.put(`/api/templates/${name}`, template)
      return true
    }, { 400: 'errors.ErrInvalidJson' })
  },

  importTemplate (name) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.post(`/api/templates/import/${name}`)
      return true
    })
  },

  deleteTemplate (name) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.delete(`/api/templates/${name}`)
      return true
    })
  }
}
