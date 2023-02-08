export class ServerApi {
  _api = null

  constructor(api) {
    this._api = api
  }

  async create(data) {
    const res = await this._api.post('/api/servers', data)
    return res.data.id
  }

  async list(page = 1) {
    const res = await this._api.get('/api/servers', { page })
    return res.data
  }

  async get(id, withSocket = true) {
    const res = await this._api.get(`/api/servers/${id}?perms`)
    if (withSocket) {
      return new Server(this._api, res.data)
    } else {
      return res.data
    }
  }

  async getStatus(id) {
    const res = await this._api.get(`/api/servers/${id}/status`)
    return res.data.running
  }

  async action(id, action, wait = false) {
    await this._api.post(`/api/servers/${id}/${action}?wait=${wait}`)
    return true
  }

  async start(id, wait = false) {
    return await this.action(id, 'start', wait)
  }

  async stop(id, wait = false) {
    return await this.action(id, 'stop', wait)
  }

  async kill(id, wait = false) {
    return await this.action(id, 'kill', wait)
  }

  async install(id, wait = false) {
    return await this.action(id, 'install', wait)
  }

  async sendCommand(id, command) {
    await this._api.post(`/api/servers/${id}/console`, command)
    return true
  }

  async getConsole(id, time = 0) {
    const res = await this._api.get(`/api/servers/${id}/console?time=${time}`)
    return res.data
  }

  async updateName(id, name) {
    await this._api.put(`/api/servers/${id}/name/${encodeURIComponent(name)}`)
    return true
  }

  async getDefinition(id) {
    const res = await this._api.get(`/api/servers/${id}`)
    return res.data
  }

  async updateDefinition(id, data) {
    await this._api.post(`/api/servers/${id}`, data)
    return true
  }

  async getData(id) {
    const res = await this._api.get(`/api/servers/${id}/data`)
    return res.data.data
  }

  async updateData(id, data) {
    await this._api.post(`/api/servers/${id}/data`, { data })
    return true
  }

  async getUsers(id) {
    const res = await this._api.get(`/api/servers/${id}/user`)
    return res.data
  }

  async updateUser(id, user) {
    await this._api.put(`/api/servers/${id}/user/${user.email}`, user)
    return true
  }

  async deleteUser(id, email) {
    await this._api.delete(`/api/servers/${id}/user/${email}`)
    return true
  }

  getFileUrl (id, path) {
    if (path.indexOf('/') === 0) path = path.substring(1)
    path = encodeURIComponent(path)
    path = path.replace(/%2F/g, '/')
    return `/api/servers/${id}/file/${path}`
  }

  async getFile(id, path = '', raw = false) {
    if (path.indexOf('/') === 0) path = path.substring(1)
    const res = await this._api.get(
      this.getFileUrl(id, path),
      undefined,
      undefined,
      raw ? { responseType: 'arraybuffer' } : undefined
    )

    if (raw) {
      return new TextDecoder('utf-8').decode(new Uint8Array(res.data))
    } else {
      return res.data
    }
  }

  async fileExists(id, path) {
    if (path.indexOf('/') === 0) path = path.substring(1)
    try {
      const res = await this._api._axios.get(this._api._host + this.getFileUrl(id, path))

      if (res.headers['content-disposition']) {
        return 'file'
      } else {
        return 'folder'
      }
    } catch (e) {
      if (e.response.status === 404) return false
      this._api._handleError(e)
    }
  }

  async uploadFile(id, path, content, onUploadProgress) {
    if (path.indexOf('/') === 0) path = path.substring(1)
    let blob = null
    if (content instanceof Blob || content instanceof File) {
      blob = content
    } else if (typeof content === 'string') {
      blob = new Blob([content])
    } else {
      blob = new Blob([JSON.stringify(content)])
    }

    const data = new FormData()
    data.append('file', blob)

    await this._api.put(
      this.getFileUrl(id, path),
      data,
      undefined,
      { 'Content-Type': 'multipart/formdata' },
      { onUploadProgress }
    )
    return true
  }

  async createFolder(id, path) {
    if (path.indexOf('/') === 0) path = path.substring(1)
    await this._api.put(
      this.getFileUrl(id, path),
      undefined,
      { folder: true }
    )
    return true
  }

  async archiveFile(id, destination, files) {
    if (destination.startsWith('/')) destination = destination.substring(1)
    if (!Array.isArray(files)) files = [files]
    files.map(file => file.startsWith('/') ? file.substring(1) : file)
    await this._api.post(`/api/servers/${id}/archive/${destination}`, files)
    return true
  }

  async extractFile(id, path, destination) {
    if (path.startsWith('/')) path = path.substring(1)
    await this._api.get(`/api/servers/${id}/extract/${path}`, { destination })
    return true
  }

  async deleteFile(id, path) {
    if (path.indexOf('/') === 0) path = path.substring(1)
    await this._api.delete(this.getFileUrl(id, path))
    return true
  }

  async getOAuthClients(id) {
    const res = await this._api.get(`/api/servers/${id}/oauth2`)
    return res.data
  }

  async createOAuthClient(id, name, description) {
    const res = await this._api.post(`/api/servers/${id}/oauth2`, { name, description })
    return res.data
  }

  async deleteOAuthClient(id, clientId) {
    await this._api.delete(`/api/servers/${id}/oauth2/${clientId}`)
    return true
  }

  async delete(id) {
    await this._api.delete(`/api/servers/${id}`)
    return true
  }
}

