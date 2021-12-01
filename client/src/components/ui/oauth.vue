<template>
  <v-card>
    <v-card-title v-text="$t('oauth.Clients')" />
    <v-card-subtitle>
      <div>{{ $t('oauth.' + (isServer ? 'ServerDescription' : 'AccountDescription'))}}</div>
      <a target="_blank" :href="apiDocsUrl">
        {{ $t('oauth.Docs') }}
        <v-icon class="caption">mdi-launch</v-icon>
      </a>
    </v-card-subtitle>
    <v-card-text>
      <v-list>
        <v-list-item v-for="client in clients" :key="client.client_id">
          <v-list-item-content>
            <v-list-item-title v-text="client.name || $t('oauth.UnnamedClient')" />
            <v-list-item-subtitle v-text="client.client_id" />
            <v-tooltip bottom>
              <template v-slot:activator="{ on }">
                <v-list-item-subtitle v-on="on" v-text="client.description" />
              </template>
              <span v-text="client.description" />
            </v-tooltip>
          </v-list-item-content>
          <v-list-item-icon>
            <v-btn icon @click="deleteClient(client)">
              <v-icon>mdi-delete</v-icon>
            </v-btn>
          </v-list-item-icon>
        </v-list-item>
        <v-list-item @click="newClient()">
          <v-list-item-avatar>
            <v-icon>mdi-plus</v-icon>
          </v-list-item-avatar>
          <v-list-item-content>
            <v-list-item-title v-text="$t('common.Create')" />
          </v-list-item-content>
        </v-list-item>
      </v-list>

      <ui-overlay
        v-model="creating"
        :title="$t('oauth.Create')"
        card
        closable
        @close="resetCreate()"
      >
        <v-row>
          <v-col cols="12">
            <ui-input v-model="newName" autofocus :label="$t('common.Name')" />
          </v-col>
          <v-col cols="12">
            <ui-input v-model="newDescription" :label="$t('common.Description')" @keyup.enter="createClient()" />
          </v-col>
          <v-col cols="12" md="6">
            <v-btn large block color="error" v-text="$t('common.Cancel')" @click="resetCreate()" />
          </v-col>
          <v-col cols="12" md="6">
            <v-btn large block color="success" v-text="$t('common.Create')" @click="createClient()" />
          </v-col>
        </v-row>
      </ui-overlay>
      <ui-overlay
        v-model="newClientDataOpen"
        :title="$t('oauth.Credentials')"
        card
        closable
        @close="newClientData = false"
      >
        <v-row>
          <v-col cols="12">
            <v-alert border="bottom" text type="warning" dense>
              {{ $t('oauth.NewClientWarning') }}
            </v-alert>
          </v-col>
          <v-col cols="12">
            <p class="title" v-text="$t('oauth.ClientId')" />
            <p class="body" v-text="newClientData.id" />
          </v-col>
          <v-col cols="12">
            <p class="title" v-text="$t('oauth.ClientSecret')" />
            <p class="body" v-text="newClientData.secret" />
          </v-col>
        </v-row>
      </ui-overlay>
      <ui-overlay
        v-model="deleteOpen"
        :title="$t('oauth.Delete', { clientName: deleteName })"
        card
        closable
        @close="resetDelete()"
      >
        <v-row>
          <v-col cols="12">
            <v-alert border="bottom" text type="error" dense>
              {{ $t('oauth.DeleteWarning') }}
            </v-alert>
          </v-col>
          <v-col cols="12">
            <ui-switch v-model="deleteConfirm" :label="$t('oauth.ConfirmDelete')" />
          </v-col>
          <v-col cols="12" md="6">
            <v-btn large block color="error" v-text="$t('common.Cancel')" @click="resetDelete()" />
          </v-col>
          <v-col cols="12" md="6">
            <v-btn :disabled="!deleteConfirm" large block color="success" v-text="$t('common.Delete')" @click="deleteConfirmed()" />
          </v-col>
        </v-row>
      </ui-overlay>
    </v-card-text>
  </v-card>
</template>

<script>
export default {
  props: {
    server: { type: Object, default: () => undefined }
  },
  data () {
    return {
      apiDocsUrl: location.origin + '/swagger/index.html',
      isServer: false,
      clients: [],
      creating: false,
      newName: '',
      newDescription: '',
      newClientDataOpen: false,
      newClientData: false,
      deleteOpen: false,
      deleteConfirm: false,
      deleteName: '',
      toDelete: false
    }
  },
  mounted () {
    this.isServer = !!this.server
    this.getClients()
  },
  methods: {
    async getClients (client) {
      if (this.isServer) {
        this.clients = await this.$api.getServerOAuthClients(this.server.id)
      } else {
        this.clients = await this.$api.getUserOAuthClients()
      }
    },
    newClient () {
      this.creating = true
    },
    async createClient () {
      let clientData = false
      if (this.isServer) {
        clientData = await this.$api.createServerOAuthClient(this.server.id, this.newName, this.newDescription)
      } else {
        clientData = await this.$api.createUserOAuthClient(this.newName, this.newDescription)
      }
      this.resetCreate()
      this.newClientData = clientData
      this.newClientDataOpen = !!clientData
      this.getClients()
    },
    resetCreate () {
      this.creating = false
      this.newName = ''
      this.newDescription = ''
    },
    resetDelete () {
      this.deleteOpen = false
      this.toDelete = false
      this.deleteName = false
      this.deleteConfirm = false
    },
    deleteClient (client) {
      this.deleteOpen = true
      this.toDelete = client.client_id
      this.deleteName = client.name || this.$t('oauth.UnnamedClient')
      this.deleteConfirm = false
    },
    async deleteConfirmed () {
      if (this.isServer) {
        await this.$api.deleteServerOAuthClient(this.server.id, this.toDelete)
      } else {
        await this.$api.deleteUserOAuthClient(this.toDelete)
      }
      this.resetDelete()
      this.getClients()
    }
  }
}
</script>
