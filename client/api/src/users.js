export class UserApi {
  _api = null

  constructor(api) {
    this._api = api
  }

  async list(page = 1) {
    const res = await this._api.get('/api/users', { page })
    return res.data
  }

  async search(query) {
    const res = await this._api.get('/api/users', { username: `*${query}*` })
    return res.data.users
  }

  async create(username, email, password) {
    const res = await this._api.post('/api/users', { username, email, password })
    return res.data.id
  }

  async get(id) {
    const res = await this._api.get(`/api/users/${id}`)
    return res.data
  }

  async getPermissions(id) {
    const res = await this._api.get(`/api/users/${id}/perms`)
    return res.data
  }

  async update(id, user) {
    await this._api.post(`/api/users/${id}`, user)
    return true
  }

  async updatePermissions(id, permissions) {
    await this._api.put(`/api/users/${id}/perms`, permissions)
    return true
  }

  async delete(id) {
    await this._api.delete(`/api/users/${id}`)
    return true
  }
}
