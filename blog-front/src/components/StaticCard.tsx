import styles from '@/styles/components/card.module.scss'
import React, { useMemo } from 'react'
import { CardProps } from '@/components/DynamicCard'

/**
 * 静态卡片组件
 * @param children 子元素
 * @param padding 内边距
 * @param radius 圆角
 * @param useDefaultPadding 是否使用默认的内边距
 * @param title 标题
 * @param icon 图标
 * @constructor
 */
const StaticCard: React.FC<CardProps> = ({
  children,
  padding = 0,
  radius = 12,
  useDefaultPadding = false,
  title,
  icon
}):React.ReactNode => {

  const cardStyle = useMemo((): Record<string, string> => ({
    padding: useDefaultPadding ? '' : typeof padding === 'number' ? `${padding}px` : padding,
    borderRadius: typeof radius === 'number' ? `${radius}px` : radius
  }), [useDefaultPadding, padding, radius])

  return (
    <div className={`${styles.card} animate-on-scroll w-full ${useDefaultPadding && 'sm:p-6 p-4'}`}
         style={cardStyle}
    >
      { (icon || title) && <div className="pb-2 flex gap-col-1.5 desc-text">
        { icon && <i className={`${icon} inline-block`}></i> }
        { title && <p className="text-sm">{ title }</p> }
      </div> }
      {children}
    </div>
  )
}

export default StaticCard;