'use client'

import React, { useCallback, useEffect, useRef } from 'react'
import { usePathname } from 'next/navigation'

/**
 * 处理组件滚动到可视区域后显示动画
 * 添加了指定 className 的元素在未出现在可视区区域时会隐藏显示
 * 只有出现在可视区域后才会播放动画动画效果并显示在页面中 该效果只会出现一次
 * @constructor
 */
const ScrollVisible: React.FC = (): null => {
  const pathName = usePathname();
  const elements = useRef<Set<Element>>(new Set<Element>());
  const obsServer = useRef<IntersectionObserver>();
  const updateObserveElements = useCallback((newElement: Set<Element>) => {
    const toRemove: Element[] = [], toAdd: Element[] = [];
    elements.current.forEach(item => {
      if (!newElement.has(item)) {
        toRemove.push(item);
      }
    })
    newElement.forEach(item => {
      if (!elements.current.has(item)) {
        toAdd.push(item);
      }
    })
    toRemove.forEach(item => {
      obsServer.current?.unobserve(item);
    })
    toAdd.forEach(item => {
      obsServer.current?.observe(item);
    })
    elements.current = newElement;
  }, [])

  useEffect(() => {
    if (!obsServer.current) {
      obsServer.current = new IntersectionObserver((entries, observer) => {
        entries.forEach(entry => {
          if (entry.isIntersecting) {
            entry.target.classList.add('visible');
            observer.unobserve(entry.target);
          }
        })
      }, { threshold: 0.1 })
    }
    const newElements = new Set(document.querySelectorAll('.animate-on-scroll'));
    updateObserveElements(newElements);
  }, [pathName, updateObserveElements])
  return null
}

export default ScrollVisible;