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
      <span v-text="$t('servers.Memory')" />
      <div class="flex-grow-1" />
      <v-btn-toggle v-model="mibMode" borderless dense mandatory>
        <v-btn :value="false" v-text="'MB'" />
        <v-btn :value="true"  v-text="'MiB'" />
      </v-btn-toggle>
    </v-card-title>
    <v-card-text>
      <apexchart
        :series="series"
        :options="options"
        height="250"
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
      maxPoints: 40,
      mibMode: false,
      options: {
        chart: {
          id: 'memory',
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
          curve: 'straight'
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
              if (!this.mibMode) {
                if (value < 1e3) return this.intl.format(value) + ' B'
                if (value < 1e6) return this.intl.format(value / 1e3) + ' KB'
                if (value < 1e9) return this.intl.format(value / 1e6) + ' MB'
                if (value < 1e12) return this.intl.format(value / 1e9) + ' GB'
                return this.intl.format(value / 1e12) + ' TB'
              } else {
                if (value < Math.pow(2, 10)) return this.intl.format(value) + ' B'
                if (value < Math.pow(2, 20)) return this.intl.format(value / Math.pow(2, 10)) + ' KiB'
                if (value < Math.pow(2, 30)) return this.intl.format(value / Math.pow(2, 20)) + ' MiB'
                if (value < Math.pow(2, 40)) return this.intl.format(value / Math.pow(2, 30)) + ' GiB'
                return this.intl.format(value / Math.pow(2, 40)) + ' TiB'
              }
            }
          },
          min: 0,
          max: () => this.series[0]
            ? Math.ceil(
              this.series[0].data.reduce(
                (acc, curr) => Math.max(acc, curr[1]),
                this.mibMode ? 1025 : 1000
              )
            )
            : this.mibMode ? 1025 : 1000
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

    const chartData = new Array(this.maxPoints)
    const timestamp = new Date().getTime()

    for (let i = 0; i < chartData.length; i++) {
      chartData[i] = [timestamp - ((chartData.length - i) * 3000), 0]
    }

    this.series = [{ name: this.$t('servers.Memory'), data: chartData }]
  },
  methods: {
    updateStats (data) {
      this.options = { ...this.options, chart: { ...this.options.chart, foreColor: this.$isDark() ? '#FFF' : '#000000DE' } }
      const chartData = [...((this.series[0] || {}).data || []), [new Date().getTime(), data.memory]]
      this.series = [{
        name: this.$t('servers.Memory'),
        data: chartData.length > this.maxPoints ? chartData.slice(chartData.length - this.maxPoints) : chartData
      }]
    }
  }
}
</script>
