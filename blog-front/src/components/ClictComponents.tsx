'use client'

import '@/styles/components/client-components.scss'
import React, { useMemo, useState } from 'react'

export const CommonLike: React.FC<{
  count: number,
  entityKey: string | number,
  type: 'topic' | 'article' | 'comment'
}> = ({ count, entityKey, type }) => {
  const [likeCount, setLikeCount] = useState<number>(count);

  const storageKey = useMemo<string>(() => {
    if (type === 'topic') {
      return 'TOPIC:LIKE:SET'
    } else if (type === 'article') {
      return 'ARTICLE:LIKE:SET'
    } else {
      return 'COMMENT:LIKE:SET'
    }
  }, [type])

  const likeKeys = useMemo<Record<string | number, null>>(() => {
    const stringValue = localStorage.getItem(storageKey)
    if (stringValue && stringValue.length > 0) {
      return JSON.parse(stringValue) as Record<string | number, null>;
    }
    return {};
  }, [storageKey])

  const handleLike = () => {
    likeKeys[entityKey] = null;
    localStorage.setItem(storageKey, JSON.stringify(likeKeys));
    setLikeCount(next => next + 1);
  }
  return (
    <>
      { likeKeys[entityKey] === null ? (
        <i className="inline-block text-red-5 i-tabler:heart-filled" />
      ) : (
        <i className="inline-block i-tabler:heart common-like-icon" onClick={handleLike} />
      ) }
      <span className="font-mono ml-1">{likeCount}</span>
    </>
  )
}