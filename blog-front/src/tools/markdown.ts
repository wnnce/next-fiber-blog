import markdownit from 'markdown-it';
import hljs from 'highlight.js';
import type MarkdownIt from 'markdown-it'
import katex from '@vscode/markdown-it-katex'
import { alert } from "@mdit/plugin-alert"
import { imgLazyload } from "@mdit/plugin-img-lazyload"
import { tasklist } from "@mdit/plugin-tasklist"
import anchor from 'markdown-it-anchor';
import toc from 'markdown-it-toc-done-right';

export const getArticleRender = (tocRenderCallback?: (html: string) => void): MarkdownIt => {
  const articleRender = markdownit('default', {
    highlight: (str, lang) => {
      let resultValue: string
      if (lang && hljs.getLanguage(lang)) {
        resultValue = hljs.highlight(str, { language: lang }).value;
      } else {
        resultValue = str;
      }
      return `<pre><code class="hljs ${lang && `language-${lang}`}">${resultValue}</code></pre>`
    }
  }).use(katex).use(alert).use(imgLazyload).use(tasklist).use(anchor)

  if (tocRenderCallback) {
    articleRender.use(toc, {
      callback: (html: string, _: object) => {
        tocRenderCallback(html);
      }
    });
  }
  return articleRender;
}

export const getTopicRender = () => {
  return getArticleRender();
}

export const getCommentRender = () => {
  return markdownit('default', {
    breaks: true,
    highlight: (str, lang) => {
      let resultValue: string
      if (lang && hljs.getLanguage(lang)) {
        resultValue = hljs.highlight(str, { language: lang }).value;
      } else {
        resultValue = str;
      }
      return `<pre><code class="hljs ${lang && `language-${lang}`}">${resultValue}</code></pre>`
    }
  }).disable(['image', 'heading'], true)
}