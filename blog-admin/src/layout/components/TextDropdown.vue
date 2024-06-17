<script setup lang="ts">

import type { OptionItem } from '@/assets/script/types'

interface Props {
  text: string | number;
  options: OptionItem[];
}

defineProps<Props>();

const emits = defineEmits<{
  (e: 'select', value: string | number): void
}>()

</script>

<template>
  <a-dropdown :popup-max-height="false" @select="emits('select', $event)">
    <span class="pointer">{{ text }} <icon-down/></span>
    <template #content>
      <a-doption v-for="(item, index) in options" :key="index" :value="item.value">
        {{ item.label }}
      </a-doption>
    </template>
  </a-dropdown>
</template>

<style scoped lang="scss">
.arco-dropdown-open .arco-icon-down {
  transform: rotate(180deg);
}
</style>