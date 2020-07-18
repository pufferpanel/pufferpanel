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
              <v-text-field
                v-model="node.name"
                :label="$t('common.Name')"
                outlined
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col
              cols="12"
              md="6"
            >
              <v-text-field
                v-model="node.publicHost"
                :label="$t('nodes.PublicHost')"
                type="text"
                outlined
              />
            </v-col>
            <v-col
              cols="12"
              md="6"
            >
              <v-text-field
                v-model="node.publicPort"
                :label="$t('nodes.PublicPort')"
                type="number"
                outlined
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col
              cols="12"
              md="6"
            >
              <v-text-field
                v-model="node.privateHost"
                :label="$t('nodes.PrivateHost')"
                type="text"
                outlined
              />
            </v-col>
            <v-col
              cols="12"
              md="6"
            >
              <v-text-field
                v-model="node.privatePort"
                :label="$t('nodes.PrivatePort')"
                type="number"
                outlined
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field
                v-model="node.sftpPort"
                :label="$t('nodes.SftpPort')"
                type="number"
                outlined
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
