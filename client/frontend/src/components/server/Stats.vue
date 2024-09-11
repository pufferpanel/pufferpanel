<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Chart, { _adapters, Tooltip } from 'chart.js/auto'
import 'chartjs-adapter-date-fns'
import Query from './Query.vue'

const fromCss = (el, prop) => {
  return getComputedStyle(el).getPropertyValue(prop).trim()
}

const defaultFamily = "'Helvetica Neue', 'Helvetica', 'Arial', sans-serif"

const props = defineProps({
  server: { type: Object, required: true }
})

const { t, locale } = useI18n()

const numFormat = new Intl.NumberFormat(
  [locale.value.replace('_', '-'), 'en'],
  { maximumFractionDigits: 2 }
)

const formatCpu = (value) => {
  return numFormat.format(value) + ' %'
}

const formatMemory = (value) => {
  if (!value) return numFormat.format(0) + ' B'
  if (value < Math.pow(2, 10)) return numFormat.format(value) + ' B'
  if (value < Math.pow(2, 20)) return numFormat.format(value / Math.pow(2, 10)) + ' KiB'
  if (value < Math.pow(2, 30)) return numFormat.format(value / Math.pow(2, 20)) + ' MiB'
  if (value < Math.pow(2, 40)) return numFormat.format(value / Math.pow(2, 30)) + ' GiB'
  return numFormat.format(value / Math.pow(2, 40)) + ' TiB'
}

// date-fns localization approach takes importing an object per locale
// which is both a lot of data and annoying to have to handle as that
// means manually dealing with fallbacks, dynamic imports, etc
// so instead we just force the adapter to use the browsers built in
// date localization tooling and just let date-fns do the logic parts
const intl = new Intl.DateTimeFormat(
  [locale.value.replace('_', '-'), 'en'],
  { hour: 'numeric', minute: 'numeric', second: 'numeric' }
)
_adapters._date.prototype.format = (time) => {
  return intl.format(time)
}

Tooltip.positioners.cursor = (_, eventPosition) => {
  return { x: eventPosition.x, y: eventPosition.y }
}

const cpuChartEl = ref(null)
const memoryChartEl = ref(null)
let cpuChart = null
let memoryChart = null
const cpu = []
const memory = []
const jvmHeapUsed = []
const jvmHeapAlloc = []
const jvmMetaUsed = []
const jvmMetaAlloc = []

function addData(d) {
  const x = new Date().getTime()

  cpu.push({ x, y: d.cpu })
  memory.push({ x, y: d.memory })

  if (d.jvm) {
    jvmHeapUsed.push({ x, y: d.jvm.heapUsed })
    jvmHeapAlloc.push({ x, y: d.jvm.heapTotal - d.jvm.heapUsed })
    jvmMetaUsed.push({ x, y: d.jvm.metaspaceUsed })
    jvmMetaAlloc.push({ x, y: d.jvm.metaspaceTotal - d.jvm.metaspaceUsed })

    memoryChart.show(memoryChart.data.datasets.findIndex(set => set.label === t('servers.JvmMetaUsed')))
    memoryChart.show(memoryChart.data.datasets.findIndex(set => set.label === t('servers.JvmMetaAlloc')))
    memoryChart.show(memoryChart.data.datasets.findIndex(set => set.label === t('servers.JvmHeapUsed')))
    memoryChart.show(memoryChart.data.datasets.findIndex(set => set.label === t('servers.JvmHeapAlloc')))
    memoryChart.hide(memoryChart.data.datasets.findIndex(set => set.label === t('servers.Memory')))
  }

  for (let graph of [cpu, memory, jvmHeapUsed, jvmHeapAlloc, jvmMetaUsed, jvmMetaAlloc]) {
    while (graph.length > 60) {
      graph.shift()
    }
  }

  for (let chart of [cpuChart, memoryChart]) {
    chart.options.scales.x.min = x - (60 * 1000)
    chart.update()
  }
}

