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
      return templateFromApi((await ctx.$http.get(`/api/templates/${name}`)).data)
    })
  },

  saveTemplate (name, template) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.put(`/api/templates/${name}`, templateToApi(template))
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
  },

  templateFromApiJson (template, server = false) {
    if (typeof template === 'string') {
      return templateFromApi(JSON.parse(template), server)
    } else {
      return templateFromApi(template, server)
    }
  },

  templateToApiJson (template, stringify = true) {
    return stringify ? JSON.stringify(templateToApi(template), undefined, 2) : templateToApi(template)
  }
}

const templateFromApi = (template, server = false) => {
  const { name, display, type, readme } = template

  const normalizeArrayOps = element => {
    if (element.type === 'download' && typeof element.files === 'string') element.files = [element.files]
    if (element.type === 'command' && typeof element.commands === 'string') element.commands = [element.commands]
    return element
  }

  if (!template.run) template.run = {}
  const command = template.run.command || ''
  const workingDirectory = template.run.workingDirectory || ''
  const autostart = template.run.autostart
  const autorestart = template.run.autorestart
  const autorecover = template.run.autorecover
  const stop = {}
  if (template.run.stop) {
    stop.type = 'command'
    stop.stop = template.run.stop
  } else if (template.run.stopCode) {
    stop.type = 'signal'
    stop.stop = template.run.stopCode
  } else {
    stop.type = 'command'
    stop.stop = ''
  }

  const envVars = template.run.environmentVars || {}

  const pre = (template.run.pre || []).map(normalizeArrayOps)
  const post = (template.run.post || []).map(normalizeArrayOps)
  const install = (template.install || []).map(normalizeArrayOps)

  if (!template.data) template.data = {}
  const vars = []
  for (const name in template.data) {
    if (!template.data[name].type) template.data[name].type = 'string'
    vars.push({ ...template.data[name], name })
  }

  const defaultEnv = template.environment
  const supportedEnvs = !server ? template.supportedEnvironments || [] : undefined
  const id = server ? template.id : undefined

  return {
    id,
    name,
    display,
    type,
    command,
    workingDirectory,
    stop,
    pre,
    post,
    envVars,
    vars,
    install,
    defaultEnv,
    supportedEnvs,
    readme,
    autostart,
    autorestart,
    autorecover
  }
}

const templateToApi = (template) => {
  const { id, name, display, type, command, workingDirectory, stop, pre, post, envVars, vars, install, defaultEnv, supportedEnvs, autostart, autorestart, autorecover } = template

  const convertedStop = {}
  if (stop.type === 'signal') {
    convertedStop.stopCode = Number(stop.stop)
  } else {
    convertedStop.stop = stop.stop
  }

  const convertedVars = {}
  vars.forEach(variable => {
    convertedVars[variable.name] = variable
    delete convertedVars[variable.name].name
  })

  return {
    id,
    name,
    display,
    type,
    install,
    run: {
      ...convertedStop,
      command,
      workingDirectory,
      pre,
      post,
      environmentVars: envVars,
      autostart,
      autorestart,
      autorecover
    },
    data: convertedVars,
    environment: defaultEnv,
    supportedEnvironments: supportedEnvs
  }
}
