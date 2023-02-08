import { install, uninstall } from '@github/hotkey'

let currentHotkeys = {}

function registerHotkey (el, keys, vnode) {
  install(el, keys)
  el.addEventListener('hotkey-fire', (e) => {
    if (!(vnode.props || {}).onHotkey) return
    e.preventDefault()
    if (el.offsetWidth || el.offsetHeight || el.getClientRects().length) {
      vnode.props.onHotkey(e.detail.path.join(' '))
    }
  })
}

// looking at vue internals isn't great, but i've not found a better way to do this yet
function getGroup(i, route) {
  let r = route
  if (!r) r = i.$route
  const rootId = i._.root.uid
  if (i._.uid === rootId) return 'root'
  if ((i._.type || {}).name === 'RouterView') return r.name
  return getGroup(i._.parent.ctx, r)
}

export default {
  install: (app) => {
    app.provide('hotkeys', () => {
      const res = {}
      Object.keys(currentHotkeys).map(k => {
        res[k] = Object.values(currentHotkeys[k]).flat()
      })
      return res
    })
    app.directive('hotkey', {
      beforeMount(el, binding, vnode) {
        if (!binding.value) return
        const group = getGroup(binding.instance)
        if (!currentHotkeys[group]) currentHotkeys[group] = {}
        if (!currentHotkeys[group][el]) currentHotkeys[group][el] = []
        currentHotkeys[group][el].push(binding.value)
        if (typeof binding.value === 'string') {
          registerHotkey(el, binding.value, vnode)
        }
        if (Array.isArray(binding.value)) {
          registerHotkey(el, binding.value.join(','), vnode)
        }
      },
      beforeUnmount(el, binding, vnode) {
        Object.keys(currentHotkeys).map(k => {
          if (currentHotkeys[k][el]) {
            currentHotkeys[k][el] = currentHotkeys[k][el].filter(e => e !== binding.value)
            if (currentHotkeys[k][el].length === 0) delete currentHotkeys[k][el]
            if (Object.keys(currentHotkeys[k]).length === 0) delete currentHotkeys[k]
          }
        })
        uninstall(el)
      }
    })
  }
}
