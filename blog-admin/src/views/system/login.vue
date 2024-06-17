<script setup lang="ts">
import { useRoute } from 'vue-router'
import { onMounted, reactive, ref } from 'vue'
import type { LoginForm } from '@/api/system/types'
const { redirect } = useRoute().query;

const isSavePassword = ref<boolean>(false);

const defaultFormData: LoginForm = {
  username: undefined,
  password: undefined,
  code: undefined
}
const formData = reactive<LoginForm>({ ...defaultFormData })
const loginButtonLoading = ref<boolean>(false);
const handleSubmit = async () => {
  
}
</script>

<template>
  <div class="login-container flex">
    <div class="left-background"></div>
    <div class="right-div flex justify-center item-center">
      <div class="login-card radius-xl">
        <h1 class="title">欢迎访问博客后台管理系统</h1>
        <a-form :model="formData" layout="vertical" @submit="handleSubmit">
          <a-form-item field="username" label="用户名" hide-label :rules="[{required: true, message: '用户名不能为空'}]">
            <a-input v-model="formData.username" size="large" placeholder="请输入用户名">
              <template #prefix>
                <icon-user />
              </template>
            </a-input>
          </a-form-item>
          <a-form-item field="password" label="密码" hide-label :rules="[{required: true, message: '密码不能为空'}]">
            <a-input-password v-model="formData.password" size="large" placeholder="请输入密码">
              <template #prefix>
                <icon-unlock />
              </template>
            </a-input-password>
          </a-form-item>
          <a-form-item field="code" label="验证码" hide-label :rules="[{required: true, message: '验证码不能为空'}]">
            <div class="code-input flex">
              <a-input v-model="formData.code" size="large" placeholder="请输入验证码">
                <template #prefix>
                  <icon-safe />
                </template>
              </a-input>
              <div class="verify-code">验证码</div>
            </div>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" html-type="submit" size="large" long :loading="loginButtonLoading">登录</a-button>
          </a-form-item>
        </a-form>
        <div class="save-password">
          <input id="save-password" type="checkbox" v-model="isSavePassword" />
          <label for="save-password">
            保存密码
          </label>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.login-container {
  height: 100vh;
  width: 100%;
  background-color: var(--background-color);
  .left-background {
    width: 30%;
    flex-shrink: 1;
    background-image: url("/images/login-bg.png");
    background-repeat: no-repeat;
    background-size: 80%;
    background-position: center;
    background-color: #87CEFA;
  }
  .right-div {
    flex: 1;
    .login-card {
      row-gap: var(--space-sm);
      padding: var(--space-xxl);
      box-shadow: 0 0 24px var(--shadom-color);
      background-color: white;
      .title {
        font-weight: 600;
        padding: var(--space-sm) 0 var(--space-xl) 0;
      }
      .code-input {
        column-gap: var(--space-md);
        .verify-code {
          flex-shrink: 1;
          width: 100px;
          background-color: #1E90FF;
        }
      }
      .save-password {
        font-size: 14px;
      }
    }
  }
}
</style>