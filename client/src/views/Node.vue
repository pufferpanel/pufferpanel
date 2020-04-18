<template>
  <v-row>
    <v-col cols="12" md="6" offset-md="3">
      <v-card>
        <v-card-title v-text="$t('nodes.Edit')" />
        <v-card-text class="mt-6">
          <v-row>
            <v-col>
              <v-text-field :label="$t('common.Name')" v-model="node.name" outlined hide-details />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" md="6">
              <v-text-field :label="$t('nodes.PublicHost')" v-model="node.publicHost" type="text" outlined hide-details />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field :label="$t('nodes.PublicPort')" v-model="node.publicPort" type="number" outlined hide-details />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" md="6">
              <v-text-field :label="$t('nodes.PrivateHost')" v-model="node.privateHost" type="text" outlined hide-details />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field :label="$t('nodes.PrivatePort')" v-model="node.privatePort" type="number" outlined hide-details />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field :label="$t('nodes.SftpPort')" v-model="node.sftpPort" type="number" outlined hide-details />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-btn v-text="$t('nodes.Update')" large block color="primary" @click="updateNode" />
            </v-col>
            <v-col cols="12">
              <v-btn v-text="$t('nodes.Delete')" block color="error" @click="deleteNode" />
            </v-col>
            <v-col cols="12" md="6">
              <v-btn v-text="$t('nodes.SaveConfig')" :disabled="loadingDeploy" text block @click="downloadConfig()" />
            </v-col>
            <v-col cols="12" md="6">
              <v-btn v-text="$t('nodes.SavePublicKey')" :disabled="loadingDeploy" text block @click="download(deployment.publicKey, 'public.pem')" />
            </v-col>
            <v-col cols="12">
              <span v-html="markdown($t('nodes.DeploymentInstruction'))" />
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script>
import { handleError } from '@/utils/api'
import { typeNode } from '@/utils/types'
import markdown from '@/utils/markdown'

export default {
  data () {
    return {
      loading: true,
      loadingDeploy: true,
      node: {},
      deployment: {
        clientId: '',
        clientSecret: '',
        publicKey: ''
      },
      configTemplate: {
        logs: '/var/log/pufferpanel',
        web: {},
        token:{
            public: '/etc/pufferpanel/public.pem',
        },
        panel: {
          enable: false
        },
        daemon: {
          auth: {
            url: location.protocol + '//' + location.host
          },
          data: {
            cache: '/var/lib/pufferpanel/cache',
            servers: '/var/lib/pufferpanel/servers'
          },
          sftp: {},
        }
      }
    }
  },
  mounted () {
    this.loadData()
  },
  methods: {
    loadData () {
      const ctx = this
      ctx.$http.get(`/api/nodes/${ctx.$route.params.id}`).then(response => {
        ctx.node = response.data
        ctx.loading = false
      }).catch(handleError(ctx))
      ctx.$http.get(`/api/nodes/${ctx.$route.params.id}/deployment`).then(response => {
        ctx.deployment = response.data
        ctx.loadingDeploy = false
      }).catch(handleError(ctx))
    },
    updateNode () {
      const ctx = this
      ctx.$http.put(`/api/nodes/${ctx.$route.params.id}`, typeNode(ctx.node)).then(response => {
        ctx.$toast.success(ctx.$t('nodes.UpdateSuccess'))
      }).catch(handleError(ctx))
    },
    deleteNode () {
      const ctx = this
      ctx.$http.delete(`/api/nodes/${ctx.$route.params.id}`).then(response => {
        ctx.$router.push({ name: 'Nodes' })
      }).catch(handleError(ctx))
    },
    downloadConfig () {
      const config = { ...this.configTemplate }
      config.daemon.auth.clientId = this.deployment.clientId
      config.daemon.auth.clientSecret = this.deployment.clientSecret
      config.daemon.sftp.host = `0.0.0.0:${this.node.sftpPort}`
      config.web.host = `0.0.0.0:${this.node.privatePort}`
      this.download(JSON.stringify(config, undefined, 2), 'config.json')
    },
    download (content, filename, contentType) {
      if (!contentType) contentType = 'application/octet-stream'
      var a = document.createElement('a')
      var blob = new Blob([content], { type: contentType })
      a.href = window.URL.createObjectURL(blob)
      a.download = filename
      a.click()
    },
    markdown
  }
}
</script>
