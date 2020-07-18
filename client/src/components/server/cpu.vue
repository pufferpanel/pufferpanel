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
    <v-card-title v-text="$t('servers.CPU')" />
    <v-card-text>
      <apexchart
        :options="options"
        :series="series"
        height="300"
      />
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
      intl: new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }),
      maxPoints: 20,
      options: {
        chart: {
          id: 'cpu',
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
            formatter: value => this.intl.format(Math.round(value * 100) / 100) + '%'
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
      const chartData = [...((this.series[0] || {}).data || []), [new Date().getTime(), Math.round(data.cpu * 100) / 100]]
      this.series = [{
        name: this.$t('servers.CPU'),
        data: chartData.length > this.maxPoints ? chartData.slice(chartData.length - this.maxPoints) : chartData
      }]
    },
    isDark
  }
}
</script>
