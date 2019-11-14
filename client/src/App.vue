<template>
  <v-app>
    <v-app-bar
      app
      dark
      clipped-left
      color="primary"
    >
      <v-app-bar-nav-icon
        v-if="loggedIn && $vuetify.breakpoint.smAndDown"
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
            <v-list-item-title v-text="$t('common.Account')" />
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
            <v-list-item-title v-text="$t('common.Servers')" />
          </v-list-item-content>
        </v-list-item>

        <v-list-item
          v-if="isAdmin()"
          :to="{name: 'Nodes'}"
          link
        >
          <v-list-item-icon>
            <v-icon>mdi-server-network</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title v-text="$t('common.Nodes')" />
          </v-list-item-content>
        </v-list-item>

        <v-list-item
          v-if="isAdmin()"
          :to="{name: 'Users'}"
          link
        >
          <v-list-item-icon>
            <v-icon>mdi-account-multiple</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title v-text="$t('common.Users')" />
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
              <v-list-item-title v-text="$t('common.Logout')" />
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
          @logged-in="loggedIn = true"
          v-else
        />
      </v-container>
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
      minified: false
    }
  },
  mounted () {
    this.loadConfig()

    this.$vuetify.theme.dark = isDark()

    if ((Cookies.get('puffer_auth') || '')) {
      this.loggedIn = true
    } else {
      this.loggedIn = false
    }
  },
  methods: {
    loadConfig () {
      const vue = this
      this.$http.get('/api/config').then(function (response) {
        vue.appConfig = { ...vue.appConfig, ...response.data }
      })
    },
    toggleDark () {
      doToggleDark(this.$vuetify)
    },
    logout () {
      Cookies.remove('puffer_auth')
      this.loggedIn = false
      this.$router.push({ name: 'Login' })
    }
  }
}
</script>
