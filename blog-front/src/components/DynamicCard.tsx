'use client'

import styles from '@/styles/components/card.module.scss'
import React, { useMemo, useRef } from 'react'
import type { MouseEvent } from 'react';
import { throttleTimer } from '@/tools/utils'

export interface CardProps {
  children: React.ReactNode;
  color?: string;
  padding?: number | string;
  // 是否使用默认内边距
  useDefaultPadding?: boolean;
  radius?: number | string;
  // 标题
  title?: number | string | boolean;
  // 图标
  icon?: string;
  // 盒子晃动的阈值
  multiple?: number;
}

/**
 * 动态卡片组件
 * @param children 子元素
 * @param color 背景动效颜色
 * @param padding 内边距
 * @param radius 圆角
 * @param useDefaultPadding 是否使用默认内边距
 * @param title 卡片标题
 * @param icon 卡片图标
 * @param multiple 图片晃动除以的阈值 越大晃动幅度越小
 * @constructor
 */
const DynamicCard: React.FC<CardProps> = (
  {
    children,
    color = 'rgb(70, 185, 82)',
    padding = 0,
    radius = 12,
    useDefaultPadding = true,
    title,
    icon,
    multiple = 20
  }
): React.ReactNode => {
  const cardRef = useRef<HTMLDivElement>(null);
  const lightRef = useRef<HTMLDivElement>(null);

  // 外层卡片样式
  const cardStyle = useMemo((): Record<string, string> => ({
    padding: useDefaultPadding ? '' : typeof padding === 'number' ? `${padding}px` : padding,
    borderRadius: typeof radius === 'number' ? `${radius}px` : radius,
  }), [useDefaultPadding, padding, radius])

  // 下层光源样式
  const lightStyle = useMemo((): Record<string, string> => ({
    backgroundColor: color
  }), [color])

  /**
   * 鼠标滑动事件处理 使用基于timer的节流函数 10ms内只执行一次
   */
  const handleMouseMove = throttleTimer((event: MouseEvent<HTMLDivElement>) => {
    if (!cardRef.current || !lightRef.current) {
      return;
    }
    const cardElement = cardRef.current, lightElement = lightRef.current;
    lightElement.style.display = 'block';
    const box = cardElement.getBoundingClientRect();
    const { clientX, clientY } = event;
    const { x, y, width, height } = box;
    const calcX = (clientX - x - (width / 2)) / multiple;
    const calcY = (clientY - y - (height / 2)) / multiple * -1;
    cardElement.style.transform = `rotateX(${calcX}deg) rotateY(${calcY}deg)`;
    lightElement.style.left = clientX - x - 100 + 'px';
    lightElement.style.top = clientY - y - 100 + 'px';
  }, 10)

  /**
   * 鼠标滑出事件处理
   * @param event 鼠标滑出事件
   */
  const handleMouseLeave =  (event: MouseEvent<HTMLDivElement>) => {
    if (!cardRef.current || !lightRef.current) {
      return;
    }
    cardRef.current.style.transform = 'rotateX(0) rotateY(0)';
    lightRef.current.style.display = 'none';
  }
  return (
    <div className="animate-on-scroll">
      <div ref={cardRef} className={`${styles.card} w-full h-full ${useDefaultPadding && 'sm:p-6 p-4'}`}
           style={cardStyle}
           onMouseMove={handleMouseMove}
           onMouseLeave={handleMouseLeave}
      >
        { (icon || title) && <div className="pb-2 flex gap-col-1.5 items-center desc-text">
          { icon && <i className={`${icon} inline-block`}></i> }
          { title && <p className="text-sm">{ title }</p> }
        </div> }

        <div ref={lightRef} className={`${styles.light} hidden`} style={lightStyle} />
        {children}
      </div>
    </div>
  )
}

export default DynamicCard;