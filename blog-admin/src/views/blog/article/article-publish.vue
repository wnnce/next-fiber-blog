<script setup lang="ts">
import Editor from '@/components/Editor.vue'
import { onMounted, ref } from 'vue'
import { useLocalStorage } from '@/hooks/local-storage'
import { useArcoMessage } from '@/hooks/message'
import ArticleForm from '@/views/blog/article/article-form.vue'

const { successMessage, errorMessage } = useArcoMessage();

const editorRef = ref();

const contentDraftCacheKey = "article:content:draft";
const draftModelShow = ref<boolean>(false);
const handleSaveDraft = () => {
  useLocalStorage().set<string>(contentDraftCacheKey, editorRef.value.getMarkdownValue(), undefined)
  successMessage('保存草稿成功');
}
const handleDeleteDraft = () => {
  useLocalStorage().remove(contentDraftCacheKey);
  successMessage('删除成功')
  draftModelShow.value = false;
}
const handleReaderDraft = () => {
  const draftContent = useLocalStorage().get<string>(contentDraftCacheKey);
  editorRef.value.setValue(draftContent);
  draftModelShow.value = false;
}

const formRef = ref();
const handleNext = () => {
  const articleContent = editorRef.value.getMarkdownValue();
  if (!articleContent || articleContent.trim().length === 0) {
    errorMessage('文章内容不能为空');
    return;
  }
  formRef.value.show(undefined, articleContent);
}

onMounted(() => {
  setTimeout(() => {
    const draftContent = useLocalStorage().get<string>(contentDraftCacheKey);
    if (draftContent) {
      draftModelShow.value = true;
    }
  }, 500)
})
</script>

<template>
  <div class="table-card">
    <div class="publish-option flex justify-between">
      <a-popconfirm content="是否确认清空文章内容？" type="error" position="bl"
                    :ok-button-props="{ status: 'danger' }"
                    @ok="editorRef.clear()"
      >
        <a-button status="danger">
          <template #icon><icon-eraser /></template>清空
        </a-button>
      </a-popconfirm>

      <div class="flex" style="column-gap: 16px">
        <a-popconfirm content="保存草稿后，可以在下一次进入该页面时继续编辑" type="info" position="lt"
                      @ok="handleSaveDraft"
        >
          <a-button>
            <template #icon><icon-history /></template>保存草稿
          </a-button>
        </a-popconfirm>
        <a-button type="primary" @click="handleNext">
          <template #icon><icon-right/></template>下一步
        </a-button>
      </div>
    </div>
    <Editor ref="editorRef" :min-height="700" />
    <article-form ref="formRef" @reload="editorRef.clear()" />
    <a-modal title="提示" message-type="info" width="300px" v-model:visible="draftModelShow">
      读取到还有未编辑完成的草稿，是否继续编辑？
      <template #footer>
        <a-button type="primary" size="small" @click="handleReaderDraft">确定</a-button>
        <a-button status="danger" size="small" @click="handleDeleteDraft">删除草稿</a-button>
        <a-button size="small" @click="draftModelShow = false">取消</a-button>
      </template>
    </a-modal>
  </div>
</template>

<style scoped lang="scss">

</style>