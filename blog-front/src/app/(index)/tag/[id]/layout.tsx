import '@/styles/components/label-info.scss'
import React from 'react'
import StaticCard from '@/components/StaticCard'
import { queryTag } from '@/lib/api'
import RichImage from '@/components/RichImage'
import Link from 'next/link'

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
        <div className="text-center info-text">
          <i className="inline-block text-24 text-red-4 i-tabler:exclamation-circle" />
          <p className="text-center mt-4">标签参数错误</p>
        </div>
      </StaticCard>
    )
  }
  const { data: tag } = await queryTag(tagId)
  if (!tag) {
    return (
      <StaticCard padding="1.5rem">
        <div className="text-center">
          <i className="inline-block info-text text-24 i-tabler:error-404" />
          <p className="text-center mt-4 info-text">当前标签不存在</p>
        </div>
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