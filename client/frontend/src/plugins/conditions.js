import init, { resolve_if } from 'pufferpanel-conditions/conditions.js'

export default {
  install: async (app) => {
    await init()
    app.config.globalProperties.$conditions = resolve_if
    app.provide('conditions', resolve_if)
    window.pufferpanel.conditions = resolve_if
  }
}
