<script setup lang="ts">

import { computed, onMounted, shallowReactive } from 'vue'
import { type DayStats, type IndexStats, siteApi } from '@/api/blog/site'
import HomeEcharts from '@/components/HomeEcharts.vue'

const indexStats = shallowReactive<IndexStats>({
  toDayAccess: 0,
  toDayComment: 0,
  totalAccess: 0,
  totalArticle: 0,
  totalComment: 0,
  totalTopic: 0,
  totalUser: 0,
  articleTotalView: 0,
  accessArray: [],
  articleArray: [],
  commentArray: [],
  userArray: [],
});

const queryIndexStats = async () => {
  const result = await siteApi.indexStats()
  if (result.code === 200) {
    Object.assign(indexStats, result.data)
  }
}

const accessDayStats = computed(() => {
  console.log(indexStats.accessArray)
  return formatDayStats(indexStats.accessArray);
})
const articleDayStats = computed(() => {
  return formatDayStats(indexStats.articleArray);
})
const commentDayStats = computed(() => {
  return formatDayStats(indexStats.commentArray);
})
const userDayStats = computed(() => {
  return formatDayStats(indexStats.userArray);
})
const formatDayStats = (origin: DayStats[]): DayStats[] => {
  if (origin.length  >= 7) {
    return origin
  }
  const dayStatsMap = new Map<string, number>();
  origin.forEach(item => dayStatsMap.set(item.dateItem, item.countItem));
  const firstDate = new Date();
  const dateList: DayStats[] = [];
  for (let i = 7; i >= 1; i--) {
    const previousDate = new Date(firstDate.getTime() - i * 24 * 60 * 60 * 1000);
    const year = previousDate.getFullYear();
    const month = String(previousDate.getMonth() + 1).padStart(2, '0');
    const day = String(previousDate.getDate()).padStart(2, '0');
    const dateString = `${year}-${month}-${day}`;
    dateList.push({
      dateItem: dateString,
      countItem: dayStatsMap.get(dateString) || 0,
    })
  }
  return dateList;
}

onMounted(() => {
  queryIndexStats();
})

</script>

<template>
  <div class="index-container">
    <div class="header-card-list flex">
      <div class="index-card">
        <img src="/icons/current-access.svg" loading="lazy" alt="logo">
        <div>
          <span>{{ indexStats.toDayAccess }}</span>
          <label class="info-text">今日访问</label>
        </div>
      </div>
      <div class="index-card">
        <img src="/icons/current-comment.svg" loading="lazy" alt="logo">
        <div>
          <span>{{ indexStats.toDayComment }}</span>
          <label class="info-text">新增评论数</label>
        </div>
      </div>
      <div class="index-card">
        <img src="/icons/total-access.svg" loading="lazy" alt="logo">
        <div>
          <span>{{ indexStats.totalAccess }}</span>
          <label class="info-text">总访问量</label>
        </div>
      </div>
      <div class="index-card">
        <img src="/icons/total-comment.svg" loading="lazy" alt="logo">
        <div>
          <span>{{ indexStats.totalComment }}</span>
          <label class="info-text">总评论数</label>
        </div>
      </div>
      <div class="index-card">
        <img src="/icons/total-user.svg" loading="lazy" alt="logo">
        <div>
          <span>{{ indexStats.totalUser }}</span>
          <label class="info-text">总用户数</label>
        </div>
      </div>
      <div class="index-card">
        <img src="/icons/total-article.svg" loading="lazy" alt="logo">
        <div>
          <span>{{ indexStats.totalArticle }}</span>
          <label class="info-text">总文章数</label>
        </div>
      </div>
      <div class="index-card">
        <img src="/icons/total-topic.svg" loading="lazy" alt="logo">
        <div>
          <span>{{ indexStats.totalTopic }}</span>
          <label class="info-text">总动态数</label>
        </div>
      </div>
      <div class="index-card">
        <img src="/icons/total-article-view.svg" loading="lazy" alt="logo">
        <div>
          <span>{{ indexStats.articleTotalView }}</span>
          <label class="info-text">总阅读量</label>
        </div>
      </div>
    </div>
    <div class="chart-list flex">
      <div class="index-card">
        <HomeEcharts echart-title="近7日访问量" series-name="访问量" :data="accessDayStats" />
      </div>
      <div class="index-card">
        <HomeEcharts echart-title="近7日新增评论" series-name="新增评论数" :data="commentDayStats"
                     line-color="#27AE60"
                     area-start-color="#3CB371"
                     area-end-color="#A6E3D7"
        />
      </div>
      <div class="index-card">
        <HomeEcharts echart-title="近7日文章发布字数" series-name="发布字数" :data="articleDayStats"
                     line-color="#F89B30"
                     area-start-color="#F8C630"
                     area-end-color="#F0E5A7"
        />
      </div>
      <div class="index-card">
        <HomeEcharts echart-title="近7日新增用户" series-name="用户数量" :data="userDayStats"
                     line-color="#9C27B0"
                     area-start-color="#9B59B6"
                     area-end-color="#D2B4DE"
        />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.index-container {
  .index-card {
    min-width: 20rem;
    display: flex;
    background-color: var(--card-color);
    padding: 1.5rem;
    border-radius: 0.75rem;
    align-items: center;
    gap: 1.5rem;
    > img {
      height: 72px;
      width: 72px;
    }
    > div {
      padding: 0.25rem 0;
      height: 100%;
      display: flex;
      flex-direction: column;
      justify-content: space-between;
      > span {
        font-size: 2.25rem;
      }
    }
  }
  .header-card-list {
    flex-wrap: wrap;
    gap: 1rem;
    > div {
      flex: 1;
    }
  }
  .chart-list {
    margin-top: 1rem;
    flex-wrap: wrap;
    gap: 1rem;
    > .index-card {
      width: calc(50% - 0.5rem);
      height: 320px;
    }
  }
}

@media screen and (max-width: 1250px) {
  .index-container > .chart-list > .index-card {
    width: 100%;
  }
}

</style>