import '@/styles/layouts/index-layout.scss';
import React from 'react'
import AuthorLink from '@/components/AuthorLink'
import Bulletin from '@/components/Bulletin'
import Categories from '@/components/Categories'
import TagsWordCloud from '@/components/TagsWordCloud'
import HotArticles from '@/components/HotArticles'
import SiteAbstract from '@/components/SiteAbstract'

const IndexLayout: React.FC<{
  children: React.ReactNode
}> = ({ children }) => {
  return (
    <div className="dynamic-container min-h-screen">
      <div className="layout-container py-8 px-4">
        <div className="layout-left flex flex-col gap-row-4">
          <AuthorLink />
          <Bulletin />
          <Categories />
          <TagsWordCloud />
        </div>
        <div className="layout-content">
          {children}
        </div>
        <div className="layout-right flex flex-col gap-row-4">
          <HotArticles />
          <SiteAbstract />
        </div>
      </div>
    </div>
  )
}

export default IndexLayout;