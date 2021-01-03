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
  <v-card>
    <v-card-title>
      <span v-text="$t('servers.Settings')" />
    </v-card-title>
    <v-card-text>
      <v-row>
        <v-col v-if="Object.keys(items).length === 0">
          <span v-text="$t('servers.NoEditableVars')" />
        </v-col>
        <v-col
          v-for="(item, name) in items"
          :key="name"
          cols="12"
        >
          <ui-variable-input v-model="items[name]" />
        </v-col>
      </v-row>
      <v-row v-if="Object.keys(items).length > 0">
        <v-col>
          <v-btn
            block
            color="success"
            @click="save"
            v-text="$t('common.Save')"
          />
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script>
import markdown from '@/utils/markdown'

export default {
  props: {
    server: { type: Object, default: () => {} }
  },
  data () {
    return {
      items: {}
    }
  },
  mounted () {
    this.loadData()
  },
  methods: {
    async loadData () {
      this.items = await this.$api.getServerData(this.server.id)
    },
    async save () {
      await this.$api.updateServerData(this.server.id, this.items)
      this.$toast.success(this.$t('common.Saved'))
    },
    markdown
  }
}
</script>
