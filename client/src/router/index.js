/**
 * Vue Router
 *
 * @library
 *
 * https://router.vuejs.org/en/
 */

// Lib imports
import Vue from 'vue'
import Router from 'vue-router'
// Routes
import paths from './paths'

import Cookies from 'js-cookie'

function route (path, view, name, meta) {
  return {
    name: name || view,
    path,
    component: resolve => import(`@/views/${view}.vue`).then(resolve),
    meta: meta
  }
}

function checkLoginState (next) {
  const cookie = Cookies.get('puffer_auth') || ''

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
  routes: paths
    .map(path => route(path.path, path.view, path.name, path.meta))
    .concat([
      { path: '/', redirect: 'server' },
      { path: '', redirect: 'server' },
      { path: '*', redirect: 'errors/404' }
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

export default router
