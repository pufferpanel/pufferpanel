<template>
  <v-row>
    <v-col cols="12" md="6" offset-md="3">
      <v-card>
        <v-card-title v-text="$t('nodes.Edit')" />
        <v-card-text class="mt-6">
          <v-row>
            <v-col>
              <v-text-field :label="$t('common.Name')" v-model="node.name" outlined />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" md="6">
              <v-text-field :label="$t('nodes.PublicHost')" v-model="node.publicHost" type="text" outlined />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field :label="$t('nodes.PublicPort')" v-model="node.publicPort" type="number" outlined />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" md="6">
              <v-text-field :label="$t('nodes.PrivateHost')" v-model="node.privateHost" type="text" outlined />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field :label="$t('nodes.PrivatePort')" v-model="node.privatePort" type="number" outlined />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field :label="$t('nodes.SftpPort')" v-model="node.sftpPort" type="number" outlined />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-btn v-text="$t('nodes.Update')" large block color="primary" @click="updateNode" />
            </v-col>
            <v-col cols="12">
              <v-btn v-text="$t('nodes.Delete')" block color="error" @click="deleteNode" />
            </v-col>
            <v-col cols="12">
              <v-btn v-text="$t('nodes.DeploymentData')" text block :to="`/api/nodes/${node.id}/deployment`" target="_blank" />
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
      ctx.$http.get(`/api/nodes/${ctx.$route.params.id}`).then(response => {
        ctx.node = response.data
        ctx.loading = false
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
    }
  }
}
</script>
