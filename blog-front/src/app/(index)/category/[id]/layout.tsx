import '@/styles/components/label-info.scss'
import React from 'react'
import StaticCard from '@/components/StaticCard'
import { queryCategory } from '@/lib/api'
import RichImage from '@/components/RichImage'
import Link from 'next/link'
import { HotLabel, TopLabel } from '@/components/Labels'
import Empty from '@/components/Empty'

const CategoryInfo: React.FC<{
  params: {
    id: string
  },
  children: React.ReactNode
}> = async ({ params, children }) => {
  const categoryId = parseInt(params.id);
  if (!categoryId || isNaN(categoryId) || categoryId <= 0) {
    return (
      <StaticCard padding="1.5rem">
        <Empty text="分类参数错误" icon="i-tabler:exclamation-circle" iconClassName="text-24 text-red-4" />
      </StaticCard>
    )
  }
  const { data: category } = await queryCategory(categoryId);
  if (!category) {
    return (
      <StaticCard padding="1.5rem">
        <Empty text="当前分类不存在" icon="i-tabler:error-404" iconClassName="text-24" />
      </StaticCard>
    )
  }
  return (
    <>
      <StaticCard>
        <div className="label-info-container">
          <RichImage src={category.coverUrl} style={{ height: '100% !important', width: '100% !important' }}  fill />
          <div className="label-info-content flex flex-col">
            <div className="flex-1 flex justify-center items-center">
              <h1 className="main-text font-bold text-10 ">{category.categoryName}</h1>
            </div>
            <div className="label-info-summary">
            <p className="desc-text text-xs"><Link className="a-hover-line-text-sm" href="/categorys">CATEGORIES</Link> / {category.categoryName}</p>
              <p className="info-text text-sm mt-2">{category.description}</p>
              <div className="mt-4 desc-text text-xs flex justify-between">
                <ul className="list-none flex gap-col-2">
                  {category.isTop && <li><TopLabel /></li>}
                  {category.isHot && <li><HotLabel /></li>}
                  <li className="flex items-center">
                    <i className="inline-block i-tabler:eye mr-1 text-sm" />
                    {category.viewNum}
                  </li>
                  <li className="flex items-center">
                    <i className="inline-block i-tabler:article mr-1 text-sm" />
                    0
                  </li>
                </ul>
                <p className="text-end">{category.createTime}</p>
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

export default CategoryInfo;