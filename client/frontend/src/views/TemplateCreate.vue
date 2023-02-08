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

let unmodified = `{
    "name": "",
    "display": "",
    "type": "",
    "data": {},
    "install": [],
    "run": {
        "command": "",
        "stop": "",
        "workingDirectory": "",
        "pre": [],
        "post": [],
        "environmentVars": {}
    },
    "environment": {
        "type": "standard"
    },
    "supportedEnvironments": [
        {
            "type": "standard"
        }
    ]
}`
const template = ref(unmodified)
const valid = ref({
  general: false,
  run: false
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

onMounted(() => {
  if (route.query.copy) {
    const t = sessionStorage.getItem('copiedTemplate')
    if (t) {
      template.value = t
      sessionStorage.removeItem('copiedTemplate')
      setTimeout(() => {
        unmodified = template.value
      }, 50)
    }
  }
})

async function save() {
  if (!canSave()) return
  const name = JSON.parse(template.value).name

  const exists = await api.template.exists('local', name)
  if (!exists) {
    await api.template.save(name, template.value)
    toast.success(t('templates.Saved'))
    router.push({ name: 'TemplateView', params: { id: name } })
  } else {
    events.emit(
      'confirm',
      t('templates.ConfirmOverwrite'),
      {
        text: t('templates.Overwrite'),
        icon: 'check',
        action: async () => {
          await api.template.save(name, template.value)
          toast.success(t('templates.Saved'))
          router.push({ name: 'TemplateView', params: { id: name } })
        }
      }
    )
  }
}

function canSave() {
  return Object.values(valid.value).filter(e => e === false).length === 0
}
</script>

<template>
  <div class="templatecreate">
    <div>
      <tabs anchors>
        <tab id="general" :title="t('templates.General')" icon="general" hotkey="t g">
          <general v-model="template" id-editable @valid="valid.general = $event" />
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
      <btn color="primary" :disabled="!canSave()" @click="save()"><icon name="save" />{{ t('common.Save') }}</btn>
    </div>
  </div>
</template>
