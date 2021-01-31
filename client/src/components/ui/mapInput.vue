<template>
  <div>
    <v-row
      v-for="(entry, i) in inputValue"
      :key="i"
    >
      <v-col
        cols="12"
        md="6"
      >
        <ui-input
          :ref="'key-' + entry.key"
          :value="entry.key"
          :label="keyLabel"
          @input="onKeyInput(entry.key, $event)"
        />
      </v-col>
      <v-col
        cols="12"
        md="6"
      >
        <ui-input
          :value="entry.value"
          :label="valueLabel"
          icon-behind="mdi-close"
          @click:append-outer="remove(entry.key)"
          @input="onValueInput(entry.key, $event)"
        />
      </v-col>
    </v-row>
    <v-row v-if="add">
      <v-col
        cols="12"
        md="6"
      >
        <ui-input
          value=""
          autofocus
          :label="keyLabel"
          @input="addKey($event)"
        />
      </v-col>
      <v-col
        cols="12"
        md="6"
      >
        <ui-input
          v-model="addValue"
          :label="valueLabel"
        />
      </v-col>
    </v-row>
    <v-row v-else>
      <v-col cols="12">
        <v-btn
          text
          block
          @click="add = true"
          v-text="addLabel ? addLabel : $t('common.Add')"
        />
      </v-col>
    </v-row>
  </div>
</template>

<script>
export default {
  props: {
    keyLabel: { type: String, default: () => undefined },
    valueLabel: { type: String, default: () => undefined },
    addLabel: { type: String, default: () => undefined },
    value: { type: Object, required: true }
  },
  data () {
    return {
      add: false,
      addValue: ''
    }
  },
  computed: {
    inputValue () {
      return Object.keys(this.value).map((elem, i) => {
        return { key: elem, value: this.value[elem], position: i }
      })
    }
  },
  methods: {
    onKeyInput (key, event) {
      const changed = {}
      Object.keys(this.value).map(elem => {
        if (elem === key) {
          changed[event] = this.value[elem]
        } else {
          changed[elem] = this.value[elem]
        }
      })
      this.$emit('input', changed)
    },
    onValueInput (key, event) {
      const changed = { ...this.value }
      changed[key] = event
      this.$emit('input', changed)
    },
    remove (key) {
      const changed = { ...this.value }
      delete changed[key]
      this.$emit('input', changed)
    },
    addKey (key) {
      const changed = { ...this.value }
      changed[key] = this.addValue
      this.addValue = ''
      this.add = false
      this.$emit('input', changed)
      this.$nextTick(() => {
        this.$refs['key-' + key][0].focus()
      })
    }
  }
}
</script>
