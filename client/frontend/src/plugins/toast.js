const toast = {
  show(type, message, cb, btnText) {
    const el = document.createElement('div')
    el.textContent = message
    el.classList.add('toast', 'in', type)
    if (cb) {
      const action = document.createElement('button')
      if (btnText) action.textContent = btnText
      action.addEventListener('click', cb)
      action.classList.add('button', 'toast', type)
      el.appendChild(action)
    }
    document.getElementById('toasts').appendChild(el)
    setTimeout(() => {
      el.classList.remove('in')
      setTimeout(() => {
        el.classList.add('out')
        setTimeout(() => el.remove(), 200)
      }, 5000)
    }, 200)
  },
  error(message, cb, btnText) {
    toast.show('error', message, cb, btnText)
  },
  success(message, cb, btnText) {
    toast.show('success', message, cb, btnText)
  }
}

export default {
  install: (app) => {
    app.config.globalProperties.$toast = toast
    app.provide('toast', toast)
  }
}
