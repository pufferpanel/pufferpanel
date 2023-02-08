export default (config) => {
  return {
    install: (app) => {
      app.config.globalProperties.$config = config
      app.provide('config', config)

      document.title = config.branding.name
    }
  }
}
