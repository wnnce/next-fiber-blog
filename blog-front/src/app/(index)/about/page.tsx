import React from 'react'
import useSiteConfiguration from '@/hooks/site-configuration'
import DynamicCard from '@/components/DynamicCard'
import SimpleMarkdown from '@/components/SimpleMarkdown'

const About: React.FC = async () => {
  const [about] = await useSiteConfiguration().queryConfigs('about');
  return (
    <DynamicCard padding="1.5rem" title="ABOUT" icon="i-tabler:user" multiple={80}>
      <h1 className="text-center font-bold text-4xl mb-8">关于我</h1>
      <SimpleMarkdown markdown={about.value.toString()} />
    </DynamicCard>
  )
}

export default About;