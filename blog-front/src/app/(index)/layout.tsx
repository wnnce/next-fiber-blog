import '@/styles/layouts/index-layout.scss';
import React from 'react'
import AuthorConcat from '@/components/AuthorConcat'
import Bulletin from '@/components/Bulletin'
import Categories from '@/components/Categories'
import HotArticles from '@/components/HotArticles'
import SiteAbstract from '@/components/SiteAbstract'
import TagsWordCloud from '@/components/TagsWordCloud'
import Archives from '@/components/Archives'

const IndexLayout: React.FC<{
  children: React.ReactNode
}> = ({ children }) => {
  return (
    <div className="dynamic-container min-h-screen">
      <div className="layout-container py-4 px-2 sm:py-8 sm:px-4">
        <div className="layout-left flex flex-col gap-row-4">
          <AuthorConcat />
          <Bulletin />
          <Categories />
          <TagsWordCloud />
        </div>
        <div className="layout-content">
          {children}
        </div>
        <div className="layout-right flex flex-col gap-row-4">
          <HotArticles />
          <Archives />
          <SiteAbstract />
        </div>
      </div>
    </div>
  )
}

export default IndexLayout;