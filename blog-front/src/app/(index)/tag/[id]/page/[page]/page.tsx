import React from 'react'
import DynamicCard from '@/components/DynamicCard'
import { pageLabelArticle } from '@/lib/api'
import LabelArticleList from '@/components/LabelArticleList'
import ServerPagination from '@/components/ServerPagination'

const TagArticlePage: React.FC<{
  params: {
    id: string
    page: string
  }
}> = async ({ params }) => {
  const { id, page } = params;
  const articlePage = parseInt(page);
  const tagId = parseInt(id);
  if (!tagId || isNaN(tagId) || tagId <= 0) {
    return ;
  }
  if (!articlePage || isNaN(articlePage) || articlePage < 1) {
    return (
      <DynamicCard padding="1.5rem" title="POSTS" icon="i-tabler:news">
        <div className="text-center info-text">
          <i className="text-center inline-block text-red-4 text-24 i-tabler:exclamation-circle" />
          <p className="text-center mt-4">文章页码不存在</p>
        </div>
      </DynamicCard>
    )
  }
  const { data: pageData } = await pageLabelArticle({ tagId, page: articlePage, size: 20 });
  if ( !pageData || pageData.records.length === 0) {
    return (
      <DynamicCard padding="1.5rem" title="POSTS" icon="i-tabler:news">
        <div className="text-center desc-text">
          <i className="text-center inline-block text-24 i-tabler:template" />
          <p className="text-center mt-4">还没有关联文章哦...</p>
        </div>
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
                        targetUrlPrefix={`/tag/${tagId}/page/`}
                        className="mt-4 animate-on-scroll"
      />
    </>

  )
}

export default TagArticlePage;