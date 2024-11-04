import 'github-markdown-css/github-markdown-dark.css'
import 'highlight.js/styles/atom-one-dark.min.css'
import '@/styles/components/markdown.scss'
import '@/styles/components/topic.scss'
import React from 'react'
import { pageTopic } from '@/lib/api'
import ServerPagination from '@/components/ServerPagination'
import StaticCard from '@/components/StaticCard'
import { Timeline, TimeLineItem } from '@/components/Timeline'
import RichImage from '@/components/RichImage'
import { TopicLike } from '@/components/client/CommonLike'
import useMarkdownParse from '@/hooks/markdown'

const TopicPage: React.FC<{
  params: {
    page: string
  }
}> = async ({ params }) => {
  const { page } = params;
  const numberPage = parseInt(page);
  if (!numberPage || isNaN(numberPage) || numberPage <= 0 ) {
    throw new Error('动态参数错误')
  }
  const { data: topicPage } = await pageTopic({ page: numberPage, size: 10 });

  const topicRender = useMarkdownParse().topicRender();

  return (
    <>
      <StaticCard useDefaultPadding title="TOPICS" icon="i-tabler:world" multiple={40}>
        <h1 className="text-center font-bold sm:text-4xl text-2xl mb-8">
          我的动态
        </h1>
        <Timeline>
          { topicPage.records.map(topic => (
            <TimeLineItem key={topic.topicId} time={topic.createTime}>
              <div className="topic-list-li-content markdown-body topic-markdown"
                   dangerouslySetInnerHTML={{ __html: topicRender.render(topic.content) }}
              />
              { (topic.imageUrls && topic.imageUrls.length > 0) && (
                topic.mode === 1 ? (
                  <div className="flex flex-wrap gap-1 mt-4">
                    {topic.imageUrls.map(image => (
                      <RichImage src={image} width={100} height={100} fill thumbnail preview radius={4} key={image} />
                    ))}
                  </div>
                ) : (
                  <div className="topic-flow-photo mt-4">
                    { topic.imageUrls.map(image => (
                      <RichImage src={image} key={image} width={0} height={0} imageClassName="w-full h-auto" thumbnail preview />
                    )) }
                  </div>
                )
              )}
              <ul className="flex gap-col-4 mt-4 desc-text relative">
                <li className="flex items-center">
                  <TopicLike topicId={topic.topicId} count={topic.voteUp} />
                </li>
                <li className="flex items-center">
                  <i className="inline-block i-tabler:message-chatbot" />
                  <span className="font-mono ml-1">0</span>
                </li>
                { topic.isTop && (
                  <li className="flex items-center">
                    <i className="i-tabler:pin-filled inline-block primary-color" title="置顶" />
                  </li>
                )}
                { topic.isHot && (
                  <li className="flex items-center">
                    <i className="i-tabler:flame inline-block text-orange-5" title="热门" />
                  </li>
                )}
                <li className="flex items-center absolute right-0 bottom-0">
                  <i className="inline-block text-sm i-tabler:location" />
                  <span className="ml-1 text-xs">中国-重庆</span>
                </li>
              </ul>
            </TimeLineItem>
          ))}
        </Timeline>

      </StaticCard>
      <ServerPagination className="mt-4 w-full animate-on-scroll" current={topicPage.current} size={topicPage.size}
                        total={topicPage.total}
                        targetUrlPrefix="/topic/page/"
      />
    </>
  )
}

export default TopicPage;