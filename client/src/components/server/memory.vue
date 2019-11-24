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
    <v-card-title v-text="$t('servers.Memory')" />
    <v-card-text class="pt-4">
      <line-chart
        :chart-data="datacollection"
        :options="options"
      />
    </v-card-text>
  </v-card>
</template>

<script>
import LineChart from './LineChart.js'
import { isDark } from '@/utils/dark'

export default {
  components: {
    LineChart
  },
  data () {
    return {
      maxPoints: 20,
      memory: Array.apply(null, Array(this.maxPoints)).map(Number.prototype.valueOf, 0),
      label: Array.apply(null, Array(this.maxPoints)).map(String.prototype.valueOf, ''),
      datacollection: {
        datasets: [
          {
            backgroundColor: isDark() ? this.$vuetify.theme.themes.dark.accent : this.$vuetify.theme.themes.light.accent,
            data: []
          }
        ]
      },
      options: {
        scales: {
          xAxes: [{
            ticks: {
              fontColor: isDark() ? this.$vuetify.theme.themes.dark.tertiary : this.$vuetify.theme.themes.light.tertiary
            },
            gridLines: {
              color: isDark() ? this.$vuetify.theme.themes.dark.tertiary : this.$vuetify.theme.themes.light.tertiary
            }
          }],
          yAxes: [{
            ticks: {
              beginAtZero: true,
              fontColor: isDark() ? this.$vuetify.theme.themes.dark.tertiary : this.$vuetify.theme.themes.light.tertiary
            },
            gridLines: {
              color: isDark() ? this.$vuetify.theme.themes.dark.tertiary : this.$vuetify.theme.themes.light.tertiary
            }
          }]
        },
        legend: {
          display: false
        },
        responsive: true,
        animation: false
      }
    }
  },
  mounted () {
    const root = this
    this.$socket.addEventListener('message', function (event) {
      const data = JSON.parse(event.data)
      if (data === 'undefined') {
        return
      }
      if (data.type === 'stat') {
        root.updateStats(data.data)
      }
    })
  },
  methods: {
    updateStats (data) {
      if (this.memory.length === this.maxPoints) {
        this.memory.shift()
        this.label.shift()
      }

      this.label.push(new Date().toLocaleTimeString())
      this.memory.push(data.memory)

      const newMemory = []
      const newLabel = []
      for (let i = 0; i < this.memory.length; i++) {
        newMemory[i] = this.memory[i]
        newLabel[i] = this.label[i]
      }

      this.datacollection = {
        labels: newLabel,
        datasets: [
          {
            backgroundColor: isDark() ? this.$vuetify.theme.themes.dark.accent : this.$vuetify.theme.themes.light.accent,
            data: newMemory
          }
        ]
      }
    }
  }
}
</script>
