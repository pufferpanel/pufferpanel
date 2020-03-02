<template>
  <v-card>
    <v-card-title v-text="$t('servers.Console')" />
    <v-card-text>
      <v-textarea
        id="console"
        v-model="console"
        rows="20"
        background-color="secondary"
        readonly
        solo
        flat
        no-resize
        class="console"
      />
      <v-text-field
        v-model="consoleCommand"
        outlined
        placeholder="Command..."
        append-icon="mdi-send"
        append-outer-icon="mdi-pause"
        @click:append="sendCommand"
        @click:append-outer="popoutConsole"
        @keyup.enter="sendCommand"
      />
      <v-overlay :value="consolePopup">
        <v-card
          :dark="isDark()"
          :light="!isDark()"
          class="d-flex flex-column"
          height="90vh"
          width="90vw"
        >
          <v-card-title>
            <span v-text="$t('servers.Console')" />
            <div class="flex-grow-1" />
            <v-btn
              icon
              @click="consolePopup = false"
            >
              <v-icon>mdi-close</v-icon>
            </v-btn>
          </v-card-title>
          <v-card-text
            id="popup"
            class="flex-grow-1"
          >
            <v-textarea
              id="popupText"
              ref="popup"
              v-model="consoleReadonly"
              style="height: 100%"
              background-color="secondary"
              solo
              flat
              hide-details
              no-resize
              readonly
              class="console"
            />
          </v-card-text>
        </v-card>
      </v-overlay>
    </v-card-text>
  </v-card>
</template>

<script>
import { isDark } from '@/utils/dark'

export default {
  data () {
    return {
      console: '',
      consoleReadonly: '',
      maxConsoleLength: 10000,
      buffer: [],
      refreshInterval: null,
      consoleCommand: '',
      consolePopup: false
    }
  },
  mounted () {
    const ctx = this
    this.$socket.addEventListener('message', event => {
      const data = JSON.parse(event.data)
      if (data === 'undefined') {
        return
      }
      if (data.type === 'console') {
        ctx.parseConsole(data.data)
      }
    })
    this.$socket.addEventListener('open', event => {
      this.$socket.sendObj({ type: 'replay' })
    })
    this.refreshInterval = setInterval(this.updateConsole, 1000)
  },
  beforeDestroy () {
    if (this.refreshInterval !== null) {
      clearInterval(this.refreshInterval)
    }
  },
  methods: {
    parseConsole (data) {
      const ctx = this

      if (data.logs instanceof Array) {
        data.logs.forEach(element => {
          ctx.buffer.push(element)
        })
      } else {
        this.buffer.push(data.logs)
      }
    },
    popoutConsole () {
      this.consoleReadonly = this.console
      this.consolePopup = true
      this.$nextTick(() => {
        this.$refs.popup.$el.style.height = '100%'
        this.$refs.popup.$el.children[0].style.height = '100%'
        this.$refs.popup.$el.children[0].children[0].style.height = '100%'
        this.$refs.popup.$el.children[0].children[0].children[0].style.height = '100%'
        this.$el.querySelector('#popupText').style.height = '100%'
        this.$el.querySelector('#popupText').scrollTop = this.$el.querySelector('#popupText').scrollHeight
      })
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

      const textArea = this.$el.querySelector('#console')
      this.$nextTick(() => {
        textArea.scrollTop = textArea.scrollHeight
      })
    },
    sendCommand () {
      if (this.consoleCommand.length === 0) {
        return
      }

      this.$socket.sendObj({ type: 'console', command: this.consoleCommand })

      this.consoleCommand = ''
    },
    isDark
  }
}
</script>

<style>
  .v-textarea.console textarea {
    line-height: 1.25;
  }
  #popup .v-input__slot {
    height: 100%;
  }
</style>
