import EventEmitter from 'events'
import Cookies from 'js-cookie'
import axios from 'axios'
import { VBtn, VSpacer } from 'vuetify/lib'
import { ServersApi } from './servers'
import { NodesApi } from './nodes'
import { UsersApi } from './users'
import { TemplatesApi } from './templates'
import parseTheme from '@/utils/theme'

// most functions are supplied by mixins defined below this class definition
class ApiClient extends EventEmitter {
  _ctx = null
  _sockets = {}

  isLoggedIn () {
    return this.getToken() !== ''
  }

  getToken () {
    return Cookies.get('puffer_auth') || ''
  }

  updateCookie (token) {
    Cookies.set('puffer_auth', token, { sameSite: 'strict' })
  }

  saveLoginData (token, scopes = [], silent = false) {
    this.updateCookie(token)
    localStorage.setItem('scopes', JSON.stringify(scopes))
    if (!silent) this.emit('login')
  }

  logout (reason) {
    Cookies.remove('puffer_auth')
    this.emit('logout', reason)
  }

  register (username, email, password) {
    return this.withErrorHandling(async ctx => {
      const res = (await ctx.$http.post('/auth/register', { username, email, password })).data
      const hasLogin = res.token && res.token !== ''
      if (hasLogin) this.saveLoginData(res.token, res.scopes || [])
      return hasLogin
    })
  }

  login (email, password, options = {}) {
    return this.withErrorHandling(async ctx => {
      const res = (await ctx.$http.post('/auth/login', { email, password })).data
      this.saveLoginData(res.session, res.scopes || [], options.silent)
      return true
    }, { noToast: options.silent || options.noToast })
  }

  reauth () {
    return this.withErrorHandling(async ctx => {
      const res = (await ctx.$http.post('/auth/reauth')).data
      this.updateCookie(res.session)
      return true
    }, { noToast: true })
  }

  getConfig () {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.get('/api/config')).data
    })
  }

  getTheme (theme) {
    return this.withErrorHandling(async ctx => {
      const res = (await ctx.$http.get(`/theme/${theme}.tar`, { responseType: 'arraybuffer' })).data
      return parseTheme(res)
    })
  }

  getSelf () {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.get('/api/self')).data
    })
  }

  updateSelf (username, email, password, options = {}) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.put('/api/self', { username, email, password })
      return true
    }, options)
  }

  updatePassword (oldPassword, newPassword, options = {}) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.put('/api/self', { password: oldPassword, newPassword })
      return true
    }, options)
  }

  getSetting (key, options = {}) {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.get(`/api/settings/${key}`)).data.value
    }, options)
  }

  setSetting (key, value, options = {}) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.put(`/api/settings/${key}`, { value })
      if (key === 'panel.settings.companyName') {
        this.emit('panelTitleChanged', value)
      }
      return true
    }, options)
  }

  getUserSettings (options = {}) {
    return this.withErrorHandling(async ctx => {
      const map = {};
      (await ctx.$http.get('/api/userSettings')).data.map(elem => {
        map[elem.key] = elem.value
      })
      return map
    }, options)
  }

  setUserSetting (key, value, options = {}) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.put(`/api/userSettings/${key}`, { value })
      return true
    }, options)
  }

  async withErrorHandling (f, errorOptions) {
    try {
      return await f(this._ctx)
    } catch (err) {
      return this.handleError(err, errorOptions)
    }
  }

  handleError (error, options) {
    if (!options) options = {}

    // eslint-disable-next-line no-console
    console.error('ERROR', error)

    if (options.noToast || error.response.status === 401) return false

    let msg = 'errors.ErrUnknownError'
    if (error && error.response && error.response.data.error) {
      if (error.response.data.error.code) {
        msg = 'errors.' + error.response.data.error.code
      } else {
        msg = error.response.data.error.msg
      }
    }

    if (options[error.response.status] !== undefined) msg = options[error.response.status]

    const detailsAction = {
      timeout: 6000,
      slot: [
        this._ctx.$createElement('div', { attrs: { class: 'flex-grow-1' } }, [
          this._ctx.$createElement('span', [this._ctx.$t('errors.ErrUnknownError')]),
          this._ctx.$createElement(VSpacer, []),
          this._ctx.$createElement(VBtn, {
            props: { text: true, right: true },
            on: {
              click: () => window.pufferpanel.showError(error)
            }
          }, [this._ctx.$t('common.Details')])
        ])
      ]
    }

    const errUnknown = msg === 'errors.ErrUnknownError'

    this._ctx.$toast.error(errUnknown ? undefined : this._ctx.$t(msg), errUnknown ? detailsAction : undefined)
    return false
  }
}

export default function (Vue) {
  // add mixins for resources into the ApiClient
  Object.assign(ApiClient.prototype, ServersApi, NodesApi, UsersApi, TemplatesApi)

  Vue.prototype.$api = new ApiClient()
  Vue.prototype.$http = axios.create()
  if (process.env.NODE_ENV !== 'production') {
    window.pufferpanel.api = Vue.prototype.$api
  }

  // automagically add auth token to api requests
  Vue.prototype.$http.interceptors.request.use(request => {
    if (request.url.startsWith('/api') || request.url.startsWith('/proxy')) {
      request.headers[request.method].Authorization = 'Bearer ' + Vue.prototype.$api.getToken()
    }
    return request
  }, error => {
    return Promise.reject(error)
  })

  // handle 401 api responses as session timeout
  Vue.prototype.$http.interceptors.response.use(response => {
    return response
  }, error => {
    if (((error || {}).response || {}).status === 401) {
      localStorage.setItem('reauth_reason', 'session_timed_out')
      Vue.prototype.$api.logout('session_timed_out')
    }
    return Promise.reject(error)
  })

  Vue.prototype.hasScope = scope => {
    const scopes = localStorage.getItem('scopes') || '[]'
    return JSON.parse(scopes).indexOf(scope) !== -1
  }
  Vue.prototype.hasAuth = () => Vue.prototype.$api.getToken() !== ''
  Vue.prototype.isAdmin = () => Vue.prototype.hasScope('servers.admin')

  Vue.mixin({
    beforeCreate () {
      if (Vue.prototype.$api._ctx == null) Vue.prototype.$api._ctx = this
    }
  })
}
