<script setup lang="ts">

import { useRoute } from 'vue-router'
import { onMounted, ref, watch } from 'vue'
import { useLocalUserStore } from '@/stores/user'

const route = useRoute();

const routePaths = ref<string[]>(['首页']);

const updateBreadcrumb = () => {
  const menuId = route.name ? route.name.toString() : undefined;
  routePaths.value = ['首页'];
  if (!menuId) {
    return;
  }
  const parentIds = useLocalUserStore().getMenuParentIdListMap().get(menuId);
  if (!parentIds || parentIds.length === 0) {
    return;
  }
  const names = parentIds.concat(menuId).map(id => queryMenuName(id))
  routePaths.value.push(...names);
}

const queryMenuName = (id: string): string => {
  const menu = useLocalUserStore().menuList.find(item => item.menuId.toString() === id)
  return menu ? menu.menuName : '';
}

watch(route, () => {
  updateBreadcrumb();
})

onMounted(() => {
  updateBreadcrumb();
})
</script>

<template>
  <a-breadcrumb>
    <a-breadcrumb-item v-for="(name, index) in routePaths" :key="index">
      {{ name }}
    </a-breadcrumb-item>
  </a-breadcrumb>
</template>

<style scoped lang="scss">

</style>