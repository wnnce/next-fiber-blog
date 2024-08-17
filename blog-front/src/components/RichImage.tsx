'use client'

import '@/styles/components/rich-image.scss';
import React, { CSSProperties, useState } from 'react'
import Image, { ImageLoaderProps } from 'next/image'
import { sliceThumbnailImageUrl } from '@/tools/utils'

// 图片的显示模式
declare type ImageMode = 'contain' | 'cover' | 'fill' | 'none' | 'scale-down';
// 图片加载状态 loading:加载中 success:加载完成  error:加载失败
declare type ImageState = 'loading' | 'success' | 'error';

interface Props {
  // 图片链接
  src: string,
  // 是否只显示缩略图
  thumbnail?: boolean;
  // 宽度
  width?: number;
  // 高度
  height?: number;
  // 是否填充容器
  fill?: boolean;
  // 是否懒加载
  lazy?: boolean;
  // 图片描述
  alt?: string;
  // 圆角
  radius?: number | string;
  // 图片显示模式
  mode?: ImageMode;
  // 图片加载完成的事件
  onDone?: () => void;
  // 类名列表
  className?: string,
  style?: CSSProperties,
}

const RichImage: React.FC<Props> = ({
  src,
  thumbnail = false,
  width = 100,
  height = 100,
  fill = false,
  lazy = true,
  alt = 'image',
  radius = 0,
  mode = 'cover',
  onDone,
  className,
  style
}): React.ReactNode => {
  const _radius = typeof radius === 'number' ? `${radius}px` : radius;
  const [imageStatus, setImageStatus] = useState<ImageState>('loading');
  const [blurMaskShow, setBlurMaskShow] = useState<boolean>(true);
  const thumbnailUrl = sliceThumbnailImageUrl(src, Math.min(height, width));

  const handleLoadDone = () => {
    setImageStatus('success');
    if (!thumbnail) {
      setTimeout(() => {
        setBlurMaskShow(false);
      }, 500)
    }
    onDone && onDone();
  }

  const imageLoader = ({width, src, quality}: ImageLoaderProps) => {
    if (src.startsWith('/b-oss')) {
      return 'https://file.qiniu.vnc.ink' + src.substring(6);
    }
    return src
  }
  return (
    <div className={ className ? `${className} rich-image` : 'rich-image' } style={{
      borderRadius: _radius,
      height: `${height}px`,
      width: `${width}px`,
      ...style,
    }}>
      {
        thumbnail ? (
          // 加载缩略图时 显示加载中状态
          imageStatus === 'loading' && <div className="rich-image-mask loading-mask flex justify-center items-center">
            <i className="inline-block i-tabler-rotate-clockwise-2 animate-spin text-lg" />
          </div>
        ) : (
          // 不加载所缩略图 你们开启图片模糊加载 使用缩略图和 blur 制作一个模糊遮罩 当图片加载完成后移除遮罩
          blurMaskShow && (
            <>
              <div className="rich-image-mask blur-mask" style={{
                backdropFilter: `blur(${imageStatus === 'loading' ? 4 : 0}px)`
              }} />
              <div className="rich-image-mask thumb-image-mask" style={{
                backgroundImage: `url(https://file.qiniu.vnc.ink${thumbnailUrl.substring(6)})`
              }} />
            </>
          )
        )
      }
      <Image
        src={thumbnail ? thumbnailUrl : src}
        alt={ alt }
        height={fill ? undefined : height}
        width={fill ? undefined : width}
        fill={fill}
        loading={lazy ? 'lazy' : undefined}
        placeholder="empty"
        loader={imageLoader}
        onLoad={handleLoadDone}
        style={{
          borderRadius: radius,
          objectFit: mode
        }}
      />
    </div>
  )
}

export default RichImage;