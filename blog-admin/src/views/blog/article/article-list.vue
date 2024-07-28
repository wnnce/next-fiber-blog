<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import RightOperate from '@/components/RightOperate.vue'
import { useArcoMessage } from '@/hooks/message'
import { constant } from '@/assets/script/constant'
import LoadImage from '@/components/LoadImage.vue'
import LinksForm from '@/views/blog/links/links-form.vue'
import type { Article, ArticleQueryForm, ArticleUpdateFrom } from '@/api/blog/article/types'
import { articleApi } from '@/api/blog/article'
import ArticleForm from '@/views/blog/article/article-form.vue'

const { successMessage, loading } = useArcoMessage();

const tableLoading = ref<boolean>(false);
const recordTotal = ref<number>(0);
const tableData = ref<Article[]>([]);
const defaultQueryForm: ArticleQueryForm = {
  page: 1,
  size: 10,
  title: undefined,
  categoryId: undefined,
  tagId: undefined,
  status: undefined,
  createTimeBegin: undefined,
  createTimeEnd: undefined
}
const dateRange = ref<string[]>([]);
const queryForm = reactive<ArticleQueryForm>({ ...defaultQueryForm })
const queryTableData = async () => {
  tableLoading.value = true;
  try {
    const result = await articleApi.pageArticle(queryForm);
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

const handleDelete = async (record: Article) => {
  const loadingMsg = loading('数据删除中')
  try {
    const result = await articleApi.deleteArticle(record.articleId);
    if (result.code === 200) {
      successMessage('删除成功');
      await queryTableData();
    }
  } finally {
    loadingMsg.close();
  }
}

const handleUpdateStatus = async (form: ArticleUpdateFrom) => {
  const result = await articleApi.updateSelective(form);
  if (result.code === 200) {
    successMessage('更新成功');
    return true;
  }
  return false;
}

const formRef = ref();
const showForm = (record?: Article) => {
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
        <label>文章标题</label>
        <a-input v-model="queryForm.title" placeholder="请输入文章标题" />
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
        <a-button type="primary">
          <template #icon><icon-plus /></template>
          新增
        </a-button>
      </div>
      <right-operate @refresh="handleRefresh"/>
    </div>
    <div class="article-list flex flex-column">
      <transition-group name="switch">
        <div class="article-item flex justify-between" v-for="item in tableData" :key="item.articleId">
          <div class="item-info flex">
            <div class="item-cover flex item-center">
              <load-image :src="item.coverUrl" thumbnail width="120px" height="70px" radius="8px" />
            </div>
            <div class="item-text">
              <h3 class="title">{{ item.title }}</h3>
              <p class="summary info-text">{{ item.summary }}</p>
            </div>
            <div class="flex item-center info-text">
              <ul>
                <li><icon-eye/><span>{{ item.viewNum }}</span></li>
                <li><icon-share-alt/><span>{{ item.shareNum }}</span></li>
                <li><icon-thumb-up/><span>{{ item.voteUp }}</span></li>
                <li><icon-message/><span>0</span></li>
                <li><icon-clock-circle/><span>{{ item.createTime }}</span></li>
              </ul>
            </div>
          </div>
          <div class="item-option desc-text flex item-center">
            <ul>
              <li>
                <a-switch v-model="item.isTop" unchecked-text="未置顶" checked-text="置顶"
                          :before-change="newValue => handleUpdateStatus({ articleId: item.articleId, isTop: Boolean(newValue) })"
                />
              </li>
              <li>
                <a-switch v-model="item.isHot" unchecked-text="非热门" checked-text="热门"
                          :before-change="newValue => handleUpdateStatus({ articleId: item.articleId, isHot: Boolean(newValue) })"
                />
              </li>
              <li>
                <a-switch v-model="item.isComment" unchecked-text="评论关" checked-text="评论开"
                          :before-change="newValue => handleUpdateStatus({ articleId: item.articleId, isComment: Boolean(newValue) })"
                />
              </li>
              <li>
                <a-switch :model-value="item.isPrivate" disabled unchecked-text="公开" checked-text="私密" />
              </li>
              <li>
                <a-switch :unchecked-value="1" :checked-value="0" v-model="item.status"
                          unchecked-text="禁用" checked-text="正常"
                          :before-change="newValue => handleUpdateStatus({ articleId: item.articleId, status: Number(newValue) })"
                />
              </li>
              <li>
                <a-button type="text" shape="round">
                  <template #icon>
                    <icon-edit />
                  </template>
                </a-button>
                <a-button type="text" shape="round" @click="showForm(item)">
                  <template #icon>
                    <icon-settings />
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
      </transition-group>
    </div>
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
    <article-form ref="formRef" @reload="handleQuery" />
  </div>
</template>

<style scoped lang="scss">
.article-list {
  row-gap: var(--space-sm);
  .article-item {
    > div {
      column-gap: var(--space-sm);
    }
    .item-text {
      min-width: 200px;
      max-width: 300px;
      h3, p {
        overflow: hidden;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        text-overflow: ellipsis;
      }
      h3 {
        margin-bottom: var(--space-xs);
        line-height: 130%;
      }
      p {
        line-height: 120%;
      }
    }
    ul {
      list-style: none;
      padding: 0;
    }
    li {
      display: inline;
      margin: 0 var(--space-xs);
    }
    .item-info {
      li > span {
        margin-left: 6px;
      }
    }
  }
}
</style>