const Notify = {
  install: function (Vue) {
    Vue.prototype.$notify = function (text, color = 'info') {
      this.$emit('notify', { text, color })
    }
  }
}

export default Notify
