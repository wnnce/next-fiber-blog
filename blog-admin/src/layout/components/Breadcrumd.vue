<script setup lang="ts">

import { useRoute } from 'vue-router'
import { computed } from 'vue'
import { useLocalUserStore } from '@/stores/user'

const route = useRoute();

const routePaths = computed((): string[] => {
  const currentId = route.name ? route.name.toString() : undefined;
  const baseName = ['首页'];
  if (!currentId) {
    return baseName
  }
  const parentIds = useLocalUserStore().getMenuParentIdListMap().get(currentId);
  const finalIds: string[] = parentIds ? parentIds.concat(currentId) : [currentId];
  const names = finalIds.map(id => queryMenuName(id))
  return baseName.concat(names);
})

const queryMenuName = (id: string): string => {
  const menu = useLocalUserStore().menuListMap.get(id);
  return menu ? menu.menuName : '';
}
</script>

<template>
  <a-breadcrumb>
    <template v-for="(name, index) in routePaths" :key="index">
      <a-breadcrumb-item>{{ name }}</a-breadcrumb-item>
    </template>
  </a-breadcrumb>
</template>

<style scoped lang="scss">

</style>