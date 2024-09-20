import '@/styles/components/article-page.scss';
import React from 'react'
import RichImage from '@/components/RichImage'
import DynamicCard from '@/components/DynamicCard'
import ServerPagination from '@/components/ServerPagination'

const ArticlePage: React.FC<{
  params: {
    page: number
  }
}> = async ({ params }) => {
  const nums: number[] = [1, 2, 3, 4, 5, 6, 7];
  return (
    <>
      <ul className="flex flex-col list-none gap-row-4">
        {
          nums.map(num => (
            <li key={num} className="article-page-item">
              <article>
                <DynamicCard>
                  <div className="overflow-hidden" style={{ width: '100%', aspectRatio: '16 / 6' }}>
                    <RichImage className="article-item-cover" src="/b-oss/images/2024/0704/7d8e2b2e0aa20f139b8788aa38066403.webp" fill
                               style={{ width: '100%', height: '100%' }} />
                  </div>
                  <div className="p-4 flex flex-col gap-row-2">
                    <p className="desc-text text-xs">
                      POSTED <time dateTime="2023-12-12"
                                   title="发布于2023-12-11">2023-12-12</time> &nbsp;&nbsp;&nbsp;&nbsp;WORDS <strong>112233</strong>
                    </p>
                    <h2 className="sm:text-xl text-lg">
                      <a href="#" className="a-hover-line-text-md">
                        测试标题测试标题测试标题测试标题测试标测试标题
                      </a>
                    </h2>
                    <div className="flex gap-col-8 text-xs desc-text">
                      <ul className="list-none flex gap-col-2">
                        <i className="inline-block i-tabler:category text-sm" />
                        <li>
                          <a href="#" className="a-hover-line-text-sm">前端</a>
                        </li>
                        <li>
                          <a href="#" className="a-hover-line-text-sm">后端</a>
                        </li>
                      </ul>
                      <ul className="list-none flex gap-col-2">
                        <i className="inline-block i-tabler:tag text-sm" />
                        <li>
                          <a href="#" className="a-hover-line-text-sm">Vue</a>
                        </li>
                        <li>
                          <a href="#" className="a-hover-line-text-sm">React</a>
                        </li>
                      </ul>
                    </div>
                    <p className="desc-text text-sm line-clamp-2 line-height-relaxed">
                      简介测试简介测试简介测试简介测试简介测试简介测试简介测试简介测试简介测试简介测试
                      简介测试简介测试简介测试简介测试简介测试简介测试简介测试简介测试简介测试简介测试
                      简介测试简介测试简介测试简介测试简介测试简介测试简介测试简介测试简介测试简介测试
                    </p>
                    <ul className="list-none flex gap-col-4 info-text font-mono text-xs">
                      <li className="flex items-center">
                        <i className="inline-block i-tabler:eye mr-1 text-sm" />
                        112233
                      </li>
                      <li className="flex items-center">
                        <i className="inline-block i-tabler:message-chatbot mr-1 text-sm" />
                        112233
                      </li>
                      <li className="flex items-center">
                        <i className="inline-block i-tabler:thumb-up mr-1 text-sm" />
                        112233
                      </li>
                    </ul>
                  </div>
                </ DynamicCard>
              </article>
            </li>
          ))
        }
      </ul>
      <ServerPagination className="mt-4 w-full animate-on-scroll" current={4} size={5} total={20} />
    </>
  )
}

export default ArticlePage;