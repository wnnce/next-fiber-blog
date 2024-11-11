'use client'

import React, { useEffect, useState } from 'react'
import { throttle } from '@/tools/utils'

const ArticleDynamicScrollTop: React.FC = () => {
  const [ readProgress, setReadProgress ] = useState<string>('0');

  const handleScroll = throttle(() => {
    const { scrollHeight, scrollTop, offsetHeight } = document.body;
    setReadProgress(_ => ((scrollTop + offsetHeight) / scrollHeight * 100).toFixed(0));
  })

  const handleScrollTop = () => {
    document.body.scrollTo({
      left: 0,
      top: 0,
      behavior: 'smooth'
    })
  }

  useEffect(() => {
    document.body.addEventListener('scroll', handleScroll)
    return () => {
      document.body.removeEventListener('scroll', handleScroll)
    }
  }, [])

  return (
    <>
      <style jsx>{`
        .scroll-top-button {
          color: white;
          background: linear-gradient(to top, rgb(var(--primary-color)), rgb(var(--primary-color))) no-repeat left bottom;
          background-size: 100% 0;
          padding: 0.25rem 0.25rem 0 0.25rem;
          button {
            padding: 0;
            margin: 0;
          }
        }
      `}</style>
      <div className="scroll-top-button" style={{ backgroundSize: `100% ${readProgress}%` }} >
        <button title="回到顶部" className="i-tabler:arrow-up" onClick={handleScrollTop} />
      </div>
    </>
  )
}

export default ArticleDynamicScrollTop;