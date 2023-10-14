export default {
  install: (app) => {
    app.directive('click-outside', {
      beforeMount(el, binding) {
        el.clickOutsideEvent = function (event) {
          if (!(el == event.target || el.contains(event.target))) {
            binding.value(event)
          }
        }
        document.body.addEventListener('click', el.clickOutsideEvent)
      },
      beforeUnmount(el) {
        document.body.removeEventListener('click', el.clickOutsideEvent)
      }
    })
  }
}
