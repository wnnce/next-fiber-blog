import 'highlight.js/styles/atom-one-dark.min.css'
import 'github-markdown-css/github-markdown-dark.css'
import '@/styles/components/markdown.scss'

import React from 'react'
import useSiteConfiguration from '@/hooks/site-configuration'
import DynamicCard from '@/components/DynamicCard'
import StaticCard from '@/components/StaticCard'
import Comment from '@/components/comment/Comment'
import useMarkdownParse from '@/hooks/markdown'

const About: React.FC = async () => {
  const [about] = await useSiteConfiguration().queryConfigs('about');
  const { articleRender } = useMarkdownParse();
  const aboutHtml = articleRender().render(about.value.toString());
  return (
    <section className="flex flex-col gap-row-4">
      <DynamicCard padding="1.5rem" title="ABOUT" icon="i-tabler:user" multiple={80}>
        <h1 className="text-center font-bold text-4xl mb-8">关于我</h1>
        <div className="markdown-body about-markdown" dangerouslySetInnerHTML={{__html: aboutHtml}} />
      </DynamicCard>
      <StaticCard padding="1.5rem">
        <Comment type={3} />
      </StaticCard>
    </section>
  )
}

export default About;