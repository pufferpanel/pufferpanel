<template>
  <div>
    <ui-switch
      v-if="value.type === 'boolean'"
      :value="value.value"
      :required="value.required"
      :label="value.display"
      @input="onInput($event)"
    >
      <template slot="message">
        <!-- eslint-disable-next-line vue/no-v-html -->
        <div v-html="markdown(value.desc)" />
      </template>
    </ui-switch>
    <ui-select
      v-else-if="value.type === 'option'"
      :value="value.value"
      item-text="display"
      item-value="value"
      :items="value.options"
      :label="value.display"
      @input="onInput($event)"
    >
      <template slot="message">
        <!-- eslint-disable-next-line vue/no-v-html -->
        <div v-html="markdown(value.desc)" />
      </template>
    </ui-select>
    <ui-input
      v-else
      :type="value.type === 'integer' ? 'number' : 'text'"
      :required="value.required"
      :label="value.display"
      :value="value.value"
      @input="onInput($event)"
    >
      <template slot="message">
        <!-- eslint-disable-next-line vue/no-v-html -->
        <div v-html="markdown(value.desc)" />
      </template>
    </ui-input>
  </div>
</template>

<script>
import markdown from '@/utils/markdown'

export default {
  props: {
    value: { type: Object, required: true }
  },
  computed: {
    look () {
      return this.$vuetify.theme.options.inputStyle.split('-')[0]
    },
    flat () {
      return this.$vuetify.theme.options.inputStyle.split('-').indexOf('flat') !== -1
    }
  },
  methods: {
    onInput (event) {
      this.$emit('input', { ...this.value, value: event })
    },
    markdown
  }
}
</script>
