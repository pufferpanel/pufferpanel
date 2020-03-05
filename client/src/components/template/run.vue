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
  <div>
    <v-row>
      <v-col cols="12">
        <v-text-field v-model="value.command" :label="$t('templates.Command')" outlined hide-details />
      </v-col>
    </v-row>
  </div>
</template>

<script>
export default {
  props: {
    value: { type: Object, default: () => {} }
  }
}
</script>
