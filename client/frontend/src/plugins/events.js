import { createNanoEvents } from 'nanoevents'

export default {
  install: (app) => {
    const emitter = createNanoEvents()
    app.config.globalProperties.$events = emitter
    app.provide('events', emitter)
  }
}
