// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
// Components
import './components'
import BootstrapVue from 'bootstrap-vue'
import Vuetify from 'vuetify'
// Plugins
import './plugins'
// Sync router with store
import { sync } from 'vuex-router-sync'
// Application imports
import App from './App'
import i18n from '@/i18n'
import router from '@/router'
import store from '@/store'
import { library } from '@fortawesome/fontawesome-svg-core'
import { faCheckCircle, faTimesCircle } from '@fortawesome/free-regular-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import VueNativeSocket from 'vue-native-websocket'

library.add(faCheckCircle, faTimesCircle)

Vue.component('font-awesome-icon', FontAwesomeIcon)

Vue.use(BootstrapVue)
Vue.use(Vuetify)
Vue.use(VueNativeSocket, 'ws://localhost:1234', {
  connectManually: true,
  reconnection: true,
  reconnectionDelay: 5000,
  format: 'json'
})

// Sync store with router
sync(store, router)

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  i18n,
  router,
  store,
  render: h => h(App)
}).$mount('#app')