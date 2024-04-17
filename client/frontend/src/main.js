import { createApp } from 'vue'
import makeI18n from '@/plugins/i18n'
import api, {apiClient} from '@/plugins/api'
import clickOutside from '@/plugins/clickOutside'
import conditions from '@/plugins/conditions'
import configPlugin from '@/plugins/config'
import events from '@/plugins/events'
import hotkeys from '@/plugins/hotkeys'
import theme from '@/plugins/theme'
import toast from '@/plugins/toast'
import userSettings from '@/plugins/userSettings'
import validators from '@/plugins/validators'
import makeRouter from '@/router'
import App from "@/App.vue"

const checkEnv = !!import.meta.env.VITE_CHECK_ENV
if (/app\.github\.dev/.test(window.location.host) && checkEnv) {
  const err = document.createElement('div')
  err.style.border = '16px solid red'
  err.style.backgroundColor = '#f8d7da'
  err.style.color = '#58151c'
  err.style.textAlign = 'center'
  err.style.display = 'flex'
  err.style.flexDirection = 'column'
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
  document.getElementById('hideApp').remove()
  throw new Error('github codespaces detected')
}

if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register('/sw.js', { scope: '/' })
}

window.pufferpanel = {}

async function mountApp(config) {
  createApp(App)
    .use(api)
    .use(events)
    .use(configPlugin(config))
    .use(userSettings(apiClient))
    .use(clickOutside)
    .use(conditions)
    .use(await makeI18n())
    .use(theme(config))
    .use(toast)
    .use(validators)
    .use(makeRouter(apiClient))
    .use(hotkeys)
    .mount('#app')
}

apiClient
  .getConfig()
  .then(config => {
    mountApp(config)
  })
  .catch(error => {
    console.log(error)
    mountApp({
      branding: {
        name: 'PufferPanel'
      },
      themes: {
        active: 'PufferPanel',
        available: ['PufferPanel']
      },
      registrationEnabled: true
    })
  })
