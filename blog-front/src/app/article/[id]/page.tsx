'use server'

import '@/styles/components/article-page.scss'
import 'github-markdown-css/github-markdown-dark.css'
import '@/styles/components/markdown.scss'
import 'highlight.js/styles/atom-one-dark.min.css'
import React from 'react'
import { redirect } from 'next/navigation'
import { queryArticle } from '@/lib/api'
import StaticCard from '@/components/StaticCard'
import DynamicCard from '@/components/DynamicCard'
import RichImage from '@/components/RichImage'
import { formatDateTime } from '@/tools/utils'
import Link from 'next/link'
import useMarkdownParse from '@/hooks/markdown'

const ArticlePage: React.FC<{
  params: {
    id: string
  }
}> = async ({ params }) => {
  const numberId = parseInt(params.id);
  if (!numberId || isNaN(numberId) || numberId < 0) {
    redirect('/404');
  }
  const { data: article } = await queryArticle(numberId);
  if (!article) {
    redirect('/404');
  }

  const articleRender = useMarkdownParse().articleRender()

  const articleHtml = articleRender.render(article.content);

  const readerTime = function() {
    if (article.wordCount <= 400) {
      return '1分钟';
    }
    const min = Math.ceil(article.wordCount / 400);
    if (min < 60) {
      return `${min}分钟`
    }
    const hours = Math.floor(min / 60);
    const num = min & 60;
    if (num === 0) {
      return `${hours}小时`
    }
    return `${hours}小时${num}分钟`
  }();

  return (
    <div className="dynamic-container">
      <div className="flex gap-4 py-8 px-4">
        <div className="article-toc hidden lg:block lg:w-80">
          <DynamicCard padding="1.5rem">
            112233
          </DynamicCard>
        </div>
        <div className="article-content flex-1">
          <StaticCard>
            <div className="article-header h-80 xl:h-100">
              <RichImage className="header-cover" src={article.coverUrl} fill lazy />
              <div className="header-summary flex flex-col justify-end p-6 gap-row-2">
                <h1 className="text-9 main-text">{article.title}</h1>
                <ul className="list-none flex flex-wrap gap-col-4 text-xs info-text">
                  <li className="flex gap-col-2">
                    <i className="text-sm inline-block i-tabler:clock" />
                    <time dateTime={article.createTime}>{`POSTED ${formatDateTime(article.createTime)}`}</time>
                  </li>
                  <li className="flex gap-col-2">
                    <i className="text-sm inline-block i-tabler:history" />
                    <time dateTime={article.updateTime}>{`UPDATE ${formatDateTime(article.updateTime)}`}</time>
                  </li>
                </ul>
                <ul className="list-none flex flex-wrap gap-col-4 text-xs info-text">
                  <li className="flex gap-col-2">
                    <i className="text-sm inline-block i-tabler:file-word" />
                    <span>{`WORDS ${article.wordCount}`}</span>
                  </li>
                  <li className="flex gap-col-2">
                    <i className="text-sm inline-block i-tabler:hourglass-high" />
                    <time dateTime={article.updateTime}>{`READ ${readerTime}`}</time>
                  </li>
                  <li className="flex gap-col-2">
                    <i className="text-sm inline-block i-tabler:eye" />
                    <time dateTime={article.updateTime}>{`VIEW ${article.viewNum}`}</time>
                  </li>
                </ul>
                <div className="flex gap-col-4 gap-row-2 text-xs info-text flex-wrap">
                  {(article.categories && article.categories.length) > 0 && (
                    <ul className="list-none flex gap-col-2 flex-wrap">
                      <i className="inline-block i-tabler:category text-sm" />
                      {article.categories.map(item => (
                        <li key={item.categoryId}>
                          <Link href="#" className="a-hover-line-text-sm">{item.categoryName}</Link>
                        </li>
                      ))}
                    </ul>
                  )}
                  {(article.tags && article.tags.length > 0) && (
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
              </div>
            </div>
            <div className="p-4 md:p-6 markdown-body article-markdown" dangerouslySetInnerHTML={{ __html: articleHtml }}></div>
          </StaticCard>
        </div>
      </div>
    </div>
  )
}

export default ArticlePage;