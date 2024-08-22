import { listCategory } from '@/lib/api'
import React from 'react'
import DynamicCard from '@/components/DynamicCard'
import RichImage from '@/components/RichImage'
import Link from 'next/link'

const Categories: React.FC = async () => {
  const {data: categories} = await listCategory();
  return (
    <DynamicCard padding="1.5rem" title="CATRGORIES" icon="i-tabler:category">
      <section className="text-sm">
        <ul className="list-none">
          {categories.map(category => (
            <li key={category.categoryId}>
              <Link title={category.categoryName} role="button" aria-label={`分类-${category.categoryName}`}
                 className="flex justify-between items-center my-4 pr-4 page-categories-item"
                 href={`/category/${category.categoryId}/page/1`}
              >
                <div className="flex gap-col-4 items-center">
                  <RichImage src={category.coverUrl} thumbnail width={80} height={50} radius={8} fill />
                  <span className="a-hover-line-text-sm">{category.categoryName}</span>
                </div>
                <span className="count-text info-text text-xs">{category.articleNum}</span>
              </Link>
              {(category.children && category.children.length > 0) &&
                <ul className="list-none category-children-ul pl-4">
                  {category.children.map(item => (
                    <li key={item.categoryId}>
                      <Link title={item.categoryName} role="button" aria-label={`分类-${item.categoryName}`}
                         className="flex justify-between items-center my-4 pr-4 page-categories-item"
                         href={`/category/${category.categoryId}/page/1`}
                      >
                        <div className="flex gap-col-4 items-center">
                          <RichImage src={item.coverUrl} thumbnail width={80} height={50} radius={8} />
                          <span className="a-hover-line-text-sm">{item.categoryName}</span>
                        </div>
                        <span className="count-text info-text text-xs">{item.articleNum}</span>
                      </Link>
                    </li>
                  ))}
                </ul>
              }
            </li>
          ))}
        </ul>
      </section>
    </DynamicCard>
  )
}

export default Categories;