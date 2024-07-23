<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { SelectOptionData } from '@arco-design/web-vue'
import { useDict } from '@/hooks/dict'

export declare type DictValueType = 'string' | 'number' | 'boolean';

const { queryDict } = useDict();

interface Props {
  dictKey: string;
  width?: string | number;
  placeholder?: string;
  type?: DictValueType;
}

const props = withDefaults(defineProps<Props>(), {
  width: '100%',
  placeholder: '请选择数据',
  type: 'string'
})

const modelValue = defineModel<string | number | boolean | undefined>('modelValue', {
  required: true
})

const options = ref<SelectOptionData[]>([]);

const queryDictValues = async () => {
  const dictValues = await queryDict(props.dictKey);
  if (dictValues && dictValues.length > 0) {
    options.value = dictValues.map(item => {
      let newValue: string | number | boolean | undefined = undefined;
      if (props.type === 'string') {
        newValue = item.value;
      } else if (props.type === 'number') {
        newValue = Number(item.value);
      } else {
        newValue = Boolean(item.value);
      }
      return {
        label: item.label,
        value: newValue
      }
    })
  }
}

onMounted(() => {
  queryDictValues();
})
</script>

<template>
  <a-select allow-clear v-model="modelValue" :placeholder="placeholder">
    <a-option v-for="(item, index) in options" :label="item.label" :value="item.value" :key="index" />
  </a-select>
</template>

<style scoped lang="scss">

</style>