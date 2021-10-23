<template>
  <v-card>
    <v-card-title ref="title">
      <span v-text="$t('servers.Console')" />
      <div class="flex-grow-1" />
      <v-btn
        v-hotkey="'c p'"
        icon
        @click="togglePaused()"
      >
        <v-badge
          color="error"
          :value="paused && hasNewLines"
          overlap
          dot
        >
          <v-icon>
            {{ paused ? 'mdi-play' : 'mdi-pause' }}
          </v-icon>
        </v-badge>
      </v-btn>
    </v-card-title>
    <v-card-text>
      <span style="display: none;">
        <v-chip
          ref="daemonChip"
          :color="$vuetify.theme.options.console.daemonChip"
          x-small
          label
        >
          DAEMON
        </v-chip>
      </span>
      <div
        ref="consoleEl"
        class="console"
        :style="'min-height: 25em; max-height: 36em;'"
      >
      </div>
      <ui-input
        v-if="server.permissions.sendServerConsole || isAdmin()"
        ref="cmdInput"
        v-model="consoleCommand"
        v-hotkey="'c i'"
        class="pt-2"
        placeholder="Command..."
        end-icon="mdi-send"
        @click:append="sendCommand"
        @keyup.enter="sendCommand"
      />
    </v-card-text>
  </v-card>
</template>

<script>
import AnsiUp from 'ansi_up'

const DAEMON_MESSAGE_REGEX = /^(&nbsp;|&gt;|\s)*\[DAEMON]/
const CONSOLE_REFRESH_TIME = 500
const CONSOLE_MEMORY_ALLOWED = 1024 * 1024 * 4 // 4MB
// due to pausing we have 2 copies of the console,
// so each should only get half the allowed memory
const CONSOLE_MEMORY_ALLOWED_PER_BUFFER = CONSOLE_MEMORY_ALLOWED / 2

const ansiup = new AnsiUp()
let lines = []

export default {
  props: {
    server: { type: Object, default: () => {} }
  },
  data () {
    return {
      console: [],
      consoleCommand: '',
      lastConsoleTime: 0,
      paused: false,
      hasNewLines: false,
      interval: null,
      listener: null
    }
  },
  mounted () {
    // ansi up keeps state a little aggressively, so force a reset
    ansiup.ansi_to_html('\u001b[m')

    this.interval = setInterval(() => {
      if (this.paused) return
      if (this.hasNewLines) this.$refs.consoleEl.replaceChildren(...lines.map(line => line.el))

      const consoleEl = this.$refs.consoleEl
      this.$nextTick(() => {
        if (this.paused) return
        if (!this.hasNewLines) return
        consoleEl.scrollTop = consoleEl.scrollHeight
        this.hasNewLines = false
      })
    }, CONSOLE_REFRESH_TIME)

    this.listener = event => {
      if ('epoch' in event) {
        this.lastConsoleTime = event.epoch
      } else {
        this.lastConsoleTime = Math.floor(Date.now() / 1000)
      }

      this.parseConsole(event)
    }

    this.$api.addServerListener(this.server.id, 'console', this.listener)

    this.$api.startServerTask(this.server.id, () => {
      if (this.$api.serverConnectionNeedsPolling(this.server.id)) {
        this.$api.requestServerConsoleReplay(this.server.id, this.lastConsoleTime)
      }
    }, 5000)

    this.$api.requestServerConsoleReplay(this.server.id)
  },
  beforeDestroy () {
    this.$api.removeServerListener(this.server.id, 'console', this.listener)
    this.interval && clearInterval(this.interval)
    this.interval = null
    lines = []
  },
  methods: {
    parseConsole (data) {
      let newLines = (Array.isArray(data.logs) ? data.logs.join('') : data.logs).replaceAll('\r\n', '\n')
      const endOnNewline = newLines.endsWith('\n')
      newLines = newLines.split('\n').map(line => {
        return ansiup.ansi_to_html(line) + '\n'
      })

      if (!endOnNewline && newLines.length > 0) {
        const line = newLines[newLines.length - 1]
        newLines[newLines.length - 1] = line.substring(0, line.length - 1)
      } else if (newLines.length > 0) {
        newLines.pop()
      }

      if (lines.length !== 0 && !lines[lines.length - 1].textEl.innerHTML.endsWith('\n')) {
        lines[lines.length - 1].textEl.innerHTML += newLines.shift()
        lines[lines.length - 1].crHandled = false
      }

      newLines = newLines.map((line) => {
        const isDaemonMessage = DAEMON_MESSAGE_REGEX.test(line)
        const el = document.createElement('span')
        if (isDaemonMessage) el.append(this.$refs.daemonChip.$el.cloneNode(true))
        const textEl = document.createElement('span')
        textEl.className = 'consoleLine'
        textEl.innerHTML = isDaemonMessage ? line.replace(DAEMON_MESSAGE_REGEX, '') : line
        el.append(textEl)
        return {
          el,
          textEl,
          crHandled: false,
          size: data.length * 2
        }
      })

      lines = lines.concat(newLines)

      const currentSize = lines.reduce((acc, curr) => acc + curr.size, 0)
      let freed = 0
      let toRemove = 0
      while (currentSize - freed > CONSOLE_MEMORY_ALLOWED_PER_BUFFER) {
        freed += lines[toRemove].size
        toRemove += 1
      }
      lines = lines.slice(toRemove)

      lines = lines.map(line => {
        if (line.crHandled) return line
        const endOnNewline = line.textEl.innerHTML.endsWith('\n')
        const parts = (endOnNewline ? line.textEl.innerHTML.substring(0, line.textEl.innerHTML.length - 1) : line.textEl.innerHTML).split('\r')
        let result = parts.shift()
        parts.map(part => {
          result = part + result.substring(part.length)
        })
        line.textEl.innerHTML = endOnNewline ? (result + '\n') : result
        return { ...line, crHandled: true }
      })

      this.hasNewLines = true
    },
    togglePaused () {
      this.paused = !this.paused
    },
    sendCommand () {
      this.$api.sendServerCommand(this.server.id, this.consoleCommand)
      this.consoleCommand = ''
    }
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
    word-break: break-all;
  }

  .consoleLine {
    white-space: pre;
  }
</style>
