import { ApiClient, ServerCookieSessionStore } from 'pufferpanel'

export const apiClient = new ApiClient(
  location.origin,
  new ServerCookieSessionStore()
)

export default {
  install: (app) => {
    app.config.globalProperties.$api = apiClient
    app.provide('api', apiClient)
  }
}
