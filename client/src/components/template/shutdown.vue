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
  <v-row>
    <v-col
      cols="12"
      class="pb-0"
    >
      <v-btn-toggle
        :value="value.type"
        borderless
        dense
        mandatory
        @change="$emit('input', { type: $event, stop: $event === 'command' ? '' : 0 })"
      >
        <v-btn
          value="command"
          v-text="$t('templates.stop.Command')"
        />
        <v-btn
          value="signal"
          v-text="$t('templates.stop.Signal')"
        />
      </v-btn-toggle>
    </v-col>
    <v-col cols="12">
      <ui-input
        :value="value.stop"
        :label="$t('templates.stop.Command')"
        :type="value.type === 'command' ? 'text' : 'number'"
        @input="$emit('input', { ...value, stop: $event })"
      />
    </v-col>
  </v-row>
</template>

<script>
export default {
  props: {
    value: { type: Object, default: () => {} }
  }
}
</script>
