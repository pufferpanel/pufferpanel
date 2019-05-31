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
      console: ''
    }
  },
  methods: {
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
