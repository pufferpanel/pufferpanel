<template>
  <v-row>
    <v-col cols="12" md="6" offset-md="3">
      <v-card>
        <v-card-title v-text="$t('common.AddNode')" />
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
              <v-btn v-text="$t('common.AddNode')" large block color="success" @click="submit" />
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script>
import { typeNode } from '@/utils/types'

export default {
  data () {
    return {
      node: {}
    }
  },
  methods: {
    submit () {
      const ctx = this
      ctx.$http.post('/api/nodes', typeNode(ctx.node)).then(function (response) {
        ctx.$router.push({ name: 'Nodes' })
      }).catch(function (error) {
        let msg = 'errors.ErrUnknownError'
        if (error && error.response && error.response.data.error) {
          if (error.response.data.error.code) {
            msg = 'errors.' + error.response.data.error.code
          } else {
            msg = error.response.data.error.msg
          }
        }

        ctx.$toast.error(ctx.$t(msg))
      })
    }
  }
}
</script>
