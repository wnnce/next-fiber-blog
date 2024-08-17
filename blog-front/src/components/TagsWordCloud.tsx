'use client'

import '@/styles/components/tags-word-cloud.scss'
import React, { useEffect } from 'react'
import DynamicCard from '@/components/DynamicCard'
import TagCloud, { TagCloudOptions } from 'TagCloud'

const createWordCloudTag = (text: string, color: string): string => {
  return `<a href="https://baidu.com" class="word-cloud-item" style="color: ${color};" target="_blank" title="${text}">${text}</a>`;
}

/**
 * 标签词云组件
 * @constructor
 */
const TagsWordCloud: React.FC = (): React.ReactNode => {
  const wordCloudOption: TagCloudOptions = {
    useContainerInlineStyles: false,
    useHTML: true,
  }
  const texts: string[] = [
    createWordCloudTag('前端', 'red'), createWordCloudTag('后端', 'yellow')
  ]
  useEffect(() => {
    const wordCloud = TagCloud('.tags-word-cloud', texts, wordCloudOption);
    return () => {
      wordCloud.destroy();
    }
  })
  return (
    <DynamicCard padding="1.5rem" title="TAGS" icon="i-tabler:tags">
      <div className="tags-word-cloud"></div>
    </DynamicCard>
  )
}

export default TagsWordCloud;