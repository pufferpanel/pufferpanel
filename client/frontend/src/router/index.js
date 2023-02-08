import { createRouter, createWebHistory } from 'vue-router'
import makeRoutes from './routes'
import defaultRoute from './defaultRoute'

export default function(api) {
  const router = createRouter({
    history: createWebHistory(),
    linkExactActiveClass: 'active',
    routes: makeRoutes(api),
    scrollBehavior (to, from, savedPosition) {
      if (savedPosition) {
        return savedPosition
      }
      if (to.hash) {
        return { selector: to.hash }
      }
      return { left: 0, top: 0 }
    }
  })

  router.beforeEach((to, from) => {
    if (to.meta.noAuth && api.auth.isLoggedIn()) {
        return defaultRoute(api)
    }

    if (!to.meta.noAuth && !api.auth.isLoggedIn()) {
      sessionStorage.setItem('returnTo', JSON.stringify({
        name: to.name,
        params: to.params,
        hash: to.hash,
        query: to.query
      }))
      return { name: 'Login' }
    }

    return true
  })

  return router
}
