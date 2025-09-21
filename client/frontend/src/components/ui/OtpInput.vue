<script setup>
import { ref, onMounted } from 'vue'

const props = defineProps({
  digits: { type: Number, default: () => 6 },
  disabled: { type: Boolean, default: () => false }
})
const emit = defineEmits(['update:modelValue', 'complete'])

const inputs = ref([])

function onInput(n, e) {
  if (!e.data && e.target.value.length === props.digits) {
    // data isn't set but there's a full code in the input when a password manager directly sets field value
    [...e.target.value].map((c, i) => {
      if (inputs.value[i]) inputs.value[i].value = c
    })
    if (inputs.value[n]) inputs.value[n].blur()
    emit('update:modelValue', inputs.value.map(e => e.value).join(''))
    emit('complete', inputs.value.map(e => e.value).join(''))
    return
  } else if (e.data === null) {
    // data is null when deleting chars
    emit('update:modelValue', inputs.value.map(e => e.value).join(''))
    return
  } else if (e.data === undefined && e.target.value.length === 1) {
    // likely a password manager detecting multiple fields and direct settings each field to a single digit of the code
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

function onBackspace(n, e) {
  if (e.target.value === '' && inputs.value[n - 1]) {
    inputs.value[n - 1].value = ''
    inputs.value[n - 1].focus()
  }
}

function clearAll() {
  inputs.value.map(e => e.value = '')
  if (inputs.value[0]) inputs.value[0].focus()
}
</script>

<template>
  <div class="otp-input">
    <input
      v-for="n in Array(props.digits).keys()"
      :key="n"
      :ref="pushInput"
      :name="n === 0 ? 'totp' : 'ingnore'"
      :autofill="n === 0 ? 'one-time-code' : 'none'"
      :disabled="props.disabled"
      type="number"
      @input="onInput(n, $event)"
      @keydown.exact.backspace="onBackspace(n, $event)"
      @keydown.ctrl.backspace="clearAll"
    />
    <!--
      apparently proton pass tries to handle multi field totp inputs on its end, however it sadly seems to have
      a heart attach mid-attempt and just refuses to fill in the last digit of the code, adding another input
      seems to confuse it into just throwing the full code into the first input though, which we handle for
      other password managers as well, that's why this invisible to humans input exists
    -->
    <input type="number" name="fix_proton_pass" autofill="none" style="position:absolute;left:-99999px;width:1px" />
  </div>
</template>
