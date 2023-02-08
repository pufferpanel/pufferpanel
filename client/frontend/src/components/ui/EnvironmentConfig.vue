<script>
const fields = {
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
      hint: 'env.docker.BindingsHint',
      keyLabel: 'env.docker.HostPath',
      valueLabel: 'env.docker.ContainerPath',
      default: {}
    },
    {
      name: 'portBindings',
      type: 'portBindings',
      label: 'env.docker.portBindings',
      hint: 'env.docker.PortBindingsHint',
      default: []
    }
  ],
  // to not throw up when server creation cant select a valid env
  unsupported: []
}
</script>

<script setup>
import { onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import KeyValueInput from '@/components/ui/KeyValueInput.vue'
import PortMappingInput from '@/components/ui/PortMappingInput.vue'
import Suggestion from '@/components/ui/Suggestion.vue'
import TextField from '@/components/ui/TextField.vue'

const props = defineProps({
  noFieldsMessage: { type: String, default: () => undefined },
  modelValue: {
    type: Object,
    validator: val => fields[val.type] !== undefined,
    required: true
  }
})

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()

onMounted(() => {
  const defaults = {}
  fields[props.modelValue.type].map(field => {
    if (!props.modelValue[field.name]) {
      defaults[field.name] = field.default
    }
  })
  if (Object.keys(defaults).length > 0)
    emit('update:modelValue', { ...props.modelValue, ...defaults })
})

function onInput(field, event) {
  emit('update:modelValue', { ...props.modelValue, [field]: event })
}

function getLabel(field) {
  return field.label ? t(field.label) : t(`env.${props.modelValue.type}.${field.name}`)
}
</script>

<template>
  <div class="environment-config">
    <div v-if="noFieldsMessage && fields[modelValue.type].length === 0" v-text="noFieldsMessage" />
    <div v-for="field in fields[modelValue.type]" :key="field.name" class="field">
      <key-value-input v-if="field.type === 'map'" :model-value="modelValue[field.name] || field.default" :label="getLabel(field)" :hint="field.hint ? t(field.hint) : undefined" :key-label="t(field.keyLabel)" :value-label="t(field.valueLabel)" @update:modelValue="onInput(field.name, $event)" />
      <port-mapping-input v-else-if="field.type === 'portBindings'" :model-value="modelValue[field.name] || field.default" :label="getLabel(field)" :hint="field.hint ? t(field.hint) : undefined" @update:modelValue="onInput(field.name, $event)" />
      <suggestion v-else-if="field.type === 'text' && field.options" :model-value="modelValue[field.name] || field.default" :label="getLabel(field)" :options="field.options" :hint="field.hint ? t(field.hint) : undefined" @update:modelValue="onInput(field.name, $event)" />
      <text-field v-else-if="field.type === 'text'" :model-value="modelValue[field.name] || field.default" :label="getLabel(field)" :hint="field.hint ? t(field.hint) : undefined" @update:modelValue="onInput(field.name, $event)" />
      <span v-else v-text="`${field.type} not yet implemented`" />
    </div>
  </div>
</template>
