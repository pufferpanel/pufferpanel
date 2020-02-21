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
import { handleError } from '@/utils/api'

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
  beforeDestroy () {
    this.$disconnect()
    if (this.statRequest) {
      clearInterval(this.statRequest)
    }
  },
  methods: {
    loadServer () {
      const ctx = this
      this.$http.get(`/api/servers/${this.$route.params.id}?perms`).then(response => {
        ctx.server = response.data.server
        ctx.server.permissions = response.data.permissions
        const url = `${window.location.protocol === 'http:' ? 'ws' : 'wss'}://${window.location.host}/daemon/socket/${ctx.server.id}`
        ctx.$connect(url)
        ctx.statRequest = setInterval(ctx.callStats, 3000)
      }).catch(handleError(ctx))
    },
    callStats () {
      this.$socket.sendObj({ type: 'stat' })
    }
  }
}
</script>
