import Vue from 'vue'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'

const opts = {
  theme: {
    options: {
      customProperties: true
    },
    themes: {
      light: {
        primary: '#07a7e3',
        secondary: '#e4e4e4',
        tertiary: '#888',
        accent: '#65a5f8'
      },
      dark: {
        primary: '#07a7e3',
        secondary: '#535353',
        tertiary: '#999',
        accent: '#65a5f8'
      }
    }
  },
  icons: {
    iconfont: 'mdi'
  }
}

Vue.use(Vuetify, opts)

export default new Vuetify(opts)