class Server {
  _expectClose = false
  _connectionAttemptsFailed = 0
  _connectionAttemptsMax = 5
  _connectionFailReset = null
  _socket = null
  _api = null
  _tasks = []
  _preOpenMessages = []
  _emitter = null
  readyState = WebSocket.CONNECTING

  constructor(api, serverData) {
    // inlined https://github.com/ai/nanoevents because just depending on it breaks nodejs somehow...
    this._emitter = {
      events: {},
      emit(event, ...args) {
        let callbacks = this.events[event] || []
        for (let i = 0, length = callbacks.length; i < length; i++) {
          callbacks[i](...args)
        }
      },
      on(event, cb) {
        this.events[event]?.push(cb) || (this.events[event] = [cb])
        return () => {
          this.events[event] = this.events[event]?.filter(i => cb !== i)
        }
      }
    }

    this.id = serverData.server.id
    this.ip = serverData.server.id
    this.name = serverData.server.name
    this.node = serverData.server.node
    this.port = serverData.server.port
    this.type = serverData.server.type
    this.permissions = api.auth.hasScope('servers.admin') ? {
      editServerData: true,
      editServerUsers: true,
      installServer: true,
      putServerFiles: true,
      sendServerConsole: true,
      sftpServer: true,
      startServer: true,
      stopServer: true,
      viewServerConsole: true,
      viewServerFiles: true,
      viewServerStats: true
    } : serverData.permissions
    this._api = api
    this._openSocket()
    this.emit('open')
  }

  on(event, cb) {
    return this._emitter.on(event, cb)
  }

  emit(event, data) {
    return this._emitter.emit(event, data)
  }

  _openSocket() {
    let host = this._api._host
    if (!host && typeof window !== 'undefined') {
      host = window.location.host
    }
    if (!host) throw new Error('cannot determine host to connect to')
    const protocol = host.indexOf('https://') === 0 ? 'wss' : 'ws'
    if (host.indexOf('http://') === 0) host = host.substr(7)
    if (host.indexOf('https://') === 0) host = host.substr(8)

    this._socket = new WebSocket(`${protocol}://${host}/api/servers/${this.id}/socket`)
    this.readyState = this._socket.readyState

    this._socket.addEventListener('open', e => this._onOpen(e))
    this._socket.addEventListener('message', e => this._onMessage(e))
    this._socket.addEventListener('close', e => this._onClose(e))
  }

  _onOpen(e) {
    this.readyState = this._socket.readyState
    this.emit('socket-open', e)

    this._preOpenMessages.forEach(msg => this.send(msg))
  }

  _onMessage(e) {
    this.readyState = this._socket.readyState
    const event = JSON.parse(e.data)

    this.emit('message', event)
    this.emit(event.type, event.data)
  }

  _onClose(e) {
    this.readyState = this._socket.readyState
    this.emit('socket-close', e)

    clearTimeout(this._connectionFailReset)
    if (this._expectClose) {
      this._cleanup()
    } else {
      console.warn('socket closed', e)
      if (this._connectionAttemptsFailed === this._connectionAttemptsMax) {
        // emit an error once after a certain number of failed retries
        // then keep retrying without emitting more errors until it works
        this._onError({ msg: 'Socket closed unexpectedly', event: e })
      }
      this._connectionAttemptsFailed += 1
      this._connectionFailReset = setTimeout(() => { this._connectionAttemptsFailed = 0 }, 30000)
      setTimeout(() => this._openSocket(), 5000)
    }
  }

