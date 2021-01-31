<template>
  <v-switch
    :id="id"
    :input-value="value"
    :dense="dense"
    :disabled="disabled"
    :label="label"
    :hide-details="hideDetails && !showHint"
    :persistent-hint="showHint"
    :hint="hintValue"
    :error-messages="errorMessages"
    :prepend-icon="icon"
    :append-icon="endIcon"
    :placeholder="placeholder"
    :required="required"
    :name="name"
    @change="$emit('input', $event ? true : false)"
    v-on="listeners"
  >
    <slot
      v-for="(_, slotName) in $slots"
      :slot="slotName"
      :name="slotName"
    />
  </v-switch>
</template>

<script>
export default {
  props: {
    dense: { type: Boolean, default: () => undefined },
    disabled: { type: Boolean, default: () => undefined },
    endIcon: { type: String, default: () => undefined },
    errorMessages: { type: String, default: () => undefined },
    hint: { type: String, default: () => undefined },
    icon: { type: String, default: () => undefined },
    id: { type: String, default: () => undefined },
    label: { type: String, default: () => undefined },
    name: { type: String, default: () => undefined },
    placeholder: { type: String, default: () => undefined },
    required: { type: Boolean, default: () => undefined },
    value: { type: undefined, default: () => '', required: true }
  },
  computed: {
    listeners () {
      const { input, ...listeners } = this.$listeners
      return listeners
    },
    hideDetails () {
      const hasError =
        this.errorMessages !== undefined && this.errorMessages !== ''
      return !hasError
    },
    showHint () {
      const hasHint = this.hint !== undefined && this.hint !== ''
      const hasMsgSlot =
        this.$slots.message !== undefined &&
        this.$slots.message !== '' &&
        this.$slots.message !== []
      return hasHint || hasMsgSlot
    },
    hintValue () {
      // set hint to '_' if only the slot has content to force vuetify to
      // display the hint without needing to double define it everywhere
      return this.hint ? this.hint : this.$slots.message ? '_' : undefined
    }
  }
}
</script>
