'use client'

import '@/styles/components/client-components.scss'
import '@/styles/animate.css'
import React, { ReactPortal, useEffect, useState } from 'react'
import { CSSTransition } from 'react-transition-group'
import ReactDOM from 'react-dom'

const Search: React.FC = () => {
  const [ maskVisible, setMaskVisible ] = useState<boolean>(false);
  const [ bodyVisible, setBodyVisible ] = useState<boolean>(false);
  const [ portal, setPortal ] = useState<ReactPortal>();

  useEffect(() => {
    const portalElement = ReactDOM.createPortal((
      <SearchBody maskVisible={maskVisible} bodyVisible={bodyVisible} onClose={() => {
        setBodyVisible(false)
        setMaskVisible(false)
      }} />
    ), document.body)
    setPortal(portalElement)
    return () => {
      setPortal(undefined)
    }
  }, [maskVisible, bodyVisible])

  return (
    <div className="header-search">
      <button className="i-tabler-search cursor-pointer" onClick={() => {
        setMaskVisible(true);
        setTimeout(() => {
          setBodyVisible(true);
        }, 100)
      }} />
      { portal }
    </div>
  )
}

const SearchBody: React.FC<{
  maskVisible: boolean;
  bodyVisible: boolean;
  onClose: () => void;
}> = ({ maskVisible, bodyVisible, onClose }) => {
  return (
    <CSSTransition timeout={300} in={maskVisible} classNames="mask-fade" unmountOnExit>
      <div className="fixed top-0 left-0 right-0 bottom-0 z-999 flex justify-center p-8" onClick={onClose}>
        <CSSTransition timeout={500} in={bodyVisible} unmountOnExit classNames="search-fade">
          <div className="search-body rounded-md search-fade h-min max-w-160 w-full p-4 md:p-6 bg-white mt-12 flex flex-col gap-row-4"
               onClick={() => {}}
          >
            <p className="md:text-xl search-title">文章搜索</p>
            <input className="search-input text-sm" placeholder="请输入待搜索的文章标题或内容..." />
          </div>
        </CSSTransition>
      </div>
    </CSSTransition>
  )
}

export default Search;