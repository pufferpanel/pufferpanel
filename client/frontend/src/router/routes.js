import defaultRoute from './defaultRoute'

export default (api) => [
  {
    path: '/',
    redirect: () => defaultRoute(api)
  },
  {
    path: '/auth/login',
    component: () => import('@/views/Login.vue'),
    name: 'Login',
    meta: {
      noAuth: true
    }
  },
  {
    path: '/auth/register',
    component: () => import('@/views/Registration.vue'),
    name: 'Register',
    meta: {
      noAuth: true
    }
  },
  {
    path: '/auth/invite',
    component: () => import('@/views/Invite.vue'),
    name: 'Invite',
    meta: {
      noAuth: true
    }
  },
  {
    path: '/servers',
    component: () => import('@/views/ServerList.vue'),
    name: 'ServerList',
    meta: {
      tkey: 'servers.Servers',
      permission: 'servers.view',
      icon: 'server',
      hotkey: 'g s'
    }
  },
  {
    path: '/servers/new',
    component: () => import('@/views/ServerCreate.vue'),
    name: 'ServerCreate'
  },
  {
    path: '/servers/view/:id',
    component: () => import('@/views/ServerView.vue'),
    name: 'ServerView'
  },
  {
    path: '/nodes',
    component: () => import('@/views/NodeList.vue'),
    name: 'NodeList',
    meta: {
      tkey: 'nodes.Nodes',
      permission: 'nodes.view',
      icon: 'node',
      hotkey: 'g n'
    }
  },
  {
    path: '/nodes/new',
    component: () => import('@/views/NodeCreate.vue'),
    name: 'NodeCreate'
  },
  {
    path: '/nodes/view/:id',
    component: () => import('@/views/NodeView.vue'),
    name: 'NodeView'
  },
  {
    path: '/users',
    component: () => import('@/views/UserList.vue'),
    name: 'UserList',
    meta: {
      tkey: 'users.Users',
      permission: 'users.view',
      icon: 'users',
      hotkey: 'g u'
    }
  },
  {
    path: '/users/new',
    component: () => import('@/views/UserCreate.vue'),
    name: 'UserCreate'
  },
  {
    path: '/users/view/:id',
    component: () => import('@/views/UserView.vue'),
    name: 'UserView'
  },
  {
    path: '/templates',
    component: () => import('@/views/TemplateList.vue'),
    name: 'TemplateList',
    meta: {
      tkey: 'templates.Templates',
      permission: 'templates.view',
      icon: 'template',
      hotkey: 'g t'
    }
  },
  {
    path: '/templates/new',
    component: () => import('@/views/TemplateCreate.vue'),
    name: 'TemplateCreate'
  },
  {
    path: '/templates/view/:repo/:id',
    component: () => import('@/views/TemplateView.vue'),
    name: 'TemplateView'
  },
  {
    path: '/settings',
    component: () => import('@/views/Settings.vue'),
    name: 'Settings',
    meta: {
      tkey: 'settings.Settings',
      permission: 'servers.admin',
      icon: 'settings',
      hotkey: 'g c'
    }
  },
  {
    path: '/self',
    component: () => import('@/views/Self.vue'),
    name: 'Self'
  }
]
