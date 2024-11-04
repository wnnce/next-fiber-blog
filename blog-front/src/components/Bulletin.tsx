import '@/styles/components/bulletin.scss';
import React from 'react'
import DynamicCard from '@/components/DynamicCard'
import { listNoticeByType } from '@/lib/api'
import { Swiper, SwiperSlide } from '@/components/Swiper'

/**
 * 博客公告栏组件
 * @constructor
 */
const Bulletin: React.FC = async () => {
  const { data: notices } = await listNoticeByType(2);
  return (
    <DynamicCard title="BULLETINS" icon="i-tabler:bell-exclamation">
      <section className="mx--6 bulletin-swiper">
        <Swiper autoPlay showButton={false} showDot={false}>
          { notices.map(notice => (
            <SwiperSlide className="bulletin-swiper-item" key={notice.noticeId}>
              <h3 className="text-xl main-text mb-2">
                {notice.title}
              </h3>
              <p className="info-text text-sm">
                {notice.message}
              </p>
              { notice.level === 2 && (
                <label className="waring-swiper-tag">
                  WARN
                </label>
              )}
              { notice.level === 3 && (
                <label className="danger-swiper-tag">
                  DANGER
                </label>
              )}
            </SwiperSlide>
          ))}
        </Swiper>
      </section>
    </DynamicCard>
  )
};

export default Bulletin;