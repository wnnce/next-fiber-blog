<script setup lang="ts">
import { useDict } from '@/hooks/dict'
import { computed, onMounted, ref } from 'vue'
import type { DictValue } from '@/api/system/dict/types'

const { queryDict } = useDict();

const props = defineProps<{
  dictKey: string;
  value: string | number | boolean;
}>();

const dictLabel = computed((): string => {
  const valueString = props.value.toString();
  const find = dictValues.value.find(item => item.value === valueString);
  return find ? find.label : valueString;
})

const dictValues = ref<DictValue[]>([]);

const queryDictValues = async () => {
  const queryValues = await queryDict(props.dictKey);
  if (queryValues && queryValues.length > 0) {
    dictValues.value = queryValues;
  }
}

onMounted( () => {
  queryDictValues();
})
</script>

<template>
  <span>{{ dictLabel }}</span>
</template>

<style scoped lang="scss">

</style>