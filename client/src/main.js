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

if (/app\.github\.dev/.test(window.location.host)) {
  const err = document.createElement('div')
  err.style.border = '16px solid red'
  err.style.backgroundColor = '#f8d7da'
  err.style.color = '#58151c'
  err.style.textAlign = 'center'
  err.style.display = 'flex'
  err.style.flexDirection = 'column'
  err.style.margin = '8px'
  err.innerHTML = `
<h1 style="margin:1em">!! IMPORTANT NOTICE !!</h1>
<h2 style="margin:1em">Usage of GitHub Codespaces for hosting is not permitted</h2>
<p style="text-wrap:balance;max-width:75%;align-self:center;line-height:1.5;margin:1em">
  Please read the <a href="https://docs.github.com/en/site-policy/github-terms/github-terms-for-additional-products-and-features#codespaces">GitHub Codespaces Terms of Service</a>.<br/>
  These explicitly state that the usage of GitHub Codespaces for
    <code style="background-color:white;border:1px solid;">any other activity unrelated to the development or testing of the software project associated with the repository where GitHub Codespaces is initiated</code>
  is not permitted.
</p>
<blockquote style="text-wrap:balance;max-width:75%;border-left:8px solid;align-self:center;background-color:white;margin:1em">
  In order to prevent violations of these limitations and abuse of GitHub Codespaces, GitHub may monitor your use of GitHub Codespaces.
  Misuse of GitHub Codespaces may result in termination of your access to Codespaces, restrictions in your ability to use GitHub Codespaces,
  or the disabling of repositories created to run Codespaces in a way that violates these Terms.
</blockquote>
<h4 style="text-wrap:balance;max-width:75%;align-self:center;margin:1em">
  The PufferPanel team does not tolerate or support the use of PufferPanel in order to facilitate violations against the GitHub Terms of Service
</h4>
`
  document.getElementById('app').appendChild(err)
  throw new Error('github codespaces detected')
}

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
