<script setup lang="ts">

import type { ECharts } from 'echarts'
import { onMounted, onUnmounted, watch } from 'vue'
import type { EChartsOption } from 'echarts/types/dist/echarts'
import * as echarts from 'echarts'
import { throttle } from '@/assets/script/util'

interface Props {
  min?: number | string;
  max?: number | string;
  value: number | string;
  color?: string;
  width?: string | number;
  fontSize?: number;
}

const props = withDefaults(defineProps<Props>(), {
  min: 0,
  max: 100,
  color: 'red',
  width: 24,
  fontSize: 16
})

watch(props, () => {
  updateProcessChart();
})

const processChartId = new Date().getTime().toString() + (Math.random() * 100).toFixed(0);

let processChart!: ECharts;

const updateProcessChart = () => {
  processChart.setOption<EChartsOption>({
    series: [
      {
        max: props.max,
        data: [
          { value: props.value }
        ]
      }
    ]
  })
}

const handleEchartsResize = throttle(() => {
  processChart && (processChart.resize());
})

const initProcessChart = () => {
  processChart = echarts.init(document.getElementById(processChartId), null, {
    renderer: 'canvas'
  });
  processChart.setOption<EChartsOption>({
    grid: {
      left: 0,
      right: 0,
      top: 0,
      bottom: 0
    },
    series: [
      {
        type: 'gauge',
        startAngle: 190,
        endAngle: -10,
        radius: '130%',
        center: ['50%', '65%'],
        min: 0,
        max: props.max,
        splitNumber: 0,
        itemStyle: {
          color: props.color,
        },
        progress: {
          show: true,
          width: props.width
        },
        pointer: {
          show: false
        },
        axisLine: {
          lineStyle: {
            width: props.width
          }
        },
        splitLine: {
          show: false
        },
        axisLabel: {
          show: false
        },
        anchor: {
          show: false
        },
        title: {
          show: false
        },
        detail: {
          valueAnimation: true,
          width: '60%',
          lineHeight: 30,
          borderRadius: 8,
          offsetCenter: [0, 0],
          fontSize: props.fontSize,
          fontWeight: 'bolder',
          formatter: '{value} MB',
          color: 'inherit'
        }
      }
    ]
  })
  updateProcessChart();
}

onMounted(() => {
  initProcessChart();
  window.addEventListener('resize', handleEchartsResize);
})

onUnmounted(() => {
  window.removeEventListener('resize', handleEchartsResize);
})

</script>

<template>
  <div :id="processChartId" style="height: 100%; width: 100%">

  </div>
</template>

<style scoped lang="scss">

</style>