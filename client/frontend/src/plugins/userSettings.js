let api = null
let events = null

const settings = {}

const settingsApi = {
  getAll: () => {
    return {...settings}
  },
  get: (key) => {
    return settings[key]
  },
  set: async (key, value) => {
    if (!api.auth.isLoggedIn()) return

    settings[key] = value
    return await api.settings.setUserSetting(key, value)
  },
  refresh: async () => {
    Object.keys(settings).map(key => {
      delete settings[key]
    })

    if (api.auth.isLoggedIn()) {
      const res = await api.settings.getUserSettings()
      Object.keys(res).map(key => {
        settings[key] = res[key]
      })
    }

    events.emit('userSettingsReloaded')
  }
}

export default (apiClient) => {
  return {
    install: (app) => {
      api = apiClient
      events = app.config.globalProperties.$events

      settingsApi.refresh()
      app.config.globalProperties.$userSettings = settingsApi
      app.provide('userSettings', settingsApi)

      events.on('login', settingsApi.refresh)
      events.on('logout', settingsApi.refresh)
    }
  }
}
