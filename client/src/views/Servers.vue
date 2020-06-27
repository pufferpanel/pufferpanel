<template>
  <v-container>
    <h1 v-text="$t('servers.Servers')" />
    <v-row>
      <v-col>
        <v-sheet elevation="1" class="mb-8">
          <div class="pt-2 text-center text--disabled" v-if="servers.length === 0" v-text="$t('servers.NoServers')" />
          <v-list two-line>
            <div v-for="(server, index) in servers" :key="server.id">
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
          <v-row ref="lazy" v-if="pagination.page < pagination.pageCount" v-intersect="lazyLoad">
            <v-col cols="2" offset="5">
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
import { handleError } from '@/utils/api'

export default {
  data () {
    return {
      headers: [
        {
          text: this.$t('common.Name'),
          value: 'name',
          sortable: true
        },
        {
          text: this.$t('nodes.Node'),
          value: 'node',
          sortable: true
        },
        {
          text: this.$t('common.Address'),
          value: 'address',
          sortable: true
        },
        {
          text: this.$t('common.Online'),
          value: 'online',
          sortable: true
        }
      ],
      servers: [],
      loading: false,
      recheckLazy: false,
      pagination: {
        page: 0,
        rowsPerPage: 10,
        pageCount: 1
      },
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
  watch: {
    recheckLazy (newVal) {
      const rect = /* document.getElementById('lazy') */ this.$refs.lazy.getBoundingClientRect()
      const viewHeight = Math.max(document.documentElement.clientHeight, window.innerHeight)
      if (
        !(rect.bottom < 0 || rect.top - viewHeight >= 0) &&
        this.pagination.page < this.pagination.pageCount
      ) {
        this.loadNextPage()
      }
    }
  },
  methods: {
    lazyLoad (entries, observer, isIntersecting) {
      if (isIntersecting) {
        this.loadNextPage()
      }
    },
    loadNextPage () {
      if (this.loading) return
      this.loading = true
      const ctx = this
      const { page, rowsPerPage } = this.pagination
      this.$http.get('/api/servers', {
        params: {
          page: page + 1,
          limit: rowsPerPage
        }
      }).then(response => {
        for (const i in response.data.servers) {
          const server = response.data.servers[i]

          let serverInList = false

          ctx.servers.forEach(elem => {
            if (server.id === elem.id) {
              serverInList = true
            }
          })

          if (!serverInList) {
            let ip = ''

            if (server.ip && server.ip !== '' && server.ip !== '0.0.0.0') {
              ip = server.ip
            } else {
              ip = server.node.publicHost
            }

            if (server.port) {
              ip += ':' + server.port
            }

            ctx.servers.push({
              id: server.id,
              name: server.name,
              node: server.node.name,
              address: ip,
              online: false,
              nodeAddress: server.node.publicHost + ':' + server.node.publicPort
            })
          }
        }
        ctx.pagination.page = page + 1
        const paging = response.data.paging
        ctx.pagination.pageCount = Math.ceil(paging.total / ctx.pagination.rowsPerPage)
        ctx.loading = false
        ctx.recheckLazy = true
      })
        .catch(handleError(ctx))
        .finally(() => {
          ctx.pollServerStatus()
        })
    },
    pollServerStatus () {
      for (const i in this.servers) {
        const server = this.servers[i]
        this.$http.get('/proxy/daemon/server/' + server.id + '/status').then(response => {
          if (response.data) {
            if (response.data.running) {
              server.online = true
            } else {
              server.online = false
            }
          }
        })
      }
    },
    rowSelected (item) {
      this.$router.push({ name: 'Server', params: { id: item.id } })
    },
    updatePage (newPage) {
      this.pagination.page = newPage
    }
  }
}
</script>
