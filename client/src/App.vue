<template>
  <v-app>
    <core-header v-if="!$route.meta.noHeader"/>

    <b-container fluid>
      <b-row>
        <core-drawer v-if="!$route.meta.noSidebar"/>

        <main
          role="main"
          class="col-md-9 ml-sm-auto col-lg-10 px-4">
          <!--<core-view />-->
        </main>
      </b-row>
    </b-container>
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
