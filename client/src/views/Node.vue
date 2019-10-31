<template>
  <v-row>
    <v-col cols="12" md="6" offset-md="3">
      <v-card>
        <v-card-title v-text="$t('common.EditNode')" />
        <v-card-text class="mt-6">
          <v-row>
            <v-col>
              <v-text-field :label="$t('common.Name')" v-model="node.name" outlined />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" md="6">
              <v-text-field :label="$t('common.PublicHost')" v-model="node.publicHost" type="text" outlined />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field :label="$t('common.PublicPort')" v-model="node.publicPort" type="number" outlined />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" md="6">
              <v-text-field :label="$t('common.PrivateHost')" v-model="node.privateHost" type="text" outlined />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field :label="$t('common.PrivatePort')" v-model="node.privatePort" type="number" outlined />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field :label="$t('common.SftpPort')" v-model="node.sftpPort" type="number" outlined />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-btn v-text="$t('common.UpdateNode')" large block color="primary" @click="updateNode" />
            </v-col>
            <v-col cols="12">
              <v-btn v-text="$t('common.DeleteNode')" block color="error" @click="deleteNode" />
            </v-col>
            <v-col cols="12">
              <v-btn v-text="$t('common.NodeDeploymentData')" text block :to="`/api/nodes/${node.id}/deployment`" target="_blank" />
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script>
export default {
  data () {
    return {
      loading: true,
      node: {}
    }
  },
  mounted () {
    this.loadData()
  },
  methods: {
    loadData () {
      const ctx = this
      ctx.$http.get(`/api/nodes/${ctx.$route.params.id}`).then(function (response) {
        if (response.status >= 200 && response.status < 300) {
          ctx.node = response.data
          ctx.loading = false
        }
      }).catch(function (error) {
        let msg = 'errors.ErrUnknownError'
        if (error && error.response && error.response.data.error) {
          if (error.response.data.error.code) {
            msg = 'errors.' + error.response.data.error.code
          } else {
            msg = error.response.data.error.msg
          }
        }

        ctx.$notify(ctx.$t(msg), 'error')
      })
    },
    updateNode () {
      const ctx = this
      ctx.$http.put(`/api/nodes/${ctx.$route.params.id}`, ctx.node).then(function (response) {
        ctx.$notify(ctx.$t('common.NodeUpdateSuccess'), 'success')
      }).catch(function () {
        ctx.$notify(ctx.$t('common.NodeUpdateError'), 'error')
      })
    },
    deleteNode () {
      const ctx = this
      ctx.$http.delete(`/api/nodes/${ctx.$route.params.id}`).then(function (response) {
        if (response.status >= 200 && response.status < 300) ctx.$router.push({ name: 'Nodes' })
      }).catch(function (error) {
        let msg = 'errors.ErrUnknownError'
        if (error && error.response && error.response.data.error) {
          if (error.response.data.error.code) {
            msg = 'errors.' + error.response.data.error.code
          } else {
            msg = error.response.data.error.msg
          }
        }

        ctx.$notify(ctx.$t(msg), 'error')
      })
    }
  }
}
</script>
