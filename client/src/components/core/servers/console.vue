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
      consoleTracker: null,
      statTracker: null
    }
  },
  methods: {
    openConnection () {
      try {
        this.consoleTracker = new WebSocket('/daemon/server/' + this.server.id + '/console')
        this.consoleTracker.onopen = function () {
          this.consoleTracker.onmessage = function (event) {
            let data = JSON.parse(event.data)
            if (data === 'undefined') {
              return
            }
            switch (data.type) {
              case 'console': {
                this.parseConsole(data)
                break
              }
              case 'stat': {
                this.parseStats(data)
                break
              }
            }

            this.statTracker = setInterval(this.callStats, 1000)
          }
        }
        this.consoleTracker.onerror = function (event) {
          console.log(event)
        }
      } catch (ex) {
        console.log(ex)
        this.consoleTracker = null
      }
    },
    callStats () {

    },
    parseStats (data) {

    },
    parseConsole (data) {
    
    }
  },
  mounted () {
    this.openConnection()
  },
  beforeDestroy: function () {
    if (this.consoleTracker) {
      if (this.consoleTracker instanceof WebSocket) {
        this.consoleTracker.close()
      } else {
        clearInterval(this.consoleTracker)
      }
    }

    if (this.statTracker) {
      clearInterval(this.statTracker)
    }
  }
}
</script>