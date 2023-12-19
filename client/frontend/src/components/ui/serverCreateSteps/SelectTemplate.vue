<script setup>
import { ref, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import Overlay from '@/components/ui/Overlay.vue'
import markdown from '@/utils/markdown.js'

const { t } = useI18n()
const api = inject('api')
const emit = defineEmits(['selected', 'back'])
const templatesByRepo = ref([])
const showing = ref(false)
const currentTemplate = ref({})

const props = defineProps({
  arch: { type: String, required: true },
  env: { type: String, required: true },
  os: { type: String, required: true }
})

function templateEnvMatches(template) {
  if (!Array.isArray(template.supportedEnvironments)) {
    if (!template.environment) return false // neither supported nor default env defined, ignore
    template.supportedEnvironments = [template.environment]
  }
  if (template.supportedEnvironments.filter(e => e.type === props.env).length > 0) return true
  return false
}

function templateOsMatches(template) {
  if (!template.requirements || !template.requirements.os) return true
  return template.requirements.os === props.os
}

function templateArchMatches(template) {
  if (!template.requirements || !template.requirements.arch) return true
  return template.requirements.arch === props.arch
}

async function load() {
  const repos = await api.template.listAllTemplates()
  const res = []
  Object.keys(repos).sort((a, b) => repos[a].id > repos[b].id).map(repo => {
    if (repos[repo].templates.length === 0) return
    const templates = repos[repo].templates.filter(template => {
      return templateEnvMatches(template) &&
        templateOsMatches(template) &&
        templateArchMatches(template)
    })
    res.push({ ...repos[repo], templates })
  })
  templatesByRepo.value = res
}

onMounted(async () => {
  load()
})

async function show(repo, template) {
  currentTemplate.value = await api.template.get(repo, template)
  if (currentTemplate.value.readme) {
    showing.value = true
  } else {
    // no readme, skip readme popup
    emit('selected', currentTemplate.value)
  }
}

function choice(confirm) {
  showing.value = false
  if (confirm) emit('selected', currentTemplate.value)
}
</script>

<template>
  <div class="select-template">
    <h2 v-text="t('servers.SelectTemplate')" />
    <div v-for="repo in templatesByRepo" :key="repo.id" class="list">
      <h3 class="list-header" v-text="repo.name" />
      <div v-for="template in repo.templates" :key="template.name" class="list-item template" @click="show(repo.id, template.name)">
        <span class="title" v-text="template.display" />
      </div>
    </div>
    <btn color="error" @click="emit('back')"><icon name="back" />{{ t('common.Back') }}</btn>

    <overlay v-model="showing" :title="currentTemplate.display" closable>
      <!-- eslint-disable-next-line vue/no-v-html -->
      <div dir="ltr" class="readme" v-html="markdown(currentTemplate.readme)" />
      <div class="actions">
        <btn color="error" @click="choice(false)"><icon name="close" />{{ t('common.Cancel') }}</btn>
        <btn color="primary" @click="choice(true)"><icon name="check" />{{ t('servers.SelectThisTemplate') }}</btn>
      </div>
    </overlay>
  </div>
</template>
