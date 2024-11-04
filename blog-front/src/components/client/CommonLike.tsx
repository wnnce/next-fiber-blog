'use client'

import '@/styles/components/client-components.scss'
import React, { useCallback, useEffect, useMemo, useState } from 'react'
import useMessage from '@/components/message'
import { articleVoteUp, topicVoteUp } from '@/lib/client-api'

export const TopicLike: React.FC<{
  topicId: number;
  count: number;
}> = ({ topicId, count }) => {

  const message = useMessage();

  const handleLike = async (key: string | number, done: () => void) => {
    const loadingMessage = message.showLoading('处理中...');
    try {
      const result = await topicVoteUp(topicId);
      if (result.code === 200) {
        message.showSuccess('点赞成功');
        done();
      }
    } finally {
      loadingMessage.close();
    }
  }

  return (
    <CommonLike count={count} entityKey={topicId} type="topic" onLike={handleLike} />
  )
}

export const ArticleLike: React.FC<{
  articleId: number;
  count: number;
}> = ({ articleId, count }) => {

  const message = useMessage();

  const handleLike = async (key: string | number, done: () => void) => {
    const loadingMessage = message.showLoading('处理中...');
    try {
      const result = await articleVoteUp(articleId);
      if (result.code === 200) {
        message.showSuccess('点赞成功');
        done();
      }
    } finally {
      loadingMessage.close();
    }
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
  const [isLike, setIsLike] = useState<boolean>(false);

  const storageKey = useMemo<string>(() => {
    if (type === 'topic') {
      return 'TOPIC:LIKE:SET'
    } else if (type === 'article') {
      return 'ARTICLE:LIKE:SET'
    } else {
      return 'COMMENT:LIKE:SET'
    }
  }, [type])

  const handleLike = () => {
    const stringValue = localStorage.getItem(storageKey)
    let newLikeKeys: Record<string | number, null> = {}
    if (stringValue && stringValue.length > 0) {
      newLikeKeys = JSON.parse(stringValue) as Record<string | number, null>
    }
    newLikeKeys[entityKey] = null
    localStorage.setItem(storageKey, JSON.stringify(newLikeKeys));
    setLikeCount(prev => prev + 1);
    setIsLike(true);
  }

  const handleLocalStoreChange = useCallback((event: StorageEvent) => {
    if (!event.key || event.key != storageKey) {
      return
    }
    const stringValue = localStorage.getItem(storageKey)
    if (stringValue && stringValue.length > 0) {
      const likeKeys = (JSON.parse(stringValue) as Record<string | number, null>)
      if (likeKeys[entityKey] === null && !isLike) {
        setLikeCount(prev => prev + 1)
        setIsLike(true)
      }
    }
  }, [storageKey, entityKey, isLike])

  useEffect(() => {
    const stringValue = localStorage.getItem(storageKey)
    if (stringValue && stringValue.length > 0) {
      const likeKeys = (JSON.parse(stringValue) as Record<string | number, null>)
      likeKeys[entityKey] === null && setIsLike(true)
    }
    window.addEventListener('storage', handleLocalStoreChange)
    return () => {
      window.removeEventListener('storage', handleLocalStoreChange)
    }
  }, [storageKey, entityKey, handleLocalStoreChange])

  return (
    <button className={`desc-text flex items-start common-like-button ${className || ''}`}
            disabled={isLike}
            onClick={!isLike ? () => {
              onLike(entityKey, handleLike)
            } : undefined}
    >
      {isLike ? (
        <i className="inline-block text-sm text-red-5 i-tabler:thumb-up-filled" />
      ) : (
        <i className="inline-block text-sm i-tabler:thumb-up common-like-icon" />
      )}
      {!hideText && <span className="text-xs ml-0.5">{likeCount}</span>}
    </button>
  )
}