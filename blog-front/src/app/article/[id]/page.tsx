'use server'

import '@/styles/components/article-page.scss'
import 'github-markdown-css/github-markdown-dark.css'
import '@/styles/components/markdown.scss'
import 'highlight.js/styles/atom-one-dark.min.css'
import cardStyles from '@/styles/components/card.module.scss';
import React from 'react'
import { redirect } from 'next/navigation'
import { queryArticle } from '@/lib/api'
import RichImage from '@/components/RichImage'
import { formatDateTime } from '@/tools/utils'
import Link from 'next/link'
import { getArticleRender } from '@/tools/markdown'
import { ArticleLike } from '@/components/client/CommonLike'
import Comment from '@/components/comment/Comment'
import StaticCard from '@/components/StaticCard'
import ArticleDynamicScrollTop from '@/components/client/ArticleDynamicScrollTop'
import ArticleDrawerToc from '@/components/client/ArticleDrawerToc'

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

  let articleTocHtml: string = '';
  const articleRender = getArticleRender((html: string) => {
    articleTocHtml = html;
  })
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
    const num = min % 60;
    if (num === 0) {
      return `${hours}小时`
    }
    return `${hours}小时${num}分钟`
  }();

  return (
    <div className="dynamic-container">
      <div className="sm:py-8 py-4 sm:px-4 block lg:gap-4 lg:flex">
        <div className="article-slide-option text-xl flex flex-col gap-row-1">
          <div className="toc-button lg:hidden">
            <ArticleDrawerToc tocHtml={articleTocHtml} />
          </div>
          <ArticleDynamicScrollTop />
        </div>
        <div className="flex-1 flex gap-row-4 flex-col">
          <div className={`${cardStyles.card} article-content`}>
            <div className="article-header h-80 xl:h-100">
              <RichImage className="header-cover" src={article.coverUrl} fill lazy />
              <div className="header-summary flex flex-col justify-end p-6 gap-row-2">
                <h1 className="text-9 main-text">{article.title}</h1>
                <ul className="list-none flex flex-wrap gap-col-4 text-xs main-text">
                  <li className="flex gap-col-2">
                    <i className="text-sm inline-block i-tabler:clock" />
                    <time dateTime={article.createTime}>{`POSTED ${formatDateTime(article.createTime)}`}</time>
                  </li>
                  <li className="flex gap-col-2">
                    <i className="text-sm inline-block i-tabler:history" />
                    <time dateTime={article.updateTime}>{`UPDATE ${formatDateTime(article.updateTime)}`}</time>
                  </li>
                </ul>
                <ul className="list-none flex flex-wrap gap-col-4 text-xs main-text">
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
                <div className="flex gap-col-4 gap-row-2 text-xs main-text flex-wrap">
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
            <div className="p-4 md:p-6 markdown-body article-markdown"
                 dangerouslySetInnerHTML={{ __html: articleHtml }}></div>
            <div className="article-footer sm:gap-row-3 gap-row-1.5 sm:p-2 p-1 flex flex-col desc-text">
              <ul className="flex footer-options py-2 list-none">
                <li><ArticleLike articleId={article.articleId} count={article.voteUp} /></li>
                <li><i className="inline-block i-tabler:qrcode"></i></li>
              </ul>
              <p className="sm:text-sm text-3">{ `本文使用 ${article.protocol || 'CC BY-NC-SA 4.0'} 许可协议，转载请注明出处` }</p>
              { article.tips && <p className="sm:text-sm text-3">{ `TIPS: ${article.tips}` }</p> }
            </div>
          </div>
          <div className="related-article sm:px-0 px-2">
            <h2 className="mb-4">关联文章</h2>
            <div className="related-article-list flex gap-col-4">
              <Link href={`/article/${article.articleId}`}>
                <div className="related-article-item">
                  <RichImage className="item-cover" src={article.coverUrl} fill lazy />
                  <div className="item-content">
                    <h3>{ article.title }</h3>
                  </div>
                </div>
              </Link>
            </div>
          </div>
          <StaticCard useDefaultPadding>
            <Comment type={1} articleId={article.articleId} />
          </StaticCard>
        </div>
        <div className="article-toc animate-on-scroll hidden lg:block lg:w-80">
          <div className="toc-nav-list" dangerouslySetInnerHTML={{ __html: articleTocHtml }} />
        </div>
      </div>
    </div>
  )
}

export default ArticlePage;