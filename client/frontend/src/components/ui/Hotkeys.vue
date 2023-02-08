<script setup>
import { ref, inject, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'

const getHotkeys = inject('hotkeys')
const { t } = useI18n()
const route = useRoute()
const ctxOverrides = {
  'TemplateView': 'Template',
  'TemplateCreate': 'Template'
}
const hotkeys = ref(currentHotkeys())

function getContext() {
  return ctxOverrides[route.name] || route.name
}

function currentHotkeys() {
  const h = getHotkeys()
  const res = {
    global: h.root.flat().filter(e => e !== '?' && e !== 'Shift+?' && e !== 'Escape'),
    contextual: (h[route.name] || [])
      .flat()
      .filter(e => e !== 'Escape')
      .filter(e => !/^. \d \d$/.test(e))
      .sort(),
    context: getContext()
  }

  return res
}

watch(
  () => route.name,
  async newName => {
    setTimeout(() => {
      hotkeys.value = currentHotkeys()
    }, 500)
  }
)

onMounted(() => {
  hotkeys.value = currentHotkeys()
})
</script>

<template>
  <div class="hotkey-list">
    <div v-if="hotkeys.contextual.length > 0" class="contextual">
      <h3 v-text="t(`hotkeys.${hotkeys.context}.Title`)" />
      <div v-for="keys in hotkeys.contextual" :key="keys" class="hotkey">
        <span class="keys">
          <span v-for="key in keys.split(' ')" :key="key" class="key" v-text="key" />
        </span>
        <span class="description" v-text="t(`hotkeys.${hotkeys.context}.${keys}`)" />
      </div>
    </div>
    <div class="global">
      <h3 v-text="t('hotkeys.Global.Title')" />
      <div v-for="keys in hotkeys.global" :key="keys" class="hotkey">
        <span class="keys">
          <span v-for="key in keys.split(' ')" :key="key" class="key" v-text="key" />
        </span>
        <span class="description" v-text="t(`hotkeys.Global.${keys}`)" />
      </div>
    </div>
  </div>
</template>
