export class SettingsApi {
  _api = null

  constructor(api) {
    this._api = api
  }

  async get(key) {
    const res = await this._api.get(`/api/settings/${key}`)
    return res.data.value
  }

  async set(data) {
    const res = await this._api.post('/api/settings/', data)
    return res.data.value
  }

  async getUserSettings() {
    const res = await this._api.get('/api/userSettings/')
    const map = {}
    res.data.map(e => {
      map[e.key] = e.value
    })
    return map
  }

  async setUserSetting(key, value) {
    await this._api.put(`/api/userSettings/${key}`, { value })
    return true
  }
}
