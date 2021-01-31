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
          />
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
  </v-container>
</template>

<script>
export default {
  data () {
    return {
      loading: true,
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
      if (this.hasScope('nodes.edit') || this.isAdmin()) this.$router.push({ name: 'Node', params: { id: item.id } })
    }
  }
}
</script>
