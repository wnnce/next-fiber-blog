import '@/styles/layouts/header.scss';
import React from 'react'
import Image from 'next/image'
import Link from 'next/link'
import HeaderDrawerMenu from '@/layouts/HeaderDrawerMenu'
import Search from '@/components/client/Search'
import { querySiteConfigs } from '@/tools/site-configuration'

export interface HeaderProps {
  navList: {
    name: string,
    url: string,
    icon: React.ReactNode
  }[]
}

/**
 * 顶部导航组件
 * @param navList 导航列表
 * @constructor
 */
const Header: React.FC<HeaderProps> = async ({ navList }) => {
  const [ logoItem ] = await querySiteConfigs('logo');
  return (
    <header className="p-4 w-full text-sm header">
      <div className="flex justify-between items-center dynamic-container">
        <div className="sm:hidden">
          <HeaderDrawerMenu navList={navList} />
        </div>
        <div className="nav-div flex">
          <div className="logo">
            <Link href="/">
              <Image
                src={ logoItem ? process.env.NEXT_PUBLIC_QINIU_IMAGE_DOMAIN + logoItem.value.toString().substring(6) : '/images/logo.svg' }
                alt="logo" width="100" height="60"
              />
            </Link>
          </div>
          <nav className="hidden sm:block">
            <ul className="list-none">
              { navList.map(item => {
                return (
                  <li className="hover-line-text inline-block text-center" key={item.name}>
                    <Link className="inline-block" href={item.url}>{item.name}</Link>
                  </li>
                )
              })}
            </ul>
          </nav>
        </div>
        <Search />
      </div>
    </header>
  )
}

export default Header;