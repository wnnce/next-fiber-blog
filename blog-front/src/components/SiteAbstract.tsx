import React from 'react'
import DynamicCard from '@/components/DynamicCard'

const SiteAbstract: React.FC = (): React.ReactNode => {
  return (
    <DynamicCard padding="1.5rem" title="ABSTRACT" icon="i-tabler:presentation">
      <div className="flex gap-row-2 flex-col text-sm info-text">
        <p className="flex justify-between">
          <span>已运行时间:</span>
          <span>1年1月3天</span>
        </p>
        <p className="flex justify-between">
          <span>总字数:</span>
          <span>31.1k</span>
        </p>
        <p className="flex justify-between">
          <span>总评论数:</span>
          <span>112233</span>
        </p>
        <p className="flex justify-between">
          <span>总访客数:</span>
          <span>112233</span>
        </p>
        <p className="flex justify-between">
          <span>总访问量:</span>
          <span>11223344</span>
        </p>
      </div>
    </DynamicCard>
  )
}

export default SiteAbstract;