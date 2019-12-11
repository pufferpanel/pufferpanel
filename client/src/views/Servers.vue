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
            :loading="loading"
            :items-per-page="pagination.rowsPerPage"
            :page.sync="pagination.page"
            :server-items-length="totalServers"
            @click:row="rowSelected"
            @page-count="updatePage"
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
          <div class="text-center py-2">
            <v-pagination
              v-model="pagination.page"
              :length="pagination.pageCount"
            />
          </div>
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
      error: null,
      loading: true,
      totalServers: 0,
      pagination: {
        page: 1,
        rowsPerPage: 10,
        pageCount: 1
      },
      task: null
    }
  },
  watch: {
    pagination: {
      handler () {
        this.loadData()
      },
      deep: true
    }
  },
  mounted () {
    this.pagination.page = 1
    this.loadData()
    this.task = setInterval(this.pollServerStatus, 30 * 1000)
  },
  beforeDestroy: function () {
    if (this.task != null) {
      clearInterval(this.task)
    }
  },
  methods: {
    loadData () {
      const vue = this
      vue.loading = true
      const { page, rowsPerPage } = this.pagination
      vue.servers = []
      this.$http.get('/api/servers', {
        params: {
          page: page,
          limit: rowsPerPage
        }
      }).then(function (response) {
        for (const i in response.data.servers) {
          const server = response.data.servers[i]

          let serverInList = false

          vue.servers.forEach(function (elem) {
            if (server.id === elem.id) {
              serverInList = true
            }
          })

          if (!serverInList) {
            let ip = ''

            if (server.ip && server.ip !== '' && server.ip !== '0.0.0.0') {
              ip = server.ip
              if (server.port) {
                ip += ':' + server.port
              }
            } else {
              ip = server.node.publicHost
            }

            vue.servers.push({
              id: server.id,
              name: server.name,
              node: server.node.name,
              address: ip,
              online: false,
              nodeAddress: server.node.publicHost + ':' + server.node.publicPort
            })
          }
        }
        const paging = response.data.paging
        vue.totalServers = paging.total
        vue.pagination.pageCount = Math.ceil(paging.total / vue.pagination.rowsPerPage)
      }).catch(function (error) {
        let msg = 'errors.ErrUnknownError'
        if (error && error.response && error.response.data.error) {
          if (error.response.data.error.code) {
            msg = 'errors.' + error.response.data.error.code
          } else {
            msg = error.response.data.error.msg
          }
        }

        vue.$toast.error(vue.$t(msg))
      }).then(function () {
        vue.loading = false
        vue.pollServerStatus()
      })
    },
    pollServerStatus () {
      for (const i in this.servers) {
        const server = this.servers[i]
        this.$http.get('/daemon/server/' + server.id + '/status').then(function (response) {
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
