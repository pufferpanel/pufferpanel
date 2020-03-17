<template>
  <v-container>
    <h1 v-text="$t('servers.Servers')" />
    <v-row>
      <v-col>
        <v-sheet elevation="1" class="pt-2">
          <v-data-table
            hide-default-footer
            style="cursor: pointer;"
            :headers="headers"
            :items="servers"
            :items-per-page="100"
            @click:row="rowSelected"
          >
            <template v-slot:item.online="{ item }">
              <v-icon
                v-if="item.online"
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
            </template>
          </v-data-table>
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
        this.$http.get('/daemon/server/' + server.id + '/status').then(response => {
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
