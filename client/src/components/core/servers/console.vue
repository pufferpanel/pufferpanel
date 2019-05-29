<template>
  <div>

  </div>
</template>

<script>
export default {
  props: {
    server: Object
  },
  data() {
    return {
      console: '',
      connection: null,
      statTracker: null
    }
  },
  methods: {
    openConnection() {
      try {
        let base = location.protocol === 'https' ? 'wss://' : 'ws:/' + location.host;
        this.connection = new WebSocket(base + '/daemon/server/' + this.server.id + '/console')
        this.connection.onopen = function () {
          this.connection.onmessage = function (event) {
            let data = JSON.parse(event.data)
            if (data === 'undefined') {
              return
            }
            switch (data.type) {
              case 'console': {
                this.parseConsole(data.data)
                break
              }
              case 'stat': {
                this.parseStats(data.data)
                break
              }
            }

            this.statTracker = setInterval(this.callStats, 10000)
          }
        }
        this.connection.onerror = function (event) {
          console.log(event)
        }
      } catch (ex) {
        console.log(ex)
        this.connection = null
      }
    },
    callStats() {
      if (this.connection) {
        this.connection.send(JSON.stringify({
          "type": "statsRequest"
        }))
      } else {
        this.createRequest().get("/daemon/server/" + server.id + "/stats", {timeout: 1000}).then(function (response) {
          console.log(data)
          this.parseStats(response.data.data)
        }).catch(function (error) {
          console.log(error)
        })
      }
    },
    parseStats(data) {

    },
    parseConsole(data) {
      console = console + data.logs
    }
  },
  mounted() {
    this.openConnection()
  },
  beforeDestroy: function () {
    if (this.connection) {
      if (this.connection instanceof WebSocket) {
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