<template>
  <v-footer
    id="core-footer"
    absolute
    height="82"
  >
    <div class="footer-items">
      <span
        v-for="link in links"
        :key="link.name"
      >
        <a
          :href="link.Link"
          class="tertiary--text footer-links">{{ link.name }}</a>
      </span>
    </div>
  </v-footer>
</template>

<script>
export default {
  data: () => ({
    links: [
      { name: 'Servers', Link: '/servers' }
    ]
  })
}
</script>

<style>
#core-footer {
  z-index: 0;
}
</style>
