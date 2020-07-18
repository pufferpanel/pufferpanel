<template>
  <v-container>
    <h1 v-text="$t('users.Users')" />
    <v-row>
      <v-col>
        <v-list
          two-line
          elevation="1"
        >
          <div
            v-for="(user, index) in users"
            :key="user.id"
          >
            <v-list-item :to="(hasScope('users.edit') || isAdmin()) ? {name: 'User', params: {id: user.id}} : undefined">
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
import { handleError } from '@/utils/api'

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
      ctx.$http.get('/api/users').then(response => {
        response.data.users.forEach(user => {
          ctx.users.push(user)
        })
        ctx.loading = false
      }).catch(handleError(ctx))
    }
  }
}
</script>
