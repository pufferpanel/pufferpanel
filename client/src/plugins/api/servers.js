import { ServerSocket } from './serverSocket'

// used as a mixin for the ApiClient, use of `this` refers to the ApiClient instance
export const ServersApi = {
  startServerConnection (id) {
    if (this._sockets[id]) return
    this._sockets[id] = new ServerSocket(this, id)
    return Promise.resolve()
  },

  serverConnectionNeedsPolling (id) {
    if (!this._sockets[id]) throw new Error('connection not opened')
    return this._sockets[id].needsPolling()
  },

  addServerListener (id, event, listener) {
    if (!this._sockets[id]) throw new Error('connection not opened')
    this._sockets[id].on(event, listener)
  },

  removeServerListener (id, event, listener) {
    if (!this._sockets[id]) return
    this._sockets[id].off(event, listener)
  },

  startServerTask (id, f, interval) {
    if (!this._sockets[id]) throw new Error('connection not opened')
    this._sockets[id].startTask(f, interval)
  },

  stopServerTask (id, ref) {
    if (!this._sockets[id]) return
    this._sockets[id].stopTask(ref)
  },

  sendToServer (id, message) {
    if (!this._sockets[id]) throw new Error('connection not opened')
    this._sockets[id].send(message)
  },

  closeServerConnection (id) {
    if (this._sockets[id]) {
      this._sockets[id].on('close', () => {
        delete this._sockets[id]
      })
      this._sockets[id].close()
    }
  },

  // websocket actions

  requestServerStats (id) {
    this.sendToServer(id, { type: 'stat' })
  },

  requestServerStatus (id) {
    this.sendToServer(id, { type: 'status' })
  },

  requestServerConsoleReplay (id, since = 0) {
    this.sendToServer(id, { type: 'replay', since })
  },

  sendServerCommand (id, command) {
    this.sendToServer(id, { type: 'console', command })
  },

  requestServerFile (id, path) {
    this.sendToServer(id, { type: 'file', action: 'get', path })
  },

  requestServerFolderCreation (id, path) {
    this.sendToServer(id, { type: 'file', action: 'create', path })
  },

  requestServerFileDeletion (id, path) {
    this.sendToServer(id, { type: 'file', action: 'delete', path })
  },

  async sendServerAction (id, action) {
    if (action === 'restart') {
      await this.serverAction(id, 'stop', true)
      action = 'start'
    }

    this.sendToServer(id, { type: action })
  },

  // http actions

  getServers (page, pageSize = 10) {
    return this.withErrorHandling(async ctx => {
      const res = (await ctx.$http.get('/api/servers', { params: { page, limit: pageSize } })).data
      const servers = res.servers.map(server => {
        let ip = ''

        if (server.ip && server.ip !== '' && server.ip !== '0.0.0.0') {
          ip = server.ip
        } else {
          ip = server.node.publicHost
        }

        if (server.port) {
          ip += ':' + server.port
        }

        return {
          id: server.id,
          name: server.name,
          node: server.node.name,
          address: ip
        }
      })
      return { servers, pages: Math.ceil(res.paging.total / pageSize) }
    })
  },

  createServer (data) {
    return this.withErrorHandling(async ctx => {
      const res = (await ctx.$http.post('/api/servers', data)).data
      return res.id
    })
  },

  getServer (id) {
    return this.withErrorHandling(async ctx => {
      const res = (await ctx.$http.get(`/api/servers/${id}?perms`)).data
      return { ...res.server, permissions: res.permissions }
    })
  },

  getServerDefinition (id) {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.get(`/proxy/daemon/server/${id}`)).data
    })
  },

  getServerData (id) {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.get(`/proxy/daemon/server/${id}/data`)).data.data
    })
  },

  getServerTasks (id) {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.get(`/proxy/daemon/server/${id}/tasks`)).data.tasks
    })
  },

  getServerSocketUrl (id) {
    const protocol = window.location.protocol === 'http:' ? 'ws' : 'wss'
    return `${protocol}://${window.location.host}/proxy/daemon/socket/${id}`
  },

  getServerStatus (id, options) {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.get(`/proxy/daemon/server/${id}/status`)).data.running
    }, options)
  },

  getServerStats (id) {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.get(`/proxy/daemon/server/${id}/stats`)).data
    })
  },

  getServerConsole (id, time = 0) {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.get(`/proxy/daemon/server/${id}/console?time=${time}`)).data
    })
  },

  getServerUsers (id) {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.get(`/api/servers/${id}/user`)).data
    })
  },

  getServerFileUrl (id, path) {
    if (path.indexOf('/') === 0) path = path.substring(1)
    path = encodeURIComponent(path)
    path = path.replace(/%2F/g, '/')
    return `/proxy/daemon/server/${id}/file/${path}`
  },

  downloadServerFile (id, path, asUtf8 = false) {
    return this.withErrorHandling(async ctx => {
      const data = (await ctx.$http.get(
        this.getServerFileUrl(id, path),
        asUtf8 ? { responseType: 'arraybuffer' } : undefined
      )).data
      if (asUtf8) {
        return new TextDecoder('utf-8').decode(new Uint8Array(data))
      } else {
        return data
      }
    })
  },

  uploadServerFile (id, path, content, onUploadProgress) {
    return this.withErrorHandling(async ctx => {
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

      await ctx.$http.put(
        this.getServerFileUrl(id, path),
        data,
        { headers: { 'Content-Type': 'multipart/form-data' }, onUploadProgress }
      )
      return true
    })
  },

  createServerFolder (id, path) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.put(this.getServerFileUrl(id, path), undefined, { params: { folder: true } })
      return true
    })
  },

  deleteServerFile (id, path) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.delete(this.getServerFileUrl(id, path))
      return true
    })
  },

  archiveServerFiles (id, destination, files) {
    if (destination.startsWith('/')) destination = destination.substring(1)
    if (!Array.isArray(files)) files = [files]
    files.map(file => {
      return file.startsWith('/') ? file.substring(1) : file
    })

    return this.withErrorHandling(async ctx => {
      await ctx.$http.post(`/proxy/daemon/server/${id}/archive/${destination}`, files)
      return true
    })
  },

  extractServerFile (id, path, destination) {
    if (path.startsWith('/')) path = path.substring(1)
    return this.withErrorHandling(async ctx => {
      await ctx.$http.get(`/proxy/daemon/server/${id}/extract/${path}`, { params: { destination } })
      return true
    })
  },

  serverAction (id, action, wait = false) {
    return this.withErrorHandling(async ctx => {
      await this._ctx.$http.post(`/proxy/daemon/server/${id}/${action}?wait=${wait}`)
    })
  },

  serverCommand (id, command) {
    return this.withErrorHandling(async ctx => {
      await this._ctx.$http.post(`/proxy/daemon/server/${id}/console`, command)
    })
  },

  updateServerName (id, newName) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.put(`/api/servers/${id}/name/${encodeURIComponent(newName)}`)
      return true
    })
  },

  updateServerDefinition (id, data) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.post(`/proxy/daemon/server/${id}`, data)
      return true
    })
  },

  updateServerData (id, data) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.post(`/proxy/daemon/server/${id}/data`, { data })
      return true
    })
  },

  createServerTask (id, task) {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.post(`/proxy/daemon/server/${id}/tasks`, task)).id
    })
  },

  runServerTask (id, taskId) {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.post(`/proxy/daemon/server/${id}/tasks/${taskId}/run`)).id
    })
  },

  editServerTask (serverId, taskId, task) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.put(`/proxy/daemon/server/${serverId}/tasks/${taskId}`, task)
      return true
    })
  },

  deleteServerTask (serverId, taskId) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.delete(`/proxy/daemon/server/${serverId}/tasks/${taskId}`)
      return true
    })
  },

  updateServerUser (id, user) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.put(`/api/servers/${id}/user/${user.email}`, user)
      return true
    })
  },

  reloadServer (id) {
    this.withErrorHandling(async ctx => {
      await ctx.$http.post(`/proxy/daemon/server/${id}/reload`)
      return true
    })
  },

  deleteServer (id) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.delete(`/api/servers/${id}`)
      return true
    })
  },

  deleteServerUser (id, email) {
    return this.withErrorHandling(async ctx => {
      await ctx.$http.delete(`/api/servers/${id}/user/${email}`)
      return true
    })
  },

  getServerOAuthClients (id) {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.get(`/api/servers/${id}/oauth2`)).data
    })
  },

  createServerOAuthClient (id, name, description) {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.post(`/api/servers/${id}/oauth2`, { name, description })).data
    })
  },

  deleteServerOAuthClient (id, clientId) {
    return this.withErrorHandling(async ctx => {
      return (await ctx.$http.delete(`/api/servers/${id}/oauth2/${clientId}`)).data
    })
  }
}
