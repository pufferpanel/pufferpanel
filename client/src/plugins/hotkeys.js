import { install, uninstall } from '@github/hotkey'

function registerHotkey (el, keys, vnode) {
  install(el, keys)
  el.addEventListener('hotkey-fire', (e) => {
    if (!vnode.child.$listeners.hotkey) return
    e.preventDefault()
    if (el.offsetWidth || el.offsetHeight || el.getClientRects().length) {
      vnode.child.$emit('hotkey', e.detail.path.join(' '))
    }
  })
}

export default function (Vue) {
  Vue.directive('hotkey', {
    bind: function (el, binding, vnode) {
      if (typeof binding.value === 'string') {
        registerHotkey(el, binding.value, vnode)
      }
      if (Array.isArray(binding.value)) {
        binding.value.map((v) => {
          registerHotkey(el, v, vnode)
        })
      }
    },
    unbind: function (el, binding, vnode) {
      uninstall(el)
    }
  })
}
