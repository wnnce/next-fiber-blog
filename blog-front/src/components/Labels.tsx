import '@/styles/components/labels.scss';
import React, { CSSProperties } from 'react'

interface LabelProps {
  className?: string;
  style?: CSSProperties;
}

/**
 * 热门标签
 * @constructor
 */
export const HotLabel: React.FC<LabelProps> = ({ className, style }): React.ReactNode => {
  return (
    <label className={`${className ? className : ''} font-mono info-text bg-red-5 hot-label`}
           style={style} />
  )
}

/**
 * 置顶标签
 * @constructor
 */
export const TopLabel: React.FC<LabelProps> = ({ className, style }): React.ReactNode => {
  return (
    <label className={`${className ? className : ''} font-mono info-text top-label`}
           style={style} />
  )
}