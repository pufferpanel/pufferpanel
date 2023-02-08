<script setup>
import { ref, onUpdated } from 'vue'
import { useI18n } from 'vue-i18n'
import { generateOperatorLabel } from '@/utils/operators.js'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import Operator from './Operator.vue'
import Overlay from '@/components/ui/Overlay.vue'

const props = defineProps({
  addLabel: { type: String, default: () => undefined },
  modelValue: { type: Array, required: true }
})

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()

const edit = ref({})
const editIndex = ref(0)
const editOpen = ref(false)
const model = ref([...props.modelValue])

function update() {
  emit('update:modelValue', model.value)
}

function add() {
  editIndex.value = model.value.length
  edit.value = { type: 'command', commands: [] }
  editOpen.value = true
}

function startEdit(i) {
  editIndex.value = i
  edit.value = model.value[i]
  editOpen.value = true
}

function cancelEdit() {
  editOpen.value = false
  editIndex.value = 0
  edit.value = {}
}

function confirmEdit() {
  editOpen.value = false
  model.value[editIndex.value] = edit.value
  editIndex.value = 0
  edit.value = {}
  update()
}

function swap(i1, i2) {
  const x = model.value[i1]
  model.value[i1] = model.value[i2]
  model.value[i2] = x
  update()
}

function remove(i) {
  model.value.splice(i, 1)
  update()
}
</script>

<template>
  <div class="operators">
    <ul class="list">
      <li v-for="(step, index) in model" :key="index" class="list-item clickable">
        <div class="list-item-content" @click="startEdit(index)">
          <span v-text="generateOperatorLabel(t, step)" />
          <btn :disabled="index === 0" variant="icon" @click.stop="swap(index, index-1)"><icon name="up" /></btn>
          <btn :disabled="index === model.length - 1" variant="icon" @click.stop="swap(index, index+1)"><icon name="down" /></btn>
          <btn variant="icon" @click.stop="remove(index)"><icon name="remove" /></btn>
        </div>
      </li>
    </ul>
    <btn variant="text" @click="add()"><icon name="plus" />{{ addLabel || t('common.Add') }}</btn>
    <overlay v-model="editOpen">
      <operator v-model="edit" />
      <div class="actions">
        <btn color="error" @click="cancelEdit()"><icon name="close" />{{ t('common.Cancel') }}</btn>
        <btn color="primary" @click="confirmEdit()"><icon name="save" />{{ t('common.Save') }}</btn>
      </div>
    </overlay>
  </div>
</template>
