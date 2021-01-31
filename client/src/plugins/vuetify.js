import Vue from 'vue'
import Vuetify, { VSnackbar } from 'vuetify/lib'
import 'vuetify/dist/vuetify.min.css'
import VuetifyToast from 'vuetify-toast-snackbar'

const opts = {
  components: { VSnackbar },
  theme: {
    options: {
      inputStyle: 'outlined',
      customProperties: true
    },
    themes: {
      light: {
        primary: '#07a7e3',
        anchor: '#07a7e3',
        accent: '#65a5f8'
      },
      dark: {
        primary: '#3b8db8',
        anchor: '#07a7e3',
        accent: '#65a5f8'
      }
    }
  },
  icons: {
    iconfont: 'mdi'
  }
}

const vuetify = new Vuetify(opts)

Vue.use(Vuetify, opts)
Vue.use(VuetifyToast, { $vuetify: vuetify.framework, x: 'center', y: 'top', timeout: 2500, queueable: true })

export default vuetify
