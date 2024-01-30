<script setup>
import { ref, onUpdated } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Dropdown from '@/components/ui/Dropdown.vue'
import Icon from '@/components/ui/Icon.vue'
import Overlay from '@/components/ui/Overlay.vue'
import EnvironmentConfig from '@/components/ui/EnvironmentConfig.vue'

const props = defineProps({
  modelValue: { type: String, required: true }
})

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()
const template = ref(JSON.parse(props.modelValue))
const edit = ref({})
const editing = ref(false)
const adding = ref(false)
const editIndex = ref(0)
const envs = [
  { value: 'host', label: t('env.host.name') },
  { value: 'docker', label: t('env.docker.name') }
]
const envDefaults = {
  host: { type: 'host' },
  docker: { type: 'docker', image: 'pufferpanel/generic' }
}

function update() {
  emit('update:modelValue', JSON.stringify(template.value, undefined, 4))
}

function editEnvChanged(newEnv) {
  edit.value = envDefaults[newEnv]
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

function unsupportedEnvs() {
  return envs.filter(env => {
    return template.value.supportedEnvironments.filter(e => e.type === env.value).length === 0
  })
}

function addEnv() {
  edit.value = envDefaults[unsupportedEnvs()[0].value]
  adding.value = true
  editing.value = true
}

function editEnv(index) {
  edit.value = template.value.supportedEnvironments[index]
  editIndex.value = index
  editing.value = true
}

function cancelEdit() {
  editing.value = false
  adding.value = false
  edit.value = {}
  editIndex.value = 0
}

function confirmEdit() {
  if (adding.value === true) {
    template.value.supportedEnvironments.push(edit.value)
  } else {
    template.value.supportedEnvironments[editIndex.value] = edit.value
  }
  cancelEdit()
  update()
}

function removeEnv(index) {
  const env = template.value.supportedEnvironments[index]
  template.value.supportedEnvironments.splice(index, 1)
  if (template.value.environment.type === env.type) template.value.environment = template.value.supportedEnvironments[0]
  update()
}
</script>

<template>
  <div class="environments">
    <ul class="list">
      <li v-for="(env, index) in template.supportedEnvironments" :key="env.type" class="list-item clickable">
        <div class="list-item-content" @click="editEnv(index)">
          <span v-text="t(`env.${env.type}.name`)" />
          <btn :disabled="template.supportedEnvironments.length < 2" variant="icon" @click.stop="removeEnv(index)"><icon name="remove" /></btn>
        </div>
      </li>
    </ul>
    <btn :disabled="unsupportedEnvs().length === 0" variant="text" @click="addEnv()"><icon name="plus" />{{ t('templates.AddEnvironment') }}</btn>
    <overlay v-model="editing">
      <dropdown v-if="adding" v-model="edit.type" :options="unsupportedEnvs()" :label="t('templates.Environment')" @update:modelValue="editEnvChanged" />
      <environment-config v-model="edit" :no-fields-message="t('env.NoEnvFields')" />
      <div class="actions">
        <btn color="error" @click="cancelEdit()"><icon name="close" />{{ t('common.Cancel') }}</btn>
        <btn color="primary" @click="confirmEdit()"><icon name="check" />{{ t('common.Confirm') }}</btn>
      </div>
    </overlay>
  </div>
</template>
