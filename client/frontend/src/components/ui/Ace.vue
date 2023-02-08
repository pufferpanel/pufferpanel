<script setup>
import { ref, onMounted, onUnmounted, onUpdated, nextTick } from 'vue'

const props = defineProps({
  id: { type: String, default: () => 'editor' },
  file: { type: String, default: () => undefined },
  mode: { type: String, default: () => undefined },
  modelValue: { type: String, required: true }
})

const emit = defineEmits(['update:modelValue'])

const el = ref(null)

let editor = null

onMounted(() => {
  if (!window.ace) {
    const ace = document.createElement('script')
    ace.src = '/js/ace/ace.js'
    ace.type = 'text/javascript'

    ace.onload = () => {
      const ml = document.createElement('script')
      ml.src = '/js/ace/ext-modelist.js'
      ml.type = 'text/javascript'

      ml.onload = () => {
        init()
      }

      document.head.append(ml)
    }

    document.head.append(ace)
  } else {
    init()
  }
})

onUnmounted(() => {
  if (editor) editor.destroy()
})

onUpdated(() => {
  if (editor && editor.session.getValue() !== props.modelValue)
    editor.session.setValue(props.modelValue)
})

function init() {
  nextTick(() => {
    editor = window.ace.edit(props.id)
    let theme = getComputedStyle(el.value).getPropertyValue('--ace-theme').trim() || 'monokai'
    editor.setTheme(`ace/theme/${theme}`)
    if (props.mode) {
      editor.session.setMode(`ace/mode/${props.mode}`)
    } else if (props.file) {
      const modelist = ace.require('ace/ext/modelist')
      const mode = modelist.getModeForPath(props.file).mode
      editor.session.setMode(mode)
    }
    editor.session.setValue(props.modelValue)
    editor.on('change', () => {
      emit('update:modelValue', editor.session.getValue())
    })
    editor.focus()
    editor.commands.addCommand({
      name: 'unfocusEditor',
      bindKey: { win: 'Escape', mac: 'Escape' },
      exec: (editor) => {
        editor.blur()
      }
    })
  })
}
</script>

<template>
  <div :id="id" ref="el" v-hotkey class="ace" />
</template>
