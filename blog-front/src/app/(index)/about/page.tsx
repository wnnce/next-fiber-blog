import 'highlight.js/styles/atom-one-dark.min.css'
import 'github-markdown-css/github-markdown-dark.css'
import '@/styles/components/markdown.scss'

import React from 'react'
import DynamicCard from '@/components/DynamicCard'
import StaticCard from '@/components/StaticCard'
import Comment from '@/components/comment/Comment'
import { getArticleRender } from '@/tools/markdown'
import { querySiteConfigs } from '@/tools/site-configuration'

const About: React.FC = async () => {
  const [about] = await querySiteConfigs('about');
  const aboutHtml = getArticleRender().render(about.value.toString());
  return (
    <section className="flex flex-col gap-row-4">
      <DynamicCard title="ABOUT" icon="i-tabler:user" multiple={80}>
        <h1 className="text-center font-bold text-4xl mb-8">关于我</h1>
        <div className="markdown-body about-markdown" dangerouslySetInnerHTML={{ __html: aboutHtml }}/>
      </DynamicCard>
      <StaticCard useDefaultPadding>
        <Comment type={3} />
      </StaticCard>
    </section>
  )
}

export default About;