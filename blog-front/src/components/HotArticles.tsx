import '@/styles/components/hot-article.scss';
import React from 'react'
import DynamicCard from '@/components/DynamicCard'
import { listHotArticle } from '@/lib/api'
import Link from 'next/link'

const HotArticles: React.FC = async () => {
  const { data: hotArticles } = await listHotArticle();
  return (
    <DynamicCard title="HOTS" icon="i-tabler:flame">
      <section className="mt-2">
        <ol className="list-none flex flex-col gap-row-4 text-sm">
          { hotArticles.map(article => (
            <li className="line-clamp-1 hot-article-item hover-line-text" key={article.articleId}>
              <Link href={`/article/${article.articleId}`} title={article.title}>
                { article.title }
                <i className="i-tabler:chevron-right text-xl arrow-icon"></i>
                <i className="text-lg line-icon"></i>
              </Link>
            </li>
          ))}
        </ol>
      </section>
    </DynamicCard>
  )
}

export default HotArticles;