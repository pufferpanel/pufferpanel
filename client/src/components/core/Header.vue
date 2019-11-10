<template>
  <b-navbar
    type="dark"
    variant="dark"
    class="fixed-top">

    <router-link to="/">
      <b-navbar-brand>PufferPanel</b-navbar-brand>
    </router-link>
    <!--<b-collapse
      id="nav_collapse"
      is-nav>
      <b-navbar-nav class="ml-auto">
        <b-nav-item-dropdown
          :text="$t('common.Language')"
          right>
          <b-dropdown-item href="#">EN</b-dropdown-item>
        </b-nav-item-dropdown>
        <b-nav-item-dropdown
          :text="$t('common.Account')"
          right>
          <router-link :to="{name: 'Account'}"><b-dropdown-item v-text="$t('common.Settings')"></b-dropdown-item></router-link>
          <router-link :to="{name: 'Logout'}"><b-dropdown-item v-text="$t('common.Logout')"></b-dropdown-item></router-link>
        </b-nav-item-dropdown>
      </b-navbar-nav>
    </b-collapse>-->
  </b-navbar>
</template>
