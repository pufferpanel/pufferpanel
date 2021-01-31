import EventEmitter from 'events'

export class ServerSocket extends EventEmitter {
  _expectClose = false
  _connectionAttemptsFailed = 0
  _connectionAttemptsMax = 5
  _connectionFailReset = null
  _socket = null
  _api = null
  _id = null
  _tasks = []
  _preOpenMessages = []
  readyState = WebSocket.CONNECTING

  constructor (api, id) {
    super()

    if (process.env.NODE_ENV !== 'production') {
      if (!window.pufferpanel.socket) window.pufferpanel.socket = {}
      window.pufferpanel.socket[id] = this

      this.simulateSocketError = (allowRetry = true) => {
        if (!allowRetry) this._connectionAttemptsFailed = this._connectionAttemptsMax
        this._socket.close()
      }
    }

    this._id = id
    this._api = api
    this._openSocket()
    this.emit('open')
  }

  _openSocket () {
    this._socket = new WebSocket(this._api.getServerSocketUrl(this._id))
    this.readyState = this._socket.readyState

    this._socket.addEventListener('open', e => this._onOpen(e))
    this._socket.addEventListener('message', e => this._onMessage(e))
    this._socket.addEventListener('close', e => this._onClose(e))
  }

  _onOpen (e) {
    this.readyState = this._socket.readyState
    this.emit('socket-open', e)

    this._preOpenMessages.forEach(msg => this.send(msg))
  }

  _onMessage (e) {
    this.readyState = this._socket.readyState
    const event = JSON.parse(e.data)

    this.emit('message', event)
    this.emit(event.type, event.data)
  }

  _onClose (e) {
    this.readyState = this._socket.readyState
    this.emit('socket-close', e)

    clearTimeout(this._connectionFailReset)
    if (this._expectClose) {
      this._cleanup()
    } else {
      if (this._connectionAttemptsFailed >= this._connectionAttemptsMax) {
        this._onError({ msg: 'Socket closed unexpectedly', event: e })
      } else {
        this._connectionAttemptsFailed += 1
        setTimeout(() => { this._connectionAttemptsFailed = 0 }, 30000)
        setTimeout(() => this._openSocket(), 500 + (500 * this._connectionAttemptsFailed))
      }
    }
  }

  _onError (e) {
    // eslint-disable-next-line no-console
    console.error('SOCKET ERROR', e)

    this.readyState = this._socket.readyState
    this.emit('error', e)
  }

  _cleanup () {
    this._tasks.forEach(task => clearInterval(task))
  }

  startTask (f, interval) {
    f()
    return this._tasks.push(setInterval(f, interval))
  }

  stopServerTaks (ref) {
    for (const task in this._tasks) {
      if (task === ref) clearInterval(ref)
    }
  }

  needsPolling () {
    const state = this._socket.readyState
    return state === WebSocket.CLOSING || state === WebSocket.CLOSED
  }

  async send (message) {
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
        case 'start':
        case 'stop':
        case 'kill':
        case 'install': {
          this._api.serverAction(this._id, message.type)
          break
        }
        case 'stat': {
          const running = await this._api.getServerStatus(this._id)
          let stats = { cpu: 0, memory: 0 }

          if (running) stats = await this._api.getServerStats(this._id)

          this._onMessage({ data: JSON.stringify({ data: { running }, type: 'status' }) })
          this._onMessage({ data: JSON.stringify({ data: stats, type: 'stat' }) })
          break
        }
        case 'status': {
          const running = await this._api.getServerStatus(this._id)
          this._onMessage({ data: JSON.stringify({ data: { running }, type: 'status' }) })
          break
        }
        case 'replay': {
          const data = await this._api.getServerConsole(this._id, message.since)
          this._onMessage({ data: JSON.stringify({ data, type: 'console' }) })
          break
        }
        case 'file': {
          switch (message.action) {
            case 'get': {
              const files = await this._api.downloadServerFile(this._id, message.path)
              this._onMessage({ data: JSON.stringify({ data: { files, path: message.path }, type: 'file' }) })
              break
            }
            case 'create': {
              await this._api.createServerFolder(this._id, message.path)
              this.send({ type: 'file', action: 'get', path: message.path })
              break
            }
            case 'delete': {
              await this._api.deleteServerFile(this._id, message.path)
              let fetchPath = message.path.split('/')
              fetchPath.pop()
              fetchPath = fetchPath.join('/') || '/'
              this.send({ type: 'file', action: 'get', path: fetchPath })
              break
            }
          }
          break
        }
        case 'console': {
          this._api.serverCommand(this._id, message.command)
          break
        }
        default: {
          // eslint-disable-next-line no-console
          console.error('SOCKET SEND', 'got unknown message', message)
          break
        }
      }
    }
  }

  close () {
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
