'use client'

import '@/styles/components/client-components.scss'
import '@/styles/animate.css'
import React, {
  ChangeEvent,
  CompositionEvent,
  ReactPortal,
  useEffect,
  useRef,
  useState
} from 'react'
import { CSSTransition, TransitionGroup } from 'react-transition-group'
import ReactDOM from 'react-dom'
import useMessage from '@/components/message'
import { SimpleArticle } from '@/lib/types'
import { searchArticle } from '@/lib/client-api'
import Empty from '@/components/Empty'
import Link from 'next/link'

const Search: React.FC = () => {
  const [ maskVisible, setMaskVisible ] = useState<boolean>(false);
  const [ bodyVisible, setBodyVisible ] = useState<boolean>(false);
  const [ portal, setPortal ] = useState<ReactPortal>();

  useEffect(() => {
    const portalElement = ReactDOM.createPortal((
      <SearchBody maskVisible={maskVisible} bodyVisible={bodyVisible} onClose={() => {
        setBodyVisible(false)
        setMaskVisible(false)
      }} />
    ), document.body)
    setPortal(portalElement)
    return () => {
      setPortal(undefined)
    }
  }, [maskVisible, bodyVisible])

  return (
    <div className="header-search">
      <button className="i-tabler-search cursor-pointer" onClick={() => {
        setMaskVisible(true);
        setTimeout(() => {
          setBodyVisible(true);
        }, 100)
      }} />
      { portal }
    </div>
  )
}

const SearchBody: React.FC<{
  maskVisible: boolean;
  bodyVisible: boolean;
  onClose: () => void;
}> = ({ maskVisible, bodyVisible, onClose }) => {
  const [ articleList, setArticleList ] = useState<SimpleArticle[] | undefined>(undefined);
  const composition = useRef<boolean>(false);
  const timeRef = useRef<number | undefined | NodeJS.Timeout>(undefined);
  const message = useMessage();
  const handleInputChange = (event: ChangeEvent<HTMLInputElement>) => {
    if (composition.current) {
      return
    }
    if (timeRef.current) {
      clearTimeout(timeRef.current)
    }
    timeRef.current = setTimeout(() => {
      (event.target.value && event.target.value.length > 0) ? handleSearch(event.target.value) : setArticleList(undefined);
    }, 1000)
  }

  const handleSearch = async (keyword: string) => {
    const loadingMessage = message.showLoading('搜索中...');
    try {
      const result = await searchArticle(keyword);
      if (result.code === 200) {
        setArticleList(result.data || []);
      }
    } finally {
      loadingMessage.close();
    }
  }

  const handleClose = () => {
    onClose();
    composition.current = false;
    setArticleList(undefined);
  }

  return (
    <CSSTransition timeout={300} in={maskVisible} classNames="mask-fade" unmountOnExit>
      <div className="fixed top-0 left-0 right-0 bottom-0 z-999 flex justify-center p-8" onClick={handleClose}>
        <CSSTransition timeout={500} in={bodyVisible} unmountOnExit classNames="search-fade">
          <div className="search-body rounded-md search-fade h-min max-w-160 w-full p-4 md:p-6 bg-white mt-12 flex flex-col gap-row-4"
               onClick={(event) => event.stopPropagation()}
          >
            <p className="md:text-xl search-title">文章搜索</p>
            <input className="search-input text-sm" placeholder="请输入待搜索的文章标题或内容..."
                   onChange={handleInputChange}
                   onCompositionStart={() => {
                     composition.current = true;
                   }}
                   onCompositionEnd={(event: CompositionEvent<HTMLInputElement>) => {
                     composition.current = false
                     handleSearch((event.target as HTMLInputElement).value)
                   }}
            />
            { articleList && articleList.length === 0 ? (
              <Empty iconSize="4rem" textSize="0.75rem" text="没有符合条件的数据" />
            ) : (
              <div className="search-result">
                <TransitionGroup component="ul" className="list-none flex flex-col gap-row-2">
                  { articleList?.map(item => (
                    <CSSTransition timeout={300} key={item.articleId} classNames="list-fade">
                      <li>
                        <h3 className="info-text" onClick={handleClose}>
                          <Link className="a-hover-line-text-sm" href={`/article/${item.articleId}`}
                                dangerouslySetInnerHTML={{ __html: item.title }}
                          />
                        </h3>
                        <p className="desc-text text-sm line-clamp-2"
                           dangerouslySetInnerHTML={{ __html: item.summary }}
                        />
                      </li>
                    </CSSTransition>
                  ))}
                </TransitionGroup>
              </div>
            )}
          </div>
        </CSSTransition>
      </div>
    </CSSTransition>
  )
}

export default Search;