import '@/styles/components/article-page.scss';
import React from 'react'
import RichImage from '@/components/RichImage'
import DynamicCard from '@/components/DynamicCard'
import ServerPagination from '@/components/ServerPagination'
import { pageArticle } from '@/lib/api'
import Empty from '@/components/Empty'
import { formatDateTime } from '@/tools/utils'
import Link from 'next/link'
import { querySiteConfigs } from '@/tools/site-configuration'

const ArticlePage: React.FC<{
  params: {
    page: string
  }
}> = async ({ params }) => {
  const numberPage = parseInt(params.page);
  if (!numberPage || isNaN(numberPage) || numberPage <= 0 ) {
    throw new Error('页码参数错误')
  }
  const [ articleSizeItem ] = await querySiteConfigs('articleSize');
  const { data: articlePage } = await pageArticle({ page: numberPage,  size: articleSizeItem ? articleSizeItem.value : 5})
  if (articlePage.records.length === 0) {
    return (
      <DynamicCard>
        <Empty textClassName="text-sm" text="还没有非置顶文章哦..." iconSize="6rem" />
      </DynamicCard>
    )
  }
  return (
    <>
      <ul className="flex flex-col list-none gap-row-4">
        {
          articlePage.records.map(article => (
            <li key={article.articleId} className="article-page-item">
              <article>
                <DynamicCard multiple={40} useDefaultPadding={false}>
                  <div className="overflow-hidden" style={{ width: '100%', aspectRatio: '16 / 6' }}>
                    <RichImage className="article-item-cover" src={article.coverUrl} fill
                               style={{ width: '100%', height: '100%' }} />
                  </div>
                  <div className="p-4 flex flex-col sm:gap-row-2 gap-row-1">
                    <p className="desc-text text-xs">
                      POSTED <time dateTime={article.createTime}>
                        {formatDateTime(article.createTime)}
                      </time> &nbsp;&nbsp;&nbsp;&nbsp;WORDS <strong>{article.wordCount}</strong>
                    </p>
                    <h2 className="sm:text-xl text-lg">
                      {article.isHot && <i className="i-tabler-flame inline-block text-xl sm:text-2xl translate-y-1 text-orange-5" title="热门" />}
                      <Link href={`/article/${article.articleId}`} className="a-hover-line-text-md">
                        {article.title}
                      </Link>
                    </h2>
                    <div className="flex gap-col-8 gap-row-2 text-xs desc-text flex-wrap">
                      { (article.categories && article.categories.length) > 0 && (
                        <ul className="list-none flex gap-col-2 flex-wrap">
                          <i className="inline-block i-tabler:category text-sm" />
                          {article.categories.map(item => (
                            <li key={item.categoryId}>
                              <Link href="#" className="a-hover-line-text-sm">{item.categoryName}</Link>
                            </li>
                          ))}
                        </ul>
                      )}
                      { (article.tags && article.tags.length > 0) && (
                        <ul className="list-none flex gap-col-2 flex-wrap">
                          <i className="inline-block i-tabler:tag text-sm" />
                          {article.tags.map(item => (
                            <li key={item.tagId}>
                              <Link href="#" className="a-hover-line-text-sm">{item.tagName}</Link>
                            </li>
                          ))}
                        </ul>
                      )}
                    </div>
                    <p className="desc-text text-sm line-clamp-2 line-height-relaxed">
                      {article.summary}
                    </p>
                    <ul className="list-none flex gap-col-4 info-text text-xs">
                      <li className="flex items-center">
                        <i className="inline-block i-tabler:eye mr-1 text-sm" />
                        { article.viewNum }
                      </li>
                      <li className="flex items-center">
                        <i className="inline-block i-tabler:message-chatbot mr-1 text-sm" />
                        { article.commentNum }
                      </li>
                      <li className="flex items-center">
                        <i className="inline-block i-tabler:thumb-up mr-1 text-sm" />
                        { article.voteUp }
                      </li>
                    </ul>
                  </div>
                </ DynamicCard>
              </article>
            </li>
          ))
        }
      </ul>
      <ServerPagination className="mt-4 w-full animate-on-scroll" current={articlePage.current} size={articlePage.size}
                        total={articlePage.total}
                        targetUrlPrefix="/page/"
      />
    </>
  )
}

export default ArticlePage;