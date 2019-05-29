import Vue from 'vue'

// Lib imports
import axios from 'axios'

Vue.prototype.$http = createClient()
Vue.prototype.axios = Vue.prototype.$http

/**
 * @deprecated Use $http instead
 * @returns {AxiosInstance}
 */
Vue.prototype.createRequest = function () {
  return Vue.prototype.$http
}

function createClient() {
  return axios.create()
}