/**
 * Vue Router
 *
 * @library
 *
 * https://router.vuejs.org/en/
 */

// Lib imports
import Vue from 'vue'
import VueAnalytics from 'vue-analytics'
import Router from 'vue-router'
import Meta from 'vue-meta'
// Routes
import paths from './paths'

import Cookies from 'js-cookie'

function route (path, view, name, meta) {
  return {
    name: name || view,
    path,
    component: (resolve) => import(
      `@/views/${view}.vue`
      ).then(resolve),
    meta: meta
  }
}

function checkLoginState (next) {
  let cookie = Cookies.get('puffer_auth')
  if (cookie === undefined) {
    cookie = ''
  }
  if (cookie === '') {
    next('/auth/login')
  } else {
    next()
  }
}

Vue.use(Router)

// Create a new router
const router = new Router({
  mode: 'history',
  routes: paths.map(path => route(path.path, path.view, path.name, path.meta)).concat([
    { path: '/', redirect: 'Servers' },
    { path: '', redirect: 'Servers' },
    {
      path: '*', component: (resolve) => import(
        `@/views/errors/404.vue`
        ).then(resolve)
    }
  ]),
  scrollBehavior (to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    if (to.hash) {
      return { selector: to.hash }
    }
    return { x: 0, y: 0 }
  }
})

router.beforeEach((to, from, next) => {
  if (to.matched.some(r => r.meta.noAuth)) {
    next()
  } else {
    checkLoginState(next)
  }
})

Vue.use(Meta)

// Bootstrap Analytics
// Set in .env
// https://github.com/MatteoGabriele/vue-analytics
if (process.env.GOOGLE_ANALYTICS) {
  Vue.use(VueAnalytics, {
    id: process.env.GOOGLE_ANALYTICS,
    router,
    autoTracking: {
      page: process.env.NODE_ENV !== 'development'
    }
  })
}

export default router
