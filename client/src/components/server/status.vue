<!--
  - Copyright 2019 Padduck, LLC
  -  Licensed under the Apache License, Version 2.0 (the "License");
  -  you may not use this file except in compliance with the License.
  -  You may obtain a copy of the License at
  -          http://www.apache.org/licenses/LICENSE-2.0
  -  Unless required by applicable law or agreed to in writing, software
  -  distributed under the License is distributed on an "AS IS" BASIS,
  -  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  -  See the License for the specific language governing permissions and
  -  limitations under the License.
  -->

<template>
<v-tooltip bottom>
  <template v-slot:activator="{ on }">
    <v-icon v-on="on" dense :color="online ? 'success' : 'error'">mdi-brightness-1</v-icon>
  </template>
  <span v-text="online ? $t('common.Online') : $t('common.Offline')" />
</v-tooltip>
</template>

<script>
export default {
  props: {
    server: { type: Object, default: function () { return {} } }
  },
  data () {
    return {
      online: false,
      interval: null
    }
  },
  mounted () {
    const ctx = this
    this.getStatus(ctx)
    this.interval = setInterval(ctx.getStatus, 5000, ctx)
  },
  beforeDestroy () {
    clearInterval(this.interval)
  },
  methods: {
    getStatus (ctx) {
      ctx.$http.get(`/daemon/server/${ctx.server.id}/status`).then(function (response) {
        ctx.online = response.data.running
      }).catch(function (error) {
        let msg = 'errors.ErrUnknownError'
        if (error && error.response && error.response.data.error) {
          if (error.response.data.error.code) {
            msg = 'errors.' + error.response.data.error.code
          } else {
            msg = error.response.data.error.msg
          }
        }

        ctx.$toast.error(ctx.$t(msg))
      })
    }
  }
}
</script>
