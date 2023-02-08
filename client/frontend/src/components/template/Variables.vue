<script setup>
import { ref, onUpdated } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Dropdown from '@/components/ui/Dropdown.vue'
import Icon from '@/components/ui/Icon.vue'
import KeyValueInput from '@/components/ui/KeyValueInput.vue'
import Overlay from '@/components/ui/Overlay.vue'
import TextField from '@/components/ui/TextField.vue'
import Toggle from '@/components/ui/Toggle.vue'

const props = defineProps({
  modelValue: { type: String, required: true }
})

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()

const template = ref(JSON.parse(props.modelValue))
const editOpen = ref(false)
const edit = ref({})

function update() {
  emit('update:modelValue', JSON.stringify(template.value, undefined, 4))
}

function add() {
  edit.value = {
    name: '',
    display: '',
    desc: '',
    type: 'string',
    value: '',
    required: false,
    userEdit: false,
    options: {}
  }
  editOpen.value = true
}

function startEdit(name) {
  edit.value = template.value.data[name]
  edit.value.name = name
  edit.value.oldName = name
  if (!edit.value.display) edit.value.display = ''
  if (!edit.value.desc) edit.value.desc = ''
  if (!edit.value.type) edit.value.type = 'string'
  if (!edit.value.value) edit.value.value = ''
  if (!edit.value.required) edit.value.required = false
  if (!edit.value.userEdit) edit.value.userEdit = false
  if (!edit.value.options) {
    edit.value.options = {}
  } else {
    const o = {}
    edit.value.options.map(e => {
      o[e.value] = e.display
    })
    edit.value.options = o
  }
  editOpen.value = true
}

function cancelEdit() {
  editOpen.value = false
  edit.value = {}
}

function confirmEdit() {
  editOpen.value = false
  if (edit.value.oldName && edit.value.oldName !== edit.value.name) {
    delete template.value.data[edit.value.oldName]
  }
  const name = edit.value.name
  delete edit.value.name
  if (edit.value.oldName) delete edit.value.oldName
  if (edit.value.type === 'boolean' || Object.keys(edit.value.options).length === 0) {
    delete edit.value.options
  } else {
    const o = []
    Object.keys(edit.value.options).map(k => {
      o.push({ value: k, display: edit.value.options[k] })
    })
    edit.value.options = o
  }
  template.value.data[name] = edit.value
  edit.value = {}
  update()
}

function remove(name) {
  delete template.value.data[name]
  update()
}

function nameValid() {
  return edit.value.name.trim() !== ''
}

onUpdated(() => {
  try {
    const u = JSON.parse(props.modelValue)
    // reserializing to avoid issues due to formatting
    if (JSON.stringify(template.value) !== JSON.stringify(u))
      template.value = u
  } catch {
    // expected failure caused by json editor producing invalid json during modification
  }
})

const types = [
  { value: 'string', label: t('templates.variables.types.String') },
  { value: 'boolean', label: t('templates.variables.types.Boolean') },
  { value: 'integer', label: t('templates.variables.types.Number') },
  { value: 'options', label: t('templates.variables.types.Options') }
]

const supportsOptions = {
  string: true,
  options: true
}
</script>

<template>
  <div class="variables">
    <div class="hint" v-text="t('templates.description.Variables')" />
    <ul class="list">
      <li v-for="(item, name) in template.data" :key="name" class="list-item clickable" @click="startEdit(name)">
        <div class="list-item-content">
          <span v-text="item.display" />
          <btn variant="icon" @click.stop="remove(name)"><icon name="remove" /></btn>
        </div>
      </li>
    </ul>
    <btn variant="text" @click="add()"><icon name="plus" />{{ t('templates.AddVariable') }}</btn>

    <overlay v-model="editOpen">
      <text-field v-model="edit.name" :label="t('common.Name')" />
      <text-field v-model="edit.display" :label="t('templates.Display')" />
      <text-field v-model="edit.desc" :label="t('templates.variables.Description')" />
      <dropdown v-model="edit.type" :label="t('templates.variables.Type')" :options="types" />
      <text-field v-model="edit.value" :label="t('templates.variables.Value')" />
      <toggle v-model="edit.required" :label="t('templates.variables.Required')" />
      <toggle v-model="edit.userEdit" :label="t('templates.variables.UserEdit')" />
      <key-value-input v-if="supportsOptions[edit.type]" v-model="edit.options" :label="t('templates.variables.Options')" />

      <div class="actions">
        <btn color="error" @click="cancelEdit()"><icon name="close" />{{ t('common.Cancel') }}</btn>
        <btn color="primary" :disabled="!nameValid()" @click="confirmEdit()"><icon name="check" />{{ t('common.Apply') }}</btn>
      </div>
    </overlay>
  </div>
</template>
