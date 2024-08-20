'use client'

import '@/styles/components/tags-word-cloud.scss'
import React, { useEffect, useState } from 'react'
import DynamicCard from '@/components/DynamicCard'
import TagCloud from 'TagCloud'
import { listTag } from '@/lib/api'

const createWordCloudTag = (text: string, color: string): string => {
  return `<a href="https://baidu.com" class="word-cloud-item" style="color: ${color};" target="_blank" title="${text}">${text}</a>`;
}

/**
 * 标签词云组件
 * @constructor
 */
const TagsWordCloud: React.FC = (): React.ReactNode => {
  const [texts, setTexts] = useState<string[]>([]);
  useEffect(() => {
    const queryData = async () => {
      const result = await listTag();
      if (result.code === 200 && result.data) {
        const texts: string[] = result.data.map(tag => createWordCloudTag(tag.tagName, tag.color))
        setTexts(texts)
      }
    }
    queryData();
  }, [])
  useEffect( () => {
    const wordCloud = TagCloud('.tags-word-cloud', texts, {
      useContainerInlineStyles: false,
      useHTML: true,
    });
    return () => {
      wordCloud.destroy();
    }
  }, [texts])
  return (
    <DynamicCard padding="1.5rem" title="TAGS" icon="i-tabler:tags">
      <div className="tags-word-cloud"></div>
    </DynamicCard>
  )
}

export default TagsWordCloud;