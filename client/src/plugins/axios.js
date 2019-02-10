import Vue from 'vue'

// Lib imports
import axios from 'axios'
import Cookies from 'js-cookie'

Vue.prototype.$http = axios
Vue.prototype.axios = axios

Vue.prototype.createRequest = function () {
  let cookie = Cookies.get('puffer_auth')
  if (cookie === undefined) {
    cookie = ''
  }
  if (cookie === '') {
    return axios.create()
  } else {
    return axios.create({
      headers: { 'Authorization': 'Bearer ' + cookie }
    })
  }
}
