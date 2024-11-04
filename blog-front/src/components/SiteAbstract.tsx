import React from 'react'
import DynamicCard from '@/components/DynamicCard'
import { querySiteStats } from '@/lib/api'
import { formatWordCount } from '@/tools/utils'

const SiteAbstract: React.FC = async () => {
  const { data: stats } = await querySiteStats();
  return (
    <DynamicCard title="ABSTRACT" icon="i-tabler:presentation">
      <div className="flex gap-row-2 flex-col text-sm info-text">
        <p className="flex justify-between">
          <span>已运行时间:</span>
          <span>1年1月3天</span>
        </p>
        <p className="flex justify-between">
          <span>总字数:</span>
          <span>{ formatWordCount(stats.wordTotal || 0) }</span>
        </p>
        <p className="flex justify-between">
          <span>总评论数:</span>
          <span>{ stats.commentCount || 0 }</span>
        </p>
        <p className="flex justify-between">
          <span>总访客数:</span>
          <span>{ stats.visitorCount || 0 }</span>
        </p>
        <p className="flex justify-between">
          <span>总访问量:</span>
          <span>{ stats.accessCount || 0 }</span>
        </p>
      </div>
    </DynamicCard>
  )
}

export default SiteAbstract;