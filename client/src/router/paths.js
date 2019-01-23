/**
 * Define all of your application routes here
 * for more information on routes, see the
 * official documentation https://router.vuejs.org/en/
 */

export default [
  {
    path: '/server',
    view: 'Servers'
  },
  {
    path: '/error/404',
    name: 'Error',
    view: 'errors/404',
    meta: {
      noAuth: true,
      noSidebar: true,
      noFooter: true,
      noHeader: true
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
