let _isDark = null
let _ctx = null

export default function (Vue) {
  Vue.prototype.$isDark = () => _isDark
  Vue.prototype.$toggleDark = () => Vue.prototype.$setDark(!_isDark)
  Vue.prototype.$setDark = async (newValue) => {
    if (newValue !== false && newValue !== true) {
      if (newValue === 'false') {
        newValue = false
      } else if (newValue === 'true') {
        newValue = false
      } else {
        return
      }
    }
    if (Vue.prototype.$api.isLoggedIn()) await Vue.prototype.$api.setUserSetting('dark', `${newValue}`)
    localStorage.setItem('dark', `${newValue}`)
    _isDark = newValue
    _ctx.$vuetify.theme.dark = _isDark
  }

  if (process.env.NODE_ENV !== 'production') {
    window.pufferpanel.isDark = Vue.prototype.$isDark
    window.pufferpanel.setDark = Vue.prototype.$setDark
    window.pufferpanel.toggleDark = Vue.prototype.$toggleDark
  }

  Vue.mixin({
    async created () {
      if (_ctx == null) {
        _ctx = this
        const loaded = Vue.prototype.$api.isLoggedIn() ? (await Vue.prototype.$api.getUserSettings()).dark : null
        if (loaded === 'true' || loaded === 'false') {
          _isDark = loaded === 'true'
        } else {
          _isDark = (localStorage.getItem('dark') || '') === 'true'
        }
        _ctx.$vuetify.theme.dark = _isDark
      }
    }
  })
}
