<template>
  <v-container>
    <h1 v-text="$t('servers.Servers')" />
    <v-row>
      <v-col>
        <v-sheet
          elevation="1"
          class="mb-8"
        >
          <div
            v-if="servers.length === 0"
            class="pt-2 text-center text--disabled"
            v-text="$t('servers.NoServers')"
          />
          <v-list two-line>
            <div
              v-for="(server, index) in servers"
              :key="server.id"
            >
              <v-list-item :to="{ name: 'Server', params: { id: server.id } }">
                <v-list-item-avatar>
                  <v-icon
                    v-if="server.online"
                    color="success"
                  >
                    mdi-check-circle
                  </v-icon>
                  <v-icon
                    v-else
                    color="error"
                  >
                    mdi-alert-circle
                  </v-icon>
                </v-list-item-avatar>
                <v-list-item-content>
                  <v-list-item-title v-text="server.name" />
                  <v-list-item-subtitle v-text="server.address + ' @ ' + server.node" />
                </v-list-item-content>
              </v-list-item>
              <v-divider v-if="index !== servers.length - 1" />
            </div>
          </v-list>
          <v-row
            v-if="page < pageCount"
            ref="lazy"
            v-intersect="lazyLoad"
          >
            <v-col
              cols="2"
              offset="5"
            >
              <v-progress-circular
                indeterminate
                class="mr-2"
              />
              <span v-text="$t('common.Loading')" />
            </v-col>
          </v-row>
          <v-btn
            v-show="hasScope('servers.create') || isAdmin()"
            color="primary"
            bottom
            right
            fixed
            fab
            dark
            large
            :to="{name: 'AddServer'}"
          >
            <v-icon>mdi-plus</v-icon>
          </v-btn>
        </v-sheet>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  data () {
    return {
      servers: [],
      loading: false,
      page: 0,
      pageCount: 1,
      task: null
    }
  },
  mounted () {
    this.task = setInterval(this.pollServerStatus, 30 * 1000)
  },
  beforeDestroy () {
    if (this.task != null) {
      clearInterval(this.task)
    }
  },
  methods: {
    recheckLazy () {
      const rect = this.$refs.lazy.getBoundingClientRect()
      const viewHeight = Math.max(document.documentElement.clientHeight, window.innerHeight)
      if (
        !(rect.bottom < 0 || rect.top - viewHeight >= 0) &&
        this.page < this.pageCount
      ) {
        this.loadNextPage()
      }
    },
    lazyLoad (entries, observer, isIntersecting) {
      if (isIntersecting) {
        this.loadNextPage()
      }
    },
    async loadNextPage () {
      if (this.loading) return
      this.loading = true
      const { servers, pages } = await this.$api.getServers(this.page + 1)
      servers.filter(server => {
        return this.servers.filter(elem => server.id === elem.id).length === 0
      }).forEach(server => this.servers.push(server))
      this.page = this.page + 1
      this.pageCount = pages
      this.loading = false
      this.recheckLazy()
      this.pollServerStatus()
    },
    async pollServerStatus () {
      for (const i in this.servers) {
        const servers = [...this.servers]
        servers[i].online = await this.$api.getServerStatus(servers[i].id)
        this.servers = servers
      }
    }
  }
}
</script>
