<template>
  <div>
    <server-render
      v-if="server"
      :server="server"
    />
    <v-row v-else>
      <v-col cols="5" />
      <v-col cols="2">
        <v-progress-circular
          indeterminate
          class="mr-2"
        />
        <strong v-text="$t('common.Loading')" />
      </v-col>
    </v-row>
  </div>
</template>

<script>
export default {
  data () {
    return {
      server: null,
      recover: null,
      statRequest: null
    }
  },
  mounted () {
    this.server = this.loadServer()
  },
  beforeDestroy: function () {
    this.$disconnect()
    if (this.statRequest) {
      clearInterval(this.statRequest)
    }
  },
  methods: {
    loadServer () {
      const vue = this
      this.$http.get(`/api/servers/${this.$route.params.id}?perms`).then(function (response) {
        vue.server = response.data.server
        vue.server.permissions = response.data.permissions
        const url = `${window.location.protocol === 'http:' ? 'ws' : 'wss'}://${window.location.host}/daemon/server/${vue.server.id}/socket`
        vue.$connect(url)
        vue.statRequest = setInterval(vue.callStats, 3000)
      }).catch(function (error) {
        let msg = 'errors.ErrUnknownError'
        if (error && error.response && error.response.data.error) {
          if (error.response.data.error.code) {
            msg = 'errors.' + error.response.data.error.code
          } else {
            msg = error.response.data.error.msg
          }
        }

        vue.$notify(vue.$t(msg), 'error')
      })
    },
    callStats () {
      this.$socket.sendObj({ type: 'stat' })
    }
  }
}
</script>
