<script setup>
import { ref, computed, onUpdated } from 'vue'
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
if (Array.isArray(template.value.groups)) {
  template.value.groups.sort((a, b) => a.order > b.order ? 1 : -1)
}
const editOpen = ref(false)
const edit = ref({})
const addGroupOpen = ref(false)
const newGroup = ref({
  display: '',
  description: ''
})
const changeVarGroupOpen = ref(false)
const changeVarGroup = ref({})

const grouplessVars = computed(() => {
  if (Array.isArray(template.value.groups)) {
    return Object.keys(template.value.data).filter(varname => {
      return template.value.groups.map(g => g.variables).flat().indexOf(varname) === -1
    })
  } else {
    return Object.keys(template.value.data)
  }
})

const groupSelectOptions = computed(() => {
  return (template.value.groups || []).map(group => {
    return { value: group.order, label: group.display }
  })
})

function update() {
  emit('update:modelValue', JSON.stringify(template.value, undefined, 4))
}

function add(group) {
  edit.value = {
    name: '',
    display: '',
    desc: '',
    type: 'string',
    value: '',
    required: false,
    userEdit: false,
    options: {},
    group
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
  if (edit.value.group != undefined) {
    template.value.groups = template.value.groups.map(group => {
      if (group === edit.value.group) {
        group.variables.push(name)
      }
      return group
    })
    delete edit.value.group
  }
  template.value.data[name] = edit.value
  edit.value = {}
  update()
}

function remove(name) {
  delete template.value.data[name]
  template.value.groups = template.value.groups.map(group => {
    const variables = group.variables.filter(v => v !== name)
    return { ...group, variables }
  })
  update()
}

function isFirstGroup(order) {
  let min = Infinity
  template.value.groups.map(group => {
    if (group.order < min) min = group.order
  })
  return min === order
}

function getLastGroup() {
  let max = -Infinity
  template.value.groups.map(group => {
    if (group.order > max) max = group.order
  })
  return max === -Infinity ? 0 : max
}

function isLastGroup(order) {
  return getLastGroup(order) === order
}

function moveGroupUp(order) {
  template.value.groups = template.value.groups.map(group => {
    if (group.order === order) {
      group.order = order - 1
    } else if (group.order === (order - 1)) {
      group.order = order
    }
    return group
  }).sort((a, b) => a.order > b.order ? 1 : -1)
  update()
}

function moveGroupDown(order) {
  template.value.groups = template.value.groups.map(group => {
    if (group.order === order) {
      group.order = order + 1
    } else if (group.order === (order + 1)) {
      group.order = order
    }
    return group
  }).sort((a, b) => a.order > b.order ? 1 : -1)
  update()
}

function resetGroupAdd() {
  addGroupOpen.value = false
  newGroup.value.display = ''
  newGroup.value.description = ''
}

function addGroup() {
  if (!Array.isArray(template.value.groups))
    template.value.groups = []
  template.value.groups.push({
    variables: [],
    order: getLastGroup() + 1,
    display: newGroup.value.display,
    description: newGroup.value.description
  })
  resetGroupAdd()
  update()
}

function removeGroup(group) {
  template.value.groups = template.value.groups.filter(g => g !== group)
  update()
}

function changeGroup(variable, currentGroup) {
  changeVarGroup.value = {
    variable,
    currentGroup,
    selected: currentGroup ? { value: currentGroup.order, label: currentGroup.display } : {}
  }
  changeVarGroupOpen.value = true
}

function resetChangeVarGroup() {
  changeVarGroupOpen.value = false
  changeVarGroup.value = {}
}

function confirmChangeVarGroup() {
  template.value.groups = template.value.groups.map(group => {
    if (group === changeVarGroup.value.currentGroup) {
      group.variables = group.variables.filter(v => v != changeVarGroup.value.variable)
    } else if (group.order === changeVarGroup.value.selected.value && group.display === changeVarGroup.value.selected.label) {
      group.variables.push(changeVarGroup.value.variable)
    }
    return group
  })
  resetChangeVarGroup()
  update()
}

function nameValid() {
  return edit.value.name.trim() !== ''
}

onUpdated(() => {
  try {
    const u = JSON.parse(props.modelValue)
    // reserializing to avoid issues due to formatting
    if (JSON.stringify(template.value) !== JSON.stringify(u)) {
      template.value = u
      if (Array.isArray(template.value.groups)) {
        template.value.groups.sort((a, b) => a.order > b.order ? 1 : -1)
      }
    }
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
    <div v-if="template.groups && template.groups.length > 0">
      <div v-for="group in template.groups" :key="group.order">
        <div class="group-header">
          <div class="title">
            <h3 v-text="group.display" />
            <div class="hint" v-text="group.description" />
          </div>
          <btn variant="icon" @click="add(group)"><icon name="plus" /></btn>
          <btn :disabled="isFirstGroup(group.order)" variant="icon" @click="moveGroupUp(group.order)"><icon name="up" /></btn>
          <btn :disabled="isLastGroup(group.order)" variant="icon" @click="moveGroupDown(group.order)"><icon name="down" /></btn>
          <btn variant="icon" @click="removeGroup(group)"><icon name="close" /></btn>
        </div>
        <ul class="list">
          <li v-for="name in group.variables" :key="name" class="list-item clickable" @click="startEdit(name)">
            <div class="list-item-content">
              <span v-text="template.data[name].display" />
              <btn variant="icon" @click.stop="changeGroup(name, group)"><icon name="select-group" /></btn>
              <btn variant="icon" @click.stop="remove(name)"><icon name="remove" /></btn>
            </div>
          </li>
        </ul>
        <hr v-if="!isLastGroup(group.order) || grouplessVars.length > 0" />
      </div>
      <div v-if="grouplessVars.length > 0">
        <div class="group-header">
          <h3 class="title" v-text="t('templates.NoGroup')" />
          <btn variant="icon" @click="add()"><icon name="plus" /></btn>
        </div>
        <ul class="list">
          <li v-for="name in grouplessVars" :key="name" class="list-item clickable" @click="startEdit(name)">
            <div class="list-item-content">
              <span v-text="template.data[name].display" />
              <btn variant="icon" @click.stop="changeGroup(name)"><icon name="select-group" /></btn>
              <btn variant="icon" @click.stop="remove(name)"><icon name="remove" /></btn>
            </div>
          </li>
        </ul>
      </div>
    </div>
    <div v-else>
      <ul class="list">
        <li v-for="(item, name) in template.data" :key="name" class="list-item clickable" @click="startEdit(name)">
          <div class="list-item-content">
            <span v-text="item.display" />
            <btn variant="icon" @click.stop="remove(name)"><icon name="remove" /></btn>
          </div>
        </li>
      </ul>
      <btn variant="text" @click="add()"><icon name="plus" />{{ t('templates.AddVariable') }}</btn>
    </div>
    <btn variant="text" @click="addGroupOpen = true"><icon name="plus" />{{ t('templates.AddVariableGroup') }}</btn>

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

    <overlay v-model="addGroupOpen">
      <text-field v-model="newGroup.display" :label="t('templates.Display')" />
      <text-field v-model="newGroup.description" :label="t('templates.variables.Description')" />
      <div class="actions">
        <btn color="error" @click="resetGroupAdd()"><icon name="close" />{{ t('common.Cancel') }}</btn>
        <btn color="primary" :disabled="newGroup.display.trim() === ''" @click="addGroup()"><icon name="check" />{{ t('common.Apply') }}</btn>
      </div>
    </overlay>

    <overlay v-model="changeVarGroupOpen" class="select-var-group" :title="t('templates.SelectNewVarGroup')">
      <dropdown v-model="changeVarGroup.selected" :options="groupSelectOptions" :object="true" />

      <div class="actions">
        <btn color="error" @click="resetChangeVarGroup()"><icon name="close" />{{ t('common.Cancel') }}</btn>
        <btn color="primary" @click="confirmChangeVarGroup()"><icon name="check" />{{ t('common.Apply') }}</btn>
      </div>
    </overlay>
  </div>
</template>
