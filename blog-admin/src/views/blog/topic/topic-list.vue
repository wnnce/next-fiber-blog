<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import RightOperate from '@/components/RightOperate.vue'
import { useArcoMessage } from '@/hooks/message'
import { constant } from '@/assets/script/constant'
import type { Topic, TopicQueryForm, TopicUpdateForm } from '@/api/blog/topic/types'
import { topicApi } from '@/api/blog/topic'
import TopicForm from '@/views/blog/topic/topic-form.vue'
import DictSelect from '@/components/DictSelect.vue'
import DictLabel from '@/components/DictLabel.vue'
import MarkdownPreview from '@/components/MarkdownPreview.vue'
import LoadImage from '@/components/LoadImage.vue'

const { successMessage, loading } = useArcoMessage();

const tableLoading = ref<boolean>(false);
const recordTotal = ref<number>(0);
const tableData = ref<Topic[]>([]);
const defaultQueryForm: TopicQueryForm = {
  page: 1,
  size: 10,
  location: undefined,
  status: undefined,
  createTimeBegin: undefined,
  createTimeEnd: undefined
}
const dateRange = ref<string[]>([]);
const queryForm = reactive<TopicQueryForm>({ ...defaultQueryForm })
const queryTableData = async () => {
  tableLoading.value = true;
  try {
    const result = await topicApi.pageTopic(queryForm);
    if (result && result.code === 200) {
      const { total, records } = result.data;
      recordTotal.value = total;
      tableData.value = records;
    }
  } finally {
    tableLoading.value = false;
  }
}

const queryLoading = ref<boolean>(false);
const handleQuery = async () => {
  queryForm.page = 1;
  queryLoading.value = true;
  try {
    await queryTableData();
  } finally {
    queryLoading.value = false;
  }
}

const handleReset = () => {
  Object.assign(queryForm, defaultQueryForm);
  dateRange.value = [];
}

const handleRefresh = () => {
  queryTableData();
}

const handleDateChange = () => {
  if (!dateRange.value) {
    queryForm.createTimeBegin = undefined;
    queryForm.createTimeEnd = undefined;
  } else {
    const [ begin, end ] = dateRange.value;
    queryForm.createTimeBegin = begin;
    queryForm.createTimeEnd = end;
  }
}

const handleDelete = async (record: Topic) => {
  const loadingMsg = loading('数据删除中')
  try {
    const result = await topicApi.deleteTopic(record.topicId);
    if (result.code === 200) {
      successMessage('删除成功');
      await queryTableData();
    }
  } finally {
    loadingMsg.close();
  }
}

const handleUpdateSelective = async (form: TopicUpdateForm) => {
  const result = await topicApi.updateSelective(form);
  if (result.code === 200) {
    successMessage('更新成功');
    return true;
  }
  return false;
}

const formRef = ref();
const showForm = (record?: Topic) => {
  formRef.value.show(record);
}

onMounted(() => {
  queryTableData();
})
</script>

