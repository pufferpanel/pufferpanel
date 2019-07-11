import Vue from 'vue'
// Lib imports
import axios from 'axios'

Vue.prototype.$http = axios.create()
Vue.prototype.axios = Vue.prototype.$http
