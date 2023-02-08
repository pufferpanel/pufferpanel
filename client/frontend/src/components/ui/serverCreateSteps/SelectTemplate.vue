<script setup>
import { ref, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import Overlay from '@/components/ui/Overlay.vue'
import markdown from '@/utils/markdown.js'

const { t } = useI18n()
const api = inject('api')
const toast = inject('toast')
const emit = defineEmits(['selected'])
const templatesByRepo = ref([])
const showing = ref(false)
const currentTemplate = ref({})

async function load() {
  templatesByRepo.value = await api.template.listAllTemplates()
}

onMounted(async () => {
  load()
})

async function show(repo, template) {
  currentTemplate.value = await api.template.get(repo, template)
  showing.value = true
}

function choice(confirm) {
  showing.value = false
  if (confirm) emit('selected', currentTemplate.value)
}
</script>

<template>
  <div class="select-template">
    <h2 v-text="t('servers.SelectTemplate')" />
    <div v-for="(templates, repo) in templatesByRepo" :key="repo" class="list">
      <h3 class="list-header" v-text="repo" />
      <div v-for="template in templates" :key="template.name" class="list-item template" @click="show(repo, template.name)">
        <span v-text="template.display" />
      </div>
    </div>

    <overlay v-model="showing" :title="currentTemplate.display" closable>
      <!-- eslint-disable-next-line vue/no-v-html -->
      <div v-if="currentTemplate.readme" dir="ltr" class="readme" v-html="markdown(currentTemplate.readme)" />
      <h2 v-else v-text="t('servers.ConfirmTemplateChoice')" />
      <div class="actions">
        <btn color="error" @click="choice(false)"><icon name="close" />{{ t('common.Cancel') }}</btn>
        <btn color="primary" @click="choice(true)"><icon name="check" />{{ t('servers.SelectThisTemplate') }}</btn>
      </div>
    </overlay>
  </div>
</template>
