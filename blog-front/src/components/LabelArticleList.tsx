import React from 'react'
import { Article } from '@/lib/types'
import RichImage from '@/components/RichImage'
import Link from 'next/link'
import { HotLabel, TopLabel } from '@/components/Labels'
import { formatDateTime } from '@/tools/utils'

const LabelArticleList: React.FC<{
  articles:Article[]
}> = ({ articles }) => {
  return (
    <ul className="list-none flex flex-col gap-row-4">
      { articles.map(article => (
        <div className="flex gap-col-4" key={article.articleId}>
          <RichImage src={article.coverUrl} width={80} height={80} fill thumbnail radius={12} />
          <div className="flex flex-col justify-between py-4">
            <p className="text-xs desc-text flex items-center">
              <i className="inline-block i-tabler:clock mr-1" />
              <time dateTime={article.createTime}>{ formatDateTime(article.createTime) }</time>
            </p>
            <h3>
              { article.isTop && <TopLabel className="mr-2" /> }
              { article.isHot && <HotLabel className="mr-2" /> }
              <Link className="a-hover-line-text-md" href="#">{article.title}</Link>
            </h3>
          </div>
        </div>
      )) }
    </ul>
  )
}

export default LabelArticleList;