<template>
  <b-container>

  </b-container>
</template>

<script>
import servers from '@/plugins/servers.js'

export default {
  props: {
    server: Object
  },
  methods: {
    loadConsole() {
      servers.connect(this.server);
    }
  }
}
</script>

<style scoped>

</style>