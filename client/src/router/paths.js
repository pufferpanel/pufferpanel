/**
 * Define all of your application routes here
 * for more information on routes, see the
 * official documentation https://router.vuejs.org/en/
 */

export default [
  {
    path: '/dashboard',
    // Relative to /src/views
    view: 'Dashboard'
  },
  {
    path: '/user-profile',
    name: 'User Profile',
    view: 'UserProfile'
  },
  {
    path: '/table-list',
    name: 'Table List',
    view: 'TableList'
  },
  {
    path: '/typography',
    view: 'Typography'
  },
  {
    path: '/icons',
    view: 'Icons'
  },
  {
    path: '/maps',
    view: 'Maps'
  },
  {
    path: '/notifications',
    view: 'Notifications'
  },
  {
    path: '/upgrade',
    name: 'Upgrade to PRO',
    view: 'Upgrade'
  },
  {
    path: '/404',
    name: 'Error',
    view: 'Error-404',
    meta: {
      noAuth: true
    }
  },
  {
    path: '/auth/login',
    name: 'Login',
    view: 'Login',
    meta: {
      noAuth: true,
      noSidebar: true,
      noFooter: true,
      noHeader: true
    }
  },
  {
    path: '/auth/register',
    name: 'Register',
    view: 'Register',
    meta: {
      noAuth: true,
      noSidebar: true,
      noFooter: true,
      noHeader: true
    }
  }
]
