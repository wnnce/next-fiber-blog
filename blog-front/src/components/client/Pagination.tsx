'use client'

import '@/styles/components/client-components.scss'
import React, { useMemo } from 'react'

const Pagination: React.FC<{
  page: number;
  pages: number;
  className?: string;
  onChange: (newPage: number) => void;
}> = ({ page, pages, className, onChange }) => {
  const preButtonIsDisabled = useMemo<boolean>(() => {
    return page === 1
  }, [page]);

  const nextButtonIsDisabled = useMemo<boolean>(() => {
    return page >= pages
  }, [page, pages])

  return (
    <div className={`flex justify-between desc-text text-xs client-pagination ${className || ''}`}>
      <button className={`page-button px-2 py-1 ${preButtonIsDisabled ? 'cursor-not-allowed' : 'active-button'}`}
              disabled={preButtonIsDisabled} onClick={() => onChange(page - 1)}
      >
        Pre
      </button>
      <span>{`第${page}页，共${pages}页`}</span>
      <button className={`page-button px-2 py-1 ${nextButtonIsDisabled ? 'cursor-not-allowed' : 'active-button'}`}
              disabled={nextButtonIsDisabled} onClick={() => onChange(page + 1)}
      >
        Next
      </button>
    </div>
  )
}

export default Pagination;