<template>
  <v-app>
    <core-header v-if="!$route.meta.noHeader"/>

    <b-container fluid>
      <b-row>
        <b-col cols="12">
          <main role="main">
            <router-view/>
          </main>
        </b-col>
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
  @import '~vuetify/dist/vuetify.min.css';
  @import 'styles/pufferpanel.css';
</style>
