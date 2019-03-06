<template>
  <v-app>
    <core-drawer v-if="!$route.meta.noSidebar"/>

    <core-view />
  </v-app>
</template>

<style lang="scss">
  @import '~bootstrap/scss/bootstrap';
  @import '~bootstrap-vue/src/index.scss';
</style>

<style lang="css">
  @import '~animate.css/animate.css';
  @import '~selectize/dist/css/selectize.css';
  @import 'styles/pufferpanel.css';
</style>
