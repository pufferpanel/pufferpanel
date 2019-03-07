<template>
  <div>
    <core-header v-if="!$route.meta.noHeader"/>

    <b-container fluid>
      <b-row>
        <core-drawer v-if="!$route.meta.noSidebar"/>

        <main
          v-if="!$route.meta.noBase"
          role="main"
          class="col-md-9 ml-sm-auto col-lg-10 px-4">
          <b-container>
            <router-view/>
          </b-container>
        </main>
        <router-view v-else />
      </b-row>
    </b-container>
  </div>
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