const chartOptions = (mode) => {
  const options = {
    responsive: true,
    aspectRatio: (ctx) => {
      if (ctx.chart.canvas) {
        return parseFloat(fromCss(ctx.chart.canvas.parentElement, 'aspect-ratio')) || 2
      } else return 2
    },
    parsing: false,
    locale: locale.value.split('_')[0] || 'en',
    interaction: {
      mode: 'x',
      intersect: false
    },
    animations: {
      y: {
        duration: 0
      }
    },
    plugins: {
      tooltip: {
        position: 'cursor',
        usePointStyle: true,
        callbacks: {
          label: (ctx) => {
            return ' ' + ctx.dataset.label + ': ' + ctx.chart.scales[ctx.dataset.yAxisID].options.ticks.callback(ctx.parsed.y)
          },
          labelPointStyle: () => 'circle'
        },
        itemSort: (a, b) => {
          return a.datasetIndex < b.datasetIndex
        },
        multiKeyBackground: (ctx) => fromCss(ctx.chart.canvas, 'background-color') || '#fff',
        padding: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-tooltip-padding') || 6,
        backgroundColor: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-tooltip-background-color') || 'rgba(0, 0, 0, 0.8)',
        cornerRadius: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-tooltip-corner-radius') || 6,
        borderColor: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-tooltip-border-color') || 'rgba(0, 0, 0, 0)',
        borderWidth: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-tooltip-border-width') || 0,
        titleAlign: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-tooltip-title-align') || 'left',
        titleColor: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-tooltip-font-color') || '#fff',
        bodyColor: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-tooltip-font-color') || '#fff',
        bodySpacing: (ctx) => parseInt(fromCss(ctx.chart.canvas, '--chartjs-tooltip-body-spacing')) || 2,
        titleFont: {
          family: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-tooltip-title-font-family') || fromCss(ctx.chart.canvas, '--chartjs-font-family') || defaultFamily,
          size: (ctx) => parseInt(fromCss(ctx.chart.canvas, '--chartjs-tooltip-title-font-size') || fromCss(ctx.chart.canvas, '--chartjs-font-size')) || 12,
          weight: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-tooltip-title-font-weight') || 'bold'
        },
        bodyFont: {
          family: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-tooltip-body-font-family') || fromCss(ctx.chart.canvas, '--chartjs-font-family') || defaultFamily,
          size: (ctx) => parseInt(fromCss(ctx.chart.canvas, '--chartjs-tooltip-body-font-size') || fromCss(ctx.chart.canvas, '--chartjs-font-size')) || 12,
          weight: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-tooltip-body-font-weight') || fromCss(ctx.chart.canvas, '--chartjs-font-weight')
        }
      },
      legend: {
        display: false
      },
      title: {
        display: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-title-display') == 'true',
        align: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-title-align') || 'center',
        color: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-title-font-color') || fromCss(ctx.chart.canvas, '--chartjs-color') || fromCss(ctx.chart.canvas, 'color'),
        position: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-title-position') || 'top',
        font: {
          family: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-title-font-family') || fromCss(ctx.chart.canvas, '--chartjs-font-family') || defaultFamily,
          size: (ctx) => parseInt(fromCss(ctx.chart.canvas, '--chartjs-title-font-size') || fromCss(ctx.chart.canvas, '--chartjs-font-size')) || 12,
          weight: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-title-font-weight') || 'bold'
        }
      }
    },
    scales: {
      x: {
        type: 'timeseries',
        min: new Date().getTime() - (60 * 1000),
        grid: {
          color: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-grid-color')
        },
        ticks: {
          min: 12,
          source: 'data',
          color: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-axis-font-color') || fromCss(ctx.chart.canvas, '--chartjs-color') || fromCss(ctx.chart.canvas, 'color'),
          font: {
            family: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-axis-font-family') || fromCss(ctx.chart.canvas, '--chartjs-font-family') || defaultFamily,
            size: (ctx) => parseInt(fromCss(ctx.chart.canvas, '--chartjs-axis-font-size') || fromCss(ctx.chart.canvas, '--chartjs-font-size')) || 12,
            weight: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-axis-font-weight') || fromCss(ctx.chart.canvas, '--chartjs-font-weight')
          }
        }
      }
    },
    elements: {
      line: {
        tension: 0.3
      },
      point: {
        pointStyle: false,
        hoverRadius: 20
      }
    }
  }

  if (mode === 'memory') {
    options.plugins.title.text = t('servers.Memory')
    options.scales.memory = {
      type: 'linear',
      position: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-axis-y-position') || 'left',
      min: 0,
      suggestedMax: 1024 * 1024,
      grid: {
        display: true,
        color: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-grid-color')
      },
      ticks: {
        callback: formatMemory,
        color: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-axis-font-color') || fromCss(ctx.chart.canvas, '--chartjs-color') || fromCss(ctx.chart.canvas, 'color'),
        font: {
          family: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-axis-font-family') || fromCss(ctx.chart.canvas, '--chartjs-font-family') || defaultFamily,
          size: (ctx) => parseInt(fromCss(ctx.chart.canvas, '--chartjs-axis-font-size') || fromCss(ctx.chart.canvas, '--chartjs-font-size')) || 12,
          weight: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-axis-font-weight') || fromCss(ctx.chart.canvas, '--chartjs-font-weight')
        }
      }
    }
  }

  if (mode === 'cpu') {
    options.plugins.title.text = t('servers.CPU')
    options.scales.cpu = {
      type: 'linear',
      position: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-axis-y-position') || 'left',
      min: 0,
      suggestedMax: 100,
      grid: {
        color: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-grid-color')
      },
      ticks: {
        callback: formatCpu,
        color: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-axis-font-color') || fromCss(ctx.chart.canvas, '--chartjs-color') || fromCss(ctx.chart.canvas, 'color'),
        font: {
          family: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-axis-font-family') || fromCss(ctx.chart.canvas, '--chartjs-font-family') || defaultFamily,
          size: (ctx) => parseInt(fromCss(ctx.chart.canvas, '--chartjs-axis-font-size') || fromCss(ctx.chart.canvas, '--chartjs-font-size')) || 12,
          weight: (ctx) => fromCss(ctx.chart.canvas, '--chartjs-axis-font-weight') || fromCss(ctx.chart.canvas, '--chartjs-font-weight')
        }
      }
    }
  }

  return options
}

