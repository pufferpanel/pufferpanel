// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
// Components
import './components'
// Plugins
import vuetify from './plugins/vuetify'
import api from './plugins/api'
import dark from './plugins/dark'
import ace from './plugins/ace'
import hotkeys from './plugins/hotkeys'
// Application imports
import App from './App'
import i18n from '@/i18n'
import router from '@/router'
import '@/styles/pufferpanel.css'
// iconfont
import '@mdi/font/css/materialdesignicons.min.css'

if ('serviceWorker' in navigator) {
  if (process.env.NODE_ENV === 'production') {
    navigator.serviceWorker.register('/service-worker.js', { scope: '/' })
  } else {
    navigator.serviceWorker.register('/service-worker-dev.js', { scope: '/' })
  }
}

window.pufferpanel = {}

Vue.use(api)
Vue.use(dark)
Vue.use(ace)
Vue.use(hotkeys)

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  i18n,
  router,
  vuetify,
  render: h => h(App)
}).$mount('#app')
