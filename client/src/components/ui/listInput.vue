<template>
  <div>
    <v-row
      v-for="(entry, i) in value"
      :key="i"
    >
      <v-col cols="12">
        <ui-input
          :value="entry"
          icon-behind="mdi-close"
          @click:append-outer="remove(i)"
          @input="onInput(i, $event)"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12">
        <v-btn
          text
          block
          @click="add()"
          v-text="addLabel ? addLabel : $t('common.Add')"
        />
      </v-col>
    </v-row>
  </div>
</template>

<script>
export default {
  props: {
    addLabel: { type: String, default: () => undefined },
    value: { type: Array, required: true }
  },
  methods: {
    onInput (index, event) {
      const changed = [...this.value]
      changed[index] = event
      this.$emit('input', changed)
    },
    remove (index) {
      const changed = [...this.value]
      changed.splice(index, 1)
      this.$emit('input', changed)
    },
    add () {
      const changed = [...this.value]
      changed.push('')
      this.$emit('input', changed)
    }
  }
}
</script>
