<template>
  <v-container>
    <h1 v-text="$t('common.User')" />
    <v-row>
      <v-col>
            <v-list two-line elevation="1">
              <div v-for="(user, index) in users">
                <v-list-item :to="{name: 'User', params: {id: user.id}}">
                  <v-list-item-content>
                    <v-list-item-title v-text="user.username" />
                    <v-list-item-subtitle v-text="user.email" />
                  </v-list-item-content>
                </v-list-item>
                <v-divider v-if="index !== users.length - 1" />
              </div>
            </v-list>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  data () {
    return {
      loading: true,
      users: []
    }
  },
  mounted () {
    this.loadData()
  },
  methods: {
    loadData () {
      const ctx = this
      ctx.loading = true
      ctx.users = []
      ctx.$http.get('/api/users').then(function (response) {
        response.data.users.forEach(function (user) {
          ctx.users.push(user)
        })
        ctx.loading = false
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
