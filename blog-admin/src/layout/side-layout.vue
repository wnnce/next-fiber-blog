<script setup lang="ts">
import { IconApps, IconBug, IconBulb, } from '@arco-design/web-vue/es/icon';
import { useAppConfigStore } from '@/stores/app-config'

const configStore = useAppConfigStore();

</script>

<template>
  <div class="side-div">
    <a-menu :style="{ height: '100%', width: configStore.sideWidth }" show-collapse-button>
      <a-sub-menu key="0">
        <template #icon><icon-apps></icon-apps></template>
        <template #title>Navigation 1</template>
        <a-menu-item key="0_0">Menu 1</a-menu-item>
        <a-menu-item key="0_1">Menu 2</a-menu-item>
        <a-menu-item key="0_2" disabled>Menu 3</a-menu-item>
      </a-sub-menu>
      <a-sub-menu key="1">
        <template #icon><icon-bug></icon-bug></template>
        <template #title>Navigation 2</template>
        <a-menu-item key="1_0">Menu 1</a-menu-item>
        <a-menu-item key="1_1">Menu 2</a-menu-item>
        <a-menu-item key="1_2">Menu 3</a-menu-item>
      </a-sub-menu>
      <a-sub-menu key="2">
        <template #icon><icon-bulb></icon-bulb></template>
        <template #title>Navigation 3</template>
        <a-menu-item-group title="Menu Group 1">
          <a-menu-item key="2_0">Menu 1</a-menu-item>
          <a-menu-item key="2_1">Menu 2</a-menu-item>
        </a-menu-item-group>
        <a-menu-item-group title="Menu Group 2">
          <a-menu-item key="2_2">Menu 3</a-menu-item>
          <a-menu-item key="2_3">Menu 4</a-menu-item>
        </a-menu-item-group>
      </a-sub-menu>
    </a-menu>
  </div>

</template>

<style scoped lang="scss">
.side-div {
  box-shadow: 8px 0 8px -8px var(--shadom-color);
  height: 100%;
  border-top: 2px solid var(--border-color);
  .logo {
    padding: 24px 0;
  }
}
</style>