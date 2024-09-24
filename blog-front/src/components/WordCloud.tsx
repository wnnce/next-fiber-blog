'use client'

import '@/styles/components/tags-word-cloud.scss'
import React, { useEffect } from 'react'
import DynamicCard from '@/components/DynamicCard'
import TagCloud from 'TagCloud'

/**
 * 标签词云组件
 * @constructor
 */
export const WordCloud: React.FC<{
  children: React.ReactNode
}> = ({ children }): React.ReactNode => {
  useEffect( () => {
    const wordContainer = document.querySelector<HTMLDivElement>('.hidden-word-cloud-list');
    if (!wordContainer) {
      return ;
    }
    const wordList = wordContainer.querySelectorAll<HTMLSpanElement>('.temp-word-cloud-li-item');
    const texts: string[] = []
    wordList.forEach(item => {
      const firstElement = item.firstElementChild
      firstElement && firstElement.classList.add('word-cloud-item');
      texts.push(item.innerHTML);
    })
    const wordCloud = TagCloud('.tags-word-cloud', texts, {
      useContainerInlineStyles: false,
      useHTML: true,
    });
    return () => {
      wordCloud.destroy();
    }
  }, [children])
  return (
    <DynamicCard padding="1.5rem" title="TAGS" icon="i-tabler:tags">
      <div className="tags-word-cloud"></div>
      <ul className="hidden-word-cloud-list hidden">
        { children }
      </ul>
    </DynamicCard>
  )
}

export const WordCloudItem: React.FC<{
  children: React.ReactNode
}> = ({ children }) => {
  return (
    <li className="temp-word-cloud-li-item">{ children }</li>
  )
}