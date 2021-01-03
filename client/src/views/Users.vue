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
        <v-row
          v-if="page < pageCount"
          ref="lazy"
          v-intersect="lazyLoad"
        >
          <v-col
            cols="2"
            offset="5"
          >
            <v-progress-circular
              indeterminate
              class="mr-2"
            />
            <span v-text="$t('common.Loading')" />
          </v-col>
        </v-row>
        <v-btn
          v-show="hasScope('users.edit') || isAdmin()"
          color="primary"
          bottom
          right
          fixed
          fab
          dark
          large
          :to="{name: 'AddUser'}"
        >
          <v-icon>mdi-plus</v-icon>
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  data () {
    return {
      loading: false,
      users: [],
      page: 0,
      pageCount: 1
    }
  },
  methods: {
    recheckLazy () {
      const rect = this.$refs.lazy.getBoundingClientRect()
      const viewHeight = Math.max(document.documentElement.clientHeight, window.innerHeight)
      if (
        !(rect.bottom < 0 || rect.top - viewHeight >= 0) &&
        this.page < this.pageCount
      ) {
        this.loadNextPage()
      }
    },
    lazyLoad (entries, observer, isIntersecting) {
      if (isIntersecting) {
        this.loadNextPage()
      }
    },
    async loadNextPage () {
      if (this.loading) return
      this.loading = true
      const { users, pages } = await this.$api.getUsers(this.page + 1)
      users.filter(user => {
        return this.users.filter(elem => user.id === elem.id).length === 0
      }).forEach(user => this.users.push(user))
      this.page = this.page + 1
      this.pageCount = pages
      this.loading = false
      this.recheckLazy()
    }
  }
}
</script>
