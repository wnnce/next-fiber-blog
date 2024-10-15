<script setup lang="ts">
import type { User } from '@/api/blog/user/types'
import { shallowRef, ref, computed } from 'vue'
import LoadImage from '@/components/LoadImage.vue'

const modalShow = ref<boolean>(false);
const show = (record: User) => {
  userInfo.value = record
  modalShow.value = true;
}

const onClose = () => {
  userInfo.value = undefined
}

const userInfo = shallowRef<User>()

const upgradeProgress = computed(() => {
  if (!userInfo.value) {
    return 0
  }
  const total = 100 << userInfo.value.level
  return userInfo.value.expertise / total;
})

defineExpose({
  show
})

</script>

<template>
  <a-modal title="站点用户详情" v-model:visible="modalShow" @close="onClose" :footer="false"
           unmount-on-close
           width="600px"
  >
    <div class="flex user-info-container" v-if="userInfo">
      <div class="info-left">
        <LoadImage :src="userInfo.avatar" radius="50%" :width="96" :height="96" />
        <h2>{{userInfo.username}}</h2>
        <div class="expertise flex info-text">
          <span>Lv {{userInfo.level}}</span>
          <a-progress :percent="upgradeProgress" :show-text="false"/>
          <span>Lv {{userInfo.level + 1}}</span>
        </div>
        <div class="labels" v-if="userInfo.labels">
          <a-tag v-for="item in userInfo.labels" color="orange" bordered :key="item">
            {{ item }}
          </a-tag>
        </div>
      </div>
      <div class="info-right">
        <div><span class="desc-text">昵称：</span><span>{{userInfo.nickname}}</span></div>
        <div><span class="desc-text">邮箱：</span><span>{{userInfo.email}}</span></div>
        <div><span class="desc-text">创建时间：</span><span>{{userInfo.createTime}}</span></div>
        <div><span class="desc-text">注册IP：</span><span>{{userInfo.registerIp}}</span></div>
        <div><span class="desc-text">注册地址：</span><span>{{userInfo.registerLocation}}</span></div>
        <div><span class="desc-text">个人站点：</span><a class="link-text" :href="userInfo.link" :title="userInfo.link" target="_blank">{{ userInfo.link }}</a></div>
        <div v-if="userInfo.summary"><span class="desc-text">个人简介：</span><span>{{userInfo.summary}}</span></div>
      </div>
    </div>
  </a-modal>
</template>

<style scoped lang="scss">
.user-info-container {
  padding: 0.25rem 1.5rem;
  column-gap: 1rem;
  > div {
    display: flex;
    flex-direction: column;
    row-gap: 0.5rem;
  }
  .info-left {
    width: 30%;
    h2 {
      font-size: 2rem;
      font-weight: bold;
    }
    .expertise {
      column-gap: 0.5rem;
      align-items: center;
      span {
        flex-shrink: 0;
      }
    }
    .labels {
      display: flex;
      column-gap: 0.5rem;
      flex-wrap: wrap;
    }
  }
  .info-right {
    font-size: 0.9rem;
  }
}
</style>