<template>
  <server-type-generic :server="server" />
</template>

<script>
export default {
  props: {
    server: { type: Object, default: () => {} }
  }
}
</script>

<style scoped>

</style>
