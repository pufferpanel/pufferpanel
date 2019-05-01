<template>
  <b-navbar
    type="dark"
    variant="dark"
    class="fixed-top">

    <router-link to="/">
      <b-navbar-brand>PufferPanel</b-navbar-brand>
    </router-link>
    <b-collapse
      id="nav_collapse"
      is-nav>
      <b-navbar-nav class="ml-auto">
        <b-nav-item-dropdown
          text="Language"
          right>
          <b-dropdown-item href="#">EN</b-dropdown-item>
        </b-nav-item-dropdown>
        <b-nav-item-dropdown
          text="Account"
          right>
          <router-link to="/account"><b-dropdown-item>Settings</b-dropdown-item></router-link>
          <router-link to="/auth/logout"><b-dropdown-item>Logout</b-dropdown-item></router-link>
        </b-nav-item-dropdown>
      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>
