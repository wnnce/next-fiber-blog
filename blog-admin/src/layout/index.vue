<script setup lang="ts">
import SideLayout from '@/layout/side-layout.vue'
import HeaderLayout from '@/layout/header-layout.vue'
import PageTabs from '@/layout/components/PageTabs.vue'
import { useLocalUserStore } from '@/stores/user'

</script>

<template>
  <main class="container flex flex-column">
    <div class="header">
      <header-layout />
    </div>
    <div class="content flex">
      <div class="content-side">
        <side-layout />
      </div>
      <div class="content-main">
        <div class="page-tabs">
          <page-tabs />
        </div>
        <div class="main-div">
          <router-view v-slot="{ Component }">
            <transition name="switch" mode="out-in">
              <keep-alive :include="useLocalUserStore().keepaliveInclude">
                <component :is="Component" />
              </keep-alive>
            </transition>
          </router-view>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped lang="scss">
.container {
  height: 100vh;
  width: 100%;
  background-color: var(--background-color);
  transition: color 0.5s, background-color 0.5s;
  .header {
    flex-shrink: 1;
  }
  .content {
    flex: 1;
    .content-side {
      flex-shrink: 1;
    }
    .content-main {
      flex: 1;
      .page-tabs {
        background-color: var(--card-color);
        border-top: 1px solid var(--border-color);
      }
      .main-div {
        overflow: auto;
        padding: var(--space-md);
      }
    }
  }
}
</style>