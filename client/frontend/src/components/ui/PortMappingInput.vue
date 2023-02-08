<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from './Btn.vue'
import Icon from './Icon.vue'
import Suggestion from './Suggestion.vue'
import TextField from './TextField.vue'

const props = defineProps({
  // default to random id as labels need target ids to exist exactly once
  label: { type: String, default: () => undefined },
  addLabel: { type: String, default: () => undefined },
  hint: { type: String, default: () => undefined },
  error: { type: String, default: () => undefined },
  modelValue: { type: Array, default: () => [] }
})

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()

const entries = ref([])

onMounted(() => {
  for (let entry in props.modelValue) {
    const [head, protocol] = entry.split('/')
    const [host, outsidePort, insidePort] = head.split(':')
    entries.value.push({ host, outsidePort, insidePort, protocol })
  }
})

function emitUpdate() {
  const result = []
  entries
    .value
    .filter(entry => entry.host && entry.outsidePort && entry.insidePort && entry.protocol)
    .map(entry => result.push(`${entry.host}:${entry.outsidePort}:${entry.insidePort}/${entry.protocol}`))
  emit('update:modelValue', result)
}

function addEntry() {
  entries.value.push({ host: '0.0.0.0', outsidePort: '', insidePort: '', protocol: 'tcp' })
  emitUpdate()
}

function onInput(entry, field, event) {
  entry[field] = event
  emitUpdate()
}

function removeEntry(item) {
  entries.value = entries.value.filter(entry => entry !== item)
  emitUpdate()
}
</script>

<template>
  <div class="port-mapping-input">
    <div v-if="label" class="label" v-text="label" />
    <div v-if="error" class="error" v-text="error" />
    <div v-else-if="hint" class="hint" v-text="hint" />
    <div v-for="(entry, index) in entries" :key="index" class="entry">
      <div class="fields">
        <text-field :model-value="entry.host" :label="t('common.Host')" @update:modelValue="onInput(entry, 'host', $event)" />
        <text-field :model-value="entry.outsidePort" :label="t('env.docker.OutsidePort')" @update:modelValue="onInput(entry, 'outsidePort', $event)" />
        <text-field :model-value="entry.insidePort" :label="t('env.docker.InsidePort')" @update:modelValue="onInput(entry, 'insidePort', $event)" />
        <suggestion :model-value="entry.protocol" :label="t('common.Protocol')" :options="['tcp', 'udp']" @update:modelValue="onInput(entry, 'protocol', $event)" />
      </div>
      <btn variant="icon" @click="removeEntry(entry)"><icon name="remove" /></btn>
    </div>
    <btn variant="text" @click="addEntry()"><icon name="plus" />{{ addLabel || t('common.Add') }}</btn>
  </div>
</template>
