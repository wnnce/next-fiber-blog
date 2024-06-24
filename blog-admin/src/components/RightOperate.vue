<script setup lang="ts">

const emits = defineEmits<{
  (e: 'refresh'): void
}>()

</script>

<template>
  <div class="flex right-operate">
    <a-button size="small" shape="circle" @click="emits('refresh')">
      <template #icon>
        <icon-refresh />
      </template>
    </a-button>
  </div>
</template>

<style scoped lang="scss">
.right-operate {
  column-gap: var(--space-sm);
  align-items: flex-end;
}
</style>