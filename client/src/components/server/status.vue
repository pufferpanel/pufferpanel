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
      <v-icon
        dense
        :color="online === true ? 'success' : online === false ? 'error' : 'grey'"
        v-on="on"
      >
        mdi-brightness-1
      </v-icon>
    </template>
    <span v-text="online === true ? $t('common.Online') : online === false ? $t('common.Offline') : $t('common.Unknown')" />
  </v-tooltip>
</template>

<script>
export default {
  data () {
    return {
      online: null
    }
  },
  mounted () {
    const ctx = this

    this.$socket.addEventListener('message', event => {
      const data = JSON.parse(event.data)
      if (!data) return
      if (data.type === 'status') ctx.online = data.data.running
    })

    setTimeout(() => {
      ctx.$socket.sendObj({ type: 'status' })
    }, 200)
  }
}
</script>
