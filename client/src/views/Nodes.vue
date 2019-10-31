<template>
  <v-container>
    <v-data-table
      style="cursor: pointer;"
      :items="nodes"
      :headers="headers"
      :loading="loading"
      hide-default-footer
      @click:row="rowClicked"
    />
    <v-btn
      v-show="isAdmin()"
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
          text: this.$t('common.PublicHost'),
          value: 'publicHost'
        },
        {
          text: this.$t('common.PublicPort'),
          value: 'publicPort'
        },
        {
          text: this.$t('common.PrivateHost'),
          value: 'privateHost'
        },
        {
          text: this.$t('common.PrivatePort'),
          value: 'privatePort'
        },
        {
          text: this.$t('common.SftpPort'),
          value: 'sftpPort'
        }
      ]
    }
  },
  mounted () {
    this.loadData()
  },
  methods: {
    loadData () {
      const ctx = this
      ctx.loading = true
      ctx.nodes = []
      ctx.$http.get('/api/nodes').then(function (response) {
        if (response.status >= 200 && response.status < 300) {
          response.data.forEach(function (node) {
            ctx.nodes.push(node)
          })
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

        ctx.$toast.error(ctx.$t(msg))
      })
    },
    rowClicked (item) {
      this.$router.push({ name: 'Node', params: { id: item.id } })
    }
  }
}
</script>
