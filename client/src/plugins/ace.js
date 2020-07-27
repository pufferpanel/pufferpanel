import Vue from 'vue'

Vue.component('Ace', {
  props: {
    editorId: { type: String, default: () => 'editor' },
    value: { type: String, default: () => '' },
    lang: { type: String, default: () => 'text' },
    theme: { type: String, default: () => 'monokai' },
    file: { type: String, default: () => '' },
    height: { type: String, default: () => '100%' }
  },
  data () {
    return {
      editor: Object,
      ready: false
    }
  },
  mounted () {
    if (!window.ace) {
      const ctx = this
      const ace = document.createElement('script')
      ace.src = '/ace/ace.js'
      ace.onload = () => {
        const modelist = document.createElement('script')
        modelist.src = '/ace/ext-modelist.js'
        modelist.onload = () => {
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
    initialize () {
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

      this.ready = true
      this.$emit('editorready', true)
    },
    setValue (newValue) {
      this.editor.getSession().setValue(newValue, 1)
    }
  },
  render (createElement, context) {
    const height = this.height ? this.height : '100%'
    return createElement('div', { attrs: { id: this.editorId }, style: `width:100%;height:${height};` }, [])
  }
})
