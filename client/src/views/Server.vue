<template>
  <div>
    <server-render
      v-if="server"
      :server="server"
    />
    <v-row v-else>
      <v-col
        cols="12"
        class="d-flex align-center justify-center"
      >
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
      server: null
    }
  },
  mounted () {
    this.loadServer()
  },
  beforeDestroy () {
    this.$api.closeServerConnection(this.server.id)
  },
  methods: {
    async loadServer () {
      const server = await this.$api.getServer(this.$route.params.id)
      await this.$api.startServerConnection(server.id)
      this.server = server
    }
  }
}
</script>