let task = null
let stopListener = null
onMounted(() => {
  cpuChart = new Chart(cpuChartEl.value, {
    type: 'line',
    options: chartOptions('cpu'),
    data: {
      datasets: [
        {
          label: t('servers.CPU'),
          yAxisID: 'cpu',
          borderColor: fromCss(cpuChartEl.value, '--chartjs-series-line-cpu'),
          backgroundColor: fromCss(cpuChartEl.value, '--chartjs-series-fill-cpu'),
          data: cpu
        }
      ]
    }
  })

  memoryChart = new Chart(memoryChartEl.value, {
    type: 'line',
    options: chartOptions('memory'),
    data: {
      datasets: [
      {
          label: t('servers.JvmMetaUsed'),
          yAxisID: 'memory',
          fill: 'origin',
          borderColor: fromCss(memoryChartEl.value, '--chartjs-series-line-jvm-metaspace-used'),
          backgroundColor: fromCss(memoryChartEl.value, '--chartjs-series-fill-jvm-metaspace-used'),
          stack: 'jvmMemory',
          hidden: true,
          data: jvmMetaUsed
        },
        {
          label: t('servers.JvmMetaAlloc'),
          yAxisID: 'memory',
          fill: '-1',
          borderColor: fromCss(memoryChartEl.value, '--chartjs-series-line-jvm-metaspace-allocated'),
          backgroundColor: fromCss(memoryChartEl.value, '--chartjs-series-fill-jvm-metaspace-allocated'),
          stack: 'jvmMemory',
          hidden: true,
          data: jvmMetaAlloc
        },
        {
          label: t('servers.JvmHeapUsed'),
          yAxisID: 'memory',
          fill: '-1',
          borderColor: fromCss(memoryChartEl.value, '--chartjs-series-line-jvm-heapspace-used'),
          backgroundColor: fromCss(memoryChartEl.value, '--chartjs-series-fill-jvm-heapspace-used'),
          stack: 'jvmMemory',
          hidden: true,
          data: jvmHeapUsed
        },
        {
          label: t('servers.JvmHeapAlloc'),
          yAxisID: 'memory',
          fill: '-1',
          borderColor: fromCss(memoryChartEl.value, '--chartjs-series-line-jvm-heapspace-allocated'),
          backgroundColor: fromCss(memoryChartEl.value, '--chartjs-series-fill-jvm-heapspace-allocated'),
          stack: 'jvmMemory',
          hidden: true,
          data: jvmHeapAlloc
        },
        {
          label: t('servers.Memory'),
          yAxisID: 'memory',
          borderColor: fromCss(memoryChartEl.value, '--chartjs-series-line-memory'),
          backgroundColor: fromCss(memoryChartEl.value, '--chartjs-series-fill-memory'),
          data: memory
        }
      ]
    }
  })

  stopListener = props.server.on('stat', addData)

  task = props.server.startTask(async () => {
    if (props.server.needsPolling() && props.server.hasScope('server.stats')) {
      addData(await props.server.getStats())
    }
  }, 5000)
})

onUnmounted(() => {
  if (cpuChart) cpuChart.destroy()
  if (memoryChart) memoryChart.destroy()
  if (task) props.server.stopTask(task)
  if (stopListener) stopListener()
})
</script>

<template>
  <Query :server="server" />
  <div class="chart memory">
    <canvas ref="memoryChartEl"/>
  </div>
  <div class="chart cpu">
    <canvas ref="cpuChartEl"/>
  </div>
</template>
