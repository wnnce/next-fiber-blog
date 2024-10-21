<script setup lang="ts">

import type { DayStats } from '@/api/blog/site'
import * as echarts from 'echarts'
import { onMounted, onUnmounted, watch } from 'vue'
import { EChartsOption } from 'echarts/types/dist/echarts'
import type { ECharts } from 'echarts'
import { throttle } from '@/assets/script/util'

interface Props {
  data: DayStats[];
  seriesName?: string;
  echartTitle?: string;
  lineColor?: string;
  areaStartColor?: string;
  areaEndColor?: string;
}

const props = withDefaults(defineProps<Props>(), {
  seriesName: '数据',
  echartTitle: '样式图表',
  lineColor: '#3185FC',
  areaStartColor: '#3A4DE9CC',
  areaEndColor: '#3A4DE94C',
})

const elementId = new Date().getTime().toString() + Math.random().toString()

let customEcharts!: ECharts;

let dateList: string[], countList: number[];

watch(props, () => {
  dateList = [];
  countList = [];
  props.data.forEach(item => {
    dateList.push(item.dateItem);
    countList.push(item.countItem);
  })
  updateEchartsData()
})

const updateEchartsData = () => {
  customEcharts.setOption<EChartsOption>({
    xAxis: { data: dateList, show: false, boundaryGap: false },
    series: [
      {
        name: props.seriesName,
        data: countList,
      }
    ]
  })
}

const initEcharts = () => {
  const dom = document.getElementById(elementId);
  customEcharts = echarts.init(dom, null, {
    renderer: 'canvas',
  })
  customEcharts.setOption<EChartsOption>({
    tooltip: {
      trigger: 'axis'
    },
    grid: {
      show: false,
      left: '4px',
      right: '4px',
      bottom: '12px',
    },
    xAxis: {
      data: []
    },
    yAxis: {
      splitLine: false,
      axisLabel: {
        show: false,
      }
    },
    series: [
      {
        name: props.seriesName,
        type: 'line',
        smooth: true,
        data: [],
        symbol: 'none',
        lineStyle: {
          color: props.lineColor,
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            {
              offset: 0,
              color: props.areaStartColor
            },
            {
              offset: 1,
              color: props.areaEndColor
            }
          ])
        },
      }
    ]
  })
  updateEchartsData()
}

const handleEchartsResize = throttle(() => {
  customEcharts && (customEcharts.resize());
})

onMounted(() => {
  initEcharts();
  window.addEventListener('resize', handleEchartsResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleEchartsResize)
})

</script>

<template>
  <div class="home-echart-container">
    <h3 class="info-text">{{ echartTitle }}</h3>
    <div class="echarts" :id="elementId"></div>
  </div>

</template>

<style scoped lang="scss">
.home-echart-container {
  > h3 {
    padding: 0.25rem 0 0.25rem 0.5rem;
    border-left: 6px solid var(--color-primary-light-4);
  }
  height: 100%;
  width: 100%;
  .echarts {
    width: 100%;
    height: 100%;
  }
}
</style>