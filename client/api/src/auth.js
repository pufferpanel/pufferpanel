function is2xx(status) {
  return status >= 200 && status < 300
}

function b64Decode(value) {
  const padding = '='.repeat((4 - value.length % 4) % 4)
  const base64 = (value + padding).replace(/-/g, '+').replace(/_/g, '/')

  return new Uint8Array(window.atob(base64).split('').map(c => c.charCodeAt(0)))
}

export class AuthApi {
  _api = null
  _sessionStore = null

  constructor(api, sessionStore) {
    this._api = api
    this._sessionStore = sessionStore
  }

  _handleLogin(scopes) {
    this._sessionStore.setScopes(scopes)
    return true
  }

  _decodeWebauthnChallenge(challenge) {
    if (challenge.publicKey.challenge) {
      challenge.publicKey.challenge = b64Decode(challenge.publicKey.challenge)
    }

    if (challenge.publicKey.user && challenge.publicKey.user.id) {
      challenge.publicKey.user.id = b64Decode(challenge.publicKey.user.id)
    }

    if (challenge.publicKey.allowCredentials) {
      challenge.publicKey.allowCredentials = challenge.publicKey.allowCredentials.map(cred => {
        cred.id = b64Decode(cred.id)
        return cred
      })
    }

    if (challenge.publicKey.excludeCredentials) {
      challenge.publicKey.excludeCredentials = challenge.publicKey.excludeCredentials.map(cred => {
        cred.id = b64Decode(cred.id)
        return cred
      })
    }

    return challenge
  }

  async oauth(clientId, clientSecret) {
    const res = await this._api.post('/oauth2/token', `grant_type=client_credentials&client_id=${clientId}&client_secret=${clientSecret}`)
    this._sessionStore.setToken(res.data.access_token)
    return true
  }

  async login(email, password) {
    const res = await this._api.post('/auth/login', { email, password })
    if (res.data.webauthnChallenge) {
      res.data.webauthnChallenge = this._decodeWebauthnChallenge(res.data.webauthnChallenge)
    }
    if (res.data.needsSecondFactor) return res.data
    return this._handleLogin(res.data.scopes)
  }

  async loginOtp(token) {
    const res = await this._api.post('/auth/otp', { token })
    return this._handleLogin(res.data.scopes)
  }

  async startPasskeyLogin(email, onError) {
    const res = await this._api.post('/auth/passkey', { email }, undefined, undefined, { onError })
    return this._decodeWebauthnChallenge(res.data)
  }

  async validatePasskeyLogin(data) {
    const res = await this._api.put('/auth/passkey', data)
    return this._handleLogin(res.data.scopes)
  }

  async passkeyLogin(email) {
    let attestation
    try {
      const assertion = await this.startPasskeyLogin(email, (error) => {
        if (error.status === 400 && error.response.data.error.code === 'ErrInvalidCredentials') {
          throw 'passkey only rejected'
        }
      })

      attestation = await navigator.credentials.get(assertion)
    } catch(_) {
      return false
    }

    return await this.validatePasskeyLogin(attestation)
  }

  async register(username, email, password) {
    const res = await this._api.post('/auth/register', { username, email, password })
    return this._handleLogin(res.data.scopes)
  }

  async reauth() {
    const res = await this._api.post('/auth/reauth')
    if (is2xx(res.status)) {
      this._handleLogin(res.data.scopes)
    }
  }

  getToken() {
    return this._sessionStore.getToken()
  }

  isLoggedIn() {
    return this._sessionStore.isLoggedIn()
  }

  hasScope(scope) {
    if (!this.isLoggedIn()) return false
    const scopes = this._sessionStore.getScopes()
    if (scopes !== null) {
      if (scopes.indexOf(scope) !== -1) return true
      return scopes.indexOf('admin') !== -1
    }
    return false
  }

  async logout() {
    await this._api.post('/auth/logout')
    this._sessionStore.deleteSession()
  }
}
