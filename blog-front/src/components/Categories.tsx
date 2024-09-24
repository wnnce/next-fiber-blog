import '@/styles/components/categories.scss';
import React from 'react'
import DynamicCard from '@/components/DynamicCard'
import { listCategory } from '@/lib/api'
import Link from 'next/link'

/**
 * 分类卡片组件
 * @constructor
 */
const Categories: React.FC = async () => {
  const { data: categories } = await listCategory();
  return (
    <DynamicCard padding="1.5rem" title="CATEGORIES" icon="i-tabler:category" >
      <section>
        <ul className="list-none text-sm mt-2">
          { categories.map(category => (
            <li key={category.categoryId}>
              <Link title={category.categoryName} role="button" aria-label={`分类-${category.categoryName}`}
                 className="categories-li flex justify-between items-center"
                 href={`/category/${category.categoryId}/page/1`}
              >
                <span>{category.categoryName}</span>
                <span className="count-text info-text text-xs">{category.articleNum}</span>
              </Link>
              {(category.children && category.children.length > 0) &&
                <ul className="list-none category-children-ul pl-4 my-1">
                  { category.children.map(item => (
                    <li key={item.categoryId}>
                      <Link title={category.categoryName} role="button" aria-label={`分类-${category.categoryName}`}
                         className="categories-li flex justify-between items-center"
                         href={`/category/${category.categoryId}/page/1`}
                      >
                        <span>{item.categoryName}</span>
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

export default Categories