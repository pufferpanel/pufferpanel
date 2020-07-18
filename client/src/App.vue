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
      <v-btn
        icon
        @click="toggleDark"
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
            @click="logout"
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

    <v-content>
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
        <router-view
          v-else
          @logged-in="didLogIn()"
          @show-error-details="showError"
        />
      </v-container>
      <common-language v-model="showLanguageSelect" />
      <common-overlay
        v-model="errorOverlayOpen"
        card
        closable
        :title="$t('common.ErrorDetails')"
      >
        <code v-text="errorText" />
      </common-overlay>
    </v-content>
  </v-app>
</template>

<script>
import Cookies from 'js-cookie'
import config from './config'
import { toggleDark as doToggleDark, isDark } from './utils/dark'

export default {
  data () {
    return {
      appConfig: config.defaultAppConfig,
      loggedIn: false,
      drawer: null,
      minified: false,
      reauhTask: null,
      showLanguageSelect: false,
      errorOverlayOpen: false,
      errorText: ''
    }
  },
  mounted () {
    this.loadConfig()

    this.$vuetify.theme.dark = isDark()

    if ((Cookies.get('puffer_auth') || '')) {
      this.didLogIn()
    } else {
      this.loggedIn = false
    }
  },
  methods: {
    loadConfig () {
      const ctx = this
      this.$http.get('/api/config').then(response => {
        ctx.appConfig = { ...ctx.appConfig, ...response.data }
      }).catch(error => console.log('config failed', error)) // eslint-disable-line no-console
    },
    toggleDark () {
      doToggleDark(this.$vuetify)
    },
    logout () {
      this.reauthTask && clearInterval(this.reauthTask)
      this.reauthTask = null
      Cookies.remove('puffer_auth')
      this.loggedIn = false
      this.$router.push({ name: 'Login' })
    },
    didLogIn () {
      this.loggedIn = true
      const ctx = this
      this.reauthTask = setInterval(() => {
        ctx.$http.post('/auth/reauth').then(response => {
          response.data.session && Cookies.set('puffer_auth', response.data.session)
        }).catch(error => console.log('reauth failed', error)) // eslint-disable-line no-console
      }, 1000 * 60 * 10)
    },
    showError (error) {
      const getCircularReplacer = () => {
        const seen = new WeakSet()
        return (key, value) => {
          if (key === 'password') return '<password>'
          if (typeof value === 'string') {
            try {
              const json = JSON.parse(value)
              if (typeof json === 'object' && json !== null) {
                if (Object.keys(json).indexOf('password' !== -1)) {
                  json.password = '<password>'
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
      this.errorText = JSON.stringify(error, getCircularReplacer(), 4)
      this.errorOverlayOpen = true
    }
  }
}
</script>
