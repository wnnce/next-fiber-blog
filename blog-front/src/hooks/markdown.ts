import markdownit from 'markdown-it';
import hljs from 'highlight.js';
import type MarkdownIt from 'markdown-it'
import { alert } from "@mdit/plugin-alert"
import { imgLazyload } from "@mdit/plugin-img-lazyload"
import { tasklist } from "@mdit/plugin-tasklist"
import anchor from 'markdown-it-anchor';

const useMarkdownParse = () => {
  let _articleRender: MarkdownIt;
  let _commentRender: MarkdownIt;

  const articleRender = (): MarkdownIt => {
    if (!_articleRender) {
      _articleRender = markdownit('default', {
        highlight: (str, lang) => {
          let resultValue: string
          if (lang && hljs.getLanguage(lang)) {
            resultValue = hljs.highlight(str, { language: lang }).value;
          } else {
            resultValue = str;
          }
          return `<pre><code class="hljs ${lang && `language-${lang}`}">${resultValue}</code></pre>`
        }
      }).use(alert).use(imgLazyload).use(tasklist).use(anchor, {
        callback: (token, info) => {
          console.log(info);
        }
      });
    }
    return _articleRender;
  }

  const topicRender = () => {
    return articleRender();
  }

  const commentRender = () => {
    if (!_commentRender) {
      _commentRender = markdownit('default', {
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
    return _commentRender;
  }

  return { articleRender, topicRender, commentRender }
}

export default useMarkdownParse;