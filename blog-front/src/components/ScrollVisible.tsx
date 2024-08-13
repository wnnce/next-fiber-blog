"use client"

import React, { useEffect } from 'react'

/**
 * 处理组件滚动到可视区域后显示动画
 * 添加了指定 className 的元素在未出现在可视区区域时会隐藏显示
 * 只有出现在可视区域后才会播放动画动画效果并显示在页面中 该效果只会出现一次
 * @constructor
 */
const ScrollVisible: React.FC = (): null => {
  useEffect(() => {
    const elements = document.querySelectorAll('.animate-on-scroll');
    const obsServer = new IntersectionObserver((entries, observer) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          entry.target.classList.add('visible');
          observer.unobserve(entry.target);
        }
      })
    }, { threshold: 0.1 })
    elements.forEach(item => {
      obsServer.observe(item);
    })
    return () => {
      elements.forEach(item => {
        obsServer.unobserve(item);
      })
    }
  }, [])

  return null
}

export default ScrollVisible;