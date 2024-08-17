import '@/styles/components/categories.scss';
import React from 'react'
import DynamicCard from '@/components/DynamicCard'

/**
 * 分类卡片组件
 * @constructor
 */
const Categories: React.FC = (): React.ReactNode => {
  return (
    <DynamicCard padding="1.5rem" title="CATEGORIES" icon="i-tabler:category" >
      <section>
        <ul className="list-none text-sm mt-2">
          <li>
            <a title="前端" role="button" aria-label="分类-前端"
               className="categories-li flex justify-between items-center"
               href="#"
            >
              <span>前端</span>
              <span className="count-text info-text text-xs">12</span>
            </a>
            <ul className="list-none category-children-ul pl-4 my-1">
              <li>
                <a title="前端" className="categories-li flex justify-between items-center" href="#">
                  <span>前端</span>
                  <span className="count-text info-text text-xs">12</span>
                </a>
              </li>
              <li>
                <a title="前端" className="categories-li flex justify-between items-center" href="#">
                  <span>前端</span>
                  <span className="count-text info-text text-xs">12</span>
                </a>
              </li>
            </ul>
          </li>
          <li>
            <a title="前端" className="categories-li flex justify-between items-center" href="#">
              <span>前端</span>
              <span className="count-text info-text text-xs">12</span>
            </a>
          </li>
          <li>
            <a title="前端" className="categories-li flex justify-between items-center" href="#">
              <span>前端</span>
              <span className="count-text info-text text-xs">12</span>
            </a>
            <ul className="list-none category-children-ul pl-4 my-1">
              <li>
                <a title="前端" className="categories-li flex justify-between items-center" href="#">
                  <span>前端</span>
                  <span className="count-text info-text text-xs">12</span>
                </a>
              </li>
              <li>
                <a title="前端" className="categories-li flex justify-between items-center" href="#">
                  <span>前端</span>
                  <span className="count-text info-text text-xs">12</span>
                </a>
              </li>
            </ul>
          </li>
        </ul>
      </section>
    </DynamicCard>
  )
}

export default Categories