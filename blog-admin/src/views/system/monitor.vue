<script setup lang="ts">

import { recordApi } from '@/api/system/record'
import { onMounted, shallowReactive, ref } from 'vue'
import type { ApplicationMonitor } from '@/api/system/record/types'
import MemoryProcessChart from '@/components/charts/MemoryProcessChart.vue'
import MemoryPieChart from '@/components/charts/MemoryPieChart.vue'
import DynamicChart from '@/components/charts/DynamicChart.vue'

const monitorStats = shallowReactive<ApplicationMonitor>({
  hostname: '',
  platform: '',
  platformVersion: '',
  cpuNumber: 0,
  cpuPercent: 0.00,
  memoryTotal: 0,
  memoryUsed: 0,
  memoryAvailable: 0,
  memoryUsedPercent: 0.00,
  Sys: 0,
  HeapSys: 0,
  HeapInuse: 0,
  HeapIdle: 0,
  StackSys: 0,
  StackInuse: 0,
  PauseTotalNs: 0,
  NumGC: 0,
  GCCPUFraction: 0.00,
});

const queryApplicationMonitor = async () => {
  const result = await recordApi.applicationMonitor();
  if (result.code === 200) {
    const { hostname, platform, platformVersion, cpuNumber, cpuPercent, memoryTotal, memoryUsed,
      memoryAvailable, memoryUsedPercent, Sys, HeapSys, HeapInuse, HeapIdle, StackSys,
      StackInuse, PauseTotalNs, NumGC, GCCPUFraction, LastGC } = result.data;
    Object.assign(monitorStats, { hostname, platform, platformVersion, cpuNumber, cpuPercent,
      memoryTotal, memoryUsed, memoryAvailable, memoryUsedPercent, Sys, HeapSys, HeapInuse, HeapIdle,
      StackSys, StackInuse, PauseTotalNs, NumGC, GCCPUFraction, LastGC })
    dynamicChartRef.value.pushNewValue(memoryUsedPercent, cpuPercent);
    // 3s 查询一次
    setTimeout(() => {
      queryApplicationMonitor()
    }, 3000)
  }
}

const dynamicChartRef = ref();

onMounted(() => {
  queryApplicationMonitor();
})

</script>

<template>
  <div class="monitor-container">
    <div class="monitor-card">
      <h3 class="monitor-card-title">监控信息</h3>
      <table class="stats-table">
        <tr><td>主机名称</td><td>{{ monitorStats.hostname }}</td></tr>
        <tr><td>平台</td><td>{{ monitorStats.platform }}</td></tr>
        <tr><td>平台版本</td><td>{{ monitorStats.platformVersion }}</td></tr>
        <tr><td>GC次数</td><td>{{ monitorStats.NumGC }}</td></tr>
        <tr><td>总GC暂停时间</td><td>{{ (monitorStats.PauseTotalNs / 1e6).toFixed(0) + 'ms' }}</td></tr>
      </table>
    </div>
    <div class="chart-container flex">
      <div class="monitor-card" style="height: 24rem">
        <h3 class="monitor-card-title">内存/CPU</h3>
        <DynamicChart ref="dynamicChartRef" />
      </div>
      <div class="flex process-list">
        <div class="monitor-card">
          <h4 class="monitor-card-title">堆内存</h4>
          <MemoryProcessChart :value="(monitorStats.HeapInuse / 1024 / 1024).toFixed(2)"
                        :max="(monitorStats.HeapSys / 1024 / 1024).toFixed(2)"
                        color="#27AE60"
          />
        </div>
        <div class="monitor-card">
          <h4 class="monitor-card-title">程序内存分析</h4>
          <MemoryPieChart :sys-member="monitorStats.Sys" :heap-member="monitorStats.HeapSys"
                          :stack-member="monitorStats.StackSys"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.monitor-container {
  display: flex;
  flex-direction: column;
  .monitor-card {
    background-color: var(--card-color);
    padding: 1rem;
    border-radius: 0.5rem;
    .monitor-card-title {
      padding: 0.25rem 0 0.25rem 0.75rem;
      margin-bottom: 1rem;
      border-left: 6px solid #1E90FF;
    }
  }
  .stats-table {
    width: 100%;
    tr {
      td:first-child {
        color: var(--color-text-3);
      }
      td {
        padding: 0.75rem;
      }
    }
  }
  .chart-container {
    margin-top: 1rem;
    width: 100%;
    column-gap: 1rem;
    > div {
      flex: 1;
    }
    .process-list {
      gap: 0.5rem;
      flex-wrap: wrap;
      > .monitor-card {
        padding: 0.75rem;
        width: calc(50% - 0.25rem);
        height: calc(50% - 0.25rem);
        .monitor-card-title {
          margin-bottom: 0.75rem;
        }
      }
    }
  }
}
</style>