export default {
  install: (app) => {
    app.directive('click-outside', {
      beforeMount(el, binding) {
        el.clickOutsideEvent = function (event) {
          if (!(el == event.target || el.contains(event.target))) {
            if (binding.modifiers.stop && event.type === 'click' && el.checkVisibility()) {
              event.stopPropagation()
            }
            binding.value(event)
          }
        }
        document.body.addEventListener('click', el.clickOutsideEvent, true)
        document.body.addEventListener('contextmenu', el.clickOutsideEvent, true)
      },
      beforeUnmount(el) {
        document.body.removeEventListener('click', el.clickOutsideEvent, true)
        document.body.removeEventListener('contextmenu', el.clickOutsideEvent, true)
      }
    })
  }
}
