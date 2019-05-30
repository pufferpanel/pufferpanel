<template>
  <b-card
    :title="$t('common.Console')">
    <textarea ref="console" class="form-control console" readonly="readonly" v-text="console"></textarea>
  </b-card>
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
      statTracker: null,
      consoleRecover: null
    }
  },
  methods: {
    openConnection () {
      this.consoleRecover = null
      console.log('opening connection')

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

            root.statTracker = setInterval(this.callStats, 10 * 1000)
          })
        })
        this.connection.addEventListener('error', function (event) {
          console.log(event)
          root.connection = null

          root.consoleRecover = setTimeout(root.openConnection, 10 * 1000)
          console.log('recover scheduled as ' + root.consoleRecover)
        })
      } catch (ex) {
        console.log(ex)
        this.connection = null

        this.consoleRecover = setTimeout(this.openConnection, 10 * 1000)
        console.log('recover scheduled as ' + this.consoleRecover)
      }
    },
    callStats () {
      if (this.connection) {
        this.connection.send(JSON.stringify({
          'type': 'statsRequest'
        }))
      }
    },
    parseStats (data) {

    },
    parseConsole (data) {
      let textArea = this.$refs['console']

      let msg = ''
      if (data.logs instanceof Array) {
        data.logs.forEach(function (element) {
          msg += element
        })
      } else {
        msg = data.logs
      }

      this.console = this.console + msg
      textArea.scrollTop = textArea.scrollHeight
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

    if (this.consoleRecover) {
      clearTimeout(this.consoleRecover)
    }
  }
}
</script>

<style>
  .console {
    font: 85% 'Droid Sans Mono', monospace;
    color: #333;
    height: 300px !important;
    text-wrap: normal;
    overflow-y: scroll;
    overflow-x: hidden;
    border: 0;
    resize: none
  }

  .console[readonly=readonly] {
    background: #fefefe;
    cursor: default
  }
</style>
