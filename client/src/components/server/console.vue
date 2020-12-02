<template>
  <v-card>
    <v-card-title>
      <span v-text="$t('servers.Console')" />
      <div class="flex-grow-1" />
      <v-btn
        icon
        @click="popoutConsole"
      >
        <v-icon
          :dark="isDark()"
          :light="!isDark()"
        >
          mdi-pause
        </v-icon>
      </v-btn>
    </v-card-title>
    <v-card-text>
      <!-- eslint-disable vue/no-v-html -->
      <div
        class="console"
        style="min-height: 25em; max-height: 40em;"
        v-html="console"
      />
      <!-- eslint-enable vue/no-v-html -->
      <v-text-field
        v-if="server.permissions.sendServerConsole"
        v-model="consoleCommand"
        outlined
        hide-details
        placeholder="Command..."
        append-icon="mdi-send"
        class="pt-2"
        @click:append="sendCommand"
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
            <!-- eslint-disable vue/no-v-html -->
            <div
              class="console"
              style="height: 80vh;"
              v-html="console"
            />
            <!-- eslint-enable vue/no-v-html -->
          </v-card-text>
        </v-card>
      </v-overlay>
    </v-card-text>
  </v-card>
</template>

<script>
import AnsiUp from 'ansi_up'
import { isDark } from '@/utils/dark'

const ansiup = new AnsiUp()

export default {
  props: {
    server: { type: Object, default: () => {} }
  },
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
    // ansi up keeps state a little aggressively, so force a reset
    ansiup.ansi_to_html('\u001b[m')

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
          ctx.buffer.push(ansiup.ansi_to_html(element).replaceAll('\n', '<br>'))
        })
      } else {
        this.buffer.push(ansiup.ansi_to_html(data.logs).replaceAll('\n', '<br>'))
      }
    },
    popoutConsole () {
      this.consoleReadonly = this.console
      this.consolePopup = true
      this.$nextTick(() => {
        this.$el.querySelector('#popup .console').scrollTop = this.$el.querySelector('#popup .console').scrollHeight
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

      const console = this.$el.querySelector('.console')
      this.$nextTick(() => {
        console.scrollTop = console.scrollHeight
      })
    },
    sendCommand () {
      this.$socket.sendObj({ type: 'console', command: this.consoleCommand })
      this.consoleCommand = ''
    },
    isDark
  }
}
</script>

<style>
  .console {
      overflow-y: scroll;
      font-size: 1rem;
      font-weight: 400;
      line-height: 1.25;
      font-family: 'Roboto Mono', monospace;
      color: #ddd;
      background-color: black;
      padding: 4px;
  }
</style>
