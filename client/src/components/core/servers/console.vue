<template>
  <b-card
    header-tag="header"
    footer-tag="footer">
    <h6 slot="header" class="mb-0" v-text="$t('common.Console')"></h6>
    <textarea ref="console" class="form-control console" readonly="readonly" v-text="console"></textarea>
    <b-input-group class="mt-3">
      <b-form-input v-model="consoleCommand" placeholder="Command..." @keyup.enter="sendCommand"></b-form-input>
      <b-input-group-append>
        <b-button @click="sendCommand" variant="info">Send</b-button>
      </b-input-group-append>
    </b-input-group>

    <b-btn slot="footer" v-b-modal.console-copy v-text="$t('common.Pause')" @click="popoutConsole"></b-btn>
    <b-modal id="console-copy" size="xl" v-bind:title="$t('common.Console')">
      <textarea ref="console" class="form-control console" readonly="readonly" v-text="consoleReadonly"></textarea>
    </b-modal>
  </b-card>
</template>

<script>
export default {
  data () {
    return {
      console: '',
      consoleReadonly: '',
      maxConsoleLength: 10000,
      buffer: [],
      refreshInterval: null,
      consoleCommand: ''
    }
  },
  methods: {
    parseConsole (data) {
      let vue = this

      if (data.logs instanceof Array) {
        data.logs.forEach(function (element) {
          vue.buffer.push(element)
        })
      } else {
        this.buffer.push(data.logs)
      }
    },
    popoutConsole () {
      this.consoleReadonly = this.console
    },
    updateConsole () {
      if (this.buffer.length === 0) {
        return
      }

      let msg = this.buffer.shift()
      while (this.buffer.length > 0) {
        msg += this.buffer.shift()
      }

      let newConsole = this.console + msg
      if (newConsole.length > this.maxConsoleLength) {
        newConsole = newConsole.substring(newConsole.length - this.maxConsoleLength, newConsole.length)
      }
      this.console = newConsole

      let textArea = this.$refs['console']
      this.$nextTick(function () {
        textArea.scrollTop = textArea.scrollHeight
      })
    },
    sendCommand () {
      if (this.consoleCommand.length === 0) {
        return
      }

      this.$socket.sendObj({type: 'console', command: this.consoleCommand})

      this.consoleCommand = ''
    }
  },
  mounted () {
    let root = this
    this.$socket.addEventListener('message', function (event) {
      let data = JSON.parse(event.data)
      if (data === 'undefined') {
        return
      }
      if (data.type === 'console') {
        root.parseConsole(data.data)
      }
    })
    this.refreshInterval = setInterval(this.updateConsole, 1000)
  },
  beforeDestroy () {
    if (this.refreshInterval !== null) {
      clearInterval(this.refreshInterval)
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
