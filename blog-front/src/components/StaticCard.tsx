import styles from '@/styles/components/card.module.scss'
import React, { useMemo } from 'react'
import { CardProps } from '@/components/DynamicCard'

/**
 * 静态卡片组件
 * @param children 子元素
 * @param padding 内边距
 * @param radius 圆角
 * @constructor
 */
const StaticCard: React.FC<CardProps> = ({
  children,
  padding = 0,
  radius = 12,
}):React.ReactNode => {

  const cardStyle = useMemo((): Record<string, string> => ({
    padding: typeof padding === 'number' ? `${padding}px` : padding,
    borderRadius: typeof radius === 'number' ? `${radius}px` : radius
  }), [padding, radius])

  return (
    <div className="animate-on-scroll h-100 w-100">
      <div className={`${styles.card} w-full h-full`}
           style={cardStyle}
      >
        {children}
      </div>
    </div>
  )
}

export default StaticCard;