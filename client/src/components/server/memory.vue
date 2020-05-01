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
    <v-card-text>
      <apexchart :series="series" :options="options" height="300" />
    </v-card-text>
  </v-card>
</template>

<script>
import { isDark } from '@/utils/dark'
import VueApexCharts from 'vue-apexcharts'

export default {
  components: {
    apexchart: VueApexCharts
  },
  data () {
    return {
      intl: new Intl.NumberFormat('en-US', {maximumFractionDigits: 2}),
      maxPoints: 20,
      options: {
        chart: {
          id: 'memory',
          height: 300,
          type: 'line',
          animations: {
            enabled: false
          },
          toolbar: {
            show: false
          },
          foreColor: isDark() ? '#FFF' : '#000000DE'
        },
        colors: [isDark() ? this.$vuetify.theme.themes.dark.accent : this.$vuetify.theme.themes.light.accent],
        tooltip: { theme: [isDark() ? this.$vuetify.theme.themes.dark.accent : this.$vuetify.theme.themes.light.accent] },
        dataLabels: {
          enabled: false
        },
        stroke: {
          curve: 'smooth'
        },
        markers: {
          size: 0
        },
        xaxis: {
          labels: {
            show: true,
            rotate: 0,
            formatter: value => new Date(value).toLocaleTimeString()
          },
          tooltip: {
            enabled: false
          },
          type: 'datetime'
        },
        yaxis: {
          labels: {
            show: true,
            formatter: value => {
              if (value < 1e3) return this.intl.format(value) + ' B'
              if (value < 1e6) return this.intl.format(value / 1e3) + ' KB'
              if (value < 1e9) return this.intl.format(value / 1e6) + ' MB'
              if (value < 1e12) return this.intl.format(value / 1e9) + ' GB'
              if (value < 1e15) return this.intl.format(value / 1e12) + ' TB'
            }
          },
          min: 0
        },
        legend: {
          show: true
        }
      },
      series: []
    }
  },
  mounted () {
    const ctx = this
    this.$socket.addEventListener('message', event => {
      const data = JSON.parse(event.data)
      if (data === 'undefined') {
        return
      }
      if (data.type === 'stat') {
        ctx.updateStats(data.data)
      }
    })
  },
  methods: {
    updateStats (data) {
      this.options = { ...this.options, chart: { ...this.options.chart, foreColor: isDark() ? '#FFF' : '#000000DE' } }
      const chartData = [...((this.series[0] || {}).data || []), [new Date().getTime(), data.memory]]
      this.series = [{
        name: this.$t('servers.Memory'),
        data: chartData.length > this.maxPoints ? chartData.slice(chartData.length - this.maxPoints) : chartData
      }]
    }
  }
}
</script>
