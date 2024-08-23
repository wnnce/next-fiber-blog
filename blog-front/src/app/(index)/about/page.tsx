import 'github-markdown-css/github-markdown-dark.css';
import '@/styles/components/markdown.scss';
import React from 'react'
import useSiteConfiguration from '@/hooks/site-configuration'
import DynamicCard from '@/components/DynamicCard'
import rehypeStringify from 'rehype-stringify'
import rehypePrettyCode from "rehype-pretty-code";
import remarkFrontmatter from 'remark-frontmatter'
import remarkGfm from 'remark-gfm'
import remarkParse from 'remark-parse'
import remarkRehype from 'remark-rehype'
import {unified} from 'unified'

const About: React.FC = async () => {
  const [about] = await useSiteConfiguration().queryConfigs('about');
  const markdownFile = await unified()
    .use(remarkParse)
    .use(remarkFrontmatter)
    .use(remarkGfm)
    .use(remarkRehype)
    .use(rehypePrettyCode)
    .use(rehypeStringify)
    .process(about.value.toString());
  const aboutHtml = String(markdownFile);
  return (
    <DynamicCard padding="1.5rem" title="ABOUT" icon="i-tabler:user" multiple={80}>
      <h1 className="text-center font-bold text-4xl">关于我</h1>
      <div className="mt-4 markdown-body" dangerouslySetInnerHTML={{ __html: aboutHtml }} />
    </DynamicCard>
  )
}

export default About;