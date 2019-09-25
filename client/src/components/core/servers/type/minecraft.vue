<template>
  <core-servers-type-generic :server="server" />
</template>

<script>
export default {
  props: {
    server: { type: Object, default: function () { return {} } }
  }
}
</script>

<style scoped>

</style>
