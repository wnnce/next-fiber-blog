import '@/styles/components/hot-article.scss';
import React from 'react'
import DynamicCard from '@/components/DynamicCard'

const HotArticles: React.FC = () => {
  return (
    <DynamicCard padding="1.5rem" title="HOTS" icon="i-tabler:flame">
      <section className="mt-2">
        <ol className="list-none flex flex-col gap-row-4 text-sm">
          <li className="line-clamp-1 hot-article-item hover-line-text">
            <a href="https://baidu.com" title="测试标题测试标题">
              测试标题测试标题测试标题测试标题测试标题测试标题测试标题测试标题
              <i className="i-tabler:chevron-right text-xl arrow-icon"></i>
              <i className="text-lg line-icon"></i>
            </a>
          </li>
          <li className="line-clamp-1 hot-article-item hover-line-text">
            <a href="https://baidu.com" title="测试标题测试标题">
              测试标题测试标题测试标题测试标题测试标题测试标题测试标题测试标题
              <i className="i-tabler:chevron-right text-xl arrow-icon"></i>
              <i className="text-lg line-icon"></i>
            </a>
          </li>
          <li className="line-clamp-1 hot-article-item hover-line-text">
            <a href="https://baidu.com" title="测试标题测试标题">
              测试标题测试标题测试标题测试标题测试标题测试标题测试标题测试标题
              <i className="i-tabler:chevron-right text-xl arrow-icon"></i>
              <i className="text-lg line-icon"></i>
            </a>
          </li>
          <li className="line-clamp-1 hot-article-item hover-line-text">
            <a href="https://baidu.com" title="测试标题测试标题">
              测试标题测试标题测试标题测试标题测试标题测试标题测试标题测试标题
              <i className="i-tabler:chevron-right text-xl arrow-icon"></i>
              <i className="text-lg line-icon"></i>
            </a>
          </li>
          <li className="line-clamp-1 hot-article-item hover-line-text">
            <a href="https://baidu.com" title="测试标题测试标题">
              测试标题测试标题测试标题测试标题测试标题测试标题测试标题测试标题
              <i className="i-tabler:chevron-right text-xl arrow-icon"></i>
              <i className="text-lg line-icon"></i>
            </a>
          </li>
          <li className="line-clamp-1 hot-article-item hover-line-text">
            <a href="https://baidu.com" title="测试标题测试标题">
              测试标题测试标题测试标题测试标题测试标题测试标题测试标题测试标题
              <i className="i-tabler:chevron-right text-xl arrow-icon"></i>
              <i className="text-lg line-icon"></i>
            </a>
          </li>
        </ol>
      </section>
    </DynamicCard>
  )
}

export default HotArticles;