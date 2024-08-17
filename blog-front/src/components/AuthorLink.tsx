import '@/styles/components/author-link.scss';
import React from 'react'
import DynamicCard from '@/components/DynamicCard'
import RichImage from '@/components/RichImage'

/**
 * 博客作者信息个链接组件
 * @constructor
 */
const AuthorLink: React.FC = (): React.ReactNode => {
  return (
    <DynamicCard padding="1.5rem">
      <section>
        <div className="flex justify-center">
          <RichImage className="author-avatar" src="/images/bg.webp" width={112} height={112} fill
                     radius="50%"
                     alt="author-avatar"
          />
        </div>
        <h1 className="text-center font-bold text-2xl line-height-relaxed text-wrap mt-2">CrushCola</h1>
        <p className="text-center text-sm text-wrap info-text">测试简介测试简介</p>
        <ul className="list-none flex justify-around mt-4">
          <li className="text-center">
            <span className="text-sm font-mono info-text">POSTS</span><br />
            <a href="/page/1" className="text-lg"><strong>5</strong></a>
          </li>
          <li className="text-center">
            <span className="text-sm font-mono info-text">CATEGORIES</span><br />
            <a href="/categorys" className="text-lg"><strong>5</strong></a>
          </li>
          <li className="text-center">
            <span className="text-sm font-mono info-text">TAGS</span><br />
            <a href="/tags" className="text-lg"><strong>5</strong></a>
          </li>
        </ul>
        <div className="mt-4">
          <a href="https://github.com/wnnce" target="_blank" className="main-link">
            <i className="inline-block i-tabler:brand-github mr-1" />
            Github
          </a>
          <ul className="list-none flex flex-wrap justify-center gap-row-2 mt-4 text-xl author-links-ul">
            <li className="w-20% text-center">
              <a href="https://bing.com" title="discord" target="_blank"
                 className="inline-block i-tabler:brand-discord">
                discord
              </a>
            </li>
            <li className="w-20% text-center">
              <a href="https://bing.com" title="wechat" target="_blank" className="inline-block i-tabler:brand-wechat">
                wechat
              </a>
            </li>
            <li className="w-20% text-center">
              <a href="https://bing.com" title="weibo" target="_blank" className="inline-block i-tabler:brand-weibo">
                weibo
              </a>
            </li>
            <li className="w-20% text-center">
              <a href="https://bing.com" title="twitter" target="_blank"
                 className="inline-block i-tabler:brand-twitter">
                twitter
              </a>
            </li>
            <li className="w-20% text-center">
              <a href="https://bing.com" title="telegram" target="_blank"
                 className="inline-block i-tabler:brand-telegram">
                telegram
              </a>
            </li>
          </ul>
        </div>
      </section>
    </DynamicCard>
  )
}

export default AuthorLink;