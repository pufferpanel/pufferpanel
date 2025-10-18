<script setup>
import { getCurrentInstance, onMounted, onUnmounted, ref, nextTick } from 'vue'
import Icon from '@/components/ui/Icon.vue'

const props = defineProps({
  actions: { type: Array, required: true },
  title: { type: String, required: true }
})

const source = ref(null)
const style = ref({})
const suggestedDirection = ref(null)
const menu = ref(null)

function canOpen() {
  return Array.isArray(props.actions) && props.actions.length !== 0
}

function onClick(event) {
  if (!canOpen()) return
  handleEvent(event, 'activator')
}

function onContext(event) {
  if (!canOpen()) return
  event.preventDefault()
  handleEvent(event)
}

function handleEvent(event, eventSource) {
  const bodyWidth = document.body.scrollWidth
  const bodyHeight = document.body.scrollHeight
  source.value = eventSource || event.pointerType
  nextTick(() => {
    const suggestX = (event.pageX + menu.value.clientWidth) > bodyWidth ? 'left' : 'right'
    const suggestY = (event.pageY + menu.value.clientHeight) > bodyHeight ? 'top' : 'bottom'
    suggestedDirection.value = `${suggestY} ${suggestX}`
    style.value = {
      '--x': `${event.x}px`,
      '--y': `${event.y}px`,
      '--layer-x': `${event.layerX}px`,
      '--layer-y': `${event.layerY}px`,
      '--client-x': `${event.clientX}px`,
      '--client-y': `${event.clientY}px`,
      '--page-x': `${event.pageX}px`,
      '--page-y': `${event.pageY}px`,
      '--screen-x': `${event.screenX}px`,
      '--screen-y': `${event.screenY}px`,
      '--width': `${menu.value.clientWidth}px`,
      '--height': `${menu.value.clientHeight}px`
    }
  })
}

function click(action) {
  action.action()
  close()
}

function close() {
  source.value = null
}

onMounted(() => {
  getCurrentInstance().vnode.el.parentNode.addEventListener('contextmenu', onContext, true)
})

onUnmounted(() => {
  getCurrentInstance().vnode.el.parentNode.removeEventListener('contextmenu', onContext, true)
})
</script>

<template>
  <slot name="activator" :can-open="canOpen()" :on-click="onClick" />
  <div v-if="source !== null" class="context-menu-wrapper">
    <ul ref="menu" v-click-outside.stop="close" v-hotkey="'Escape'" :class="['context-menu', `source-${source}`]" :style="style" :data-suggested-direction="suggestedDirection" @hotkey="close()">
      <slot name="title">
        <li class="context-title" v-text="title" />
      </slot>
      <slot v-for="action, index in actions" :key="index" name="action" :action="{...action, action: () => click(action)}">
        <li :class="['context-action', action.class || '']">
          <a v-hotkey="action.hotkey" @click.stop="click(action)">
            <icon v-if="action.icon" :name="action.icon" />
            <span class="label" v-text="action.label" />
          </a>
        </li>
      </slot>
    </ul>
  </div>
</template>
