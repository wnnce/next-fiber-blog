<script setup lang="ts">
import { useDict } from '@/hooks/dict'
import { onMounted, ref } from 'vue'

const { queryDict } = useDict();

const props = defineProps<{
  dictKey: string;
  value: string | number | boolean;
}>();

const dictLabel = ref<string>(props.value.toString());

const findDictLabel = async () => {
  const dictValue = await queryDict(props.dictKey);
  if (dictValue && dictValue.length > 0) {
    const valueString = props.value.toString();
    const findDict = dictValue.find(item => item.value === valueString);
    findDict && (dictLabel.value = findDict.label);
  }
}

onMounted( () => {
  findDictLabel();
})
</script>

<template>
  <span>{{ dictLabel }}</span>
</template>

<style scoped lang="scss">

</style>