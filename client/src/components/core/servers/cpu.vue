<!--
  - Copyright 2019 Padduck, LLC
  -  Licensed under the Apache License, Version 2.0 (the "License");
  -  you may not use this file except in compliance with the License.
  -  You may obtain a copy of the License at
  -  	http://www.apache.org/licenses/LICENSE-2.0
  -  Unless required by applicable law or agreed to in writing, software
  -  distributed under the License is distributed on an "AS IS" BASIS,
  -  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  -  See the License for the specific language governing permissions and
  -  limitations under the License.
  -->

<template>
  <b-card
    header-tag="header"
    footer-tag="footer">
    <h6 slot="header" class="mb-0" v-text="$t('common.CPU')"></h6>
    <line-chart :chart-data="datacollection" :options="options"></line-chart>
  </b-card>
</template>

<script>
import LineChart from './LineChart.js'

export default {
  components: {
    LineChart
  },
  data () {
    return {
      maxPoints: 20,
      cpu: Array.apply(null, Array(this.maxPoints)).map(Number.prototype.valueOf, 0),
      label: Array.apply(null, Array(this.maxPoints)).map(String.prototype.valueOf, ''),
      datacollection: {
        datasets: [
          {
            label: this.$t('common.CPU'),
            backgroundColor: '#65a5f8',
            data: []
          }
        ]
      },
      options: {
        scales: {
          yAxes: [{
            ticks: {
              beginAtZero: true
            }
          }]
        },
        responsive: false,
        animation: false
      }
    }
  },
  methods: {
    updateStats (data) {
      if (this.cpu.length === this.maxPoints) {
        this.cpu.shift()
        this.label.shift()
      }

      this.label.push(new Date().toLocaleTimeString())
      this.cpu.push(data.cpu)

      let newCpu = []
      let newLabel = []
      for (let i = 0; i < this.cpu.length; i++) {
        newCpu[i] = this.cpu[i]
        newLabel[i] = this.label[i]
      }

      this.datacollection = {
        labels: newLabel,
        datasets: [
          {
            label: this.$t('common.CPU'),
            backgroundColor: '#65a5f8',
            data: newCpu
          }
        ]
      }
    }
  },
  mounted () {
    let root = this
    this.$socket.addEventListener('message', function (event) {
      let data = JSON.parse(event.data)
      if (data === 'undefined') {
        return
      }
      if (data.type === 'stat') {
        root.updateStats(data.data)
      }
    })
  }
}
</script>
