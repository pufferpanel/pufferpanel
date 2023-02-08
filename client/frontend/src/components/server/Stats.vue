<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import Chart from 'vue3-apexcharts'

const { t } = useI18n()
const wrapper = ref(null)
const apex = ref(null)
const width = ref('100%')

const cpu = []
const mem = []

const series = [
  {
    name: t('servers.CPU'),
    data: []
  },
  {
    name: t('servers.Memory'),
    data: []
  }
]

const props = defineProps({
  server: { type: Object, required: true }
})

const numFormat = new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 })

const options = {
  chart: {
    toolbar: {
      show: false
    },
    zoom: {
      enabled: false
    },
    animations: {
      enabled: true,
      easing: 'linear',
      dynamicAnimation: {
        speed: 1000
      }
    }
  },
  dataLabels: {
    enabled: false
  },
  stroke: {
    curve: 'smooth'
  },
  markers: {
    size: 0
  },
  legend: {
    show: false
  },
  xaxis: {
    type: 'datetime',
    range: 60000, // show last 60 seconds on chart
    labels: {
      formatter: value => new Date(value).toLocaleTimeString()
    }
  },
  yaxis: [
    {
      title: {
        text: t('servers.CPU')
      },
      min: 0,
      // force nice scale with max 100 causes chart to display max 120
      // so default max to 99 to actually get the max to default to 100
      max: (dataMax) => Math.max(99, dataMax),
      forceNiceScale: true,
      labels: {
        formatter: x => x + ' %'
      }
    },
    {
      opposite: true,
      logarithmic: false,
      logBase: 2,
      title: {
        text: t('servers.Memory')
      },
      min: 0,
      max: (dataMax) => Math.max(1024*1024, dataMax),
      labels: {
        formatter: x => {
          if (x < Math.pow(2, 10)) return numFormat.format(x) + ' B'
          if (x < Math.pow(2, 20)) return numFormat.format(x / Math.pow(2, 10)) + ' KiB'
          if (x < Math.pow(2, 30)) return numFormat.format(x / Math.pow(2, 20)) + ' MiB'
          if (x < Math.pow(2, 40)) return numFormat.format(x / Math.pow(2, 30)) + ' GiB'
          return numFormat.format(x / Math.pow(2, 40)) + ' TiB'
        }
      }
    }
  ],
  tooltip: {
    followCursor: true,
    x: {
      show: false
    }
  }
}

function addData(d) {
  const ts = new Date().getTime()
  let animate = true

  cpu.push([ts, d.cpu])
  if (cpu.length > 600) {
    animate = false
    while (cpu.length > 60) {
      cpu.shift()
    }
  }

  mem.push([ts, d.memory])
  if (mem.length > 600) {
    animate = false
    while (mem.length > 60) {
      mem.shift()
    }
  }

  apex.value.chart.updateSeries([
    { data: cpu },
    { data: mem }
  ], animate)
}

let task = null
let stopListener = null

onMounted(() => {
  nextTick(() => {
    width.value = wrapper.value.clientWidth
    const observer = new ResizeObserver(elems => {
      for (const elem of elems) {
        if (wrapper.value) width.value = wrapper.value.clientWidth
      }
    })
    observer.observe(wrapper.value)

    const font = getComputedStyle(apex.value.$el).getPropertyValue('--apex-font').trim() || 'sans-serif'
    const backgroundColor = getComputedStyle(apex.value.$el).getPropertyValue('--apex-background-color').trim() || '#fff'
    const textColor = getComputedStyle(apex.value.$el).getPropertyValue('--apex-text-color').trim() || '#000'
    const seriesColors = getComputedStyle(apex.value.$el).getPropertyValue('--apex-series-colors').split(',').map(c => c.trim()).filter(c => /#[a-fA-F0-9]{3,6}/.test(c))

    options.chart.fontFamily = font
    options.chart.background = backgroundColor
    options.chart.foreColor = textColor
    if (seriesColors.length > 0) options.colors = seriesColors
  })

  stopListener = props.server.on('stat', addData)

  task = props.server.startTask(() => {
    props.server.stats()
  }, 1000)
})

onUnmounted(() => {
  if (task) props.server.stopTask(task)
  if (stopListener) stopListener()
})
</script>

<template>
  <div ref="wrapper" class="statistics">
    <h2 v-text="t('servers.Statistics')" />
    <chart ref="apex" dir="ltr" :width="width" height="500" :options="options" :series="series" />
  </div>
</template>
