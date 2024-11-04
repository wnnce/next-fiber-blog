<script setup lang="ts">

import type { ECharts } from 'echarts'
import { computed, onMounted, onUnmounted, watch } from 'vue'
import * as echarts from 'echarts'
import type { EChartsOption } from 'echarts/types/dist/echarts'

const props = defineProps<{
  sysMember: number,
  heapMember: number,
  stackMember: number,
}>();

const otherMember = computed(() => {
  return props.sysMember - props.heapMember - props.stackMember;
})

const pieChartId = new Date().getTime().toString() + (Math.random() * 100).toFixed(0);
let pieChart!: ECharts;

watch(props, () => {
  updatePieChart();
})

const updatePieChart = () => {
  pieChart.setOption<EChartsOption>({
    series: [
      {
        data: [
          { value: (props.heapMember / 1024 / 1024).toFixed(2), name: '堆内存' },
          { value: (props.stackMember / 1024 / 1024).toFixed(2), name: '栈内存' },
          otherMember.value > 0 && { value: (otherMember.value / 1024 / 1024).toFixed(2), name: '其它' }
        ]
      }
    ]
  })
}

const initPieChart = () => {
  pieChart = echarts.init(document.getElementById(pieChartId));
  pieChart.setOption<EChartsOption>({
    tooltip: {
      trigger: 'item'
    },
    legend: {
      show: false,
    },
    series: [
      {
        type: 'pie',
        // adjust the start and end angle
        startAngle: 180,
        endAngle: 360,
        radius: ['95%', '130%'],
        center: ['50%', '75%'],
        label: {
          show: false
        },
        data: [],
      }
    ]
  })
  updatePieChart();
}

const handleChartResize = () => {
  pieChart.resize();
}

onMounted(() => {
  initPieChart();
  window.addEventListener('resize', handleChartResize);
})

onUnmounted(() => {
  window.removeEventListener('resize', handleChartResize);
})

</script>

<template>
  <div :id="pieChartId" style="height: 100%; width: 100%">

  </div>
</template>

<style scoped lang="scss">

</style>