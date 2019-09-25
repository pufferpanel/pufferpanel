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
import config from '../config'

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
      this.$http.get('/api/servers/' + this.$route.params.id).then(function (response) {
        vue.server = response.data.data
        const url = config.websocketBaseUrl + '/daemon/server/' + vue.server.id + '/socket'
        vue.$connect(url)
        vue.statRequest = setInterval(vue.callStats, 3000)
      })
    },
    callStats () {
      this.$socket.sendObj({ type: 'stat' })
    }
  }
}
</script>
