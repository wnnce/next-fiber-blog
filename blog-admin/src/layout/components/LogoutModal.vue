<script setup lang="ts">

import { ref } from 'vue'
import { userApi } from '@/api/system/user'
import { useLocalUserStore } from '@/stores/user'
import { useLocalStorage } from '@/hooks/local-storage'
import { LOCAl_USER_KEY, TOKEN_KEY } from '@/assets/script/constant'
import { useArcoMessage } from '@/hooks/message'
import { useRouter } from 'vue-router'

const router = useRouter();

const modalShow = ref<boolean>(false);
const show = () => {
  modalShow.value = true;
}

const handleLogout = async () => {
  const result = await userApi.logout();
  if (result.code !== 200) {
    return false;
  }
  useLocalUserStore().clear()
  useLocalStorage().remove(TOKEN_KEY, LOCAl_USER_KEY);
  useArcoMessage().successMessage('退出登录成功');
  setTimeout(() => {
    router.push({ path: '/login' })
  }, 200)
}

defineExpose({
  show
})
</script>

<template>
  <a-modal v-model:visible="modalShow" :on-before-ok="handleLogout" unmountOnClose
           title-align="start"
           :width="360"
  >
    <template #title>
      提示
    </template>
    <div>
      确定注销退出系统并返回登录页吗?
    </div>
  </a-modal>
</template>

<style scoped lang="scss">

</style>