<template>
  <div class="sidebar-collapse">
    <nav class="bg-light sidebar sidebar-sticky" type="dark">
      <ul class="nav flex-column">
        <li class="nav-item">
          <b-link class="nav-link" to="/account" v-text="$t('common.AccountSettings')"></b-link>
        </li>
        <li class="nav-item">
          <b-link class="nav-link" :to="{name: 'Servers'}" v-text="$t('common.Servers')"></b-link>
        </li>
      </ul>
    </nav>
  </div>
</template>

<style lang="css">
  /*@media (min-width: 992px) {
    .sidebar-collapse {
      display: block !important;
    }
  }

  @media (max-width: 991px) {
    .sidebar-collapse {
      display: none !important;
    }
  }*/
</style>