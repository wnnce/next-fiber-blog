import '@/styles/components/rich-image.scss'
import { ReactPortal, useRef } from 'react'
import ReactDOM from 'react-dom'

const useImagePreview = () => {

  const elementRef= useRef<HTMLDivElement>();

  const previewImage = (src: string): ReactPortal => {
    const newContainer = document.createElement('div');
    newContainer.classList.add('rich-image-preview-container')
    document.body.appendChild(newContainer);
    elementRef.current = newContainer;
    setTimeout(() => {
      newContainer.classList.add('preview-visible')
    }, 100)
    const originSrc = src.startsWith('/b-oss') ? process.env.NEXT_PUBLIC_QINIU_IMAGE_DOMAIN + src.substring(6) : src;
    return ReactDOM.createPortal(
      <>
        <div className="preview-close-button" onClick={destroy}>
          <i className="i-tabler:x" />
        </div>
        <div className="preview-image-body">
          <img src={originSrc} alt="" className="preview-image" loading="lazy" />
        </div>
      </>
      , newContainer)
  }

  const destroy = () => {
    if (elementRef.current) {
      elementRef.current.classList.remove('preview-visible')
      setTimeout(() => {
        if (!elementRef.current) {
          return
        }
        document.body.removeChild(elementRef.current);
        elementRef.current = undefined;
      }, 500)
    }
  }

  return { previewImage, destroy }
}

export default useImagePreview;