<template>
  <div>

  </div>
</template>

<script>
export default {
  props: {
    server: Object
  },
  data () {
    return {
      console: '',
      connection: null,
      statTracker: null
    }
  },
  methods: {
    openConnection () {
      try {
        let root = this
        let base = location.protocol === 'https' ? 'wss://' : 'ws:/' + location.host
        this.connection = new WebSocket(base + '/daemon/server/' + this.server.id + '/console')
        this.connection.addEventListener('open', function () {
          root.connection.addEventListener('message', function (event) {
            let data = JSON.parse(event.data)
            if (data === 'undefined') {
              return
            }
            switch (data.type) {
              case 'console': {
                root.parseConsole(data.data)
                break
              }
              case 'stat': {
                root.parseStats(data.data)
                break
              }
            }

            root.statTracker = setInterval(this.callStats, 10000)
          })
        })
        this.connection.addEventListener('error', function (event) {
          console.log(event)
        })
      } catch (ex) {
        console.log(ex)
        this.connection = null
      }
    },
    callStats () {
      if (this.connection) {
        this.connection.send(JSON.stringify({
          'type': 'statsRequest'
        }))
      } else {
        this.$http.get('/daemon/server/' + server.id + '/stats', { timeout: 1000 }).then(function (response) {
          console.log(data)
          this.parseStats(response.data.data)
        }).catch(function (error) {
          console.log(error)
        })
      }
    },
    parseStats (data) {

    },
    parseConsole (data) {
      this.console = this.console + data.logs
    }
  },
  mounted () {
    this.openConnection()
  },
  beforeDestroy: function () {
    if (this.connection) {
      if (this.connection.close) {
        console.log('Closing websocket')
        this.connection.close()
      } else {
        clearInterval(this.connection)
      }
    }

    if (this.statTracker) {
      clearInterval(this.statTracker)
    }
  }
}
</script>