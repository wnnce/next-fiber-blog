import React from 'react'
import DynamicCard from '@/components/DynamicCard'
import { pageLabelArticle } from '@/lib/api'
import LabelArticleList from '@/components/LabelArticleList'
import ServerPagination from '@/components/ServerPagination'
import Empty from '@/components/Empty'

const CategoryArticlePage: React.FC<{
  params: {
    id: string
    page: string
  }
}> = async ({ params }) => {
  const { id, page } = params;
  const articlePage = parseInt(page);
  const categoryId = parseInt(id);
  if (!categoryId || isNaN(categoryId) || categoryId <= 0) {
    return ;
  }
  if (!articlePage || isNaN(articlePage) || articlePage < 1) {
    return (
      <DynamicCard padding="1.5rem" title="POSTS" icon="i-tabler:news">
        <Empty text="文章页码不存在" icon="i-tabler:exclamation-circle" iconClassName="text-24 text-red-4" />
      </DynamicCard>
    )
  }
  const { data: pageData } = await pageLabelArticle({ categoryId, page: articlePage, size: 20 });
  if ( !pageData || pageData.records.length === 0) {
    return (
      <DynamicCard padding="1.5rem" title="POSTS" icon="i-tabler:news">
        <Empty text="还没有关联文章哦..." icon="i-tabler:template" iconClassName="text-24 desc-text" textClassName="desc-text" />
      </DynamicCard>
    )
  }
  return (
    <>
      <DynamicCard padding="1.5rem" title="POSTS" icon="i-tabler:news">
        <LabelArticleList articles={pageData.records} />
      </DynamicCard>
      <ServerPagination current={pageData.current}
                        size={20}
                        total={pageData.total}
                        targetUrlPrefix={`/category/${categoryId}/page/`}
                        className="mt-4 animate-on-scroll"
      />
    </>

  )
}

export default CategoryArticlePage;