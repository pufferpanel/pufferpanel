<script setup>
import { ref, inject, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import Loader from '@/components/ui/Loader.vue'
import Overlay from '@/components/ui/Overlay.vue'
import TextField from '@/components/ui/TextField.vue'

const { t } = useI18n()
const api = inject('api')
const events = inject('events')
const templatesLoaded = ref(false)
const templatesByRepo = ref([])
const firstEntry = ref(null)
const addingRepo = ref(false)
const currentRepo = ref({name: '', url: '', branch: ''})

onMounted(async () => {
  loadTemplates()
})

async function loadTemplates() {
  templatesLoaded.value = false
  const templates = await api.template.listAllTemplates()
  templatesByRepo.value = templates.sort((a, b) => a.id > b.id)
  templatesLoaded.value = true
}

function setFirstEntry(ref) {
  if (!firstEntry.value) firstEntry.value = ref
}

function focusList() {
  firstEntry.value.$el.focus()
}

function removeRepo(repo) {
  events.emit(
    'confirm',
    t('templates.ConfirmDeleteRepo', { name: repo.name }),
    {
      text: t('templates.DeleteRepo'),
      icon: 'remove',
      color: 'error',
      action: async () => {
        await api.template.deleteRepo(repo.id)
        await loadTemplates()
      }
    },
    {
      color: 'primary'
    }
  )
}

async function addRepo() {
  await api.template.addRepo({...currentRepo.value, isLocal: false, id: 2})
  resetAddRepo()
  await loadTemplates()
}

function resetAddRepo() {
  currentRepo.value = {name: '', url: '', branch: ''}
  addingRepo.value = false
}

function canAddRepo() {
  if (!currentRepo.value.name || currentRepo.value.name === '') return false
  if (!currentRepo.value.url || currentRepo.value.url === '') return false
  return true
}
</script>

<template>
  <div class="templatelist">
    <h1 v-text="t('templates.Templates')" />
    <div v-hotkey="'l'" @hotkey="focusList()">
      <div v-for="repo in templatesByRepo" :key="repo.id" class="list">
        <h2 class="list-header template-repo-header">
          <span class="name">{{repo.name}}</span>
          <btn v-if="!repo.isLocal && $api.auth.hasScope('templates.repo.delete')" class="remove" variant="icon" @click="removeRepo(repo)"><icon name="remove" /></btn>
        </h2>
        <div v-for="template in repo.templates" :key="template.name" class="list-item">
          <router-link :ref="setFirstEntry" :to="{ name: 'TemplateView', params: { repo: repo.id, id: template.name } }">
            <div class="template">
              <span class="title">{{template.display}}</span>
              <span class="subline">{{template.type}}</span>
            </div>
          </router-link>
        </div>
        <div v-if="repo.error" class="template-repo-error alert error">
          {{(repo.error.code === 'ErrGeneric' && repo.error.msg) ? t(repo.error.msg) : t('errors.' + repo.error.code)}}
        </div>
        <div v-if="repo.isLocal && $api.auth.hasScope('templates.local.edit')" class="list-item">
          <router-link v-hotkey="'c'" :to="{ name: 'TemplateCreate' }">
            <div class="createLink"><icon name="plus" />{{ t('templates.New') }}</div>
          </router-link>
        </div>
      </div>
      <div v-if="templatesLoaded">
        <a v-if="$api.auth.hasScope('templates.repo.create')" class="repo createLink" @click="addingRepo = true"><icon name="plus" /> {{t('templates.AddRepo')}}</a>
      </div>
      <div v-else class="list-item">
        <loader small />
      </div>
    </div>
    <overlay v-model="addingRepo" :title="t('templates.AddRepo')" closable class="server-name" @close="resetAddRepo()">
      <div class="actions">
        <text-field v-model="currentRepo.name" :label="t('templates.RepoName')" />
        <text-field v-model="currentRepo.url" :label="t('templates.RepoUrl')" />
        <text-field v-model="currentRepo.branch" :label="t('templates.RepoBranch')" />
        <btn v-hotkey="'Escape'" color="error" @click="resetAddRepo()"><icon name="close" />{{ t('common.Cancel') }}</btn>
        <btn :disabled="!canAddRepo()" color="primary" @click="addRepo()"><icon name="save" />{{ t('templates.AddRepo') }}</btn>
      </div>
    </overlay>
  </div>
</template>
