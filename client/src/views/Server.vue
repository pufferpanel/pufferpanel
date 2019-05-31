<template>
  <div>
    <core-servers-type-generic v-if="this.server" :server="server"></core-servers-type-generic>
    <b-row v-else>
      <b-col cols="5"/>
      <b-col cols="2">
        <b-spinner class="align-middle"/>
        <strong :text="$t('common.Loading')" v-text="$t('common.Loading')"></strong>
      </b-col>
    </b-row>
  </div>
</template>

<script>
export default {
  data () {
    return {
      server: null,
      recover: null
    }
  },
  mounted () {
    this.server = this.loadServer()
  },
  methods: {
    loadServer () {
      let vue = this
      this.$http.get('/api/servers/' + this.$route.params.id).then(function (response) {
        vue.server = response.data.data
        let base = location.protocol === 'https' ? 'wss://' : 'ws:/' + location.host
        let url = base + '/daemon/server/' + vue.server.id + '/console'
        vue.$connect(url)
      })
    }
  },
  beforeDestroy: function () {
    this.$disconnect()
  }
}
</script>
