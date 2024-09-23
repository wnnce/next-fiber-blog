'use client'

import '@/styles/layouts/header.scss'
import '@/styles/animate.css'
import React, { useState } from 'react'
import { HeaderProps } from '@/layouts/Header'
import ReactDOM from 'react-dom'
import { CSSTransition, TransitionGroup } from 'react-transition-group'
import Image from 'next/image'
import Link from 'next/link'
import { throttle } from '@/tools/utils'

const HeaderDrawerMenu: React.FC<HeaderProps> = ({ navList }) => {
  const [ maskVisible, setMaskVisible ] = useState<boolean>(false);
  const [ menuVisible, setMenuVisible ] = useState<boolean>(false);

  window.addEventListener('resize', throttle(() => {
    const { innerWidth } = window;
    if (innerWidth >= 640) {
      setMenuVisible(false);
      setMaskVisible(false);
    }
  }))

  const portal = ReactDOM.createPortal((
    <CSSTransition in={maskVisible} timeout={300} classNames="mask-fade" unmountOnExit>
      <div className="fixed top-0 left-0 right-0 bottom-0 z-999" onClick={() => {
        setMenuVisible(false)
        setMaskVisible(false)
      }}>
        <CSSTransition in={menuVisible} timeout={500} classNames="drawer-fade" unmountOnExit>
          <div className="drawer-menu h-full p-4 w-48">
            <div className="flex justify-center py-4">
              <Image src="/images/logo.svg" alt="logo" width={120} height={60} objectFit="cover" />
            </div>
            <nav>
              <ul className="flex flex-col gap-row-4 mt-4 info-text">
                { navList.map(item => (
                  <li key={item.name}>
                    <Link className="text-sm a-hover-line-text-sm flex items-center" href={item.url}>
                      <span className="mr-2 text-lg">{ item.icon }</span>
                      { item.name }
                    </Link>
                  </li>
                ))}
              </ul>
            </nav>
          </div>

        </CSSTransition>
      </div>
    </CSSTransition>
  ), document.body)

  return (
    <div>
      <button className="i-tabler:menu-2 text-lg" onClick={() => {
        setMaskVisible(true)
        setTimeout(() => {
          setMenuVisible(true)
        }, 100)
      }} />
      { portal }
    </div>
  )
}

export default HeaderDrawerMenu;
