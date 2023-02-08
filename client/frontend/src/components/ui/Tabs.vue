<script>
import { ref, onMounted, onUnmounted, provide, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import Icon from './Icon.vue'

export default {
  components: {
    Icon
  },
  props: {
    anchors: { type: Boolean, default: () => false }
  },
  setup(props, { slots }) {
    const tabButtons = ref(null)
    const needsScroller = ref(false)

    const tabs = ref([])
    const activeKey = ref('')
    provide('activeKey', activeKey)

    function setActive(key) {
      activeKey.value = key

      if (props.anchors) {
        history.replaceState(history.state, '', '#' + key)
      }
    }

    function onResize() {
      needsScroller.value = tabButtons.value.scrollWidth > tabButtons.value.offsetWidth
    }

    function scroll(dir) {
      const dist = (tabButtons.value.offsetWidth / 2)
      tabButtons.value.scrollTo({
        behavior: 'smooth',
        left: tabButtons.value.scrollLeft + (dir === 'right' ? dist : dist * -1)
      })
    }

    onMounted(() => {
      window.addEventListener('resize', onResize)
      nextTick(() => onResize())

      tabs.value = slots
        .default()
        .filter(e => e && e.props && e.props.title)
        .map(e => {
          return {
            key: e.props.id || e.props.title.toLowercase().replace(/ /g, '-'),
            title: e.props.title,
            icon: e.props.icon,
            hotkey: e.props.hotkey
          }
        })

      if (props.anchors && tabs.value.length > 0 && location.hash) {
        const tab = tabs.value.find(e => e.key === location.hash.substring(1))
        if (tab) setActive(tab.key)
      }

      if (tabs.value.length > 0 && !activeKey.value) {
        setActive(tabs.value[0].key)
      }
    })

    onUnmounted(() => {
      window.removeEventListener('resize', onResize)
    })

    return { tabButtons, needsScroller, tabs, activeKey, setActive, scroll }
  }
}
</script>

<template>
  <div class="tabs">
    <div v-if="needsScroller" class="scroll-left" @click="scroll('left')" />
    <div v-if="needsScroller" class="scroll-right" @click="scroll('right')" />
    <div ref="tabButtons" class="tab-buttons">
      <div
        v-for="tab in tabs"
        :key="tab.key"
        v-hotkey="tab.hotkey"
        :class="['tab-button', tab.key === activeKey ? 'active' : 'inactive', tab.icon ? 'has-icon' : '']"
        @click="setActive(tab.key)"
      >
        <icon v-if="tab.icon" :name="tab.icon" />
        <span class="title" v-text="tab.title" />
      </div>
    </div>
    <slot />
  </div>
</template>
