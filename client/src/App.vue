<template>
  <v-app>
    <!--<core-toolbar v-if="!$route.meta.noHeader"/>-->

    <core-drawer v-if="!$route.meta.noSidebar"/>

    <core-view />
  </v-app>
</template>

<style lang="scss">
@import '@/styles/index.scss';

/* Remove in 1.2 */
.v-datatable thead th.column.sortable i {
  vertical-align: unset;
}
</style>
