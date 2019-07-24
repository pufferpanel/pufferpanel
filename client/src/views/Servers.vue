<template>
  <b-table hover selectable select-mode="single" @row-selected="rowSelected" :items="servers" :fields="fields"
           :busy="loading">
    <template slot="name" slot-scope="data">
      <strong v-text="data.value"></strong>
    </template>
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
      <strong :text="$t('common.Loading')">Loading...</strong>
    </div>

    <div slot="header">
      <b-btn :text="$t('common.AddServer')"></b-btn>
    </div>
  </b-table>
</template>

<script>
export default {
  data () {
    return {
      fields: {
        'name': {
          sortable: true,
          label: this.$t('common.Name')
        },
        'node': {
          sortable: true,
          label: this.$t('common.Node')
        },
        'address': {
          sortable: true,
          label: this.$t('common.Address')
        },
        'online': {
          sortable: true,
          label: this.$t('common.Online')
        }
      },
      servers: [],
      error: null,
      loading: true,
      totalServers: 0,
      pagination: {
        rowsPerPage: 10
      },
      task: null,
      canAddServer: this.hasScope('servers.create')
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
    this.task = setInterval(this.pollServerStatus, 30 * 1000)
  },
  methods: {
    loadData () {
      let vueData = this
      vueData.loading = true
      const { page, rowsPerPage } = this.pagination
      vueData.servers = []
      this.$http.get('/api/servers', {
        params: {
          page: page,
          limit: rowsPerPage
        }
      }).then(function (response) {
        let responseData = response.data
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
      }).catch(function (error) {
        let msg = 'errors.ErrUnknownError'
        if (error && error.response && error.response.data.error) {
          if (error.response.data.error.code) {
            msg = 'errors.' + error.response.data.error.code
          } else {
            msg = error.response.data.error.msg
          }
        }

        vueData.error = vueData.$t(msg)
      }).then(function () {
        vueData.loading = false
        vueData.pollServerStatus()
      })
    },
    pollServerStatus () {
      let vueData = this

      for (let i in this.servers) {
        let server = vueData.servers[i]
        vueData.$http.get('/daemon/server/' + server.id + '/status').then(function (response) {
          let data = response.data
          if (data) {
            let msg = data.data
            if (msg && msg.running) {
              server.online = true
            }
          }
        })
      }
    },
    rowSelected (items) {
      this.$router.push({ name: 'Server', params: { id: items[0].id } })
    }
  },
  beforeDestroy: function () {
    if (this.task != null) {
      clearInterval(this.task)
    }
  }
}
</script>
