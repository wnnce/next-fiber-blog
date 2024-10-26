'use client'

import '@/styles/components/client-components.scss'
import React, { useMemo, useState } from 'react'

export const TopicLike: React.FC<{
  topicId: number;
  count: number;
}> = ({ topicId, count }) => {
  const handleLike = (key: string | number, done: () => void) => {
    done();
  }

  return (
    <CommonLike count={count} entityKey={topicId} type="topic" onLike={handleLike} />
  )
}

export const ArticleLike: React.FC<{
  articleId: number;
  count: number;
}> = ({ articleId, count }) => {
  const handleLike = (key: string | number, done: () => void) => {
    done();
  }
  return (
    <CommonLike count={count} entityKey={articleId} type="article" onLike={handleLike} hideText />
  )
}

export const CommonLike: React.FC<{
  count: number,
  entityKey: string | number,
  type: 'topic' | 'article' | 'comment',
  onLike: (key: string | number, done: () => void) => void,
  className?: string,
  hideText?: boolean,
}> = ({ count, entityKey, type, onLike, className, hideText = false }) => {
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
    <button className={`desc-text flex items-start common-like-button ${className || ''}`}
            disabled={likeKeys[entityKey] === null}
            onClick={likeKeys[entityKey] === undefined ? () => {
              onLike(entityKey, handleLike)
            } : undefined}
    >
      {likeKeys[entityKey] === null ? (
        <i className="inline-block text-sm text-red-5 i-tabler:thumb-up-filled" />
      ) : (
        <i className="inline-block text-sm i-tabler:thumb-up common-like-icon" />
      )}
      {!hideText && <span className="text-xs ml-0.5">{likeCount}</span>}
    </button>
  )
}