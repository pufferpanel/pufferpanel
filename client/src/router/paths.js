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
    path: '/addserver',
    view: 'AddServer',
    name: 'AddServer'
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
    path: '/user/:id',
    view: 'User',
    name: 'User'
  },
  {
    path: '/adduser',
    view: 'User',
    name: 'AddUser'
  },
  {
    path: '/user',
    view: 'Users',
    name: 'Users'
  },
  {
    path: '/template/:id',
    view: 'Template',
    name: 'Template'
  },
  {
    path: '/addtemplate',
    view: 'Template',
    name: 'AddTemplate'
  },
  {
    path: '/template',
    view: 'Templates',
    name: 'Templates'
  },
  {
    path: '/node',
    view: 'Nodes',
    name: 'Nodes'
  },
  {
    path: '/node/:id',
    view: 'Node',
    name: 'Node'
  },
  {
    path: '/addnode',
    view: 'AddNode',
    name: 'AddNode'
  },
  {
    path: '/settings',
    view: 'Settings',
    name: 'Settings'
  },
  {
    path: '/errors/404',
    name: 'Errors/404',
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
  },
  {
    path: '/auth/invite',
    name: 'Invite',
    view: 'Invite',
    meta: {
      noAuth: true,
      noSidebar: true,
      noFooter: true,
      noBase: true
    }
  }
]
