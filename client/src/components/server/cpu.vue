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
import VueApexCharts from 'vue-apexcharts'

export default {
  components: {
    apexchart: VueApexCharts
  },
  props: {
    server: { type: Object, default: () => {} }
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
          foreColor: this.$isDark() ? '#FFF' : '#000000DE'
        },
        colors: [this.$isDark() ? this.$vuetify.theme.themes.dark.accent : this.$vuetify.theme.themes.light.accent],
        tooltip: { theme: [this.$isDark() ? this.$vuetify.theme.themes.dark.accent : this.$vuetify.theme.themes.light.accent] },
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
          min: 0,
          max: () => this.series[0] ? Math.ceil(this.series[0].data.reduce((acc, curr) => Math.max(acc, curr[1]), 100)) : 100
        },
        legend: {
          show: true
        }
      },
      series: []
    }
  },
  mounted () {
    this.$api.addServerListener(this.server.id, 'stat', event => {
      this.updateStats(event)
    })
  },
  methods: {
    updateStats (data) {
      this.options = { ...this.options, chart: { ...this.options.chart, foreColor: this.$isDark() ? '#FFF' : '#000000DE' } }
      const chartData = [...((this.series[0] || {}).data || []), [new Date().getTime(), Math.round(data.cpu * 100) / 100]]
      this.series = [{
        name: this.$t('servers.CPU'),
        data: chartData.length > this.maxPoints ? chartData.slice(chartData.length - this.maxPoints) : chartData
      }]
    }
  }
}
</script>
