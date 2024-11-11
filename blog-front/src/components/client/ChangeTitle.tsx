'use client'

import React, { useEffect, useRef } from 'react'
import { querySiteConfigs } from '@/tools/site-configuration'
import { traceAccess } from '@/lib/client-api'

const ChangeTitle: React.FC = () => {
  const tabTitleRef = useRef<string>('离桑的博客');
  const isTraceRef = useRef<boolean>(false);
  const queryTabTitle = async () => {
    const [ tabTitle ] = await querySiteConfigs('tabTitle');
    if (tabTitle && tabTitle.value) {
      tabTitleRef.current = tabTitle.value.toString();
      document.title = tabTitleRef.current;
    }
  }
  const trace = () => {
    if (isTraceRef.current) {
      return
    }
    isTraceRef.current = true;
    traceAccess();
  }
  const handleChangeTitle = () => {
    document.title = document.visibilityState === 'visible' ? tabTitleRef.current : '哎呀，加载失败了( •̀ ω •́ )'
  }
  useEffect(() => {
    trace();
    queryTabTitle();
    document.addEventListener('visibilitychange', handleChangeTitle);
    return () => {
      document.removeEventListener('visibilitychange', handleChangeTitle);
    }
  }, [])

  return (
    <></>
  )
}

export default ChangeTitle;