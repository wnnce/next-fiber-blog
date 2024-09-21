import markdownit from 'markdown-it';
import hljs from 'highlight.js';
import type MarkdownIt from 'markdown-it'

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