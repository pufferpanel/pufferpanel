<template>
  <v-app>
    <v-app-bar
      app
      dark
      clipped-left
      color="primary"
    >
      <v-app-bar-nav-icon
        v-if="loggedIn && $vuetify.breakpoint.mdAndDown"
        @click="drawer = !drawer"
      />
      <v-toolbar-title class="headline">
        <span
          v-if="appConfig"
          v-text="appConfig.branding.name"
        />
      </v-toolbar-title>
      <div class="flex-grow-1" />
      <v-menu v-if="appConfig.themes.available.length > 1">
        <template v-slot:activator="{ on, attrs }">
          <v-btn
            icon
            v-bind="attrs"
            v-on="on"
          >
            <v-icon>mdi-format-color-fill</v-icon>
          </v-btn>
        </template>
        <v-list
          subheader
          class="theme-options"
        >
          <v-subheader v-text="$t('common.ThemeOptions')" />
          <v-list-item @click="$toggleDark">
            <span v-text="$t('common.DarkMode')" />
            <span class="flex-grow-1" />
            <ui-switch
              :value="$vuetify.theme.dark"
              class="ml-4 mb-4"
              @input="$setDark($event)"
            />
          </v-list-item>
          <v-radio-group
            v-model="appConfig.themes.active"
            hide-details
          >
            <v-list-item
              v-for="theme in appConfig.themes.available"
              :key="theme"
              @click="setTheme(theme)"
            >
              {{ theme }}
              <span class="flex-grow-1" />
              <v-radio :value="theme" />
            </v-list-item>
          </v-radio-group>
        </v-list>
      </v-menu>
      <v-btn
        v-else
        icon
        @click="$toggleDark"
      >
        <v-icon>mdi-lightbulb</v-icon>
      </v-btn>
      <v-btn
        icon
        @click="showLanguageSelect = true"
      >
        <v-icon>mdi-earth</v-icon>
      </v-btn>
    </v-app-bar>

    <v-navigation-drawer
      v-if="loggedIn"
      v-model="drawer"
      :mini-variant="minified"
      app
      clipped
    >
      <v-list>
        <v-list-item
          :to="{name: 'Account'}"
          link
        >
          <v-list-item-icon>
            <v-icon>mdi-account</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title v-text="$t('users.Account')" />
          </v-list-item-content>
        </v-list-item>

        <v-list-item
          v-if="hasScope('servers.view') || isAdmin()"
          :to="{name: 'Servers'}"
          link
        >
          <v-list-item-icon>
            <v-icon>mdi-server</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title v-text="$t('servers.Servers')" />
          </v-list-item-content>
        </v-list-item>

        <v-list-item
          v-if="hasScope('nodes.view') || isAdmin()"
          :to="{name: 'Nodes'}"
          link
        >
          <v-list-item-icon>
            <v-icon>mdi-server-network</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title v-text="$t('nodes.Nodes')" />
          </v-list-item-content>
        </v-list-item>

        <v-list-item
          v-if="hasScope('users.view') || isAdmin()"
          :to="{name: 'Users'}"
          link
        >
          <v-list-item-icon>
            <v-icon>mdi-account-multiple</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title v-text="$t('users.Users')" />
          </v-list-item-content>
        </v-list-item>

        <v-list-item
          v-if="hasScope('templates.view') || isAdmin()"
          :to="{name: 'Templates'}"
          link
        >
          <v-list-item-icon>
            <v-icon>mdi-file-code</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title v-text="$t('templates.Templates')" />
          </v-list-item-content>
        </v-list-item>

        <v-list-item
          v-if="isAdmin()"
          :to="{name: 'Settings'}"
          link
        >
          <v-list-item-icon>
            <v-icon>mdi-cog</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title v-text="$t('settings.Settings')" />
          </v-list-item-content>
        </v-list-item>
      </v-list>
      <template v-slot:append>
        <v-list>
          <v-list-item
            v-if="!$vuetify.breakpoint.smAndDown"
            link
            @click="minified = !minified"
          >
            <v-list-item-icon>
              <v-icon v-text="minified ? 'mdi-chevron-right' : 'mdi-chevron-left'" />
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title v-text="$t('common.Collapse')" />
            </v-list-item-content>
          </v-list-item>
          <v-list-item
            link
            @click="$api.logout()"
          >
            <v-list-item-icon>
              <v-icon>mdi-logout</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title v-text="$t('users.Logout')" />
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </template>
    </v-navigation-drawer>

    <v-main>
      <v-container fluid>
        <div
          v-if="!appConfig"
          class="d-flex flex-column"
          style="width:100%; height:100%;"
        >
          <div
            class="d-flex align-self-center flex-row"
            style="height:100%;"
          >
            <v-progress-circular
              indeterminate
              size="100"
              class="align-self-center"
            />
          </div>
        </div>
        <router-view v-else />
      </v-container>
      <ui-language v-model="showLanguageSelect" />
      <ui-overlay
        v-model="errorOverlayOpen"
        card
        closable
        :title="$t('common.ErrorDetails')"
      >
        <!-- eslint-disable-next-line vue/no-v-html -->
        <div v-html="errorText" />
      </ui-overlay>
      <ui-overlay
        v-model="otpRequested"
        card
        :title="$t('users.OtpNeeded')"
      >
        <v-row>
          <v-col>
            <ui-input v-model="otpToken" autofocus @keyup.enter="confirmOtp" />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-btn v-text="$t('users.Login')" block color="success" @click="confirmOtp" />
          </v-col>
        </v-row>
      </ui-overlay>
    </v-main>
  </v-app>
