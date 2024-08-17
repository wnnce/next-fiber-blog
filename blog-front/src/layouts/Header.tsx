import '@/styles/layouts/header.scss';
import React from 'react'
import Image from 'next/image'

interface Props {
  navList: {
    name: string,
    url: string
  }[]
}

/**
 * 顶部导航组件
 * @param navList 导航列表
 * @constructor
 */
const Header: React.FC<Props> = ({ navList }): React.ReactNode => {
  return (
    <header className="p-4 w-full text-sm header">
      <div className="flex justify-between dynamic-container">
        <div className="nav-div flex">
          <div className="logo">
            <a href="#">
              <Image src="/images/logo.svg" alt="logo" width="100" height="60" />
            </a>
          </div>
          <nav>
            <ul className="list-none">
              { navList.map(item => {
                return (
                  <li className="hover-line-text inline-block text-center" key={item.name}>
                    <a className="inline-block" href={item.url}>{item.name}</a>
                  </li>
                )
              })}
            </ul>
          </nav>
        </div>
        {/* TODO 使用 Search Client组件替换 */}
        <div className="header-search i-tabler-search cursor-pointer"></div>
      </div>
    </header>
  )
}

export default Header;