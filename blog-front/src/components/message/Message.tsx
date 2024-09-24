import '@/styles/components/message.scss'
import React from 'react'
import { NoticeType } from '@/components/message/types'

const Message: React.FC<{
  text: string;
  type: NoticeType;
  onClose: () => void;
}> = ({ text, type, onClose }) => {

  let messageIconClassName: string = '';

  switch (type) {
    case 'loading': {
      messageIconClassName = 'i-tabler:loader-2 animate-spin text-blue-5';
      break;
    }
    case 'success': {
      messageIconClassName = 'i-tabler:circle-check-filled text-green-5';
      break;
    }
    case 'waring': {
      messageIconClassName = 'i-tabler:exclamation-circle-filled text-orange-4';
      break;
    }
    case 'danger': {
      messageIconClassName = 'i-tabler:xbox-x-filled text-red-5';
      break;
    }
    case 'info': {
      messageIconClassName = 'i-tabler:exclamation-circle-filled';
      break;
    }
  }

  return (
    <div className="message-item px-4 py-2 rounded-md flex gap-col-3 items-center">
      <i className={`inline-block text-lg shrink-0 ${messageIconClassName}`} />
      <span className="text-sm main-text">{ text }</span>
      <i className="inline-block text-lg i-tabler-x shrink-0 item-close-button" onClick={onClose} />
    </div>
  )
}

export default Message;