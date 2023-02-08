import Cookies from 'js-cookie'

export class InMemorySessionStore {
  _token = null
  _scopes = null

  setToken(token) {
    this._token = token
  }

  setScopes(scopes) {
    this._scopes = scopes
  }

  getToken() {
    return this._token
  }

  getScopes() {
    return this._scopes
  }

  isLoggedIn() {
    return this._token !== null
  }

  deleteSession() {
    this._token = null
    this._scopes = null
  }
}

const defaultCookieOptions = {
  domain: undefined,
  path: '/',
  secure: typeof window !== 'undefined' ? window.location.protocol === 'https' : false,
  sameSite: 'Strict'
}

const AUTH_COOKIE_NAME = 'puffer_auth'
const SCOPES_COOKIE_NAME = 'puffer_scopes'
export class ClientCookieSessionStore {
  _cookieOptions = null

  constructor(options = defaultCookieOptions) {
    this._cookieOptions = options
  }

  setToken(token) {
    Cookies.set(AUTH_COOKIE_NAME, token, this._cookieOptions)
  }

  setScopes(scopes) {
    Cookies.set(SCOPES_COOKIE_NAME, JSON.stringify(scopes), this._cookieOptions)
  }

  getToken() {
    return Cookies.get(AUTH_COOKIE_NAME) || null
  }

  getScopes() {
    const res = Cookies.get(SCOPES_COOKIE_NAME)
    if (res) return JSON.parse(res)
    return null
  }

  isLoggedIn() {
    return this.getToken() !== null
  }

  deleteSession() {
    Cookies.remove(AUTH_COOKIE_NAME, this._cookieOptions)
    Cookies.remove(SCOPES_COOKIE_NAME, this._cookieOptions)
  }
}

const AUTH_EXP_COOKIE_NAME = 'puffer_auth_expires'
export class ServerCookieSessionStore {
  _cookieOptions = null

  constructor(options = defaultCookieOptions) {
    this._cookieOptions = options
  }

  setToken(token) {
    throw new Error('It seems you want the ClientCookieSessionStore, not the ServerCookieSessionStore')
  }

  setScopes(scopes) {
    Cookies.set(SCOPES_COOKIE_NAME, JSON.stringify(scopes), this._cookieOptions)
  }

  getToken() {
    return null
  }

  getScopes() {
    const res = Cookies.get(SCOPES_COOKIE_NAME)
    if (res) return JSON.parse(res)
    return null
  }

  isLoggedIn() {
    const res = Cookies.get(AUTH_EXP_COOKIE_NAME)
    return res !== undefined
  }

  deleteSession() {
    Cookies.remove(AUTH_EXP_COOKIE_NAME, this._cookieOptions)
    Cookies.remove(SCOPES_COOKIE_NAME, this._cookieOptions)
  }
}
