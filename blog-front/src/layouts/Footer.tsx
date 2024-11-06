import '@/styles/layouts/footer.scss';
import React from 'react';
import Image from 'next/image'
import { querySiteConfigs } from '@/tools/site-configuration'
import Link from 'next/link'

const Footer: React.FC = async () => {
  const [ title, powered, icp ] = await querySiteConfigs('title', 'powered', 'icp')
  return (
    <footer className="footer">
      <div className="dynamic-container sm:flex justify-between">
        <div className="footer-left-summary">
          <Link href="/">
            <Image src="/images/logo.svg" alt="logo" width="100" height="60" />
          </Link>
          <div className="mt-4 text-xs line-height-loose">
            <p className="flex items-center">
              <i className="i-tabler-copyright inline-block mr-1"></i >
              2022-2024 By&nbsp;<strong>{ title.value }</strong>
            </p>
            { powered.value && <p>Powered By <span dangerouslySetInnerHTML={{ __html: powered.value }} /></p> }
            <p>苟利国家生死以 | 岂因福祸避趋之</p>
            {icp.value && <p><a href="https://baidu.com" target="_blank">{icp.value}</a></p>}
          </div>
        </div>
        <div
          className="footer-right-icons flex gap-col-4 justify-end flex-wrap items-center text-7 max-w-md mt-6 sm:mt-0">
          <a className="inline-block i-tabler-creative-commons" href="https://baidu.com" target="_blank" />
          <a className="inline-block i-tabler-brand-github" href="https://github.com/wnnce" target="_blank" />
          <a className="inline-block i-tabler-brand-golang" href="https://go.dev/" target="_blank" />
          <a className="inline-block i-tabler-brand-nextjs" href="https://nextjs.org/" target="_blank" />
        </div>
      </div>
    </footer>
  )
}

export default Footer