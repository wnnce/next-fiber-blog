import React, { useCallback, useState } from 'react'
import { CommentProps } from '@/components/comment/Comment'
import { Comment, Page } from '@/lib/types'
import { pageComment } from '@/lib/client-api'
import CommentItem from '@/components/comment/CommentItem'
import Empty from '@/components/Empty'
import Pagination from '@/components/client/Pagination'

interface Props extends CommentProps {
  fid?: number;
  page: Page<Comment>;
}

const CommentList: React.FC<Props> = ({ type, articleId, topicId, fid, page }) => {
  const [ data, setData ] = useState<Page<Comment>>(page);
  const [ loading, setLoading ] = useState<boolean>(false);

  const queryCommentPage = useCallback(async (page: number) => {
    setLoading(true);
    try {
      const result = await pageComment({ page: page, size: 10, commentType: type, articleId, topicId })
      if (result.code === 200) {
        setData(result.data);
      }
    } finally {
      setLoading(false);
    }

  }, [articleId, topicId, type])

  return (
    ( data.records && data.records.length > 0 ? (
      <div className="relative">
        { loading && (
          <div className="comment-loading">
            <i className="inline-block i-tabler:loader-2 text-8 animate-spin" />
          </div>
        )}
        <ul className="list-none flex flex-col gap-row-4">
          {data.records.map(comment => (
            <li key={comment.commentId}>
              <CommentItem comment={comment} articleId={articleId} topicId={topicId} />
            </li>
          ))}
        </ul>
        { data.pages > 1 && (
          <Pagination className="mt-2" page={data.current} pages={data.pages}
                      onChange={(newPage) => queryCommentPage(newPage)}
          />
        ) }
      </div>
    ) : (
      <Empty iconSize="4rem" textSize="0.85rem" text="还没有评论..." />
    ))
  )
}

export default CommentList;