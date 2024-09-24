import '@/styles/components/pagination.scss';
import React, { useMemo } from 'react'
import Link from 'next/link'

interface Props {
  current: number;
  size: number;
  total: number;
  showCount?: number;
  targetUrlPrefix?: string;
  className?: string;
}

/**
 * 服务端渲染分页组件
 * @param current 当前页数
 * @param size 每页的记录数
 * @param total 总记录数
 * @param targetUrlPrefix 需要跳转的连接前缀
 * @param className 类列表
 * @constructor
 */
const ServerPagination: React.FC<Props> = ({
  current,
  size,
  total,
  targetUrlPrefix,
  className
}): React.ReactNode => {

  const pages = useMemo<number>(() => {
    return Math.ceil(total / Math.abs(size));
  }, [total, size])

  const showPageList = useMemo<number[]>(() => {
    const pageList: number[] = [];
    if (current < pages) {
      for (let i = 0; i < 3; i++) {
        const nextPage = current + i;
        if (nextPage > pages) {
          break
        }
        pageList.push(nextPage);
      }
    } else {
      for (let i = 0; i < 3; i++) {
        const nextPage = current - i;
        if (nextPage < 1) {
          break
        }
        pageList.unshift(nextPage);
      }
    }
    return pageList;
  }, [pages, current]);

  const showFirstButton = useMemo<boolean>(() => {
    return showPageList[0] > 1;
  }, [showPageList])

  const showLastButton = useMemo<boolean>(() => {
    return showPageList[showPageList.length - 1] < pages;
  }, [showPageList, pages])

  return (
    <div className={`pagination-container flex justify-between items-center text-sm ${className}`}>
      <div>
        {current > 1 ? <Link href={`${targetUrlPrefix}${current - 1}`}
                          className="inline-block page-button"
        >
          Pre
        </Link> : <span />}
      </div>
      <div className="flex gap-col-2 items-end">
        {showFirstButton && <>
          <Link href={`${targetUrlPrefix}1`} className="inline-block page-button">1</Link>
          <i className="inline-block i-tabler:line-dashed text-lg" />
        </>}
        {showPageList.map(item => (
          <Link href={`${targetUrlPrefix}${item}`} key={item}
             className={`inline-block ${item === current ? 'active-page-button' : 'page-button'}`}
          >
            {item}
          </Link>
        ))}
        {showLastButton && <>
          <i className="inline-block i-tabler:line-dashed text-lg" />
          <Link href={`${targetUrlPrefix}${pages}`} className="inline-block page-button">{pages}</Link>
        </>}
      </div>
      <div>
        {current < pages ? <Link href={`${targetUrlPrefix}${current - 1}`} className="inline-block page-button">
          Next
        </Link> : <span />}
      </div>
    </div>
  )
}

export default ServerPagination;