import '@/styles/layouts/footer.scss';
import React from 'react';
import Image from 'next/image'

const Footer: React.FC = (): React.ReactNode => {
  return (
    <footer className="footer">
      <div className="footer-main sm:flex justify-between">
        <div className="footer-left-summary">
          <a href="#">
            <Image src="/images/logo.svg" alt="logo" width="100" height="60" />
          </a>
          <div className="mt-4 text-xs line-height-loose">
            <p className="flex items-center">
              <i className="i-tabler-copyright inline-block mr-1"></i >
              2022-2024 By&nbsp;<strong>CrushCola</strong>
            </p>
            <p>
              Powered By <a target="_blank" href="https://vercel.com">Vercel</a> | <a target="_blank" href="https://nextjs.org">NextJS</a> | <a target="_blank" href="https://go.dev">Golang</a> | <a target="_blank" href="https://gofiber.org">GoFiber</a>
            </p>
            <p>苟利国家生死以 | 岂因福祸避趋之</p>
            <p><a href="https://baidu.com" target="_blank">渝ICP备20001411号</a></p>
          </div>
        </div>
        <div className="footer-right-icons flex gap-col-4 justify-end flex-wrap items-center text-7 max-w-md mt-6 sm:mt-0">
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