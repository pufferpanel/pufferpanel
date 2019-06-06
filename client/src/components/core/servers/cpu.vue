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
  <line-chart :chart-data="dataCollection" :options="options"></line-chart>
</template>

<script>
import LineChart from './LineChart.js'

export default {
  components: {
    LineChart
  },
  props: {
    server: Object
  },
  data () {
    return {
      refreshInterval: null,
      dataCollection: {
        labels: [],
        datasets: [{
          label: 'CPU',
          data: []
        }]
      },
      options: {
        scales: {
          yAxes: [{
            ticks: {
              beginAtZero: true
            }
          }]
        }
      }
    }
  },
  methods: {
    updateStats (data) {
      this.dataCollection.labels.push(Date.now())
      this.dataCollection.datasets[0].data.push(data.cpu)
    },
    sendStatRequest () {
      this.updateStats({
        cpu: Math.floor(Math.random() * (50 - 5 + 1)) + 5
      })
    }
  },
  mounted () {
    let root = this
    this.$socket.addEventListener('message', function (event) {
      let data = JSON.parse(event.data)
      if (data === 'undefined') {
        return
      }
      if (data.type === 'stats') {
        root.updateStats(data.data)
      }
    })

    this.refreshInterval = setInterval(this.sendStatRequest, 5000)
  },
  beforeDestroy () {
    if (this.refreshInterval !== null) {
      clearInterval(this.refreshInterval)
    }
  }
}
</script>

<style>
  .small {
    max-width: 600px;
    margin:  150px auto;
  }
</style>
