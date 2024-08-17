import '@/styles/components/bulletin.scss';
import React from 'react'
import DynamicCard from '@/components/DynamicCard'

/**
 * 博客公告栏组件
 * @constructor
 */
const Bulletin: React.FC = (): React.ReactNode => {
  return (
    <DynamicCard padding="1.5rem" title="BULLETINS" icon="i-tabler:bell-exclamation">
      <section>
        <ul className="bulletin-swiper list-none">
          <li className="bulletin-swiper-item">
            <h3 className="text-xl main-text mb-2">
              测试标题测试标题1
            </h3>
            <p className="info-text text-sm">
              通知内容通知内容通知内容通知内容通知内容通知内容通知内容通知内容通知内容通知内容通知内容通知内容
            </p>
          </li>
          <li className="bulletin-swiper-item">
            <h3 className="text-xl main-text mb-2">
              测试标题测试标题2
            </h3>
            <p className="info-text text-sm">
              通知内容通知内容通知内容通知内容
            </p>
            <label className="danger-swiper-tag">
              DANGER
            </label>
          </li>
          <li className="bulletin-swiper-item">
            <h3 className="text-xl main-text mb-2">
              测试标题测试标题3
            </h3>
            <p className="info-text text-sm text-wrap">
              通知内容通知内容通知内容通知内容通知内容通知内容通知内容通知内容v通知内容通知内容通知内容通知内容
            </p>
            <label className="waring-swiper-tag">
              WARN
            </label>
          </li>
        </ul>
      </section>
    </DynamicCard>
  )
};

export default Bulletin;