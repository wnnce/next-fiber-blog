import '@/styles/components/author-link.scss';
import React from 'react'
import DynamicCard from '@/components/DynamicCard'
import RichImage from '@/components/RichImage'
import { Concat } from '@/lib/types'
import useSiteConfiguration from '@/hooks/site-configuration'
import { listConcat, querySiteStats } from '@/lib/api'

/**
 * 作者联系方式组件
 * @param concats 联系方式列表
 * @constructor
 */
const ConcatList: React.FC<{
  concats: Concat[]
}> = ({ concats }): React.ReactNode => {
  const firstMainConcat = concats.find(item => item.isMain);
  const list = !firstMainConcat ? concats : concats.filter(item => item.concatId != firstMainConcat.concatId);
  return (
    <div className="mt-4">
      { firstMainConcat && (
        <a href={firstMainConcat.targetUrl} target="_blank" className="main-link" title={firstMainConcat.name}>
          <i className="mr-1" dangerouslySetInnerHTML={{ __html: firstMainConcat.iconSvg }} />
          {firstMainConcat.name}
        </a>
      )}
      <ul className="list-none flex flex-wrap justify-center gap-row-2 mt-4 text-xl author-links-ul">
        { list.map(concat => (
          <li className="w-20% text-center" key={concat.concatId}>
            <a href={concat.targetUrl} title={concat.name} target="_blank" className="inline-block"
               dangerouslySetInnerHTML={{ __html: concat.iconSvg }} >
            </a>
          </li>
        ))}
      </ul>
    </div>
  )
}

/**
 * 博客作者信息个链接组件
 * @constructor
 */
const AuthorConcat: React.FC = async () => {
  const [
    { data: stats},
    { data: concats },
    [avatar, title, summary]
  ] = await Promise.all([
    querySiteStats(),
    listConcat(),
    useSiteConfiguration().queryConfigs('avatar', 'title', 'summary')
  ])
  return (
    <DynamicCard>
      <section>
        <div className="flex justify-center">
          <RichImage className="author-avatar" src={avatar.value.toString()} width={112} height={112} fill
                     thumbnail
                     radius="50%"
                     alt="author-avatar"
          />
        </div>
        <h1 className="text-center font-bold text-2xl line-height-relaxed text-wrap mt-2">{title.value}</h1>
        <p className="text-center text-sm text-wrap info-text">{summary.value}</p>
        <ul className="list-none flex justify-around mt-4">
          <li className="text-center">
            <span className="text-sm font-mono info-text">POSTS</span><br />
            <a href="/page/1" className="text-lg"><strong>{ stats.articleCount || 0 }</strong></a>
          </li>
          <li className="text-center">
            <span className="text-sm font-mono info-text">CATEGORIES</span><br />
            <a href="/categorys" className="text-lg"><strong>{ stats.categoryCount || 0 }</strong></a>
          </li>
          <li className="text-center">
            <span className="text-sm font-mono info-text">TAGS</span><br />
            <a href="/tags" className="text-lg"><strong>{ stats.tagCount || 0 }</strong></a>
          </li>
        </ul>
        <ConcatList concats={concats} />
      </section>
    </DynamicCard>
  )
}
export default AuthorConcat;