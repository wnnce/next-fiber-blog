import '@/styles/components/label-info.scss'
import React from 'react'
import StaticCard from '@/components/StaticCard'
import { queryTag } from '@/lib/api'
import RichImage from '@/components/RichImage'
import Link from 'next/link'
import Empty from '@/components/Empty'

const TagInfo: React.FC<{
  params: {
    id: string
  },
  children: React.ReactNode
}> = async ({ params, children }) => {
  const tagId = parseInt(params.id);
  if (!tagId || isNaN(tagId) || tagId <= 0) {
    return (
      <StaticCard padding="1.5rem">
        <Empty text="标签参数错误" icon="i-tabler:exclamation-circle" iconClassName="text-24 text-red-4" />
      </StaticCard>
    )
  }
  const { data: tag } = await queryTag(tagId)
  if (!tag) {
    return (
      <StaticCard padding="1.5rem">
        <Empty text="当前标签不存在" icon="i-tabler:error-404" iconClassName="text-24" />
      </StaticCard>
    )
  }
  return (
    <>
      <StaticCard>
        <div className="label-info-container">
          <RichImage src={tag.coverUrl} style={{ height: '100% !important', width: '100% !important' }}  fill />
          <div className="label-info-content flex flex-col">
            <div className="flex-1 flex justify-center items-center">
              <h1 className="main-text font-bold text-10 ">{tag.tagName}</h1>
            </div>
            <div className="label-info-summary">
            <p className="desc-text text-xs"><Link className="a-hover-line-text-sm" href="/tags">TAGS</Link> / {tag.tagName}</p>
              <div className="mt-4 desc-text text-xs flex justify-between">
                <ul className="list-none flex gap-col-4">
                  <li className="flex items-center">
                    <i className="inline-block i-tabler:eye mr-1 text-sm" />
                    {tag.viewNum}
                  </li>
                  <li className="flex items-center">
                    <i className="inline-block i-tabler:article mr-1 text-sm" />
                    0
                  </li>
                </ul>
                <p className="text-end">{tag.createTime}</p>
              </div>
            </div>
          </div>
        </div>
      </StaticCard>
      <div className="mt-4">
        {children}
      </div>
    </>
  )
}

export default TagInfo;