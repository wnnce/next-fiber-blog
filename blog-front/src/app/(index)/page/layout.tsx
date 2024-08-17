import '@/styles/layouts/article-page-layout.scss';
import React from 'react'
import { Swiper, SwiperSlide } from '@/components/Swiper'
import RichImage from '@/components/RichImage'
import { ArticleTopLabel } from '@/components/ArticleLabels'

const ArticlePageLayout: React.FC<{
  children: React.ReactNode
}> = ({ children }) => {
  return (
    <section>
      <div className="top-article-container animate-on-scroll">
        <Swiper draggable={false} dotActive>
          <SwiperSlide>
            <div className="top-article-item">
              <RichImage className="top-article-item-image" src="/images/bg.webp" fill />
              <div className="top-article-summary p-4 flex flex-col sm:gap-row-2 gap-row-1">
                <h2 className="sm:text-xl text-lg">
                  <a href="#" className="a-hover-line-text-md">
                    测试标题测试标题测试标题测试标题测试标
                  </a>
                </h2>
                <p className="line-clamp-2 sm:text-sm text-xs desc-text">
                  测试简介测试简介测试简介测试简介测试简介测试简介测试简介测试简介测试简介测试简介试简介测试简介
                  试简介测试简介试简介测试简介试简介测试简介试简介测试简介试简介测试简介试简介测试简介试简介测试简介
                </p>
              </div>
            </div>
          </SwiperSlide>
          <SwiperSlide>
            <div className="top-article-item">
              <RichImage className="top-article-item-image" src="/images/bg.webp" fill />
              <div className="top-article-summary p-4 flex flex-col sm:gap-row-2 gap-row-1">
                <h2 className="sm:text-xl text-lg">
                  <a href="#">
                    测试标题测试标题测试标题测试标题测试标
                  </a>
                </h2>
                <p className="line-clamp-2 sm:text-sm text-xs desc-text">
                  测试简介测试简介测试简介测试简介测试简介测试简介测试简介测试简介测试简介测试简介试简介测试简介
                  试简介测试简介试简介测试简介试简介测试简介试简介测试简介试简介测试简介试简介测试简介试简介测试简介
                </p>
              </div>
            </div>
          </SwiperSlide>
          <SwiperSlide>
            <div className="top-article-item">
              <RichImage className="top-article-item-image" src="/images/bg.webp" fill />
              <div className="top-article-summary p-4 flex flex-col sm:gap-row-2 gap-row-1">
                <h2 className="sm:text-xl text-lg">
                  <a href="#">
                    测试标题测试标题测试标题测试标题测试标
                  </a>
                </h2>
                <p className="line-clamp-2 sm:text-sm text-xs desc-text">
                  测试简介测试简介测试简介测试简介测试简介测试简介测试简介测试简介测试简介测试简介试简介测试简介
                  试简介测试简介试简介测试简介试简介测试简介试简介测试简介试简介测试简介试简介测试简介试简介测试简介
                </p>
              </div>
            </div>
          </SwiperSlide>
        </Swiper>
      </div>
      <div className="mt-4">
        {children}
      </div>
    </section>
  )
}

export default ArticlePageLayout;