<template>
  <b-col cols="1" class="sidebar-collapse">
    <nav class="bg-light sidebar">
      <div class="sidebar-sticky">
        <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted">
          <span>General</span>
        </h6>
        <ul class="nav flex-column">
          <li class="nav-item">
            <a class="nav-link" href="#">Account Settings</a>
          </li>
          <li class="nav-item">
            <a class="nav-link" href="#">Servers</a>
          </li>
        </ul>

        <footer class="footer">
          <div class="container">
            <p><a href="https://pufferpanel.com/">PufferPanel</a> is licensed under the <a
              href="https://github.com/PufferPanel/PufferPanel/blob/master/LICENSE">Apache License</a>.<br>Running
              Version v2.0.0</p>
          </div>
        </footer>
      </div>
    </nav>
  </b-col>
</template>

<style lang="css">
  @media (min-width: 992px) {
    .sidebar-collapse {
      display: block !important;
    }
  }

  @media (max-width: 991px) {
    .sidebar-collapse {
      display: none !important;
    }
  }
</style>