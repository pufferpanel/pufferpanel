import Vue from 'vue'
import Vuetify, { VSnackbar } from 'vuetify/lib'
import 'vuetify/dist/vuetify.min.css'
import VuetifyToast from 'vuetify-toast-snackbar'

const opts = {
  components: { VSnackbar },
  theme: {
    options: {
      customProperties: true
    },
    themes: {
      light: {
        primary: '#07a7e3',
        secondary: '#e4e4e4',
        tertiary: '#888',
        anchor: '#07a7e3',
        accent: '#65a5f8'
      },
      dark: {
        primary: '#3b8db8',
        secondary: '#535353',
        tertiary: '#999',
        anchor: '#07a7e3',
        accent: '#65a5f8'
      }
    }
  },
  icons: {
    iconfont: 'mdi'
  }
}

Vue.use(Vuetify, opts)
Vue.use(VuetifyToast, { x: 'center', y: 'top', timeout: 2500, queueable: true })

export default new Vuetify(opts)
