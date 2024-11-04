'use client'

import '@/styles/components/article-page.scss'
import '@/styles/animate.css'
import React, { ReactPortal, useEffect, useState } from 'react'
import ReactDOM from 'react-dom'
import { CSSTransition } from 'react-transition-group'
import { throttle } from '@/tools/utils'

const ArticleDrawerToc: React.FC<{
  tocHtml: string,
}> = ({ tocHtml }) => {
  const [ maskVisible, setMaskVisible ] = useState<boolean>(false);
  const [ tocVisible, setTocVisible ] = useState<boolean>(false);
  const [ portal, setPortal ] = useState<ReactPortal>();

  useEffect(() => {
    const handleResize = throttle(() => {
      const { innerWidth } = window;
      if (innerWidth >= 640) {
        setTocVisible(false);
        setMaskVisible(false);
      }
    })
    window.addEventListener('resize', handleResize);
    return () => {
      window.removeEventListener('resize', handleResize);
    }
  }, [])

  useEffect(() => {
    const portalElement = ReactDOM.createPortal((
      <CSSTransition in={maskVisible} timeout={300} classNames="mask-fade" unmountOnExit>
        <div className="fixed top-0 left-0 right-0 bottom-0 z-999 flex justify-end" onClick={() => {
          setTocVisible(false)
          setMaskVisible(false)
        }}>
          <CSSTransition in={tocVisible} timeout={500} classNames="drawer-right-fade" unmountOnExit>
            <div className="drawer-menu drawer-toc-nav-list h-full p-4 w-48" dangerouslySetInnerHTML={{ __html: tocHtml }}></div>
          </CSSTransition>
        </div>
      </CSSTransition>
    ), document.body)
    setPortal(portalElement)
    return () => {
      setPortal(undefined);
    }
  }, [maskVisible, tocVisible, tocHtml])

  return (
    <>
      <button title="文章目录" className="i-tabler:menu-2" onClick={() => {
        setMaskVisible(true)
        setTimeout(() => {
          setTocVisible(true)
        }, 100)
      }} />
      { portal }
    </>
  )
}

export default ArticleDrawerToc;
