import Vue from 'vue'

Vue.component('Ace', {
  props: ['editorId', 'value', 'lang', 'theme', 'file'],
  data () {
    return {
      editor: Object
    }
  },
  mounted () {
    if (!window.ace) {
      const ctx = this
      const ace = document.createElement('script')
      ace.src = '/ace/ace.js'
      ace.onload = function () {
        const modelist = document.createElement('script')
        modelist.src = '/ace/ext-modelist.js'
        modelist.onload = function () {
          ctx.initialize()
        }
        document.head.appendChild(modelist)
      }
      document.head.appendChild(ace)
    } else {
      this.initialize()
    }
  },
  methods: {
    initialize: function () {
      const modelist = window.ace.require('ace/ext/modelist')
      const mode = this.lang || (this.file ? modelist.getModeForPath(this.file).mode : 'text')
      const theme = this.theme || 'monokai'

      this.editor = window.ace.edit(this.editorId)
      if (this.value && this.value.length > 0) {
        this.editor.getSession().setValue(this.value, 1)
      }
      this.editor.getSession().setMode(mode)
      this.editor.setTheme(`ace/theme/${theme}`)

      this.editor.on('change', () => {
        this.$emit('input', this.editor.getValue())
      })
    }
  },
  render: function (createElement, context) {
    return createElement('div', { attrs: { id: this.editorId }, style: 'width:100%;height:100%;' }, [])
  }
})
