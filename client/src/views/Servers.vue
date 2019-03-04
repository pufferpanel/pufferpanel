<template>
  <v-container
    fill-height
    fluid
    grid-list-xl
  >
    <v-layout
      justify-center
      wrap
    >
      <v-flex
        md12
      >
        <material-notification
          v-if="error"
          color="error"
          v-text="error"
        />
        <material-card
          color="blue"
          title="Servers"
        >
          <v-data-table
            :headers="headers"
            :items="servers"
            :pagination.sync="pagination"
            :total-items="totalServers"
            :loading="loading"
          >
            <template
              slot="headerCell"
              slot-scope="{ header }"
            >
              <span
                class="subheading font-weight-light text-success text--darken-3"
                v-text="header.text"
              />
            </template>
            <template
              slot="items"
              slot-scope="{ item }"
            >
              <td>{{ item.name }}</td>
              <td>{{ item.node }}</td>
              <td>{{ item.address }}</td>
              <td>
                <font-awesome-icon
                  v-if="item.online"
                  :icon= "['far','check-circle']"/>
                <font-awesome-icon
                  v-if="!item.online"
                  :icon= "['far','times-circle']"/>
              </td>
            </template>
          </v-data-table>
        </material-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
export default {
  data () {
    return {
      headers: [
        {
          text: 'Name',
          value: 'name'
        },
        {
          text: 'Node',
          value: 'node'
        },
        {
          text: 'Address',
          value: 'address'
        },
        {
          text: 'Online',
          value: 'online'
        }
      ],
      servers: [
      ],
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
