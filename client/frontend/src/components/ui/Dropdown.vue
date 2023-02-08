<script>
import { ref, nextTick } from 'vue'
import Multiselect from '@vueform/multiselect'
import Icon from '@/components/ui/Icon.vue'
import markdown from '@/utils/markdown.js'

export default {
  components: {
    Icon,
    Multiselect
  },
  props: {
    // default to random id as labels need target ids to exist exactly once
    id: { type: String, default: () => (Math.random() + 1).toString(36).substring(2) },
    label: { type: String, default: () => undefined },
    labelProp: { type: String, default: () => 'label' },
    error: { type: String, default: () => undefined },
    hint: { type: String, default: () => undefined },
    options: { type: Array, default: () => [] },
    type: { type: String, default: () => 'text' },
    icon: { type: String, default: () => undefined },
    modelValue: { type: [String, Number, Object], default: () => '' }
  },
  emits: ['update:modelValue', 'change'],
  setup() {
    const ms = ref(null)
    const isOpen = ref(false)

    function select(item) {
      nextTick(() => ms.value.select(item))
    }

    return { isOpen, ms, select, markdown }
  }
}
</script>

<template>
  <div class="dropdown-wrapper">
    <div :class="['dropdown', error ? 'error' : '']">
      <icon v-if="icon" class="pre" :name="icon" />
      <multiselect :id="id" ref="ms" :model-value="modelValue" :label="labelProp"  mode="single" :can-deselect="false" :can-clear="false" :options="options" :placeholder="label" @input="$emit('update:modelValue', $event); $emit('change', $event)" @open="$nextTick(() => isOpen = true)" @close="isOpen = false">
        <template v-for="(index, name) in $slots" #[name]="data">
          <slot :name="name" v-bind="data"></slot>
        </template>
      </multiselect>
      <label v-if="label" :class="[isOpen ? 'dropdown-open' : 'dropdown-closed', modelValue === null ? 'dropdown-placeholder' : '']" :for="id" @click="ms.open()">{{ label }}</label>
    </div>
    <span v-if="error" class="error" v-text="error" />
    <!-- eslint-disable-next-line vue/no-v-html -->
    <span v-if="hint && !error" class="hint" v-html="markdown(hint)" />
  </div>
</template>
