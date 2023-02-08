<script setup>
import { ref, inject, onUpdated } from 'vue'
import { useI18n } from 'vue-i18n'
import TextField from '@/components/ui/TextField.vue'

const props = defineProps({
  idEditable: { type: Boolean, default: () => false },
  modelValue: { type: String, required: true }
})

const emit = defineEmits(['update:modelValue', 'valid'])

const api = inject('api')
const { t } = useI18n()

const template = ref(JSON.parse(props.modelValue))
let nameDebounce = null
const nameInvalid = ref(false)
const nameUnique = ref(true)
const displayInvalid = ref(false)
const typeInvalid = ref(false)

function update() {
  emit('update:modelValue', JSON.stringify(template.value, undefined, 4))
}

function nameUpdate() {
  clearTimeout(nameDebounce)
  nameDebounce = setTimeout(async () => {
    if (template.value.name && template.value.name.trim() !== '') {
      const exists = await api.template.exists('local', template.value.name)
      nameUnique.value = !exists
    } else {
      nameUnique.value = true
    }
    validate()
  }, 500)
  update()
}

function nameError() {
  if (nameInvalid.value) {
    return t('templates.errors.NameInvalid')
  } else if (!nameUnique.value) {
    return t('templates.NameNotUnique')
  } else {
    return null
  }
}

function validate() {
  nameInvalid.value = false
  displayInvalid.value = false
  typeInvalid.value = false

  if (!template.value.name || template.value.name.trim() === '') {
    nameInvalid.value = true
  }

  if (!template.value.display || template.value.display.trim() === '') {
    displayInvalid.value = true
  }

  if (!template.value.type || template.value.type.trim() === '') {
    typeInvalid.value = true
  }

  emit('valid', !nameInvalid.value && !displayInvalid.value && !typeInvalid.value)
}

onUpdated(() => {
  try {
    const u = JSON.parse(props.modelValue)

    let needsNameUpdate = false
    if (u.name && template.value.name && template.value.name !== u.name) {
      needsNameUpdate = true
    }

    // reserializing to avoid issues due to formatting
    if (JSON.stringify(template.value) !== JSON.stringify(u)) {
      template.value = u
      validate()
      if (needsNameUpdate) nameUpdate()
    }
  } catch {
    // expected failure caused by json editor producing invalid json during modification
  }
})
</script>

<template>
  <div>
    <text-field v-model="template.name" :disabled="!idEditable" :label="t('common.Name')" :hint="idEditable ? t('templates.description.Name') : undefined" :error="nameError()" @update:modelValue="nameUpdate" @blur="validate()" />
    <text-field v-model="template.display" :label="t('templates.Display')" :hint="t('templates.description.Display')" :error="displayInvalid ? t('templates.errors.DisplayInvalid') : undefined" @update:modelValue="update" @blur="validate()" />
    <text-field v-model="template.type" :label="t('templates.Type')" :hint="t('templates.description.Type')" :error="typeInvalid ? t('templates.errors.TypeInvalid') : undefined" @update:modelValue="update" @blur="validate()" />
  </div>
</template>
