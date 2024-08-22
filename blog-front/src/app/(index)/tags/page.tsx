import React from 'react'
import DynamicCard from '@/components/DynamicCard'
import { listTag } from '@/lib/api'

const Tags: React.FC = async () => {
  const { data: tags } = await listTag();
  return (
    <DynamicCard padding="1.5rem" title="TAGS" icon="i-tabler:tags">
      <section className="text-sm">
        <ul className="list-none flex flex-wrap gap-4">
          {tags.map(tag => (
            <li key={tag.tagId}>
              <a href="#" title={tag.tagName}>
                <span className="inline-block px-2 py-1 rounded-l-1"
                      style={{ backgroundColor: tag.color }}>{tag.tagName}</span>
                <span className="inline-block px-2 py-1 rounded-r-1 bg-gray-5">{tag.articleNum}</span>
              </a>
            </li>
          ))}
        </ul>
      </section>
    </DynamicCard>
  )
}

export default Tags;