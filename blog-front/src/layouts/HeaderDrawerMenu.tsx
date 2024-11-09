'use client'

import '@/styles/layouts/header.scss'
import '@/styles/animate.css'
import React, { ReactPortal, useCallback, useEffect, useState } from 'react'
import { HeaderProps } from '@/layouts/Header'
import ReactDOM from 'react-dom'
import { CSSTransition } from 'react-transition-group'
import Image from 'next/image'
import Link from 'next/link'
import { throttle } from '@/tools/utils'
import { querySiteConfigs } from '@/tools/site-configuration'

const DrawerMenuLogo: React.FC = () => {
  const [ logoUrl, setLogoUrl ] = useState<string>('/images/logo.svg')

  const queryLogoUrl = useCallback(async () => {
    const [logoItem] = await querySiteConfigs('logo')
    if (logoItem && logoItem.value) {
      setLogoUrl(process.env.NEXT_PUBLIC_QINIU_IMAGE_DOMAIN + logoItem.value.toString().substring(6))
    }
  }, [])

  useEffect(() => {
    queryLogoUrl();
  }, [queryLogoUrl])

  return (
    <Image src={logoUrl} alt="logo" width={120} height={60} objectFit="cover" />
  )
}

const HeaderDrawerMenu: React.FC<HeaderProps> = ({ navList }) => {
  const [ maskVisible, setMaskVisible ] = useState<boolean>(false);
  const [ menuVisible, setMenuVisible ] = useState<boolean>(false);
  const [ portal, setPortal ] = useState<ReactPortal>();

  useEffect(() => {
    const handleResize = throttle(() => {
      const { innerWidth } = window;
      if (innerWidth >= 640) {
        setMenuVisible(false);
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
        <div className="fixed top-0 left-0 right-0 bottom-0 z-999" onClick={() => {
          setMenuVisible(false)
          setMaskVisible(false)
        }}>
          <CSSTransition in={menuVisible} timeout={500} classNames="drawer-fade" unmountOnExit>
            <div className="drawer-menu h-full p-4 w-48">
              <div className="flex justify-center py-4">
                <DrawerMenuLogo />
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
    setPortal(portalElement)
    return () => {
      setPortal(undefined);
    }
  }, [maskVisible, menuVisible, navList])

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
