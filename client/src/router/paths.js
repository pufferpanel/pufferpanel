/**
 * Define all of your application routes here
 * for more information on routes, see the
 * official documentation https://router.vuejs.org/en/
 */

export default [
  {
    path: '/server',
    view: 'Servers',
    name: 'Servers'
  },
  {
    path: '/account',
    view: 'Account',
    name: 'Account'
  },
  {
    path: '/server/:id',
    view: 'Server',
    name: 'Server'
  },
  {
    path: '/error/404',
    name: 'Error/404',
    view: 'errors/404',
    meta: {
      noAuth: true,
      noSidebar: false,
      noFooter: true,
      noBase: true
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
      noBase: true
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
      noBase: true
    }
  }
]
