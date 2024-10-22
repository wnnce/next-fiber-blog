<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import RightOperate from '@/components/RightOperate.vue'
import { useArcoMessage } from '@/hooks/message'
import { constant } from '@/assets/script/constant'
import LoadImage from '@/components/LoadImage.vue'
import type { Link, LinkQueryForm } from '@/api/blog/link/types'
import { linkApi } from '@/api/blog/link'
import LinksForm from '@/views/blog/links/links-form.vue'
import type { Comment, CommentQueryForm, CommentUpdateForm } from '@/api/blog/comment/types'
import { commentApi } from '@/api/blog/comment'
import DictSelect from '@/components/DictSelect.vue'
import DictLabel from '@/components/DictLabel.vue'

const { successMessage, loading } = useArcoMessage();

const tableLoading = ref<boolean>(false);
const recordTotal = ref<number>(0);
const tableData = ref<Comment[]>([]);
const defaultQueryForm: CommentQueryForm = {
  page: 1,
  size: 10,
  commentType: undefined,
  articleId: undefined,
  topicId: undefined,
  fid: undefined,
  createTimeBegin: undefined,
  createTimeEnd: undefined
}
const dateRange = ref<string[]>([]);
const tableExtendKeys = ref<string[]>([]);
const queryForm = reactive<CommentQueryForm>({ ...defaultQueryForm })
const queryTableData = async () => {
  tableLoading.value = true;
  try {
    const result = await commentApi.pageComment(queryForm);
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

const handleDelete = async (record: Comment) => {
  const loadingMsg = loading('数据删除中')
  try {
    const result = await commentApi.deleteComment(record.commentId);
    if (result.code === 200) {
      successMessage('删除成功');
      await queryTableData();
    }
  } finally {
    loadingMsg.close();
  }
}

const handleUpdateStatus = async (form: CommentUpdateForm) => {
  const result = await commentApi.updateSelective(form);
  if (result.code === 200) {
    successMessage('更新成功');
    return true;
  }
  return false;
}

const formRef = ref();
const showForm = (record?: Link) => {
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
        <label>评论类型</label>
        <dict-select dict-key="comment_type" v-model="queryForm.commentType" type="number" width="200px" placeholder="选择文章类型" />
      </div>
      <div class="search-item">
        <label>文章ID</label>
        <a-input-number hide-button v-model="queryForm.articleId" placeholder="查询文章评论" />
      </div>
      <div class="search-item">
        <label>动态ID</label>
        <a-input-number hide-button v-model="queryForm.topicId" placeholder="查询动态评论" />
      </div>
      <div class="search-item">
        <label>一级评论ID</label>
        <a-input-number hide-button v-model="queryForm.fid" placeholder="查询子评论" />
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
      <div></div>
      <right-operate @refresh="handleRefresh"/>
    </div>
    <a-table :expandable="{ title: '展开', width: 60 }"
             v-model:expanded-keys="tableExtendKeys"
             column-resizable :data="tableData"
             :loading="tableLoading"
             :pagination="false"
             row-key="commentId"
    >
      <template #columns>
        <a-table-column title="ID" data-index="commentId" width="60" />
        <a-table-column title="内容" data-index="content" width="500" ellipsis />
        <a-table-column title="用户" data-index="username" />
        <a-table-column title="点赞" data-index="voteUp" align="center" />
        <a-table-column title="踩" data-index="voteDown" align="center" />
        <a-table-column title="一级评论" align="center" >
          <template #cell="{ record }">
            {{ record.fid === 0 ? '是' : '否' }}
          </template>
        </a-table-column>
        <a-table-column title="热门" :width="80" align="center">
          <template #cell="{ record }">
            <a-switch v-model="record.isHot"
                      :before-change="(newValue: boolean) => handleUpdateStatus({ isHot: newValue, commentId: record.commentId })"
            />
          </template>
        </a-table-column>
        <a-table-column title="置顶" :width="80" align="center">
          <template #cell="{ record }">
            <a-switch v-model="record.isTop"
                      :before-change="(newValue: boolean) => handleUpdateStatus({ isTop: newValue, commentId: record.commentId })"
            />
          </template>
        </a-table-column>
        <a-table-column title="折叠" :width="80" align="center">
          <template #cell="{ record }">
            <a-switch v-model="record.isColl"
                      :before-change="(newValue: boolean) => handleUpdateStatus({ isColl: newValue, commentId: record.commentId })"
            />
          </template>
        </a-table-column>
        <a-table-column title="状态" :width="80" align="center">
          <template #cell="{ record }">
            <a-switch :checked-value="0" :unchecked-value="1" v-model="record.status"
                      :before-change="(newValue: number) => handleUpdateStatus({ status: newValue, commentId: record.commentId })"
            />
          </template>
        </a-table-column>
        <a-table-column title="操作" align="center">
          <template #cell="{ record }">
            <a-button type="text" shape="circle" @click="showForm(record)">
              <template #icon><icon-edit /></template>
            </a-button>
            <a-popconfirm content="是否确认删除数据？" type="error" position="lt"
                          :ok-button-props="{ status: 'danger' }"
                          @ok="handleDelete(record)"
            >
              <a-button type="text" shape="circle" status="danger">
                <template #icon><icon-delete /></template>
              </a-button>
            </a-popconfirm>
          </template>
        </a-table-column>
      </template>
      <template #expand-row="{ record }">
        <table class="column-extend-table">
          <tr><td>内容</td><td>{{ record.content }}</td></tr>
          <tr><td>位置</td><td>{{ record.location }}</td></tr>
          <tr><td>IP</td><td>{{ record.commentIp }}</td></tr>
          <tr><td>UA</td><td>{{ record.commentUa }}</td></tr>
          <tr><td>排序</td><td>{{ record.sort }}</td></tr>
          <tr><td>评论类型</td><td><dict-label dict-key="comment_type" :value="record.commentType"/></td></tr>
          <tr><td>评论时间</td><td>{{ record.createTime }}</td></tr>
          <tr v-if="record.articleTitle"><td>文章标题</td><td>{{ record.articleTitle }}</td></tr>
          <tr v-if="record.fid"><td>一级评论ID</td><td>{{ record.fid }}</td></tr>
          <tr v-if="record.rid"><td>上级评论ID</td><td>{{ record.rid }}</td></tr>
        </table>
      </template>
    </a-table>
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
    <links-form ref="formRef" @reload="handleQuery" />
  </div>
</template>

<style scoped lang="scss">
.column-extend-table {
  tr {
    td {
      padding: 0.25rem;
      &:first-child {
        color: var(--color-text-3);
      }
      &:last-child {
        white-space: pre-wrap;
      }
    }
  }
}
</style>