  _onError(e) {
    // eslint-disable-next-line no-console
    console.error('SOCKET ERROR', e)

    this.readyState = this._socket.readyState
    this.emit('error', e)
  }

  _cleanup() {
    this._tasks.forEach(task => clearInterval(task))
  }

  startTask(f, interval) {
    f()
    const task = setInterval(f, interval)
    this._tasks.push(task)
    return task
  }

  stopTask(ref) {
    for (const task of this._tasks) {
      if (task === ref) clearInterval(ref)
    }
  }

  needsPolling() {
    const state = this._socket.readyState
    return state === WebSocket.CLOSING || state === WebSocket.CLOSED
  }

  async send(message) {
    let msg
    if (typeof message !== 'string') {
      msg = JSON.stringify(message)
    } else {
      msg = message
    }

    if (this._socket.readyState === WebSocket.CONNECTING) {
      this._preOpenMessages.push(message)
    } else if (this._socket.readyState === WebSocket.OPEN) {
      this._socket.send(msg)
    } else {
      // replicate socket behavior through http
      switch (message.type) {
        case 'status':
          const running = await this._api.server.getStatus(this.id)
          this._onMessage({ data: JSON.stringify({ data: { running }, type: 'status' }) })
          break
        case 'start':
        case 'stop':
        case 'kill':
        case 'install':
          this._api.server.action(this.id, message.type)
          break
        case 'console':
          this._api.server.sendCommand(this.id, message.command)
          break
        case 'replay':
          this._api.server.getConsole(this.id, message.since)
          break
        default:
          console.error('SOCKET SEND', 'unknown message', message)
      }
    }
  }

  status() {
    this.send({ type: 'status' })
  }

  stats() {
    this.send({ type: 'stat' })
  }

  start() {
    this.send({ type: 'start' })
  }

  stop() {
    this.send({ type: 'stop' })
  }

  kill() {
    this.send({ type: 'kill' })
  }

  install() {
    this.send({ type: 'install' })
  }

  sendCommand(command) {
    this.send({ type: 'console', command })
  }

  replayConsole(since = 0) {
    this.send({ type: 'replay', since })
  }

  async updateName(name) {
    const r = await this._api.server.updateName(this.id, name)
    this.name = name
    return r
  }

  async getDefinition() {
    return await this._api.server.getDefinition(this.id)
  }

  async updateDefinition(data) {
    return await this._api.server.updateDefinition(this.id, data)
  }

  async getData() {
    return await this._api.server.getData(this.id)
  }

  async updateData(data) {
    return await this._api.server.updateData(this.id, data)
  }

  async delete() {
    return await this._api.server.delete(this.id)
  }

  async getUsers() {
    return await this._api.server.getUsers(this.id)
  }

  async updateUser(user) {
    return await this._api.server.updateUser(this.id, user)
  }

  async deleteUser(email) {
    return await this._api.server.deleteUser(this.id, email)
  }

  getFileUrl(path) {
    return this._api.server.getFileUrl(this.id, path)
  }

  async getFile(path = '', raw = false) {
    return await this._api.server.getFile(this.id, path, raw)
  }

  async fileExists(path) {
    return await this._api.server.fileExists(this.id, path)
  }

  async uploadFile(path, content, onUploadProgress) {
    return await this._api.server.uploadFile(this.id, path, content, onUploadProgress)
  }

  async createFolder(path) {
    return await this._api.server.createFolder(this.id, path)
  }

  async archiveFile(destination, files) {
    return await this._api.server.archiveFile(this.id, destination, files)
  }

  async extractFile(path, destination) {
    return await this._api.server.extractFile(this.id, path, destination)
  }

  async deleteFile(path) {
    return await this._api.server.deleteFile(this.id, path)
  }

  async getOAuthClients() {
    return await this._api.server.getOAuthClients(this.id)
  }

  async createOAuthClient(name, description) {
    return await this._api.server.createOAuthClient(this.id, name, description)
  }

  async deleteOAuthClient(clientId) {
    return await this._api.server.deleteOAuthClient(this.id, clientId)
  }

  close() {
    this._expectClose = true
    this._cleanup()
    if (this.readyState === WebSocket.CONNECTING || this.readyState === WebSocket.OPEN) {
      this._socket.addEventListener('close', () => this.emit('close'))
      this._socket.close()
    } else {
      this.emit('close')
    }
  }
}
