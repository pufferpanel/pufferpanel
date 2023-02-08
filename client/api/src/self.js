export class SelfApi {
  _api = null

  constructor(api) {
    this._api = api
  }

  async get() {
    const res = await this._api.get('/api/self')
    return res.data
  }

  async updateDetails(username, email, password) {
    await this._api.put('/api/self', { username, email, password })
    return true
  }

  async changePassword(password, newPassword) {
    await this._api.put('/api/self', { password, newPassword })
    return true
  }

  async isOtpEnabled() {
    const res = await this._api.get('/api/self/otp')
    return res.data.otpEnabled
  }

  async startOtpEnroll() {
    const res = await this._api.post('/api/self/otp')
    return res.data
  }

  async validateOtpEnroll(token) {
    await this._api.put('/api/self/otp', { token })
    return true
  }

  async disableOtp(token) {
    await this._api.delete(`/api/self/otp/${token}`)
    return true
  }

  async getSettings() {
    const res = await this._api.get('/api/userSettings')
    const map = {}
    res.data.map(e => map[e.key] = e.value)
    return map
  }

  async updateSetting(key, value) {
    await this._api.put(`/api/userSettings/${key}`, { value })
    return true
  }

  async getOAuthClients() {
    const res = await this._api.get(`/api/self/oauth2`)
    return res.data
  }

  async createOAuthClient(name, description) {
    const res = await this._api.post(`/api/self/oauth2`, { name, description })
    return res.data
  }

  async deleteOAuthClient(clientId) {
    await this._api.delete(`/api/self/oauth2/${clientId}`)
    return true
  }
}
