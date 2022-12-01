<template>
  <v-container>
    <h1 v-text="$t('nodes.Nodes')" />
    <v-row>
      <v-col>
        <v-sheet
          elevation="1"
          class="py-2"
        >
          <v-data-table
            style="cursor: pointer;"
            :items="nodes"
            :headers="headers"
            :loading="loading"
            hide-default-footer
            @click:row="rowClicked"
          >
          </v-data-table>
          <v-btn
            v-show="hasScope('nodes.deploy') || isAdmin()"
            color="primary"
            bottom
            right
            fixed
            fab
            dark
            large
            :to="{name: 'AddNode'}"
          >
            <v-icon>mdi-plus</v-icon>
          </v-btn>
        </v-sheet>
      </v-col>
    </v-row>
    <ui-overlay
        v-model="tryingToEditLocal"
        card
        closable
        :title="$t('nodes.CannotEdit')"
    >
      <!-- eslint-disable-next-line vue/no-v-html -->
      <span style="text-align: center; overflow-y: auto" v-html="markdown($t('nodes.LocalNodeInstructions'))"/>
    </ui-overlay>
  </v-container>
</template>

<script>
import markdown from '@/utils/markdown'

export default {
  data () {
    return {
      loading: true,
      tryingToEditLocal: false,
      nodes: [],
      headers: [
        {
          text: this.$t('common.Name'),
          value: 'name'
        },
        {
          text: this.$t('nodes.PublicHost'),
          value: 'publicHost'
        },
        {
          text: this.$t('nodes.PublicPort'),
          value: 'publicPort'
        },
        {
          text: this.$t('nodes.PrivateHost'),
          value: 'privateHost'
        },
        {
          text: this.$t('nodes.PrivatePort'),
          value: 'privatePort'
        },
        {
          text: this.$t('nodes.SftpPort'),
          value: 'sftpPort'
        }
      ]
    }
  },
  mounted () {
    this.loadData()
  },
  methods: {
    async loadData () {
      this.loading = true
      this.nodes = []
      this.nodes = await this.$api.getNodes()
      this.loading = false
    },
    rowClicked (item) {
      if (this.hasScope('nodes.edit') || this.isAdmin()) {
        if (item.isLocal) {
          this.tryingToEditLocal = true
        } else {
          this.$router.push({ name: 'Node', params: { id: item.id } })
        }
      }
    },
    markdown
  }
}
</script>
