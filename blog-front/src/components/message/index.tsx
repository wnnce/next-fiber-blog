'use client'

import "@/styles/components/message.scss";
import ReactDOM from 'react-dom/client'
import React, { useEffect, useRef, useState } from 'react'
import { TransitionGroup, CSSTransition } from 'react-transition-group'
import { Notice, NoticeResult, NoticeType } from '@/components/message/types'
import Message from '@/components/message/Message'

const messageContainerFlag = "message-container"

let add: (text: string, type: NoticeType) => NoticeResult = (): NoticeResult => {
  return { key: '', close: () => {}}
};

export const MessageContainer: React.FC = () => {
  const [ notices, setNotices ] = useState<Notice[]>([]);

  const countRef = useRef<number>(0);

  const generateNoticeKey = () => {
    const current = countRef.current;
    countRef.current += 1;
    return `${new Date().getTime()}${current}`;
  }

  add = (text: string, type: NoticeType): NoticeResult => {
    const key = generateNoticeKey();
    setNotices(preNotices => {
      return [...preNotices, { text, key, type }]
    })
    if (type !== 'loading') {
      setTimeout(() => {
        remove(key);
      }, 3000)
    }
    return { key: key, close: () => { remove(key) } }
  }

  const remove = (key: string) => {
    setNotices(preNotices => preNotices.filter(item => item.key != key));
  }

  useEffect(() => {
    if (notices.length > 10) {
      remove(notices[0].key)
    }
  }, [notices])

  return (
    <ul className="list-none message-list flex flex-col items-center gap-row-2">
      <TransitionGroup component={null}>
        { notices.map(item => (
          <CSSTransition timeout={300} key={item.key} classNames="message-slide">
            <li key={item.key}>
              <Message text={item.text} type={item.type} onClose={() => {
                remove(item.key)
              }} />
            </li>
          </CSSTransition>
        ))}
      </TransitionGroup>

    </ul>
  )
}

let messageContainerInit = false
const initMessageContainer = () => {
  let container = document.getElementById(messageContainerFlag);
  if (!container) {
    container = document.createElement('div')
    container.className = messageContainerFlag
    container.id = messageContainerFlag
    document.body.append(container)
  }
  ReactDOM.createRoot(container).render(<MessageContainer />)
}

const useMessage = () => {
  const showInfo = (text: string): NoticeResult => {
    return addNotice(text, 'info')
  }
  const showSuccess = (text: string): NoticeResult => {
    return addNotice(text, 'success')
  }
  const showWaring = (text: string): NoticeResult => {
    return addNotice(text, 'waring')
  }
  const showDanger = (text: string): NoticeResult => {
    return addNotice(text, 'danger')
  }
  const showLoading = (text: string): NoticeResult => {
    return addNotice(text, 'loading')
  }

  const addNotice = (text: string, type: NoticeType): NoticeResult => {
    if (!messageContainerInit && document) {
      messageContainerInit = true;
      initMessageContainer();
    }
    return add(text, type);
  }

  return { showInfo, showSuccess, showWaring, showDanger, showLoading };
}

export default useMessage;