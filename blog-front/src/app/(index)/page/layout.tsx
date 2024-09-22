import '@/styles/layouts/article-page-layout.scss';
import React from 'react'
import { Swiper, SwiperSlide } from '@/components/Swiper'
import RichImage from '@/components/RichImage'
import { listTopArticle } from '@/lib/api'
import Link from 'next/link'

const ArticlePageLayout: React.FC<{
  children: React.ReactNode;
}> = async ({ children }) => {
  const result = await listTopArticle();
  if (!result || result.code !== 200 || !result.data || result.data.length === 0) {
    return (
     <>{ children }</>
    )
  }
  return (
    <section>
      <div className="top-article-container animate-on-scroll">
        <Swiper draggable={false} dotActive>
          { result.data.map(article => (
            <SwiperSlide key={article.articleId}>
              <div className="top-article-item">
                <RichImage className="top-article-item-image" src={article.coverUrl} fill />
                <div className="top-article-summary p-4 flex flex-col sm:gap-row-2 gap-row-1">
                  <h2 className="sm:text-xl text-lg">
                    <Link href="#" className="a-hover-line-text-md"
                          title={article.title}
                    >
                      { article.title }
                    </Link>
                  </h2>
                  <p className="line-clamp-1 sm:text-sm text-xs desc-text">{ article.summary }</p>
                </div>
              </div>
            </SwiperSlide>
          )) }
        </Swiper>
      </div>
      <div className="mt-4">
        {children}
      </div>
    </section>
  )
}

export default ArticlePageLayout;