<template>
  <div class="table-card">
    <div class="search-div">
      <div class="search-item">
        <label>地点</label>
        <a-input v-model="queryForm.location" placeholder="请输入发布地点" />
      </div>
      <div class="search-item">
        <label>状态</label>
        <dict-select dict-key="dict_status" v-model="queryForm.status" type="number"
                     width="180px"
                     placeholder="请选择状态"
        />
      </div>
      <div class="search-item">
        <label>创建时间</label>
        <a-range-picker v-model="dateRange" @change="handleDateChange" />
      </div>
      <div class="search-buttons">
        <a-button type="primary" @click="handleQuery" :loading="queryLoading">
          <template #icon><icon-search /></template>
          搜索
        </a-button>
        <a-button @click="handleReset">
          <template #icon><icon-refresh /></template>
          重置
        </a-button>
      </div>
    </div>
    <div class="flex justify-between">
      <div class="flex" style="column-gap: 12px">
        <a-button type="primary" @click="showForm(undefined)">
          <template #icon><icon-plus /></template>
          新增
        </a-button>
      </div>
      <right-operate @refresh="handleRefresh"/>
    </div>
    <a-spin :loading="tableLoading">
      <div class="topic-list flex flex-column" v-if="tableData && tableData.length > 0">
        <div class="topic-item flex flex-column" v-for="item in tableData" :key="item.topicId">
          <div class="topic-header flex justify-between item-center">
            <p class="info-text">
              <icon-clock-circle /> {{ item.createTime }}
            </p>
            <div class="header-option flex item-center">
              <ul>
                <li>
                  <a-switch v-model="item.isTop" unchecked-text="未置顶" checked-text="置顶"
                            :before-change="newValue => { handleUpdateSelective({ topicId: item.topicId, isTop: Boolean(newValue) }) }"
                  />
                </li>
                <li>
                  <a-switch v-model="item.isHot" unchecked-text="非热门" checked-text="热门"
                            :before-change="newValue => { handleUpdateSelective({ topicId: item.topicId, isHot: Boolean(newValue) }) }"
                  />
                </li>
                <li>
                  <a-switch v-model="item.status" unchecked-text="禁用" checked-text="正常"
                            :checked-value="0" :unchecked-value="1"
                            :before-change="newValue => { handleUpdateSelective({ topicId: item.topicId, status: Number(newValue) }) }"
                  />
                </li>
                <li>
                  <a-button type="text" shape="round" @click="showForm(item)">
                    <template #icon>
                      <icon-edit />
                    </template>
                  </a-button>
                  <a-popconfirm content="是否确认删除数据？" type="error" position="lt"
                                :ok-button-props="{ status: 'danger' }"
                                @ok="handleDelete(item)"
                  >
                    <a-button type="text" shape="circle" status="danger">
                      <template #icon><icon-delete /></template>
                    </a-button>
                  </a-popconfirm>
                </li>
              </ul>
            </div>
          </div>
          <div class="item-content radius-md">
            <markdown-preview :markdown="item.content" />
            <div class="content-images flex" v-if="item.imageUrls">
              <template v-for="url in item.imageUrls" :key="url">
                <load-image :src="url" thumbnail lazy width="100px" height="100px" preview />
              </template>
            </div>
          </div>
          <div class="item-footer flex desc-text item-center">
            <ul>
              <li><icon-location /><span>中国重庆</span></li>
              <li><icon-thumb-up /><span>{{ item.voteUp }}</span></li>
              <li><icon-branch /><dict-label dict-key="topic_mode" :value="item.mode" /></li>
            </ul>
          </div>
        </div>
      </div>
      <a-empty v-else />
    </a-spin>
    <div class="flex justify-end">
      <a-pagination :total="recordTotal" size="medium"
                    v-model:current="queryForm.page"
                    v-model:page-size="queryForm.size"
                    :page-size-options="constant.pageSizeOption"
                    show-page-size
                    show-total
                    @page-size-change="handleQuery"
                    @change="queryTableData"
      />
    </div>
    <topic-form ref="formRef" @reload="handleQuery" />
  </div>
</template>

<style scoped lang="scss">
.topic-list {
  .topic-item {
    border-bottom: 1px solid var(--border-color);
  }
  .topic-item:last-child {
    border-bottom: none;
  }
}
.topic-item {
  padding: var(--space-sm);
  row-gap: var(--space-sm);
  .item-content {
    padding: var(--space-mm);
    background-color: var(--color-neutral-1);
    .content-images {
      margin-top: var(--space-mm);
      column-gap: var(--space-mm);
    }
  }
  .item-footer {
    column-gap: 48px;
    padding: 0;
    > ul > li > span {
      margin: 0 6px;
    }
  }
}
ul {
  list-style: none;
  padding: 0;
  li {
    display: inline;
    margin: 0 var(--space-xs);
  }
}
</style>