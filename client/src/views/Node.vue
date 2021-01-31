<template>
  <v-row>
    <v-col
      cols="12"
      md="6"
      offset-md="3"
    >
      <v-card>
        <v-card-title v-text="$t('nodes.Edit')" />
        <v-card-text class="mt-6">
          <v-row>
            <v-col>
              <ui-input
                v-model="node.name"
                :label="$t('common.Name')"
                hide-details
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col
              cols="12"
              md="6"
            >
              <ui-input
                v-model="node.publicHost"
                :label="$t('nodes.PublicHost')"
                hide-details
              />
            </v-col>
            <v-col
              cols="12"
              md="6"
            >
              <ui-input
                v-model="node.publicPort"
                :label="$t('nodes.PublicPort')"
                type="number"
                hide-details
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col
              cols="12"
              md="6"
            >
              <ui-input
                v-model="node.privateHost"
                :label="$t('nodes.PrivateHost')"
                hide-details
              />
            </v-col>
            <v-col
              cols="12"
              md="6"
            >
              <ui-input
                v-model="node.privatePort"
                :label="$t('nodes.PrivatePort')"
                type="number"
                hide-details
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <ui-input
                v-model="node.sftpPort"
                :label="$t('nodes.SftpPort')"
                type="number"
                hide-details
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-btn
                large
                block
                color="success"
                @click="updateNode"
                v-text="$t('nodes.Update')"
              />
            </v-col>
            <v-col cols="12">
              <v-btn
                block
                color="error"
                @click="deleteNode"
                v-text="$t('nodes.Delete')"
              />
            </v-col>
            <v-col cols="12">
              <v-btn
                :disabled="loadingDeploy"
                text
                block
                @click="downloadConfig()"
                v-text="$t('nodes.SaveConfig')"
              />
            </v-col>
            <v-col cols="12">
              <!-- eslint-disable-next-line vue/no-v-html -->
              <span v-html="markdown($t('nodes.DeploymentInstruction'))" />
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script>
import markdown from '@/utils/markdown'

export default {
  data () {
    return {
      loading: true,
      loadingDeploy: true,
      node: {},
      deployment: {
        clientId: '',
        clientSecret: ''
      },
      configTemplate: {
        logs: '/var/log/pufferpanel',
        web: {},
        token: {
          public: location.protocol + '//' + location.host + '/auth/publickey'
        },
        panel: {
          enable: false
        },
        daemon: {
          auth: {
            url: location.protocol + '//' + location.host + '/oauth2/token'
          },
          data: {
            cache: '/var/lib/pufferpanel/cache',
            servers: '/var/lib/pufferpanel/servers'
          },
          sftp: {}
        }
      }
    }
  },
  mounted () {
    this.loadData()
  },
  methods: {
    async loadData () {
      this.node = await this.$api.getNode(this.$route.params.id)
      this.loading = false
      this.deployment = await this.$api.getNodeDeployment(this.$route.params.id)
      this.loadingDeploy = false
    },
    async updateNode () {
      await this.$api.updateNode(this.$route.params.id, this.node)
      this.$toast.success(this.$t('nodes.UpdateSuccess'))
    },
    async deleteNode () {
      await this.$api.deleteNode(this.$route.params.id)
      this.$router.push({ name: 'Nodes' })
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
