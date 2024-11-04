<script setup lang="ts">

import { useLocalUserStore } from '@/stores/user'
import { reactive, ref } from 'vue'
import { useArcoMessage } from '@/hooks/message'

const userStore = useLocalUserStore();

const submitButtonLoading = ref<boolean>(false);
const formData = reactive({
  nickname: userStore.userInfo.nickname,
  email: userStore.userInfo.email,
  phone: userStore.userInfo.phone,
})

const formSubmit = async () => {
  submitButtonLoading.value = true;
  setTimeout(() => {
    useArcoMessage().successMessage('更新成功');
    submitButtonLoading.value = false;
  }, 1000)
}

</script>

<template>
  <div class="personal-container">
    <div class="card user-preview">
      <h2 class="card-title">个人信息</h2>
      <div class="text-center">
        <a-avatar :size="84" :image-url="userStore.userInfo ? userStore.userInfo.avatar : ''">User</a-avatar>
        <h2>{{ userStore.userInfo.nickname }}</h2>
        <p class="desc-text">{{ userStore.userInfo.username }}</p>
      </div>
      <p><span class="desc-text">创建时间：</span>{{ userStore.userInfo.createTime }}</p>
      <p><span class="desc-text">上次登录地址：</span>{{ userStore.userInfo.lastLoginIp }}</p>
      <p><span class="desc-text">上次登录时间：</span>{{ userStore.userInfo.lastLoginTime }}</p>
    </div>
    <div class="card user-profile">
      <h2 class="card-title">修改资料</h2>
      <a-form :model="formData" layout="vertical" auto-label-width @submit-success="formSubmit">
        <a-form-item label="昵称" field="nickname">
          <a-input v-model="formData.nickname" placeholder="请输入昵称" />
        </a-form-item>
        <a-form-item label="邮箱" field="email">
          <a-input v-model="formData.email" placeholder="请输入邮箱" />
        </a-form-item>
        <a-form-item label="联系方式" field="phone">
          <a-input v-model="formData.phone" placeholder="请输入联系方式" />
        </a-form-item>
        <a-form-item>
          <a-button html-type="submit" type="primary" :loading="submitButtonLoading">
            <template #icon><icon-save /></template>
            更新个人资料
          </a-button>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<style scoped lang="scss">
.personal-container {
  display: flex;
  column-gap: 1rem;
  .card {
    background-color: var(--card-color);
    padding: 1.5rem;
    border-radius: 0.5rem;
    .card-title {
      padding: 0.25rem 0 0.25rem 0.5rem;
      border-left: 6px solid #1E90FF;
      margin-bottom: 1rem;
    }
  }
  .user-preview {
    display: flex;
    flex-direction: column;
    row-gap: 0.75rem;
    min-width: 24rem;
    width: auto;
    flex-shrink: 0;
    > div {
      margin-bottom: 1rem;
      h2 {
        font-weight: bold;
        margin: 0.5rem 0;
      }
    }
  }
  .user-profile {
    flex: 1;
    min-width: 30rem;
  }
}
</style>