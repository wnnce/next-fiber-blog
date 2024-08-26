import 'github-markdown-css/github-markdown-dark.css';
import '@/styles/components/markdown.scss';
import React, { CSSProperties } from 'react'
import { unified } from 'unified'
import remarkParse from 'remark-parse'
import remarkFrontmatter from 'remark-frontmatter'
import remarkGfm from 'remark-gfm'
import remarkRehype from 'remark-rehype'
import rehypePrettyCode from 'rehype-pretty-code'
import rehypeStringify from 'rehype-stringify'

const SimpleMarkdown: React.FC<{
  markdown: string,
  className?: string,
  style?: CSSProperties,
}> = async ({ markdown, className, style }) => {
  const markdownFile = await unified()
    .use(remarkParse)
    .use(remarkFrontmatter)
    .use(remarkGfm)
    .use(remarkRehype)
    .use(rehypePrettyCode, {
      grid: false
    })
    .use(rehypeStringify)
    .process(markdown);
  return (
    <div className={`markdown-body ${className || ''}`}
         style={style}
         dangerouslySetInnerHTML={{ __html: markdownFile.toString() }}
    />
  )
}

export default SimpleMarkdown;