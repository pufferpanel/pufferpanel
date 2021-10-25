<template>
  <v-row>
    <v-col cols="12">
      <v-card>
        <v-card-title v-text="$t('nodes.Add')" />
        <v-card-text class="mt-6">
          <v-row>
            <v-col>
              <ui-input
                v-model="node.name"
                :label="$t('common.Name')"
                @keyup.enter="submit"
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
                @keyup.enter="submit"
              />
            </v-col>
            <v-col
              cols="12"
              md="6"
            >
              <ui-input
                v-model="node.publicPort"
                :label="$t('nodes.PublicPort')"
                @keyup.enter="submit"
                type="number"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <ui-switch
                v-model="withPrivateAddress"
                class="mt-0"
                :label="$t('nodes.WithPrivateAddress')"
                :hint="$t('nodes.WithPrivateAddressHint')"
              />
            </v-col>
          </v-row>
          <v-row v-if="withPrivateAddress">
            <v-col
              cols="12"
              md="6"
            >
              <ui-input
                v-model="node.privateHost"
                :label="$t('nodes.PrivateHost')"
                @keyup.enter="submit"
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
                @keyup.enter="submit"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <ui-input
                v-model="node.sftpPort"
                :label="$t('nodes.SftpPort')"
                type="number"
                @keyup.enter="submit"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-btn
                large
                block
                color="success"
                :disabled="!canSubmit()"
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
export default {
  data () {
    return {
      withPrivateAddress: false,
      node: {
        publicPort: 8080,
        privatePort: 8080,
        sftpPort: 5657
      }
    }
  },
  methods: {
    canSubmit () {
      if (!this.node.name) return false
      if (!this.node.publicHost) return false
      if (!this.node.publicPort) return false
      if (!this.node.sftpPort) return false
      if (this.withPrivateAddress) {
        if (!this.node.privateHost) return false
        if (!this.node.privatePort) return false
      }
      return true
    },
    async submit () {
      if (!this.canSubmit()) return
      if (!this.withPrivateAddress) {
        this.node.privateHost = this.node.publicHost
        this.node.privatePort = this.node.publicPort
      }
      const created = await this.$api.createNode(this.node)
      if (!created) return
      const nodes = await this.$api.getNodes()
      if (!nodes) this.$router.push({ name: 'Nodes' })
      const id = nodes.find(node => node.name === this.node.name).id
      this.$router.push({ name: 'Node', params: { id } })
    }
  }
}
</script>
