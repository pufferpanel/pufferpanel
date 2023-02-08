<script setup>
import { ref, inject, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Btn from '@/components/ui/Btn.vue'
import Icon from '@/components/ui/Icon.vue'
import Overlay from '@/components/ui/Overlay.vue'
import TextField from '@/components/ui/TextField.vue'

const { t } = useI18n()
const api = inject('api')
const events = inject('events')

const props = defineProps({
  server: { type: String, default: () => undefined }
})

const docsUrl = location.origin + '/swagger/index.html'

const clients = ref([])
const creating = ref(false)
const newName = ref('')
const newDescription = ref('')
const created = ref(false)
const createdData = ref(null)

onMounted(() => {
  refresh()
})

async function refresh() {
  if (props.server) {
    clients.value = await api.server.getOAuthClients(props.server)
  } else {
    clients.value = await api.self.getOAuthClients()
  }
}

function startCreate() {
  newName.value = ''
  newDescription.value = ''
  creating.value = true
}

async function create() {
  if (props.server) {
    createdData.value = await api.server.createOAuthClient(props.server, newName.value, newDescription.value)
  } else {
    createdData.value = await api.self.createOAuthClient(newName.value, newDescription.value)
  }

  created.value = true
  creating.value = false
  refresh()
}

async function deleteClient(clientId, clientName) {
  events.emit(
    'confirm',
    t('oauth.ConfirmDelete', { name: clientName }),
    {
      text: t('oauth.Delete'),
      icon: 'remove',
      color: 'error',
      action: async () => {
        if (props.server) {
          clients.value = await api.server.deleteOAuthClient(props.server, clientId)
        } else {
          clients.value = await api.self.deleteOAuthClient(clientId)
        }

        refresh()
      }
    },
    {
      color: 'primary'
    }
  )
}
</script>

<template>
  <div class="oauth">
    <div class="info">
      <div v-text="t(server ? 'oauth.ServerDescription' : 'oauth.AccountDescription')" />
      <div>
        <a target="_blank" :href="docsUrl" v-text="t('oauth.Docs')" />
      </div>
    </div>
    <div v-for="(client, index) in clients" :key="client.client_id">
      <div class="oauth-client">
        <div class="details">
          <div v-text="client.name || t('oauth.UnnamedClient')" />
          <div v-text="client.client_id" />
          <div :title="client.description" v-text="client.description" />
        </div>
        <btn variant="icon" @click="deleteClient(client.client_id, client.name || t('oauth.UnnamedClient'))"><icon name="remove" /></btn>
      </div>
      <hr v-if="index < clients.length - 1" />
    </div>
    <btn color="primary" @click="startCreate()"><icon name="plus" />{{ t('oauth.Create') }}</btn>
    <overlay v-model="creating" :title="t('oauth.Create')" closable>
      <text-field v-model="newName" autofocus :label="t('common.Name')" />
      <text-field v-model="newDescription" :label="t('common.Description')" />
      <btn color="error" @click="creating = false"><icon name="close" />{{ t('common.Cancel') }}</btn>
      <btn color="primary" @click="create()"><icon name="save" />{{ t('oauth.Create') }}</btn>
    </overlay>
    <overlay v-model="created" :title="t('oauth.Credentials')" closable>
      <div class="warning" v-text="t('oauth.NewClientWarning')" />
      <div class="client-id">
        <span class="name" v-text="t('oauth.ClientId')+':'" />
        <span class="value" v-text="createdData.id" />
      </div>
      <div class="client-secret">
        <span class="name" v-text="t('oauth.ClientSecret')+':'" />
        <span class="value" v-text="createdData.secret" />
      </div>
    </overlay>
  </div>
</template>
