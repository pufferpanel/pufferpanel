<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import Overlay from '@/components/ui/Overlay.vue'
import TextField from '@/components/ui/TextField.vue'
import Status from './Status.vue'
import Controls from './Controls.vue'

const props = defineProps({
  server: { type: Object, required: true }
})

const { t } = useI18n()
const edit = ref(false)
const name = ref(props.server.name)

async function updateName() {
  await props.server.updateName(name.value)
  edit.value = false
}
</script>

<template>
  <h1 class="server-header">
    <status :server="server" />
    <span class="name">{{ server.name }}<btn v-if="server.permissions.editServerData" variant="icon" @click="edit = !edit"><icon name="edit" /></btn></span>
    <controls :server="server" />
  </h1>
  <overlay v-model="edit" :title="t('common.Name')" closable>
    <text-field v-model="name" />
    <btn color="primary" @click="updateName()"><icon name="save" />{{ t('common.Save') }}</btn>
  </overlay>
</template>