</template>

<script>
import config from './config'

export default {
  data () {
    return {
      appConfig: config.defaultAppConfig,
      userSettings: [],
      otpRequested: false,
      otpToken: '',
      loggedIn: false,
      drawer: null,
      minified: false,
      reauthTask: null,
      showLanguageSelect: false,
      errorOverlayOpen: false,
      errorText: '',
      css: document.createElement('style'),
      themeObjects: []
    }
  },
  mounted () {
    this.$api.on('requestOtp', () => { this.otpRequested = true })
    this.$api.on('login', this.didLogIn)
    this.$api.on('logout', this.logout)
    this.$api.on('panelTitleChanged', this.updatePanelTitle)

    this.css.type = 'text/css'
    document.head.appendChild(this.css)

    this.loadConfig()

    if (this.$api.isLoggedIn()) {
      this.didLogIn()
    } else {
      this.loggedIn = false
    }
    window.pufferpanel.showError = error => this.showError(error)
  },
  methods: {
    async loadTheme () {
      const theme = await this.$api.getTheme(this.appConfig.themes.active)
      this.themeObjects.map(url => URL.revokeObjectURL(url))
      this.themeObjects = []
      this.css.textContent = ''
      const newThemeObjects = {}
      theme.forEach(file => {
        switch (file.name) {
          case 'theme.json':
            this.$vuetify.theme.themes.light = {
              ...this.$vuetify.theme.themes.light,
              ...JSON.parse(file.content).colors.light
            }
            this.$vuetify.theme.themes.dark = {
              ...this.$vuetify.theme.themes.dark,
              ...JSON.parse(file.content).colors.dark
            }
            break
          case 'theme.css':
            this.css.textContent = file.content
            break
          default:
            newThemeObjects[file.name] = URL.createObjectURL(file.blob)
        }
      })
      Object.keys(newThemeObjects).map(key => {
        this.css.textContent = this.css.textContent.split(key).join(newThemeObjects[key])
        this.themeObjects.push(newThemeObjects[key])
      })
    },
    async loadConfig () {
      this.userSettings = this.$api.isLoggedIn() ? await this.$api.getUserSettings() : {}
      if (this.userSettings.dark && this.userSettings.dark !== '') this.$vuetify.theme.dark = this.userSettings.dark === 'true'
      const config = await this.$api.getConfig()
      this.appConfig = { ...this.appConfig, ...config }
      document.title = this.appConfig.branding.name
      if (this.userSettings.preferredTheme || localStorage.getItem('theme')) {
        const stored = this.userSettings.preferredTheme || localStorage.getItem('theme')
        if (this.appConfig.themes.available.indexOf(stored) !== -1) {
          this.appConfig.themes.active = stored
        }
      }
      this.loadTheme()
    },
    updatePanelTitle (newTitle) {
      this.appConfig.branding.name = newTitle
      document.title = this.appConfig.branding.name
    },
    async setTheme (newTheme) {
      if (this.loggedIn) this.$api.setUserSetting('preferredTheme', newTheme)
      localStorage.setItem('theme', newTheme)
      this.appConfig.themes.active = newTheme
      this.loadTheme()
    },
    async confirmOtp () {
      if (await this.$api.loginOtp(this.otpToken)) {
        this.otpRequested = false
        if (this.hasScope('servers.view') || this.isAdmin()) {
          this.$router.push({ name: 'Servers' })
        } else {
          this.$router.push({ name: 'Account' })
        }
      }
      this.otpToken = ''
    },
    logout (reason) {
      if (!this.loggedIn) return
      this.reauthTask && clearInterval(this.reauthTask)
      this.reauthTask = null
      this.loggedIn = false
      this.$router.push({ name: 'Login' })
      if (reason === 'session_timed_out') this.$toast.error(this.$t('errors.ErrSessionTimedOut'))
    },
    didLogIn () {
      this.loggedIn = true
      this.reauthTask = setInterval(async () => {
        await this.$api.reauth()
      }, 1000 * 60 * 10)
      this.loadConfig()
    },
    showError (error) {
      const getCircularReplacer = () => {
        const seen = new WeakSet()
        return (key, value) => {
          if (key === 'password') return '[password]'
          if (typeof value === 'string') {
            try {
              const json = JSON.parse(value)
              if (typeof json === 'object' && json !== null) {
                if (Object.keys(json).indexOf('password') !== -1) {
                  json.password = '[password]'
                }
                return JSON.stringify(json)
              } else { return value }
            } catch { return value }
          }
          if (typeof value === 'object' && value !== null) {
            if (seen.has(value)) {
              return
            }
            seen.add(value)
          }
          return value
        }
      }

      if (error.config) {
        const responseCode = error.message.substring(error.message.length - 3, error.message.length)
        let codeMessage = error.message
        switch (responseCode) {
          case '401':
            codeMessage = 'Not logged in (401)'
            break
          case '403':
            codeMessage = 'Permission denied (403)'
            break
          case '404':
            codeMessage = 'Endpoint not found, is something blocking access to the API? (404)'
            break
          case '500':
            codeMessage = 'Server error (500)'
            break
          case '502':
            codeMessage = 'Bad Gateway, is PufferPanel running? (502)'
            break
        }

        let auth = error.config.headers.Authorization
        if (auth) {
          auth = auth.substring(0, 16)
          if (auth.length === 16) auth = auth + '...'
        } else auth = 'Not present'

        let body = error.config.data
        if (body) {
          body = JSON.stringify(JSON.parse(body), getCircularReplacer(), 4)
        }

        const message = `${codeMessage}

Endpoint: ${error.config.method} ${error.config.url}

Authorization Header: ${auth}

${body ? 'Request Body: ' + body : ''}

Stacktrace: ${error.stack}`
          .replace(/>/g, '$gt;')
          .replace(/</g, '&lt;')
          .replace(/ /g, '&nbsp;')
          .replace(/\n/g, '<br>')

        this.errorText = message
      } else {
        this.errorText = JSON.stringify(error, getCircularReplacer(), 4)
          .replace(/>/g, '$gt;')
          .replace(/</g, '&lt;')
          .replace(/ /g, '&nbsp;')
          .replace(/\n/g, '<br>')
      }

      this.errorOverlayOpen = true
    }
  }
}
</script>
