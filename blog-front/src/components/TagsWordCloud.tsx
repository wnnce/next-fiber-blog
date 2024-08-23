import React from 'react'
import { listTag } from '@/lib/api'
import { WordCloud, WordCloudItem } from '@/components/WordCloud'
import Link from 'next/link'

const TagsWordCloud: React.FC = async () => {
  const { data: tags } = await listTag();
  return (
    <WordCloud>
      { tags.map(tag => (
        <WordCloudItem key={tag.tagId}>
          <Link href={`/tag/${tag.tagId}/page/1`} target="_blank" style={{color: tag.color}}
                title={tag.tagName}
          >
            { tag.tagName }
          </Link>
        </WordCloudItem>
      )) }
    </WordCloud>
  )
}

export default TagsWordCloud;