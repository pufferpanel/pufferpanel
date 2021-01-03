// used as a mixin for the ApiClient, use of `this` refers to the ApiClient instance
export const UsersApi = {
  getUsers (page, pageSize = 10) {
    return this.withErrorHandling(async ctx => {
      const res = (await ctx.$http.get('/api/users', { params: { page, limit: pageSize } })).data
      return { users: res.users, pages: Math.ceil(res.paging.total / pageSize) }
    })
  },

  searchUsers (query, cancelToken) {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.get('/api/users', { params: { username: `*${query}*` }, cancelToken })).data.users
    })
  },

  createUser (user) {
    return this.withErrorHandling(async ctx => {
      const id = (await ctx.$http.post('/api/users', user)).data.id
      await this.updateUserPermissions(id, user)
      return id
    })
  },

  getUser (id) {
    return this.withErrorHandling(async ctx => {
      const user = (await ctx.$http.get(`/api/users/${id}`)).data
      const perms = await this.getUserPermissions(id)
      return { ...user, ...perms }
    })
  },

  getUserPermissions (id) {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.get(`/api/users/${id}/perms`)).data
    })
  },

  updateUser (id, user) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.post(`/api/users/${id}`, user)
      await this.updateUserPermissions(id, user)
      return true
    })
  },

  updateUserPermissions (id, permissions) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.put(`/api/users/${id}/perms`, permissions)
      return true
    })
  },

  deleteUser (id) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.delete(`/api/users/${id}`)
      return true
    })
  }
}
