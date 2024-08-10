<script setup lang="ts">
import { onMounted, ref, shallowRef } from 'vue'
import { siteApi, type SiteConfigurationItem, siteConfigurationRequiredField } from '@/api/blog/site'
import MarkdownPreview from '@/components/MarkdownPreview.vue'
import LoadImage from '@/components/LoadImage.vue'
import ConfigurationForm from '@/views/blog/site/configuration-form.vue'

const loading = ref<boolean>(false);
const siteConfiguration = shallowRef<Record<string, SiteConfigurationItem>>({});
const querySiteConfiguration = async () => {
  const result = await siteApi.configuration();
  if (result.code === 200 && result.data) {
    const tempRecord: Record<string, SiteConfigurationItem> = {};
    siteConfigurationRequiredField.forEach(key => {
      tempRecord[key] = result.data[key];
      tempRecord[key].extend = false;
    })
    Object.keys(result.data).filter(key => !siteConfigurationRequiredField.includes(key)).map(key => {
      tempRecord[key] = result.data[key];
      tempRecord[key].extend = true;
    })
    siteConfiguration.value = tempRecord;
  }
}

const formRef = ref();
const showForm = () => {
  formRef.value.show(siteConfiguration.value);
}

onMounted(() => {
  querySiteConfiguration()
})
</script>

<template>
  <div class="table-card">
    <div class="flex justify-end">
      <a-button type="primary" @click="showForm">
        <template #icon>
          <icon-edit />
        </template>
        修改配置
      </a-button>
    </div>
    <a-spin :loading="loading">
      <div class="configuration-table">
        <div class="table-item flex" v-for="(item, key) in siteConfiguration" :key="key">
          <div class="item-key">
            {{ item.name || key }}
          </div>
          <div class="item-value">
            <template v-if="item.type === 'text' || item.type === 'number'">
              {{ item.value }}
            </template>
            <template v-else-if="item.type === 'html'">
              <div v-html="item.value"></div>
            </template>
            <template v-else-if="item.type === 'markdown'">
              <markdown-preview :markdown="item.value.toString()" />
            </template>
            <template v-else-if="item.type === 'image'">
              <load-image :src="item.value.toString()" lazy preview radius="12px" width="100px" height="100px" />
            </template>
            <template v-else-if="item.type === 'color'">
              <a-color-picker :default-value="item.value.toString()" disabled />
            </template>
          </div>
        </div>
      </div>
    </a-spin>
    <configuration-form ref="formRef" @reload="querySiteConfiguration" />
  </div>
</template>

<style scoped lang="scss">
.configuration-table {
  border: 1px solid var(--color-border-2);
  border-collapse: collapse;
  .table-item {
    > div {
      padding: var(--space-mm);
      box-shadow: 0 0 1px var(--color-border-2);
    }
    > div:first-child {
      text-align: end;
      font-weight: 500;
      color: var(--color-text-2);
      min-width: 120px;
      background-color: var(--color-neutral-1);
    }
    > div:last-child {
      flex: 1;
    }
  }
}

</style>