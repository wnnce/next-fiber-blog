<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import RightOperate from '@/components/RightOperate.vue'
import { useArcoMessage } from '@/hooks/message'
import { constant } from '@/assets/script/constant'
import LoadImage from '@/components/LoadImage.vue'
import type { Article, ArticleQueryForm, ArticleUpdateFrom } from '@/api/blog/article/types'
import { articleApi } from '@/api/blog/article'
import ArticleForm from '@/views/blog/article/article-form.vue'
import { tagApi } from '@/api/blog/tags'
import type { SelectOptionData, TreeNodeData } from '@arco-design/web-vue'
import { categoryApi } from '@/api/blog/category'
import type { Category } from '@/api/blog/category/types'
import { useLocalStorage } from '@/hooks/local-storage'
import { useRouter } from 'vue-router'
import { parseWordCount } from '../../../assets/script/util'

const { successMessage, loading } = useArcoMessage();
const router = useRouter();

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

const tagsSelectOption = ref<SelectOptionData[]>([]);
const queryTagSelectData = async () => {
  const result = await tagApi.listTag();
  const { code, data } = result;
  if (code === 200 && data) {
    tagsSelectOption.value = data.map(item => {
      return {
        label: item.tagName,
        value: item.tagId,
      }
    })
  }
}

const categoryTreeSelectOption = ref<TreeNodeData[]>([]);
const queryCategorySelectData = async () => {
  const result = await categoryApi.tree();
  const { code, data } = result;
  if (code === 200 && data) {
    categoryTreeSelectOption.value = parseCategoryToSelectOption(data);
  }
}

const parseCategoryToSelectOption = (categorys: Category[]): TreeNodeData[] => {
  if (!categorys || categorys.length === 0) {
    return [];
  }
  return categorys.map(item => {
    return {
      key: item.categoryId,
      title: item.categoryName,
      children: item.children && item.children.length > 0 ? parseCategoryToSelectOption(item.children) : undefined,
    }
  })
}

const formRef = ref();
const showForm = (record?: Article) => {
  formRef.value.show(record);
}

const routeArticlePublish = (articleId?: number) => {
  console.log(articleId);
  if (articleId) {
    useLocalStorage().set<number>('edit-article', articleId, undefined);
  }
  router.push({ path: '/article/publish' })
}

onMounted(() => {
  queryTableData();
  queryTagSelectData();
  queryCategorySelectData();
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
        <label>标签</label>
        <a-select v-model="queryForm.tagId" :options="tagsSelectOption"
                  placeholder="请选择文章标签"
                  allow-clear
                  style="width: 200px"
                  @clear="queryForm.tagId = undefined"
        />
      </div>
      <div class="search-item">
        <label>分类</label>
        <a-tree-select v-model="queryForm.categoryId" :data="categoryTreeSelectOption"
                       placeholder="请选择文章分类"
                       allow-clear
                       style="width: 240px"
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
        <a-button type="primary" @click="routeArticlePublish(undefined)">
          <template #icon><icon-plus /></template>
          新增
        </a-button>
      </div>
      <right-operate @refresh="handleRefresh"/>
    </div>
    <a-spin :loading="tableLoading">
      <div class="article-list flex flex-column" v-if="tableData && tableData.length > 0">
        <div class="article-item flex justify-between" v-for="item in tableData" :key="item.articleId">
          <div class="item-info flex">
            <div class="item-cover flex item-center">
              <load-image :src="item.coverUrl" thumbnail width="120px" height="70px" radius="8px" />
            </div>
            <div class="item-text">
              <h3 class="title">{{ item.title }}</h3>
              <p class="summary desc-text">{{ item.summary }}</p>
            </div>
            <div class="flex item-center desc-text">
              <ul>
                <li><icon-eye/><span>{{ item.viewNum }}</span></li>
                <li><icon-share-alt/><span>{{ item.shareNum }}</span></li>
                <li><icon-thumb-up/><span>{{ item.voteUp }}</span></li>
                <li><icon-message/><span>{{ item.commentNum }}</span></li>
                <li><icon-book/><span>{{ parseWordCount(item.wordCount) }}</span></li>
                <li><icon-clock-circle/><span>{{ item.createTime }}</span></li>
              </ul>
            </div>
          </div>
          <div class="item-option desc-text flex item-center">
            <ul>
              <li>
                <a-switch v-model="item.isTop" unchecked-text="未置顶" checked-text="置顶"
                          :before-change="(newValue: boolean) => handleUpdateStatus({ articleId: item.articleId, isTop: newValue })"
                />
              </li>
              <li>
                <a-switch v-model="item.isHot" unchecked-text="非热门" checked-text="热门"
                          :before-change="(newValue: boolean) => handleUpdateStatus({ articleId: item.articleId, isHot: newValue })"
                />
              </li>
              <li>
                <a-switch v-model="item.isComment" unchecked-text="评论关" checked-text="评论开"
                          :before-change="(newValue: boolean) => handleUpdateStatus({ articleId: item.articleId, isComment: newValue })"
                />
              </li>
              <li>
                <a-switch :model-value="item.isPrivate" disabled unchecked-text="公开" checked-text="私密" />
              </li>
              <li>
                <a-switch :unchecked-value="1" :checked-value="0" v-model="item.status"
                          unchecked-text="禁用" checked-text="正常"
                          :before-change="(newValue: number) => handleUpdateStatus({ articleId: item.articleId, status: newValue })"
                />
              </li>
              <li>
                <a-button type="text" shape="round" @click="routeArticlePublish(item.articleId)">
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
    <article-form ref="formRef" @reload="handleQuery" />
  </div>
</template>

<style scoped lang="scss">
.article-list {
  .article-item:last-child {
    border-bottom: none;
  }
  .article-item:hover {
    background-color: var(--color-border-1);
  }
  .article-item {
    padding: var(--space-sm);
    border-bottom: 1px solid var(--color-border-2);
    transition: background-color 300ms ease;
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