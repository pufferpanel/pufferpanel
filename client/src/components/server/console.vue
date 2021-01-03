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
      <ui-input
        v-if="server.permissions.sendServerConsole || isAdmin()"
        v-model="consoleCommand"
        class="pt-2"
        placeholder="Command..."
        end-icon="mdi-send"
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
      maxConsoleLength: 1000000,
      consoleCommand: '',
      consolePopup: false,
      lastConsoleTime: 0
    }
  },
  mounted () {
    // ansi up keeps state a little aggressively, so force a reset
    ansiup.ansi_to_html('\u001b[m')

    this.$api.addServerListener(this.server.id, 'console', event => {
      if ('epoch' in event) {
        this.lastConsoleTime = event.epoch
      } else {
        this.lastConsoleTime = Math.floor(Date.now() / 1000)
      }

      this.parseConsole(event)
    })

    this.$api.startServerTask(this.server.id, () => {
      if (this.$api.serverConnectionNeedsPolling(this.server.id)) {
        this.$api.requestServerConsoleReplay(this.server.id, this.lastConsoleTime)
      }
    }, 5000)

    this.$api.requestServerConsoleReplay(this.server.id)
  },
  methods: {
    parseConsole (data) {
      let newConsole = this.console
      const logs = Array.isArray(data.logs) ? data.logs : [data.logs]
      logs.forEach(element => {
        newConsole += ansiup.ansi_to_html(element).replace(/\n/g, '<br>')
      })

      if (newConsole.length > this.maxConsoleLength) {
        newConsole = newConsole.substring(newConsole.length - this.maxConsoleLength, newConsole.length)
      }
      this.console = newConsole

      const console = this.$el.querySelector('.console')
      this.$nextTick(() => {
        console.scrollTop = console.scrollHeight
      })
    },
    popoutConsole () {
      this.consoleReadonly = this.console
      this.consolePopup = true
      this.$nextTick(() => {
        this.$el.querySelector('#popup .console').scrollTop = this.$el.querySelector('#popup .console').scrollHeight
      })
    },
    sendCommand () {
      this.$api.sendServerCommand(this.server.id, this.consoleCommand)
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
