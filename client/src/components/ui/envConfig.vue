<template>
  <v-row v-if="fields.length !== 0">
    <v-col
      v-for="field in fields"
      :key="field.name"
      cols="12"
    >
      <h4
        v-if="field.type === 'map' || field.type === 'list' || field.headline"
        v-text="getLabel(field)"
      />

      <component
        :is="field.component"
        v-if="field.type === 'custom'"
        :value="value[field.name] || field.default"
        @input="onInput(field.name, $event)"
      />
      <ui-map-input
        v-else-if="field.type === 'map'"
        :value="value[field.name] || field.default"
        :key-label="field.keyLabel ? $t(field.keyLabel) : undefined"
        :value-label="field.valueLabel ? $t(field.valueLabel) : undefined"
        @input="onInput(field.name, $event)"
      />
      <ui-list-input
        v-else-if="field.type === 'list'"
        :value="value[field.name] || field.default"
        @input="onInput(field.name, $event)"
      />
      <ui-input-suggestions
        v-else-if="field.options !== undefined"
        :type="field.type"
        :label="getLabel(field)"
        :items="field.options"
        :value="value[field.name] || field.default"
        @input="onInput(field.name, $event)"
      />
      <ui-input
        v-else
        :type="field.type"
        :label="getLabel(field)"
        :value="value[field.name] || field.default"
        @input="onInput(field.name, $event)"
      />
    </v-col>
  </v-row>
</template>

<script>
const envs = {
  standard: [],
  tty: [],
  docker: [
    {
      name: 'image',
      type: 'text',
      label: 'templates.DockerImage',
      default: 'pufferpanel/generic'
    },
    {
      name: 'networkMode',
      type: 'text',
      options: [
        'bridge',
        'host',
        'overlay',
        'macvlan',
        'none'
      ],
      label: 'env.docker.networkMode',
      default: 'host'
    },
    {
      name: 'networkName',
      type: 'text',
      default: ''
    },
    {
      name: 'bindings',
      type: 'map',
      keyLabel: 'env.docker.HostPath',
      valueLabel: 'env.docker.ContainerPath',
      default: {}
    },
    {
      name: 'portBindings',
      type: 'custom',
      component: 'ui-docker-port-bindings',
      headline: true,
      default: []
    }
  ]
}

export default {
  props: {
    value: {
      type: Object,
      validator: val => {
        return envs[val.type] !== undefined
      },
      required: true
    }
  },
  computed: {
    fields () {
      return envs[this.value.type]
    }
  },
  mounted () {
    const defaulted = {}
    this.fields.map(elem => {
      if (!this.value[elem.name]) {
        defaulted[elem.name] = elem.default
      }
    })
    this.$emit('input', { ...this.value, ...defaulted })
  },
  methods: {
    onInput (field, event) {
      this.$emit('input', { ...this.value, [field]: event })
    },
    getLabel (field) {
      return field.label ? this.$t(field.label) : this.$t(`env.${this.value.type}.${field.name}`)
    }
  }
}
</script>
