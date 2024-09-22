import '@/styles/components/comment.scss'
import '@/styles/components/markdown.scss'
import 'github-markdown-css/github-markdown-dark.css'
import 'highlight.js/styles/atom-one-dark.min.css'
import React, { useContext, useMemo, useState } from 'react'
import { Comment } from '@/lib/types'
import Image from 'next/image'
import { formatDateTime, formatUa } from '@/tools/utils'
import CommentList from '@/components/comment/CommentList'
import { CommonLike } from '@/components/client/CommonLike'
import { LevelContext } from '@/components/comment/context/LevelContext'
import { CommentEditor } from '@/components/comment/Comment'
import useMarkdownParse from '@/hooks/markdown'

const CommentItem: React.FC<{
  comment: Comment;
  topicId?: number;
  articleId?: number;
}> = ({ comment, articleId, topicId }) => {
  const level = useContext<number>(LevelContext);
  const [ isReply, setIsReply ] = useState<boolean>(false);
  const commentRender = useMarkdownParse().commentRender();

  const commentContent = useMemo<string>(() => {
    return commentRender.render(comment.content);
  }, [ comment, commentRender ])

  const avatarSize = useMemo<number>(() => {
    const size = 48 - ((level - 1) * 12)
    return size < 12 ? 12 : size;
  }, [level])

  const result = formatUa(comment.commentUa);
  const browser = useMemo((): string => {
    const { name, version } = result.browser;
    if (!result.browser || !name) {
      return '未知'
    }
    return name + ' ' + (version ? version.substring(0, version.indexOf('.')) : '')
  }, [result.browser])

  const system = useMemo((): string => {
    const { name, version } = result.os;
    if (!result.os || !name) {
      return '未知';
    }
    if (!version) {
      return name;
    }
    const index = version.indexOf('.')
    return name + ' ' +  (index > 0 ? version.substring(0, index): version);
  }, [result.os])

  return (
    <div className="comment-item flex gap-col-3">
      <div className="item-user-avatar shrink-0">
        <Image src={comment.user.avatar} alt="avatar"
               height={avatarSize} width={avatarSize}
               objectFit="cover"
               className="rounded-md"
        />
      </div>
      <div className="item-body flex-1 flex flex-col gap-row-2">
        <div className="item-body-header flex gap-1 items-center flex-wrap desc-text">
          <a className="username text-sm" href={comment.user.link}
             title={comment.user.link}
             target="_blank"
          >
            {comment.user.nickname}
          </a>
          <span className="user-level px-1 main-text">{`Lv${comment.user.level}`}</span>
          {(comment.rid != 0 && comment.fid != 0 && comment.parentUser) && (
            <>
              <i className="inline-block i-tabler:arrow-badge-right" />
              <a className="username text-sm" href={comment.parentUser.link}
                 title={comment.parentUser.link}
                 target="_blank"
              >
                {comment.parentUser.nickname}
              </a>
            </>
          )}
          <time className="text-xs" dateTime={comment.createTime}
                title={comment.createTime}
          >
            {formatDateTime(comment.createTime)}
          </time>
          <span className="body-header-tag text-xs">重庆</span>
          <span className="body-header-tag text-xs">{browser}</span>
          <span className="body-header-tag text-xs">{system}</span>
        </div>
        <div className="item-body-content info-text markdown-body"
             dangerouslySetInnerHTML={{ __html: commentContent }}
        />
        <div className="item-options flex gap-col-3 items-center">
          <CommonLike className="option-button" count={comment.voteUp} entityKey={comment.commentId}
                      type="comment"
                      onLike={(key, done) => {
                        done()
                      }}
          />
          <button className="desc-text option-button text-xs" onClick={() => {
              setIsReply(true)
            }}
          >
            回复
          </button>
        </div>
        { isReply && (
          <CommentEditor
            fid={comment.fid === 0 ? comment.commentId : comment.fid}
            rid={(comment.fid === 0 && comment.rid === 0) ? 0 : comment.commentId}
            parentNickname={comment.user.nickname}
            onClose={() => {
              setIsReply(false);
            }}
          />
        ) }
        {(comment.children && comment.children.records.length > 0) && (
          <LevelContext.Provider value={level + 1}>
            <CommentList
              page={comment.children}
              type={comment.commentType}
              articleId={articleId}
              topicId={topicId}
              fid={comment.commentId}
            />
          </LevelContext.Provider>
        ) }
      </div>
    </div>
  )
}

export default CommentItem;