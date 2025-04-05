<script setup>
import { ref, onMounted } from 'vue'

const props = defineProps({
  digits: { type: Number, default: () => 6 },
  disabled: { type: Boolean, default: () => false }
})
const emit = defineEmits(['update:modelValue', 'complete'])

const inputs = ref([])

function onInput(n, e) {
  if (e.data === null) {
    // data is null when deleting chars
    emit('update:modelValue', inputs.value.map(e => e.value).join(''))
    return
  }

  if (e.data.length === props.digits && /^\d+$/.test(e.data)) {
    // full code was pasted
    [...e.data].map((c, i) => {
      if (inputs.value[i]) inputs.value[i].value = c
    })
    if (inputs.value[n]) inputs.value[n].blur()
    emit('update:modelValue', inputs.value.map(e => e.value).join(''))
    emit('complete', inputs.value.map(e => e.value).join(''))
    return
  }

  if (e.data.length > 1 && /^\d+$/.test(e.data)) {
    // numbers were pasted, but it's not the full code, fill inputs from current one onwards
    [...e.data].map((c, i) => {
      if (inputs.value[n + i]) inputs.value[n + i].value = c
    })
    emit('update:modelValue', inputs.value.map(e => e.value).join(''))
    if (inputs.value.map(e => e.value).join('').length === props.digits) {
      if (inputs.value[n]) inputs.value[n].blur()
      emit('complete', inputs.value.map(e => e.value).join(''))
    } else {
      if (inputs.value[n + e.data.length]) inputs.value[n + e.data.length].focus()
    }
    return
  }

  if (/^\d$/.test(e.data)) {
    // got a single digit input, update current field and jump to next one
    e.target.value = e.data
    emit('update:modelValue', inputs.value.map(e => e.value).join(''))
    if (n === props.digits - 1) {
      if (inputs.value[n]) inputs.value[n].blur()
      emit('complete', inputs.value.map(e => e.value).join(''))
    } else if (inputs.value[n + 1]) {
      inputs.value[n + 1].focus()
    }
    return
  }

  e.target.value = [...e.target.value].filter(c => /\d/.test(c))[0] || ''
}

onMounted(() => {
  if (inputs.value[0]) inputs.value[0].focus()
})

function pushInput(input) {
  if (inputs.value.indexOf(input) === -1) inputs.value.push(input)
}
</script>

<template>
  <div class="otp-input">
    <input v-for="n in Array(props.digits).keys()" :key="n" :ref="pushInput" :disabled="props.disabled" type="text" @input="onInput(n, $event)" />
  </div>
</template>
