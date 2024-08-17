'use client'

import '@/styles/components/swiper.scss';
import React, { type MouseEvent, TouchEvent, useEffect, useMemo, useRef } from 'react'
import { throttle } from '@/tools/utils'

interface SwiperProps {
  draggable?: boolean;
  showDot?: boolean;
  showButton?: boolean;
  autoPlay?: boolean;
  interval?: number;
  dotActive?: boolean;
  children?: React.ReactNode;
}

/**
 * 轮播图组件
 * @param draggable 是否支持拖拽
 * @param showDot 是否显示指示点
 * @param showButton 是否显示切换按钮
 * @param autoPlay 是否自动播放
 * @param interval 自动播放间隔
 * @param dotActive 指示点是否可用
 * @param children 子组件
 * @constructor
 */
export const Swiper: React.FC<SwiperProps> = ({
  draggable = true,
  showDot = true,
  showButton = true,
  autoPlay = true,
  interval = 4000,
  dotActive = false,
  children
}): React.ReactNode => {
  const currentIndexRef = useRef<number>(0);
  const timerRef = useRef<any>(null);
  const swiperRef = useRef<HTMLUListElement>(null);
  const dotLineRef = useRef<HTMLElement>(null);
  const isDragging = useRef<boolean>(false);
  const startClientX = useRef<number>(0);
  const totalSlide = useMemo<number>(() => {
    return React.Children.count(children)
  }, [children]);
  const dotIndexList = useMemo<number[]>(() => {
    if (!showDot) {
      return [];
    }
    return Array.from({ length: totalSlide }, (_, index) => index + 1);
  }, [showDot, totalSlide])

  useEffect(() => {
    autoPlay && startAutoPlay();
    return () => {
      stopAutoPlay();
    }
  })

  const startAutoPlay = () => {
    stopAutoPlay();
    timerRef.current = setTimeout(() => {
      currentIndexRef.current = currentIndexRef.current >= totalSlide - 1 ? 0 : currentIndexRef.current + 1;
      handleUpdateSwiperTranslateX();
      startAutoPlay();
    }, interval);
  }

  const stopAutoPlay = () => {
    timerRef.current && clearTimeout(timerRef.current);
  }

  const handleSwiperMouseEnter = () => {
    // 鼠标进入停止自动播放
    stopAutoPlay();
  }
  const handleSwiperMouseMove = throttle((event: MouseEvent<HTMLDivElement>) => {
    handleMove(event.clientX);
  })
  const handleSwiperMouseLeave = (event: MouseEvent<HTMLDivElement>) => {
    if (isDragging.current) {
      isDragging.current = false;
      handleMoveSwitch(event.clientX);
    }
    // 鼠标移出开始自动播放
    autoPlay && startAutoPlay();
  }
  const handleSwiperMouseDown = (event: MouseEvent<HTMLDivElement>) => {
    isDragging.current = true;
    startClientX.current = event.clientX;
  }
  const handleSwiperMouseUp = (event: MouseEvent<HTMLDivElement>) => {
    isDragging.current = false;
    handleMoveSwitch(event.clientX);
  }
  const handleTouchMouseStart = (event: TouchEvent<HTMLDivElement>) => {
    stopAutoPlay();
    isDragging.current = true;
    startClientX.current = event.touches[0].clientX;
  }
  const handleTouchMouseMove = throttle((event: TouchEvent<HTMLDivElement>) => {
    handleMove(event.touches[0].clientX);
  })
  const handleTouchMouseEnd = (event: TouchEvent<HTMLDivElement>) => {
    isDragging.current = false;
    handleMoveSwitch(event.changedTouches[0].clientX)
    autoPlay && startAutoPlay();
  }

  const handleMove = (clientX: number) => {
    if (!swiperRef.current || !isDragging.current) {
      return;
    }
    const diff = clientX - startClientX.current;
    if ((diff > 0 && currentIndexRef.current === 0) || (diff < 0 && currentIndexRef.current >= totalSlide - 1)) {
      return;
    }
    swiperRef.current.style.transition = 'none';
    swiperRef.current.style.transform = `translateX(calc(-${currentIndexRef.current * 100}% + ${diff}px))`
  }

  const handleMoveSwitch = (clientX: number) => {
    if (!swiperRef.current) {
      return;
    }
    const moveWidth = clientX - startClientX.current;
    if ((moveWidth > 0 && currentIndexRef.current === 0) || (moveWidth < 0 && currentIndexRef.current >= totalSlide - 1)) {
      return;
    }
    if (Math.abs(moveWidth) >= swiperRef.current.offsetWidth / 3) {
      currentIndexRef.current = currentIndexRef.current + (moveWidth < 0 ? 1 : -1);
    }
    swiperRef.current.style.transition = 'transform 500ms ease';
    handleUpdateSwiperTranslateX();
  }
  const handleSwitchPre = () => {
    stopAutoPlay();
    if (currentIndexRef.current === 0) {
      return;
    }
    currentIndexRef.current = currentIndexRef.current - 1;
    handleUpdateSwiperTranslateX();
    autoPlay && startAutoPlay();
  }

  const handleSwitchNext = () => {
    stopAutoPlay();
    if (currentIndexRef.current === totalSlide - 1) {
      return;
    }
    currentIndexRef.current = currentIndexRef.current + 1;
    handleUpdateSwiperTranslateX();
    autoPlay && startAutoPlay();
  }

  const handleUpdateSwiperTranslateX = () => {
    if (!swiperRef.current) {
      return;
    }
    const num = currentIndexRef.current * 100;
    swiperRef.current.style.transform = `translateX(-${num}%)`;
    if (dotLineRef.current && showDot) {
      dotLineRef.current.style.transform = `translateX(${num}%)`;
    }
  }

  const handleDotClick = (dotIndex: number) => {
    if (dotIndex === currentIndexRef.current) {
      return;
    }
    currentIndexRef.current = dotIndex;
    handleUpdateSwiperTranslateX();
  }

  return (
    <div className="swiper-container"
         onMouseEnter={handleSwiperMouseEnter}
         onMouseLeave={handleSwiperMouseLeave}
    >
      {/* 拖拽遮罩 */}
      {draggable &&
        <div className="swiper-mouse-mask"
             onMouseMove={handleSwiperMouseMove}
             onMouseDown={handleSwiperMouseDown}
             onMouseUp={handleSwiperMouseUp}
             onTouchStart={handleTouchMouseStart}
             onTouchEnd={handleTouchMouseEnd}
             onTouchMove={handleTouchMouseMove}
        />
      }
      {showButton && <>
        <button className="swiper-left-button" onClick={handleSwitchPre}>
          <i className="inline-block i-tabler:chevron-left"></i>
        </button>
        <button className="swiper-right-button" onClick={handleSwitchNext}>
          <i className="inline-block i-tabler:chevron-right"></i>
        </button>
      </>}
      {showDot && <div className="dot-container">
        <i className="dot-active-line" ref={dotLineRef} />
        {dotIndexList.map(index => (
          <i className={ dotActive ? 'dot-item active-dot-item' : 'dot-item' }
             onClick={dotActive ? () => {
               handleDotClick(index - 1);
             } : undefined}
             key={index}
          />
        ))}
      </div>}
      <ul className="swiper list-none" ref={swiperRef}>
        {children}
      </ul>
    </div>
  )
}

/**
 * 轮播图子组件
 * @param children
 * @constructor
 */
export const SwiperSlide: React.FC<{
  children: React.ReactNode
}> = ({ children }): React.ReactNode => {
  return <li className="swiper-slide">{ children }</li>
}