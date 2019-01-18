<template>
  <v-app>
    <core-filter v-if="$route.error" />

    <core-toolbar v-if="$route.error" />

    <core-drawer  v-if="$route.error" />

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
