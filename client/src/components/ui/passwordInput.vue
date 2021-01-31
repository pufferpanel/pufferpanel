<template>
  <ui-input
    :id="id"
    v-model="internalValue"
    :disabled="disabled"
    :label="label"
    :error-messages="errorMessages"
    icon="mdi-lock"
    :end-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
    :name="name"
    :type="showPassword ? 'text' : 'password'"
    @click:append="showPassword = !showPassword"
    v-on="listeners"
  />
</template>

<script>
export default {
  props: {
    disabled: { type: Boolean, default: () => false },
    errorMessages: { type: String, default: () => undefined },
    id: { type: String, default: () => undefined },
    label: { type: String, default: () => undefined },
    name: { type: String, default: () => undefined },
    value: { type: String, default: () => '' }
  },
  data () {
    return {
      internalValue: this.value,
      showPassword: false
    }
  },
  computed: {
    listeners () {
      const listeners = { ...this.$listeners }
      delete listeners.input
      delete listeners['click:append']
      return listeners
    }
  },
  watch: {
    internalValue (val) {
      this.$emit('input', this.internalValue)
    }
  }
}
</script>
