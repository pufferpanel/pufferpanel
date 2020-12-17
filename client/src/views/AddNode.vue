<template>
  <v-row>
    <v-col
      cols="12"
      md="6"
      offset-md="3"
    >
      <v-card>
        <v-card-title v-text="$t('nodes.Add')" />
        <v-card-text class="mt-6">
          <v-row>
            <v-col>
              <ui-input
                v-model="node.name"
                :label="$t('common.Name')"
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
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <ui-input
                v-model="node.sftpPort"
                :label="$t('nodes.SftpPort')"
                type="number"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-btn
                large
                block
                color="success"
                @click="submit"
                v-text="$t('nodes.Add')"
              />
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script>
import { typeNode } from '@/utils/types'
import { handleError } from '@/utils/api'

export default {
  data () {
    return {
      node: {}
    }
  },
  methods: {
    submit () {
      const ctx = this
      ctx.$http.post('/api/nodes', typeNode(ctx.node)).then(response => {
        ctx.$router.push({ name: 'Nodes' })
      }).catch(handleError(ctx))
    }
  }
}
</script>
