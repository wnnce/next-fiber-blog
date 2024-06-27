<script setup lang="ts">
import { useLocalUserStore } from '@/stores/user'
import { useRoute, useRouter } from 'vue-router'
import { onMounted, ref } from 'vue'

const route = useRoute();
const router = useRouter();

const activeName = ref<string>('');

router.beforeEach((to, form, next) => {
  const toName = to.name ? to.name.toString() : '';
  if (toName === activeName.value || toName === '-1') {
    next();
    return;
  }
  const toIsCache = to.meta ? to.meta.keepAlive as boolean : false;
  const keepaliveComponent = useLocalUserStore().queryKeepaliveComponent(toName);
  if (!keepaliveComponent) {
    useLocalUserStore().addKeepaliveComponent(toName, toIsCache);
  }
  activeName.value = toName;
  const formIsCache = form.meta ? form.meta.keepAlive as boolean : false;
  const formIsVisible = form.meta ? form.meta.isVisible as boolean : false;
  if (!formIsCache || !formIsVisible) {
    const formName = form.name ? form.name.toString() : '';
    useLocalUserStore().removeKeepaliveComponent(formName);
  }
  next();
})

const handleTabChange = (key: string | number) => {
  const keepaliveComponent = useLocalUserStore().queryKeepaliveComponent(key.toString());
  router.push({ path: keepaliveComponent?.path })
}
const handleTabDelete = (key: string | number) => {
  useLocalUserStore().removeKeepaliveComponent(key.toString())
  const length = useLocalUserStore().keepaliveList.length;
  if (length === 0) {
    router.push({ path: '/index' })
    return;
  }
  const lastPageTab = useLocalUserStore().keepaliveList[length - 1];
  activeName.value = lastPageTab.menuId;
  router.push({ path: lastPageTab.path })
}

const handleMounted = async () => {
  const routeName = route.name ? route.name.toString() : undefined;
  if (!routeName || routeName === '-1') {
    return;
  }
  const routeItem = useLocalUserStore().menuRouteMap.get(routeName);
  if (!routeItem) {
    return;
  }
  const toIsCache = routeItem.meta ? routeItem.meta.keepAlive as boolean : false;
  const keepaliveComponent = useLocalUserStore().queryKeepaliveComponent(routeName);
  if (!keepaliveComponent) {
    useLocalUserStore().addKeepaliveComponent(routeName, toIsCache);
  }
  activeName.value = routeName;
}

onMounted(() => {
  handleMounted()
})

</script>

<template>
  <a-tabs editable
          type="card"
          hide-content
          :active-key="activeName"
          animation
          @change="handleTabChange"
          @delete="handleTabDelete"
  >
    <a-tab-pane v-for="item in useLocalUserStore().keepaliveList" :key="item.menuId" :title="item.name" />
  </a-tabs>
</template>

<style scoped lang="scss">

</style>