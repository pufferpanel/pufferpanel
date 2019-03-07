<template>
  <b-navbar
    type="dark"
    variant="dark"
    class="fixed-top">

    <b-navbar-brand href="/">PufferPanel</b-navbar-brand>
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
          <b-dropdown-item href="/account">Settings</b-dropdown-item>
          <b-dropdown-item href="/auth/logout">Logout</b-dropdown-item>
        </b-nav-item-dropdown>
      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>
