import React from 'react'
import useSiteConfiguration from '@/hooks/site-configuration'
import DynamicCard from '@/components/DynamicCard'
import SimpleMarkdown from '@/components/SimpleMarkdown'
import StaticCard from '@/components/StaticCard'
import Comment from '@/components/comment/Comment'

const About: React.FC = async () => {
  const [about] = await useSiteConfiguration().queryConfigs('about');
  return (
    <section className="flex flex-col gap-row-4">
      <DynamicCard padding="1.5rem" title="ABOUT" icon="i-tabler:user" multiple={80}>
        <h1 className="text-center font-bold text-4xl mb-8">关于我</h1>
        <SimpleMarkdown markdown={about.value.toString()} />
      </DynamicCard>
      <StaticCard padding="1.5rem">
        <Comment type={3} />
      </StaticCard>
    </section>
  )
}

export default About;