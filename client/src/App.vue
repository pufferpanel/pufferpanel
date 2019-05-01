<template>
  <v-app>
    <core-header v-if="!$route.meta.noHeader"/>

    <b-container fluid>
      <b-row>
        <core-drawer v-if="!$route.meta.noSidebar"/>

        <b-col
          v-if="$route.meta.noBase"
          :cols="$route.meta.noSidebar ? '12' : '10'">
          <main role="main">
            <router-view/>
          </main>
        </b-col>
        <b-col v-else class="cols-md-11 cols-sm-12">
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
