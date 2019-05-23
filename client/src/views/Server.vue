<template>
  <div>
    <div v-if="this.server">
      <!--<keep-alive>-->
        <!--<core-servers-minecraft v-if="server.type === 'minecraft'"></core-servers-minecraft>
        <core-servers-generic v-else v-bind:server="server"></core-servers-generic>-->
      <core-servers-generic v-bind:server="server"></core-servers-generic>
      <!--</keep-alive>-->
    </div>
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
      server: null
    }
  },
  mounted () {
    this.server = this.loadServer()
  },
  methods: {
    loadServer () {
      let vue = this
      this.createRequest().get('/api/servers/' + this.$route.params.id).then(function (response) {
        console.log(response)
        vue.server = response.data.data
      })
    }
  }
}
</script>
