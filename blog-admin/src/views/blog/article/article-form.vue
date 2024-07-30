<script setup lang="ts">

import { onMounted, reactive, ref } from 'vue'
import { useArcoMessage } from '@/hooks/message'
import type { Result } from '@/api/request'
import type { FileItem, SelectOptionData, TreeNodeData } from '@arco-design/web-vue'
import type { FieldRule } from '@arco-design/web-vue/es/form/interface'
import ImageUpload from '@/components/ImageUpload.vue'
import DictSelect from '@/components/DictSelect.vue'
import type { Article, ArticleForm } from '@/api/blog/article/types'
import { articleApi } from '@/api/blog/article'
import type { Category } from '@/api/blog/category/types'
import { tagApi } from '@/api/blog/tags'
import { categoryApi } from '@/api/blog/category'

const emits = defineEmits<{
  (e: 'reload'): void
}>()

const { successMessage } = useArcoMessage();

const modalShow = ref<boolean>(false);
const show = (record?: Article, content?: string) => {
  if (record) {
    const { articleId, title, summary, coverUrl, categoryIds, tagIds, content, protocol, tips, password, isTop, isHot, isComment, isPrivate, sort, status } = record;
    Object.assign(formData, { articleId, title, summary, coverUrl : coverUrl || '', categoryIds, tagIds, content, protocol, tips, password, isTop, isHot, isComment, isPrivate, sort, status })
    formatAvatarToFileList();
  }
  if (!formData.content && content) {
    formData.content = content;
  }
  modalShow.value = true;
}
const onClose = () => {
  fileList.value = [];
  Object.assign(formData, defaultFormData);
}

const defaultFormData: ArticleForm = {
  articleId: undefined,
  title: undefined,
  summary: undefined,
  coverUrl: '',
  categoryIds: [],
  tagIds: [],
  protocol: undefined,
  tips: undefined,
  password: undefined,
  isTop: false,
  isHot: false,
  isComment: true,
  isPrivate: false,
  sort: 1,
  status: 0,
}
const formData = reactive<ArticleForm>({ ...defaultFormData })
const formRules: Record<string, FieldRule<any> | FieldRule<any>[]> = {
  title: { required: true, message: '文章标题不能为空' },
  summary: { required: true, message: '文章简介不能为空' },
  coverUrl: { required: true, message: '文章封面不能为空' },
  categoryIds: { required: true, message: '文章分类不能为空' },
  tagIds: { required: true, message: '文章标签不能为空' },
  password: { required: formData.isPrivate === true, message: '私密文章密码不能为空' },
  sort: { required: true, message: '显示顺序不能为空' },
  status: { required: true, message: '状态不能为空' },
}

const submitButtonLoading = ref<boolean>(false);
const formSubmit = async () => {
  submitButtonLoading.value = true;
  try {
    let result: Result<null>;
    if (formData.articleId) {
      result = await articleApi.updateArticle(formData);
    } else {
      result = await articleApi.saveArticle(formData);
    }
    if (result.code === 200) {
      successMessage(formData.articleId ? '修改成功' : '保存成功');
      emits('reload');
      modalShow.value = false;
    }
  } finally {
    submitButtonLoading.value = false;
  }
}

const fileList = ref<FileItem[]>([]);
const formatAvatarToFileList = () => {
  if (!formData.coverUrl || formData.coverUrl.trim().length === 0) {
    return;
  }
  fileList.value = [{
    uid: new Date().getTime().toString(),
    status: 'done',
    percent: 1,
    url: formData.coverUrl
  }]
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

onMounted(() => {
  queryTagSelectData();
  queryCategorySelectData();
})

defineExpose({
  show
})

</script>

<template>
  <a-modal :title="formData.articleId ? '修改文章' : '添加文章'" v-model:visible="modalShow" @close="onClose" :footer="false">
    <a-form :model="formData" auto-label-width @submit="formSubmit" :rules="formRules">
      <a-form-item label="标题" field="title">
        <a-input v-model="formData.title" placeholder="请输入文章标题" />
      </a-form-item>
      <a-form-item label="简介" field="summary">
        <a-textarea v-model="formData.summary" placeholder="请输入文章简介" />
      </a-form-item>
      <a-form-item label="封面" field="coverUrl">
        <image-upload v-model:file-list="fileList" v-model:file-url="formData.coverUrl" width="260px" height="140px" />
      </a-form-item>
      <a-form-item label="分类" field="categoryIds">
        <a-tree-select v-model="formData.categoryIds" :data="categoryTreeSelectOption" multiple placeholder="请选择文章分类" />
      </a-form-item>
      <a-form-item label="标签" field="tagIds">
        <a-select v-model="formData.tagIds" :options="tagsSelectOption" multiple placeholder="请选择文章标签" />
      </a-form-item>
      <a-form-item label="许可协议" field="protocol">
        <a-input v-model="formData.protocol" placeholder="请输入许可协议" />
      </a-form-item>
      <a-form-item label="底部提示" field="tips">
        <a-input v-model="formData.tips" placeholder="请输入文章底部提示" />
      </a-form-item>
      <div class="flex">
        <a-form-item label="置顶" field="isTop" required>
          <a-radio-group v-model="formData.isTop">
            <a-radio :value="true">是</a-radio>
            <a-radio :value="false">否</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="热门" field="isHot" required>
          <a-radio-group v-model="formData.isHot">
            <a-radio :value="true">是</a-radio>
            <a-radio :value="false">否</a-radio>
          </a-radio-group>
        </a-form-item>
      </div>
      <div class="flex">
        <a-form-item label="开启评论" field="isComment" required>
          <a-radio-group v-model="formData.isComment">
            <a-radio :value="true">是</a-radio>
            <a-radio :value="false">否</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="私密" field="isPrivate" required>
          <a-radio-group v-model="formData.isPrivate">
            <a-radio :value="true">是</a-radio>
            <a-radio :value="false">否</a-radio>
          </a-radio-group>
        </a-form-item>
      </div>
      <transition name="switch">
        <a-form-item label="文章密码" field="password" v-if="formData.isPrivate">
          <a-input-password v-model="formData.password" placeholder="请输入文章密码" />
        </a-form-item>
      </transition>
      <a-form-item label="显示顺序" field="sort">
        <a-input-number v-model="formData.sort" placeholder="请输入显示顺序" />
      </a-form-item>
      <a-form-item label="状态" field="status" required>
        <dict-select dict-key="dict_status" v-model="formData.status" type="number" />
      </a-form-item>
      <div class="flex justify-between" style="width: 100%; column-gap: 24px">
        <a-button html-type="submit" type="primary" size="large" long :loading="submitButtonLoading">
          <template #icon><icon-save /></template>
          提交
        </a-button>
        <a-button size="large" long @click="modalShow = false" :disabled="submitButtonLoading">
          <template #icon><icon-close /></template>
          取消
        </a-button>
      </div>
    </a-form>
  </a-modal>
</template>

<style scoped lang="scss">
.tree-select {
  width: 100%;
  max-height: 160px;
  overflow-y: auto;
  border: 1px solid var(--border-color);
}
</style>