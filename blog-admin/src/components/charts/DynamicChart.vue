<script setup lang="ts">

import type { ECharts } from 'echarts'
import * as echarts from 'echarts'
import { onMounted, onUnmounted } from 'vue'
import type { EChartsOption } from 'echarts/types/dist/echarts'

const memoryParents: [Date, string][] = [], cpuParents: [Date, string][] = [];

const pushNewValue = (memoryParent?: number, cpuParent?: number) => {
  if (!memoryParent && !cpuParent) {
    return;
  }
  if (memoryParent !== undefined) {
    memoryParents.push([new Date(), memoryParent.toFixed(2)])
  }
  if (cpuParent !== undefined) {
    cpuParents.push([new Date(), cpuParent.toFixed(2)])
  }
  if (memoryParents.length > 150) {
    memoryParents.shift();
  }
  if (cpuParents.length > 150) {
    cpuParents.shift();
  }
  updateChart();
}

const chartId = new Date().getTime().toString() + (Math.random() * 100).toFixed(0);
let dynamicChart!: ECharts;

const updateChart = () => {
  dynamicChart.setOption<EChartsOption>({
    series: [
      {
        name: 'cpu',
        data: cpuParents,
      },
      {
        name: 'memory',
        data: memoryParents,
      },
    ]
  })
}

const handleChartResize = () => {
  dynamicChart.resize();
}

const initChart = () => {
  dynamicChart = echarts.init(document.getElementById(chartId), null, {
    renderer: 'canvas'
  })
  dynamicChart.setOption<EChartsOption>({
    tooltip: {
      trigger: 'axis',
    },
    grid: {
      show: false,
      left: '32px',
      right: 0,
      bottom: '48px',
      top: '24px',
    },
    xAxis: {
      type: 'time',
      boundaryGap: false,
      show: false
    },
    yAxis: {
      type: 'value',
      min: 0,
      max: 100,
      splitLine: {
        show: false
      }
    },
    series: [
      {
        name: 'cpu',
        type: 'line',
        showSymbol: false,
        smooth: true,
        itemStyle: {
          color: '#FF7F50',
        },
      },
      {
        name: 'memory',
        type: 'line',
        showSymbol: false,
        smooth: true,
        itemStyle: {
          color: '#03A9F4',
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            {
              offset: 0,
              color: '#A2D5F2'
            },
            {
              offset: 1,
              color: '#cde7f6'
            }
          ])
        },
      },
    ]
  })
  updateChart();
}

onMounted(() => {
  initChart();
  window.addEventListener('resize', handleChartResize);
})

onUnmounted(() => {
  window.removeEventListener('resize', handleChartResize);
})

defineExpose({
  pushNewValue
})

</script>

<template>
  <div :id="chartId" style="height: 100%; width: 100%"></div>
</template>

<style scoped lang="scss">

</style>