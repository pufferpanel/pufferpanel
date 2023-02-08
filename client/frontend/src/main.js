import { createApp } from 'vue'
import makeI18n from '@/plugins/i18n'
import api, {apiClient} from '@/plugins/api'
import clickOutside from '@/plugins/clickOutside'
import configPlugin from '@/plugins/config'
import events from '@/plugins/events'
import hotkeys from '@/plugins/hotkeys'
import theme from '@/plugins/theme'
import toast from '@/plugins/toast'
import userSettings from '@/plugins/userSettings'
import validators from '@/plugins/validators'
import makeRouter from '@/router'
import App from "@/App.vue"

if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register('/sw.js', { scope: '/' })
}

async function mountApp(config) {
  createApp(App)
    .use(api)
    .use(events)
    .use(configPlugin(config))
    .use(userSettings(apiClient))
    .use(clickOutside)
    .use(await makeI18n())
    .use(theme(config))
    .use(toast)
    .use(validators)
    .use(makeRouter(apiClient))
    .use(hotkeys)
    .mount('#app')
}

apiClient
  .getConfig()
  .then(config => {
    mountApp(config)
  })
  .catch(error => {
    console.log(error)
    mountApp({
      branding: {
        name: 'PufferPanel'
      },
      themes: {
        active: 'PufferPanel',
        available: ['PufferPanel']
      },
      registrationEnabled: true
    })
  })
