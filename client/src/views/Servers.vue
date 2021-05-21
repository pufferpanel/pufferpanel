<template>
  <v-container>
    <h1 v-text="$t('servers.Servers')" />
    <v-row>
      <v-col>
        <div
          v-if="servers.length === 0 && !loading"
          class="pt-2 text-center text--disabled"
          v-text="$t('servers.NoServers')"
        />
        <v-list
          two-line
          elevation="1"
        >
          <draggable
            :value="servers"
            :animation="100"
            handle=".dragHandle"
            @input="updateList($event)"
          >
            <v-list-item
              v-for="server in servers"
              :key="server.id"
              :to="{ name: 'Server', params: { id: server.id } }"
            >
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
              <v-list-item-icon class="dragHandle align-self-center mb-4">
                <v-icon>mdi-drag-vertical</v-icon>
              </v-list-item-icon>
            </v-list-item>
          </draggable>
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
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import draggable from 'vuedraggable'

export default {
  components: {
    draggable
  },
  data () {
    return {
      servers: [],
      loading: false,
      page: 0,
      pageCount: 1,
      task: null
    }
  },
  async mounted () {
    const serverList = (await this.$api.getUserSettings()).serverList
    if (serverList && serverList !== '') {
      this.servers = JSON.parse(serverList)
    }

    this.task = setInterval(this.pollServerStatus, 30 * 1000)
  },
  beforeDestroy () {
    if (this.task != null) {
      clearInterval(this.task)
    }
  },
  methods: {
    updateList (event) {
      this.servers = event
      this.$api.setUserSetting('serverList', JSON.stringify(event))
    },
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
      const toDelete = []
      for (const i in this.servers) {
        try {
          const servers = [...this.servers]
          const status = await this.$api.getServerStatus(servers[i].id, {
            403: () => {
              toDelete.push(servers[i].id)
              return true
            },
            404: () => {
              toDelete.push(servers[i].id)
              return true
            }
          })
          servers[i].online = status
          this.servers = servers
        } catch {}
      }
      toDelete.map(elem => this.servers.splice(this.servers.findIndex(el => el.id === elem), 1))
      if (toDelete.length !== 0) this.updateList(this.servers)
    }
  }
}
</script>
