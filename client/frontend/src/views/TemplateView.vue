<script setup>
import { ref, inject, onMounted } from 'vue'
import { useRoute, useRouter, onBeforeRouteLeave } from 'vue-router'
import { useI18n } from 'vue-i18n'

import General from '@/components/template/General.vue'
import Variables from '@/components/template/Variables.vue'
import Install from '@/components/template/Install.vue'
import Hooks from '@/components/template/Hooks.vue'
import RunConfig from '@/components/template/RunConfig.vue'
import Environment from '@/components/template/Environment.vue'

import Ace from '@/components/ui/Ace.vue'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import Loader from '@/components/ui/Loader.vue'
import Tabs from '@/components/ui/Tabs.vue'
import Tab from '@/components/ui/Tab.vue'

const api = inject('api')
const toast = inject('toast')
const events = inject('events')
const { t } = useI18n()
const route = useRoute()
const router = useRouter()

let unmodified = null
const template = ref(null)
const valid = ref({
  general: true,
  run: true
})

onBeforeRouteLeave((to, from) => {
  if (unmodified === template.value) {
    return true
  } else {
    return new Promise((resolve, reject) => {
      events.emit(
        'confirm',
        t('common.ConfirmLeave'),
        {
          text: t('common.Discard'),
          icon: 'remove',
          color: 'error',
          action: () => { resolve(true) }
        },
        {
          color: 'primary',
          action: () => { resolve(false) }
        }
      )
    })
  }
})

onMounted(async () => {
  const res = await api.template.get(route.params.repo, route.params.id)
  delete res.readme
  template.value = JSON.stringify(res, undefined, 4)
  setTimeout(() => {
    unmodified = template.value
  }, 50)
})

async function deleteTemplate() {
  if (!canDelete()) return
  events.emit(
    'confirm',
    t('templates.ConfirmDelete', { name: JSON.parse(template.value).display }),
    {
      text: t('templates.Delete'),
      icon: 'remove',
      color: 'error',
      action: async () => {
        await api.template.delete(route.params.id)
        toast.success(t('templates.Deleted'))
        router.push({ name: 'TemplateList' })
      }
    },
    {
      color: 'primary'
    }
  )
}

function canDelete() {
  return route.params.repo === 'local'
}

async function save() {
  if (!canSave()) return
  await api.template.save(route.params.id, template.value)
  toast.success(t('templates.Saved'))
}

function canSave() {
  if (route.params.repo !== 'local') return false
  return Object.values(valid.value).filter(e => e === false).length === 0
}

function createLocalCopy() {
  const t = JSON.parse(template.value)
  delete t.name
  sessionStorage.setItem('copiedTemplate', JSON.stringify(t))
  router.push({ name: 'TemplateCreate', query: { 'copy': true } })
}
</script>

<template>
  <div class="templateview">
    <loader v-if="!template" />
    <div v-else>
      <div v-if="route.params.repo !== 'local'" class="alert info">
        <span v-text="t('templates.EditLocalOnly')" />
        <btn color="primary" @click="createLocalCopy()"><icon name="copy" />{{ t('templates.CreateLocalCopy') }}</btn>
      </div>
      <tabs anchors>
        <tab id="general" :title="t('templates.General')" icon="general" hotkey="t g">
          <general v-model="template" @valid="valid.general = $event" />
        </tab>
        <tab id="variables" :title="t('templates.Variables')" icon="variables" hotkey="t v">
          <variables v-model="template" />
        </tab>
        <tab id="install" :title="t('templates.Install')" icon="install" hotkey="t i">
          <install v-model="template" />
        </tab>
        <tab id="run" :title="t('templates.RunConfig')" icon="start" hotkey="t r">
          <run-config v-model="template" @valid="valid.run = $event" />
        </tab>
        <tab id="hooks" :title="t('templates.Hooks')" icon="hooks" hotkey="t h">
          <hooks v-model="template" />
        </tab>
        <tab id="environment" :title="t('templates.Environment')" icon="environment" hotkey="t e">
          <environment v-model="template" />
        </tab>
        <tab id="json" :title="t('templates.Json')" icon="json" hotkey="t j">
          <ace id="template-json" v-model="template" class="template-json-editor" mode="json" />
        </tab>
      </tabs>
      <div v-if="route.params.repo === 'local'" class="actions">
        <btn color="error" :disabled="!canDelete()" @click="deleteTemplate()"><icon name="remove" />{{ t('templates.Delete') }}</btn>
        <btn color="primary" :disabled="!canSave()" @click="save()"><icon name="save" />{{ t('templates.Save') }}</btn>
      </div>
    </div>
  </div>
</template>
