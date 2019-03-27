<template>
  <v-container>
    <b-table striped hover :items="servers" :fields="fields" :busy="loading">
      <template slot="online" slot-scope="data">
        <font-awesome-icon
          v-if="data.value"
          :icon="['far','check-circle']"/>
        <font-awesome-icon
          v-if="!data.value"
          :icon="['far','times-circle']"/>
      </template>

      <div slot="table-busy" class="text-center text-danger my-2">
        <b-spinner class="align-middle"/>
        <strong>Loading...</strong>
      </div>
    </b-table>
  </v-container>
</template>

<script>
export default {
  data () {
    return {
      fields: {
        'name': {
          sortable: true
        },
        'node': {
          sortable: true
        },
        'address': {
          sortable: true
        },
        'online': {
          sortable: true
        }
      },
      servers: [],
      error: null,
      loading: true,
      totalServers: 0,
      pagination: {
        rowsPerPage: 10
      }
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
    this.loadData()
    this.pollServerStatus()
    setInterval(this.pollServerStatus, 30 * 1000)
  },
  methods: {
    loadData () {
      let vueData = this
      vueData.loading = true
      const { page, rowsPerPage } = this.pagination
      vueData.servers = []
      this.createRequest().get('/api/servers', {
        params: {
          page: page,
          limit: rowsPerPage
        }
      }).then(function (response) {
        let responseData = response.data
        if (responseData.success) {
          for (let i in responseData.data) {
            let server = responseData.data[i]
            vueData.servers.push({
              id: server.id,
              name: server.name,
              node: server.node.name,
              address: server.ip ? server.ip + ':' + server.port : server.node.publicHost,
              online: false,
              nodeAddress: server.node.publicHost + ':' + server.node.publicPort
            })
          }
          let paging = responseData.metadata.paging
          vueData.totalServers = paging.total
        } else {
          vueData.error = responseData.msg
        }
      }).catch(function (error) {
        if (error.response) {
          if (error.response.status === 403) {
            vueData.error = 'You do not have permissions to view servers'
          } else {
            let data = error.response.data
            let msg = 'unknown error'
            if (data) {
              msg = error
            } else if (msg.msg) {
              msg = msg.msg
            }
            vueData.error = msg
          }
        } else {
          vueData.error = error
        }
      }).then(function () {
        vueData.loading = false
      })
    },
    pollServerStatus () {
      let http = this.createRequest()
      let vueData = this

      for (let i in this.servers) {
        let server = vueData.servers[i]
        http.get('/daemon/server/' + server.id + '/status').then(function (response) {
          let data = response.data
          if (data) {
            let msg = data.data
            if (msg && msg.running) {
              server.online = true
            }
          }
        })
      }
    }
  }
}
</script>
