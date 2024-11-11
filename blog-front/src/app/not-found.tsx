import React from 'react'
import Image from 'next/image'
import Link from 'next/link'

const CustomNotFound: React.FC = () => {
  return (
    <div className="dynamic-container pb-4 sm:py-20 sm:flex justify-center items-center">
      <div className="order-2">
        <Image src="/images/404.svg" alt="404" width={480} height={480} objectFit="cover" />
      </div>
      <div className="order-1 flex flex-col gap-row-2 sm:gap-row-4 p-4 sm:p-0">
        <p className="font-bold pb-2 text-2xl md:text-3xl xl:text-5xl">404 PAGE NOT FOUND</p>
        <p className="font-bold lg:text-lg">没有找到请求的网页或资源...</p>
        <p className="text-sm">
          发生此错误的原因可能是请求的资源不存在或该资源已被删除，<br/>
          请确保请求参数正确且请求路径存在...
        </p>
        <div className="font-bold flex gap-col-4">
          <Link href="#" className="a-hover-line-text-sm text-sm " title="刷新页面">
            <i className="inline-block i-tabler:refresh mr-1 transform-translate-y-0.5" />
            刷新页面
          </Link>
          <Link href="/" className="a-hover-line-text-sm text-sm " title="返回首页">
            <i className="inline-block i-tabler:arrow-bar-left mr-1 transform-translate-y-0.5" />
            返回首页
          </Link>
        </div>
      </div>
    </div>
  )
}

export default CustomNotFound;